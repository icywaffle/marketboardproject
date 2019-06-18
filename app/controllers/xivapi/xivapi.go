package xivapi

import (
	"context"
	"fmt"
	"log"
	"marketboardproject/app/models"
	"time"

	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/controllers/xivapi/urlstring"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Current issues.
// We need to remove outliers from the price calculations.
// We have to go into the recipes, and find those too.
func NetItemPrice(recipeID int, baseprofit *models.Profits, baseinfo *models.Recipes, baseprice *models.Prices, materialprices map[int][10]int, materialingredients map[int][]int) {

	// Hold all the database info in terms of collections, so that you can manipulate it.
	itemcollection := dbconnect("Recipes")
	pricecollection := dbconnect("Prices")
	profitcollection := dbconnect("Profits")

	// Uses the Recipe and Prices struct to hold all the information from the database.
	baseinfo = finditem(itemcollection, recipeID)
	baseprice = findprices(pricecollection, baseinfo.ItemResultTargetID)

	// This can be calculated using the two above.
	materialtotal := make(map[int]int)
	/*
		materialprices shows prices for it's total materials ->	  map[14146:[1158 3057 1000 0 0 0 0 0 1050 500]]
		materialingredients shows the ingredients to make it ->   map[14146:[14147 14148 12534 0 0 0 0 0 14 17]
		materialtotal shows the total price as a sum ->			  map[14146:6765]
	*/
	// Fills the maps
	findpricesarray(itemcollection, pricecollection, baseinfo, materialprices, materialingredients)

	for itemID, pricearray := range materialprices {
		var pricesum int
		for i := 0; i < len(pricearray); i++ {
			pricesum += pricearray[i]
		}
		materialtotal[itemID] = pricesum
	}

	// Find Profit requires all the previous information from above.
	baseprofit = findprofits(profitcollection, baseinfo, baseprice, materialprices, materialingredients, materialtotal, baseinfo.ItemResultTargetID)
}

// Force updates only a single item
func UpdateItemPrices(itemID int) {
	pricecollection := dbconnect("Prices")
	prices := findprices(pricecollection, itemID) // This will also handle cases, if the item is not in the database

	// If the entries in the database is three days old, we need to actually update the prices, by forcibly going back to the API
	// If the Added is zero, then it means that it's a vendor sold price. So there's no need to update.
	now := time.Now()
	if (now.Unix()-int64(prices.Sargatanas.Prices[0].Added)) > 3*24*60*60 && prices.Sargatanas.Prices[0].Added != 0 {

		byteValue := apipriceconnect(itemID)
		// Reupdate the prices information from grabbing from the API
		prices = database.Jsonprices(byteValue)
		database.UpdatePrices(pricecollection, *prices, itemID)
	}
}

func findsum(itemID int, ingredientarray []int, materialtotal map[int]int, materialprices map[int][10]int, materialingredients map[int][]int) int {
	var tiersum int
	// Some materials are base items, so these base items won't have a map key for prices.
	temppricearray := materialprices[itemID]
	for i := 0; i < len(ingredientarray); i++ {
		materialtotalprice, ok := materialtotal[ingredientarray[i]]
		if ok {
			// If a material also has a recipe, then we want to recursively call for it's material prices.
			inneringredientarray, innerrecipe := materialingredients[ingredientarray[i]]
			if innerrecipe {
				// We we need to redefine the materialtotalprice with the one that is found by looking at the prices of the materials within the materials.
				materialtotalprice = findsum(ingredientarray[i], inneringredientarray, materialtotal, materialprices, materialingredients)
			}
			temppricearray[i] = materialtotalprice
		}

		tiersum += temppricearray[i]
	}
	return tiersum
}

func findpricesarray(itemcollection *mongo.Collection, pricecollection *mongo.Collection, baseinfo *models.Recipes, materialprices map[int][10]int, materialingredients map[int][]int) {
	var pricearray [10]int
	for i := 0; i < len(baseinfo.IngredientNames); i++ {
		// Zero is an invalid material ID
		if baseinfo.IngredientNames[i] != 0 {
			prices := findprices(pricecollection, baseinfo.IngredientNames[i])
			// The issue is for Camphorwood Branch or buyable things
			// It's in a different layout because NPC bought things are not available in the market.
			// In these cases, if null, then we go to the item, and look for PriceMid.
			pricearray[i] = prices.Sargatanas.Prices[0].PricePerUnit * baseinfo.IngredientAmounts[i]
		} else {
			continue
		}

	}
	materialprices[baseinfo.ItemResultTargetID] = pricearray
	materialingredients[baseinfo.ItemResultTargetID] = baseinfo.IngredientNames

	// If there's a recipe, we want to go in one more materialprices, and keep appending to it.
	for i := 0; i < len(baseinfo.IngredientRecipes); i++ {
		if len(baseinfo.IngredientRecipes[i]) != 0 {
			matinfo := finditem(itemcollection, baseinfo.IngredientRecipes[i][0])
			findpricesarray(itemcollection, pricecollection, matinfo, materialprices, materialingredients)
		}
	}

}

func finditem(itemcollection *mongo.Collection, recipeID int) *models.Recipes {
	// itemresult is the info in the recipeID
	itemresult := database.Ingredientmaterials(itemcollection, recipeID)
	// If the item is not in the database, then we should add it. 0 is an invalid itemID
	if itemresult.ID == 0 {
		byteValue := apirecipeconnect(recipeID)
		// TODO : create a json struct that has all these variables.
		recipes, matIDs, amounts, matrecipes := database.Jsonitemrecipe(byteValue)
		database.InsertRecipe(itemcollection, *recipes, matIDs, amounts, matrecipes)

		itemresult = database.Ingredientmaterials(itemcollection, recipeID)
	}

	return itemresult
}

// Handles finding and updating the models.Prices
func findprices(pricecollection *mongo.Collection, itemID int) *models.Prices {
	// The find the price of the ingredient itself.
	priceresult := database.Ingredientprices(pricecollection, itemID)
	// TODO : Fix this into the Ingredientprices function instead.
	if priceresult.ItemID == 0 {
		byteValue := apipriceconnect(itemID)
		// Connects to the API and takes the market listed price
		prices := database.Jsonprices(byteValue)
		// If there is no market listed price, then it must mean that there's a vendor selling it.
		if len(prices.Sargatanas.History) == 0 && len(prices.Sargatanas.Prices) == 0 {
			// This information comes from the item page. Let's unmarshal the vendor price into the price struct.
			byteValue = apiitemconnect(itemID)
			prices = database.Jsonprices(byteValue)
		}
		database.InsertPrices(pricecollection, *prices, itemID)

		priceresult = database.Ingredientprices(pricecollection, itemID)
	}

	// If the entries in the database is seven days old, we need to actually update the prices, by forcibly going back to the API
	// If the Added is zero, then it means that it's a vendor sold price. So there's no need to update.
	now := time.Now()
	if (now.Unix()-int64(priceresult.Sargatanas.Prices[0].Added)) > 7*24*60*60 && priceresult.Sargatanas.Prices[0].Added != 0 {

		byteValue := apipriceconnect(itemID)
		// Reupdate the prices information from grabbing from the API
		priceresult = database.Jsonprices(byteValue)
		database.UpdatePrices(pricecollection, *priceresult, itemID)
	}

	return priceresult
}

// Handles updates, obtaining, creating information about the baseprofits from models.Profits
func findprofits(profitcollection *mongo.Collection, baseinfo *models.Recipes, baseprice *models.Prices, materialprices map[int][10]int, materialingredients map[int][]int, materialtotal map[int]int, itemID int) *models.Profits {
	baseprofit := database.Ingredientprofits(profitcollection, itemID)
	if baseprofit.ItemID == 0 {
		fillbaseprofits(baseprofit, profitcollection, baseinfo, baseprice, materialprices, materialingredients, materialtotal, itemID)
		database.InsertProfits(profitcollection, *baseprofit, itemID)
	}

	// If the entries in the database is more than seven days old, we need to recalculate and reinput into the database.
	now := time.Now()
	if (now.Unix()-int64(baseprofit.Added)) > 7*24*60*60 && baseprofit.Added != 0 {
		fillbaseprofits(baseprofit, profitcollection, baseinfo, baseprice, materialprices, materialingredients, materialtotal, itemID)
		database.UpdateProfits(profitcollection, *baseprofit, itemID)
	}

	return baseprofit
}

// Calculates and fills baseprofits with information from the models.Profits
func fillbaseprofits(baseprofit *models.Profits, profitcollection *mongo.Collection, baseinfo *models.Recipes, baseprice *models.Prices, materialprices map[int][10]int, materialingredients map[int][]int, materialtotal map[int]int, itemID int) {
	baseprofit.ItemID = baseinfo.ItemResultTargetID
	baseprofit.RecipeID = baseinfo.ID
	baseprofit.MarketboardPrice = baseprice.Sargatanas.Prices[0].PricePerUnit

	// Assuming here, that the base materials will always be cheaper.
	// We can analyze this more in the future.
	materialcosts := findsum(baseinfo.ItemResultTargetID, baseinfo.IngredientNames, materialtotal, materialprices, materialingredients)
	baseprofit.MaterialCosts = materialcosts

	baseprofit.Profits = baseprice.Sargatanas.Prices[0].PricePerUnit - materialcosts
	baseprofit.ProfitPercentage = (float32(baseprice.Sargatanas.Prices[0].PricePerUnit) - float32(materialcosts)) / float32(materialcosts) * 100

	now := time.Now()
	baseprofit.Added = now.Unix()

}

func apirecipeconnect(recipeID int) []byte {

	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	// This ensures that when this function is called, it does not exceed the rate limit.
	// TODO: Use a channel to rate limit instead to allow multiple users to use this.

	websiteurl := urlstring.UrlItemRecipe(recipeID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	fmt.Println("Connected to API")
	return byteValue
}

func apipriceconnect(itemID int) []byte {

	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	// This ensures that when this function is called, it does not exceed the rate limit.
	// TODO: Use a channel to rate limit instead to allow multiple users to use this.

	websiteurl := urlstring.UrlPrices(itemID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	fmt.Println("Connected to API")
	return byteValue
}

func apiitemconnect(itemID int) []byte {
	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	// This ensures that when this function is called, it does not exceed the rate limit.
	// TODO: Use a channel to rate limit instead to allow multiple users to use this.

	websiteurl := urlstring.UrlItem(itemID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	return byteValue
}

// Collection names are either "Prices" or "Recipes"
func dbconnect(collectionname string) *mongo.Collection {
	// Apply the user string mongoURI to be able to connect to the database
	// In this case, since this is backend, only the server should be allowed to connect here.
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("Marketboard").Collection(collectionname)

	return collection
}
