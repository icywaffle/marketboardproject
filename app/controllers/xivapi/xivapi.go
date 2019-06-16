package xivapi

import (
	"context"
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
func NetItemPrice(recipeID int, results *models.Result) {

	// Hold all the database info in terms of collections, so that you can manipulate it.
	itemcollection := dbconnect("Recipes")
	pricecollection := dbconnect("Prices")

	// Uses the Recipe and Prices struct to hold all the information from the database.
	baseinfo := finditem(itemcollection, recipeID)
	baseprice := findprices(pricecollection, baseinfo.ItemResultTargetID)

	//These two should be sent to the front end.
	materialprices := make(map[int][10]int) // Already multiplied by the ingredient amounts.
	materialingredients := make(map[int][]int)
	// This can be calculated using the two above.
	materialtotal := make(map[int]int)
	/*
		materialprices shows prices for it's total materials ->	  map[14146:[1158 3057 1000 0 0 0 0 0 1050 500]]
		materialingredients shows the ingredients to make it ->   map[14146:[14147 14148 12534 0 0 0 0 0 14 17]
		materialtotal shows the total price as a sum ->			  map[14146:6765]
	*/

	// Fills the maps
	findpricesarray(itemcollection, pricecollection, baseinfo, materialprices, materialingredients)

	// All these calculations below can be done in the front end javascript. This is here in the backend for reference.
	for itemID, pricearray := range materialprices {
		var pricesum int
		for i := 0; i < len(pricearray); i++ {
			pricesum += pricearray[i]
		}
		materialtotal[itemID] = pricesum
	}
	_, basecurrent := avgprices(baseinfo.ItemResultTargetID, 1, baseprice)
	results.MarketboardPrice = basecurrent

	// Assuming here, that the base materials will always be cheaper.
	// We can analyze this more in the future.
	materialcosts := findsum(baseinfo.ItemResultTargetID, baseinfo.IngredientNames, materialtotal, materialprices, materialingredients)
	results.MaterialCosts = materialcosts

	results.Profits = basecurrent - materialcosts
	results.ProfitPercentage = (float32(basecurrent) - float32(materialcosts)) / float32(materialcosts) * 100

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

func findpricesarray(itemcollection *mongo.Collection, pricecollection *mongo.Collection, baseinfo *database.Recipes, materialprices map[int][10]int, materialingredients map[int][]int) {
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

func finditem(itemcollection *mongo.Collection, recipeID int) *database.Recipes {
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
func findprices(pricecollection *mongo.Collection, itemID int) *database.Prices {
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

	return priceresult
}

func apirecipeconnect(recipeID int) []byte {
	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	// This ensures that when this function is called, it does not exceed the rate limit.
	// TODO: Use a channel to rate limit instead to allow multiple users to use this.

	websiteurl := urlstring.UrlItemRecipe(recipeID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	return byteValue
}

func apipriceconnect(itemID int) []byte {
	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	// This ensures that when this function is called, it does not exceed the rate limit.
	// TODO: Use a channel to rate limit instead to allow multiple users to use this.

	websiteurl := urlstring.UrlPrices(itemID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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

func avgprices(ingredient int, ingredientamount int, matprices *database.Prices) (int, int) {

	// Average Price History for the latest 20 entries.
	var hissum int
	for i := 0; i < len(matprices.Sargatanas.History) && i < 2; i++ {
		hissum = hissum + matprices.Sargatanas.History[i].PricePerUnit
	}

	soldaverage := hissum / 2

	// Average Price Listings for the latest 20 entries.

	var listsum int
	for i := 0; i < len(matprices.Sargatanas.Prices) && i < 2; i++ {
		listsum = listsum + matprices.Sargatanas.Prices[i].PricePerUnit
	}

	currentaverage := listsum / 2

	// Multiply by the ingredient amount.
	averagesoldcost := soldaverage * ingredientamount
	averagecurrentcost := currentaverage * ingredientamount

	return averagesoldcost, averagecurrentcost
}
