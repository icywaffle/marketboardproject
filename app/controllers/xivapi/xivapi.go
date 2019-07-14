package xivapi

import (
	"context"
	"fmt"
	"math"
	"time"

	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// We want to separate the times, just in case we only update one struct.
// Changing these times will allow us to just update our entries accordingly if the
// structs have been changed.
var UpdatedRecipesStructTime = int64(1563090742) // Last Update : 6/26/19 - 11AM
var UpdatedPricesStructTime = int64(1561493761)
var UpdatedProfitsStructTime = int64(1562657644) // Last Update : 7/9/19 - 9PM

type Collections struct {
	Prices  *mongo.Collection
	Recipes *mongo.Collection
	Profits *mongo.Collection
}

type Information struct {
	Recipes *models.Recipes
	Prices  *models.Prices
	Profits *models.Profits
}

type InnerInformation struct {
	Recipes      map[int]*models.Recipes      // Contains the inner recipes for some key = Recipe.ID
	SimplePrices map[int]*models.SimplePrices // Contains the inner prices for some key =  Item ID
	Profits      map[int]*models.Profits      // Contains the profits for the inner recipes for some key = Recipe.Id
}

type CollectionHandler interface {
	FindRecipesDocument(recipeID int) *models.Recipes
	FindPricesDocument(itemID int) *models.Prices
	FindProfitsDocument(recipeID int) *models.Profits
	SimplifyPricesDocument(recipeID int) *models.SimplePrices
	InsertRecipesDocument(recipeID int) *models.Recipes
	InsertPricesDocument(itemID int) *models.Prices
	InsertProfitsDocument(profits *models.Profits)
}

type ProfitHandler interface {
	ProfitDescCursor() []*models.Profits
}

func (coll Collections) FindRecipesDocument(recipeID int) *models.Recipes {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Recipes
	coll.Recipes.FindOne(context.TODO(), filter).Decode(&result)
	return &result
}

func (coll Collections) FindPricesDocument(itemID int) *models.Prices {
	filter := bson.M{"ItemID": itemID}
	var result models.Prices
	coll.Prices.FindOne(context.TODO(), filter).Decode(&result)
	return &result
}

// This puts less strain if we don't need history.
func (coll Collections) SimplifyPricesDocument(itemID int) *models.SimplePrices {
	filter := bson.M{"ItemID": itemID}
	var fullprices models.Prices
	coll.Prices.FindOne(context.TODO(), filter).Decode(&fullprices)
	var result models.SimplePrices
	result.ItemID = fullprices.ItemID
	result.LowestMarketPrice = fullprices.Sargatanas.Prices[0].PricePerUnit
	result.VendorPrice = fullprices.VendorPrice
	result.Added = fullprices.Added
	return &result
}

func (coll Collections) FindProfitsDocument(recipeID int) *models.Profits {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Profits
	coll.Profits.FindOne(context.TODO(), filter).Decode(&result)
	return &result
}

// Will insert a document, or update it if it's already in the collection.
func (coll Collections) InsertRecipesDocument(recipeID int) *models.Recipes {
	byteValue := database.ApiConnect(recipeID, "recipe")
	result := database.Jsonitemrecipe(byteValue)
	// These variables are not in the json file.
	now := time.Now()
	result.Added = now.Unix()
	// Testing if there's an entry in the DB
	filter := bson.M{"RecipeID": recipeID}

	var options options.CountOptions
	options.SetLimit(1)
	findcount, _ := coll.Recipes.CountDocuments(context.TODO(), filter, &options)
	if findcount < 1 {
		coll.Recipes.InsertOne(context.TODO(), result)
		fmt.Println("Inserted Recipe into Database: ", result.ID)
	} else {
		coll.Recipes.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: result},
		})
		fmt.Println("Updated Item into Recipe Collection :", result.ID)
	}

	return &result
}

func (coll Collections) InsertPricesDocument(itemID int) *models.Prices {
	byteValue := database.ApiConnect(itemID, "market/item")
	result := database.Jsonprices(byteValue)
	// ItemID is not part of the Json file.
	result.ItemID = itemID
	// If we do have an empty result, it means that we need to search for the vendor prices.
	if result.ItemID != 0 && len(result.Sargatanas.Prices) == 0 {
		byteValue = database.ApiConnect(result.ItemID, "item")
		result = database.Jsonprices(byteValue)
		// We have to rewrite the Item ID into this
		result.ItemID = itemID
	}
	//These variables are not in the json file.
	now := time.Now()
	result.Added = now.Unix()

	filter := bson.M{"ItemID": itemID}

	var options options.CountOptions
	options.SetLimit(1)
	findcount, _ := coll.Prices.CountDocuments(context.TODO(), filter, &options)
	if findcount < 1 {
		coll.Prices.InsertOne(context.TODO(), result)
		fmt.Println("Inserted Prices into Database: ", result.ItemID)
	} else {
		coll.Prices.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: result},
		})
		fmt.Println("Updated Item into Prices Collection :", result.ItemID)
	}

	return &result
}

// Uses the Recipes and Profits from Information, and returns a Profit model.
func (info Information) FillProfitsDocument(recipeID int) *models.Profits {
	var profits models.Profits
	profits.Name = info.Recipes.Name
	profits.IconID = info.Recipes.IconID
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

	// We're only looking at profits for a RECIPE. So it must have some Ingredients.
	for i := 0; i < len(info.Recipes.IngredientID); i++ {
		// Check in the map if there are recipes first.
		// If there are recipes in the map, FillProfitsDocument(that new recipeID)
	}
	// We may get multiple items per craft.
	profits.Profits = itempriceperunit*info.Recipes.AmountResult - materialcosts
	profitmaterialratio := (float64(profits.Profits) / float64(materialcosts)) //0.01
	profits.ProfitPercentage = float32(math.Ceil(profitmaterialratio*10000) / 100)

	now := time.Now()
	profits.Added = now.Unix()

	return &profits

}
func (coll Collections) InsertProfitsDocument(profits *models.Profits) {
	filter := bson.M{"RecipeID": profits.RecipeID}

	var options options.CountOptions
	options.SetLimit(1)
	findcount, _ := coll.Profits.CountDocuments(context.TODO(), filter, &options)
	if findcount < 1 {
		coll.Profits.InsertOne(context.TODO(), profits)
		fmt.Println("Inserted Profits into Database: ", profits.RecipeID)
	} else {
		coll.Profits.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: profits},
		})
		fmt.Println("Updated Item into Profit Collection :", profits.RecipeID)
	}

}

// Gives a Descending Sorted Array, of 20 items with the most profit from the DB
func (coll Collections) ProfitDescCursor() []*models.Profits {
	options := options.FindOptions{}
	options.Sort = bson.D{{Key: "ProfitPercentage", Value: -1}}
	// We can set this to be bigger later in the future
	limit := int64(20)
	options.Limit = &limit
	cursor, _ := coll.Profits.Find(context.Background(), bson.D{}, &options)

	var allprofits []*models.Profits
	for cursor.Next(context.TODO()) {
		var tempprofits models.Profits
		cursor.Decode(&tempprofits)

		allprofits = append(allprofits, &tempprofits)
	}
	defer cursor.Close(context.TODO())

	return allprofits

}

// Fills the entire map full of related recipes.
func BaseInformation(collections CollectionHandler, recipeID int, innerinfo InnerInformation) {

	// Adds a base recipe to the map
	baserecipe := collections.FindRecipesDocument(recipeID)
	innerinfo.Recipes[recipeID] = baserecipe

	// Finds the prices for the base item of a recipe.
	innerinfo.SimplePrices[baserecipe.ItemResultTargetID] = collections.SimplifyPricesDocument(baserecipe.ItemResultTargetID)
	// Also adds all the ingredients prices of current recipe into the map.
	for i := 0; i < len(baserecipe.IngredientID); i++ {
		innerinfo.SimplePrices[baserecipe.IngredientID[i]] = collections.SimplifyPricesDocument(baserecipe.IngredientID[i])
	}

	// Recursively call using the inner recipes (if they exist), to add more recipes and prices to our map
	for i := 0; i < len(baserecipe.IngredientRecipes); i++ {
		if baserecipe.IngredientRecipes[i] != nil {
			BaseInformation(collections, baserecipe.IngredientRecipes[i][0], innerinfo)
		}
	}

	// Profits require the recipes and prices maps to be completely filled first.
	innerinfo.Profits[recipeID] = collections.FindProfitsDocument(recipeID)

}

func ProfitInformation(profit ProfitHandler) []*models.Profits {

	return profit.ProfitDescCursor()
}
