package xivapi

import (
	"context"
	"fmt"
	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Prices  *mongo.Collection
	Recipes *mongo.Collection
	Profits *mongo.Collection
}

type Information struct {
	Recipes       *models.Recipes
	Prices        *models.Prices
	Profits       *models.Profits
	Matprofitmaps *models.Matprofitmaps
}

// Fills out information about the base item's Prices and Profits.
func (coll *Collections) BaseInformation(recipeID int) *Information {

	var info Information

	// Each info depends on the previous one.
	info.Recipes = coll.findrecipe(recipeID)
	info.Prices = coll.findprices(info.Recipes.ItemResultTargetID)
	info.Profits = coll.findprofits(&info)

	return &info
}

// Calls ingredient amounts and item IDs, and returns the results
func (coll *Collections) queryrecipes(recipeID int) *models.Recipes {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Recipes
	coll.Recipes.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// Call the prices from the database, and return the sold average and the current average
func (coll *Collections) queryprices(itemID int) *models.Prices {
	filter := bson.M{"ItemID": itemID}
	var result models.Prices
	coll.Prices.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

func (coll *Collections) queryprofits(recipeID int) *models.Profits {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Profits
	coll.Profits.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// Will access the API if there's no recipe in the database.
func (coll *Collections) findrecipe(recipeID int) *models.Recipes {
	// Only info.Recipes is filled.
	recipeinfo := coll.queryrecipes(recipeID)
	// If the item is not in the database, then we should add it. 0 is an invalid itemID
	if recipeinfo.ID == 0 {
		byteValue := database.ApiConnect(recipeID, "recipe")
		recipeinfo = database.Jsonitemrecipe(byteValue)

		database.InsertRecipe(coll.Recipes, *recipeinfo)

	}

	return recipeinfo
}

// Will access API if there's no prices in the database. Will also auto update.
func (coll *Collections) findprices(itemID int) *models.Prices {
	// Only info.Prices is filled
	pricesinfo := coll.queryprices(itemID)
	// If the item is not in the database, then we should add it. 0 is an invalid itemID
	if pricesinfo.ItemID == 0 {
		byteValue := database.ApiConnect(itemID, "market/item")
		pricesinfo = database.Jsonprices(byteValue)
		// If there is no market listed price, then it must mean that there's a vendor selling it.
		if len(pricesinfo.Sargatanas.History) == 0 && len(pricesinfo.Sargatanas.Prices) == 0 {
			// This information comes from the item page. Let's unmarshal the vendor price into the price struct.
			byteValue = database.ApiConnect(itemID, "item")
			pricesinfo = database.Jsonprices(byteValue)
		}
		database.InsertPrices(coll.Prices, *pricesinfo, itemID)
	}

	// If the entries in the database is seven days old, we need to actually update the prices, by forcibly going back to the API
	now := time.Now()
	if len(pricesinfo.Sargatanas.Prices) > 0 {
		// There's no need to update if we're looking at the VendorPrice.
		if pricesinfo.VendorPrice == 0 {
			// Marketboard could be empty, and we can try to check again.
			if len(pricesinfo.Sargatanas.Prices) == 0 || (now.Unix()-int64(pricesinfo.Sargatanas.Prices[0].Added)) > 7*24*60*60 {
				byteValue := database.ApiConnect(itemID, "market/item")
				pricesinfo = database.Jsonprices(byteValue)
				fmt.Println("Entry is Seven Days Old.")
				database.UpdatePrices(coll.Prices, *pricesinfo, itemID)
			}
		}

	}

	return pricesinfo
}

// Handles updates, obtaining, creating information about the baseprofits from models.Profits
func (coll *Collections) findprofits(info *Information) *models.Profits {
	// Inside the database
	profitsinfo := coll.queryprofits(info.Recipes.ID)
	info.Matprofitmaps = coll.fillprofitmaps(info.Recipes)
	// If not inside the database, call the database, and fill it up with information
	if profitsinfo.ItemID == 0 {
		profitsinfo = coll.fillbaseprofits(info, info.Matprofitmaps)
		database.InsertProfits(coll.Profits, *profitsinfo)
	}

	// If the entries in the database is more than seven days old, we need to recalculate and reinput into the database.
	now := time.Now()
	if (now.Unix()-int64(profitsinfo.Added)) > 7*24*60*60 && profitsinfo.Added != 0 {
		coll.fillbaseprofits(info, info.Matprofitmaps)
		database.UpdateProfits(coll.Profits, *profitsinfo, info.Recipes.ID)
	}

	return profitsinfo
}

// Fills the Profits struct
func (coll *Collections) fillbaseprofits(info *Information, matprofitmaps *models.Matprofitmaps) *models.Profits {

	var profits models.Profits
	profits.ItemID = info.Recipes.ItemResultTargetID
	profits.RecipeID = info.Recipes.ID

	// Some items won't have a market history, because they're from vendors.
	var itempriceperunit int
	if info.Prices.VendorPrice != 0 {
		itempriceperunit = info.Prices.VendorPrice
	} else {
		itempriceperunit = info.Prices.Sargatanas.Prices[0].PricePerUnit
	}
	profits.MarketboardPrice = itempriceperunit

	// Assuming here, that the base materials will always be cheaper.
	// We can analyze this more in the future.
	materialcosts := coll.findsum(info, matprofitmaps)
	profits.MaterialCosts = materialcosts

	profits.Profits = itempriceperunit - materialcosts
	profits.ProfitPercentage = (float32(itempriceperunit) - float32(materialcosts)) / float32(materialcosts) * 100

	now := time.Now()
	profits.Added = now.Unix()

	return &profits
}

// Uses the price collection from the database to fill the individual material maps.
func (coll *Collections) fillprofitmaps(baserecipe *models.Recipes) *models.Matprofitmaps {
	var pricearray [10]int
	var matprofitmaps models.Matprofitmaps
	matprofitmaps.Costs = make(map[int][10]int)
	matprofitmaps.Ingredients = make(map[int][]int)
	matprofitmaps.Total = make(map[int]int)
	for i := 0; i < len(baserecipe.IngredientNames); i++ {
		// Zero is an invalid material ID
		if baserecipe.IngredientNames[i] != 0 {
			matpriceinfo := coll.findprices(baserecipe.IngredientNames[i])
			// We need to deal with vendor prices, since they won't have market prices.
			if matpriceinfo.VendorPrice != 0 {
				pricearray[i] = matpriceinfo.VendorPrice * baserecipe.IngredientAmounts[i]
			} else {
				if len(matpriceinfo.Sargatanas.Prices) > 0 {
					pricearray[i] = matpriceinfo.Sargatanas.Prices[0].PricePerUnit * baserecipe.IngredientAmounts[i]
				}

			}

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
			matinfo := coll.findrecipe(baserecipe.IngredientRecipes[i][0])
			coll.fillprofitmaps(matinfo)
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

func (coll *Collections) findsum(info *Information, matprofitmaps *models.Matprofitmaps) int {
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
				materialinfo.Recipes = coll.findrecipe(info.Recipes.IngredientRecipes[i][0])
				materialinfo.Prices = coll.findprices(info.Recipes.IngredientNames[i])
				// We we need to redefine the materialtotalprice with the one that is found by looking at the prices of the materials within the materials.
				materialtotalprice = coll.findsum(materialinfo, matprofitmaps)
			}
			temppricearray[i] = materialtotalprice
		}

		tiersum += temppricearray[i]
	}
	return tiersum
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

*/
