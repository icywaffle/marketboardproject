package xivapi

import (
	"context"
	"fmt"
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
	*InnerInformation
}

type InnerInformation struct {
	InnerRecipes      map[int]*models.Recipes      // Contains the inner recipes for some key = Recipe.ID
	InnerSimplePrices map[int]*models.SimplePrices // Contains the inner prices for some key =  Item ID
	InnerProfits      map[int]*models.Profits      // Contains the profits for the inner recipes for some key = Recipe.Id
}

type CollectionHandler interface {
	FindRecipesDocument(recipeID int) (*models.Recipes, bool)
	FindPricesDocument(itemID int) (*models.Prices, bool)
	FindProfitsDocument(recipeID int) (*models.Profits, bool)
	SimplifyPricesDocument(recipeID int) (*models.SimplePrices, bool)
	InsertRecipesDocument(recipeID int) *models.Recipes
	InsertPricesDocument(itemID int) *models.Prices
	InsertProfitsDocument(profits *models.Profits)
}

type ProfitHandler interface {
	ProfitDescCursor() []*models.Profits
}

// Will return false if there's no recipe in the database.
func (coll Collections) FindRecipesDocument(recipeID int) (*models.Recipes, bool) {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Recipes
	err := coll.Recipes.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (coll Collections) FindPricesDocument(itemID int) (*models.Prices, bool) {
	filter := bson.M{"ItemID": itemID}
	var result models.Prices
	err := coll.Prices.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

// This is used for quick and easy prices
func (coll Collections) SimplifyPricesDocument(itemID int) (*models.SimplePrices, bool) {
	filter := bson.M{"ItemID": itemID}
	var fullprices models.Prices
	err := coll.Prices.FindOne(context.TODO(), filter).Decode(&fullprices)
	if err != nil {
		return nil, false
	}
	// We should only just copy the value so that the fullprices will disappear
	var result models.SimplePrices
	result.ItemID = fullprices.ItemID
	result.HistoryPrice = fullprices.Sargatanas.History[0].PricePerUnit
	result.LowestMarketPrice = fullprices.Sargatanas.Prices[0].PricePerUnit
	result.OnMarketboard = fullprices.OnMarketboard
	result.Added = fullprices.Added
	return &result, true
}

// This is used for stronger analysis functions, where we want to see trends etc.
func (coll Collections) FindProfitsDocument(recipeID int) (*models.Profits, bool) {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Profits
	err := coll.Profits.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
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

	// If we do have an empty result, it means it's not on the board.
	// XIVAPI may change this though.
	if result.ItemID != 0 && len(result.Sargatanas.Prices) == 0 {
		result.OnMarketboard = true
	} else {
		result.OnMarketboard = false
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

// Uses the Recipes and Prices from Information, and returns a Profit model.
// Will require profits from the map if the recipe depends on recipes.
func (info Information) FillProfitsDocument(recipeID int) *models.Profits {
	var profits models.Profits

	recipedoc := info.InnerRecipes[recipeID]
	pricesdoc := info.InnerSimplePrices[recipedoc.ItemResultTargetID]

	profits.RecipeID = recipeID
	profits.ItemID = recipedoc.ItemResultTargetID

	var materialcost int
	for i := 0; i < len(recipedoc.IngredientID); i++ {
		// The top of the stack should have no ingredient recipes.
		if recipedoc.IngredientRecipes[i] == nil {
			innerpricedoc := info.InnerSimplePrices[recipedoc.IngredientID[i]]
			materialcost += innerpricedoc.LowestMarketPrice * recipedoc.IngredientAmounts[i]
		} else {
			// If we do have recipes, it should already be defined in the map.
			innerprofitdoc := info.InnerProfits[recipedoc.IngredientID[i]]
			materialcost += innerprofitdoc.MaterialCosts * recipedoc.IngredientAmounts[i]
		}
	}
	profits.MaterialCosts = materialcost

	if pricesdoc.OnMarketboard {
		profits.Profits = pricesdoc.LowestMarketPrice - materialcost
	} else {
		profits.Profits = pricesdoc.HistoryPrice - materialcost
	}

	// Our profit depends on how much money we've spent going into it.
	profits.ProfitPercentage = float32(profits.Profits) / float32(materialcost)

	now := time.Now()
	unixtimenow := now.Unix()

	profits.Added = unixtimenow

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

// Uses recursion to fill the Information maps and inner information.
// A recipe w/ len(IngredientRecipes) = 0, should be at the top of the stack.
// Will handle if there are no items in the data
func BaseInformation(collections CollectionHandler, recipeID int, info Information) {

	// Adds a base recipe to the map
	baserecipe, indatabase := collections.FindRecipesDocument(recipeID)
	if !indatabase {
		baserecipe = collections.InsertRecipesDocument(recipeID)
	}
	info.InnerRecipes[recipeID] = baserecipe

	// Finds the prices for the base item of a recipe.
	info.InnerSimplePrices[baserecipe.ItemResultTargetID], indatabase = collections.SimplifyPricesDocument(baserecipe.ItemResultTargetID)
	if !indatabase {
		// It means that the prices are actually not in the database, so we just need to find them.
		collections.InsertPricesDocument(baserecipe.ItemResultTargetID)
		// This also means that we can simplify it.
		info.InnerSimplePrices[baserecipe.ItemResultTargetID], _ = collections.SimplifyPricesDocument(baserecipe.ItemResultTargetID)
	}
	// Also adds all the ingredients prices of current recipe into the map.
	for i := 0; i < len(baserecipe.IngredientID); i++ {
		info.InnerSimplePrices[baserecipe.IngredientID[i]], indatabase = collections.SimplifyPricesDocument(baserecipe.IngredientID[i])
		if !indatabase {
			collections.InsertPricesDocument(baserecipe.IngredientID[i])
			info.InnerSimplePrices[baserecipe.IngredientID[i]], _ = collections.SimplifyPricesDocument(baserecipe.IngredientID[i])
		}
	}

	// Recursively call using the inner recipes (if they exist), to add more recipes and prices to our map
	for i := 0; i < len(baserecipe.IngredientRecipes); i++ {
		if baserecipe.IngredientRecipes[i] != nil {
			// Adds all the recipes to the map.
			for j := 0; j < len(baserecipe.IngredientRecipes[i]); j++ {
				BaseInformation(collections, baserecipe.IngredientRecipes[i][j], info)
			}

		}
	}

	info.InnerProfits[recipeID] = info.FillProfitsDocument(recipeID)

}

func ProfitInformation(profit ProfitHandler) []*models.Profits {

	return profit.ProfitDescCursor()
}
