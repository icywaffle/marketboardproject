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
	Recipes      *models.Recipes
	InnerRecipes map[int]*models.Recipes // Contains the inner recipes for some key = Recipe.ID
	Prices       *models.Prices
	InnerPrices  map[int]*models.Prices // Contains the inner prices for some key =  Item ID
	Profits      *models.Profits
	InnerProfits map[int]*models.Profits // Contains the profits for the inner recipes for some key = Recipe.Id
}

type CollectionHandler interface {
	FindRecipesDocument(recipeID int) *models.Recipes
	FindPricesDocument(itemID int) *models.Prices
	FindProfitsDocument(recipeID int) *models.Profits
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

// Fills out information about a crafted recipe.
func BaseInformation(collections CollectionHandler, recipeID int) *Information {
	var info Information
	info.Recipes = collections.FindRecipesDocument(recipeID)

	info.Prices = collections.FindPricesDocument(info.Recipes.ItemResultTargetID)

	// Find profits will recursively call Baseinformation here.
	info.Profits = collections.FindProfitsDocument(recipeID)

	return &info
}

func ProfitInformation(profit ProfitHandler) []*models.Profits {

	return profit.ProfitDescCursor()
}

// Checks whether the information is filled, and inserts information into the database if not.
// Will also force insert/update if chosen to do so by boolean.
func InsertInformation(collections CollectionHandler, info Information, recipeID int, forceupdate bool) Information {
	// We need to pass a base info, because we have a map that needs to be filled.
	/*
		// This is redefined in case someone else had force updated prior to the Mutex Lock of another.
		currenttime := time.Now()
		profitstimesinceupdate := currenttime.Unix() - info.Profits.Added

		forceupdateprofits := false
		if profitstimesinceupdate > 86400/2 && forceupdate == true {
			forceupdateprofits = true
		} else {
			forceupdateprofits = false
		}
		// We have to separate the two otherwise we'll update prices needlessly, in the case that we actually have no profit.
		pricestimesinceupdate := currenttime.Unix() - info.Prices.Added
		forceupdateprices := false
		if pricestimesinceupdate > 86400/2 && forceupdate == true {
			forceupdateprices = true
		} else {
			forceupdateprices = false
		}
	*/

	// If Recipes.Added == 0, then it also means we need to insert into the database since we don't have it.
	if info.Recipes.Added < UpdatedRecipesStructTime {
		info.Recipes = collections.InsertRecipesDocument(recipeID)
	}
	//Recursively calling the function with inner recipes should handle all the recipes required.
	if _, ok := info.InnerRecipes[info.Recipes.ID]; !ok {
		info.InnerRecipes[info.Recipes.ID] = info.Recipes
	}

	// We only need to force update the prices and profit calculations afterwards.
	if info.Prices.Added < UpdatedPricesStructTime {
		// We need to pass all the items each recipe requires.
		for i := 0; i < len(info.Recipes.IngredientID); i++ {
			info.Prices = collections.InsertPricesDocument(info.Recipes.IngredientID[i])
			if _, ok := info.InnerPrices[info.Recipes.IngredientID[i]]; !ok {
				info.InnerPrices[info.Recipes.IngredientID[i]] = info.Prices
			}
		}

	}

	if info.Profits.Added < UpdatedProfitsStructTime {
		// Check from info if there are Ingredient Recipes.
		for i := 0; i < len(info.Recipes.IngredientRecipes); i++ {
			if info.Recipes.IngredientRecipes[i] != nil {
				// Recursively add the inner recipes and prices to the map.
				InsertInformation(collections, info, info.Recipes.IngredientRecipes[i][0], false)
				// Fill the profits for each item.
				info.FillProfitsDocument(info.Recipes.IngredientRecipes[i][0])
			}
		}

	}

	return info

}
