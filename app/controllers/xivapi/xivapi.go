package xivapi

import (
	"fmt"
	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/controllers/xivapi/urlstring"
	"marketboardproject/app/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Prices  *mongo.Collection
	Recipes *mongo.Collection
	Profits *mongo.Collection
}

// Fills out information about the base item's Prices and Profits.
func (coll *Collections) BaseInformation(recipeID int) (*Information, *models.Matprofitmaps) {

	var info Information

	// We need to try and separate the dependencies here.
	// We need to look through these inner functions to try and see if we can separate them here.
	info.Recipes = finditem(coll.Recipes, recipeID)

	// You can already see that this function depends on the function above it.
	// The dependencies here are too closely linked together.
	info.Prices = findprices(coll.Prices, info.Recipes.ItemResultTargetID)

	return &info, fillprofitmaps(coll, info.Recipes)
}

type Information struct {
	Recipes *models.Recipes
	Prices  *models.Prices
	Profits *models.Profits
}

func (filled *Information) ProfitInformation(coll *Collections, matprofitmaps *models.Matprofitmaps) *models.Profits {

	// Find Profit requires all the previous information from above.
	baseprofit := findprofits(coll, filled, matprofitmaps)

	return baseprofit
}

func finditem(itemcollection *mongo.Collection, recipeID int) *models.Recipes {
	// itemresult is the info in the recipeID
	itemresult := database.Ingredientmaterials(itemcollection, recipeID)
	// If the item is not in the database, then we should add it. 0 is an invalid itemID
	if itemresult.ID == 0 {
		byteValue := apirecipeconnect(recipeID)
		// TODO : create a json struct that has all these variables.
		recipes := database.Jsonitemrecipe(byteValue)

		database.InsertRecipe(itemcollection, *recipes)

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
	if len(priceresult.Sargatanas.Prices) > 0 {
		if (now.Unix()-int64(priceresult.Sargatanas.Prices[0].Added)) > 7*24*60*60 && priceresult.Sargatanas.Prices[0].Added != 0 {

			byteValue := apipriceconnect(itemID)
			// Reupdate the prices information from grabbing from the API
			priceresult = database.Jsonprices(byteValue)
			database.UpdatePrices(pricecollection, *priceresult, itemID)
		}

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
	fmt.Println("Connected to API")
	return byteValue
}

// Uses the price collection from the database to fill the individual material maps.
func fillprofitmaps(coll *Collections, baserecipe *models.Recipes) *models.Matprofitmaps {
	var pricearray [10]int
	var matprofitmaps models.Matprofitmaps
	matprofitmaps.Costs = make(map[int][10]int)
	matprofitmaps.Ingredients = make(map[int][]int)
	matprofitmaps.Total = make(map[int]int)
	for i := 0; i < len(baserecipe.IngredientNames); i++ {
		// Zero is an invalid material ID
		if baserecipe.IngredientNames[i] != 0 {
			matpriceinfo := findprices(coll.Prices, baserecipe.IngredientNames[i])
			// The issue is for Camphorwood Branch or buyable things
			// It's in a different layout because NPC bought things are not available in the market.
			// In these cases, if null, then we go to the item, and look for PriceMid.
			pricearray[i] = matpriceinfo.Sargatanas.Prices[0].PricePerUnit * baserecipe.IngredientAmounts[i]
		} else {
			continue
		}

	}
	// After receiving the price information for specific materials, we fill the maps here.
	matprofitmaps.Costs[baserecipe.ItemResultTargetID] = pricearray
	matprofitmaps.Ingredients[baserecipe.ItemResultTargetID] = baserecipe.IngredientNames

	// If there's a recipe, we want to go in one more materialprices, and keep appending to it.
	for i := 0; i < len(baserecipe.IngredientRecipes); i++ {
		if len(baserecipe.IngredientRecipes[i]) != 0 {
			matinfo := finditem(coll.Recipes, baserecipe.IngredientRecipes[i][0])
			fillprofitmaps(coll, matinfo)
		}
	}

	for itemID, pricearray := range matprofitmaps.Costs {
		var pricesum int
		for i := 0; i < len(pricearray); i++ {
			pricesum += pricearray[i]
		}
		matprofitmaps.Total[itemID] = pricesum
	}

	return &matprofitmaps
}

// Handles updates, obtaining, creating information about the baseprofits from models.Profits
// Matprofitmaps
// Information
// Profit Collection.
func findprofits(coll *Collections, info *Information, matprofitmaps *models.Matprofitmaps) *models.Profits {
	// Inside the database
	baseprofit := database.Ingredientprofits(coll.Profits, info.Recipes.ID)
	// If not inside the database, call the database, and fill it up with information
	if baseprofit.ItemID == 0 {
		fillbaseprofits(coll, info, matprofitmaps)
		database.InsertProfits(coll.Profits, *baseprofit)
	}

	// If the entries in the database is more than seven days old, we need to recalculate and reinput into the database.
	now := time.Now()
	if (now.Unix()-int64(baseprofit.Added)) > 7*24*60*60 && baseprofit.Added != 0 {
		fillbaseprofits(coll, info, matprofitmaps)
		database.UpdateProfits(coll.Profits, *baseprofit, info.Recipes.ID)
	}

	return baseprofit
}

// Fills the Profits struct
func fillbaseprofits(coll *Collections, info *Information, matprofitmaps *models.Matprofitmaps) {
	info.Profits.ItemID = info.Recipes.ItemResultTargetID
	info.Profits.RecipeID = info.Recipes.ID
	info.Profits.MarketboardPrice = info.Prices.Sargatanas.Prices[0].PricePerUnit

	// Assuming here, that the base materials will always be cheaper.
	// We can analyze this more in the future.
	materialcosts := findsum(coll, info, matprofitmaps)
	info.Profits.MaterialCosts = materialcosts

	info.Profits.Profits = info.Prices.Sargatanas.Prices[0].PricePerUnit - materialcosts
	info.Profits.ProfitPercentage = (float32(info.Prices.Sargatanas.Prices[0].PricePerUnit) - float32(materialcosts)) / float32(materialcosts) * 100

	now := time.Now()
	info.Profits.Added = now.Unix()

}

func findsum(coll *Collections, info *Information, matprofitmaps *models.Matprofitmaps) int {
	var tiersum int
	// Some materials are base items, so these base items won't have a map key for prices.
	temppricearray := matprofitmaps.Costs[info.Recipes.ItemResultTargetID]
	for i := 0; i < len(info.Recipes.IngredientNames); i++ {
		materialtotalprice, ok := matprofitmaps.Total[info.Recipes.IngredientNames[i]]
		if ok {
			// If a material also has a recipe, then we want to recursively call for it's material prices.
			_, innerrecipe := matprofitmaps.Ingredients[info.Recipes.IngredientNames[i]]
			if innerrecipe {
				var materialinfo *Information
				// We're going to need to find the information about the inner materials.
				// For now, we will deal with just the first recipe.
				materialinfo.Recipes = finditem(coll.Recipes, info.Recipes.IngredientRecipes[i][0])
				materialinfo.Prices = findprices(coll.Prices, info.Recipes.IngredientNames[i])
				// We we need to redefine the materialtotalprice with the one that is found by looking at the prices of the materials within the materials.
				materialtotalprice = findsum(coll, materialinfo, matprofitmaps)
			}
			temppricearray[i] = materialtotalprice
		}

		tiersum += temppricearray[i]
	}
	return tiersum
}

func apipriceconnect(itemID int) []byte {

	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	websiteurl := urlstring.UrlPrices(itemID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	fmt.Println("Connected to API")
	return byteValue
}
func apiitemconnect(itemID int) []byte {
	// MAX Rate limit is 20 Req/s -> 0.05s/Req, but safer to use 15req/s -> 0.06s/req
	time.Sleep(100 * time.Millisecond)
	websiteurl := urlstring.UrlItem(itemID)
	byteValue := urlstring.XiviapiRecipeConnector(websiteurl)
	fmt.Println("Connected to API")
	return byteValue
}

/*
// Force database to update the entries of prices
func ForceUpdateItemPrices(itemID int) {
	pricecollection := dbconnect("Prices")

	// Connects to the API and takes the market listed price
	byteValue := apipriceconnect(itemID)
	prices := database.Jsonprices(byteValue)
	// If there is no market listed price, then it must mean that there's a vendor selling it.
	if len(prices.Sargatanas.History) == 0 && len(prices.Sargatanas.Prices) == 0 {
		// This information comes from the item page. Let's unmarshal the vendor price into the price struct.
		byteValue = apiitemconnect(itemID)
		prices = database.Jsonprices(byteValue)
	}

	database.UpdatePrices(pricecollection, *prices, itemID)
}

// Force database to update the entries of Profits
func ForceUpdateProfits(recipeID int) {
	profitcollection := dbconnect("Profits")
	materialprices := make(map[int][10]int)
	materialingredients := make(map[int][]int)
	materialtotal := make(map[int]int)
	baseprofit, baseinfo, baseprice := NetItemPrice(recipeID, materialprices, materialingredients, materialtotal)
	fillbaseprofits(baseprofit, profitcollection, baseinfo, baseprice, materialprices, materialingredients, materialtotal)
	database.UpdateProfits(profitcollection, *baseprofit, baseinfo.ID)
}

func ForceUpdateRecipes(recipeID int) {
	itemcollection := dbconnect("Recipes")

	byteValue := apirecipeconnect(recipeID)
	recipes := database.Jsonitemrecipe(byteValue)

	database.UpdateRecipes(itemcollection, *recipes)

}

func CompareProfits() []*models.Profits {
	profitcollection := dbconnect("Profits")

	// Close the cursor once you've iterated through it.
	profitdocuments := database.Profitcomparisons(profitcollection)

	// Append all of the cursor elements to the allprofits array of information
	var allprofits []*models.Profits
	for profitdocuments.Next(context.TODO()) {
		var tempprofits models.Profits
		profitdocuments.Decode(&tempprofits)

		allprofits = append(allprofits, &tempprofits)
	}
	defer profitdocuments.Close(context.TODO())

	return allprofits
}




// Calculates and fills baseprofits with information from the models.Profits



*/
