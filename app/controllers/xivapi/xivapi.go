package xivapi

import (
	"context"
	"fmt"

	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/models"

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

type CollectionHandler interface {
	FindRecipesDocument(recipeID int) *models.Recipes
	FindPricesDocument(itemID int) *models.Prices
	FindProfitsDocument(recipeID int) *models.Profits
	InsertRecipesDocument(recipeID int) *models.Recipes
	InsertPricesDocument(itemID int) *models.Prices
}

// When testing, you create a FAKE COLLECTIONS, with a FAKE METHOD!
// When you create this FAKE METHOD, You actually return this models.Thing that you've made.
// You test that!

// In the test, build a Fake FindRecipesDocument
// Build a Fake collections, where this FindRecipesDocument will handle it.
// Then call the main Function BaseInformation(calls the fakemethod through an interface)

// This method, must be called by the main function through an interface.
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

func (coll Collections) InsertRecipesDocument(recipeID int) *models.Recipes {
	var result models.Recipes
	byteValue := database.ApiConnect(recipeID, "recipe")
	result = *database.Jsonitemrecipe(byteValue)
	coll.Recipes.InsertOne(context.TODO(), result)
	fmt.Println("Inserted Recipe into Database: ", result.ID)
	return &result
}

func (coll Collections) InsertPricesDocument(itemID int) *models.Prices {
	var result models.Prices
	byteValue := database.ApiConnect(itemID, "market/item")
	// If we do have an empty result, it means that we need to search for the vendor prices.
	if result.ItemID != 0 && len(result.Sargatanas.Prices) == 0 {
		byteValue = database.ApiConnect(result.ItemID, "item")
	}
	result = *database.Jsonprices(byteValue)
	coll.Recipes.InsertOne(context.TODO(), result)
	fmt.Println("Inserted Recipe into Database: ", result.ItemID)
	return &result
}

// Fills out information about a crafted recipe.
func BaseInformation(collections CollectionHandler, recipeID int) *Information {
	var info Information
	info.Recipes = collections.FindRecipesDocument(recipeID)
	// If the ID returns zero, then it's not in the database. We need to insert one.
	if info.Recipes.ID == 0 {
		info.Recipes = collections.InsertRecipesDocument(recipeID)
	}
	info.Prices = collections.FindPricesDocument(info.Recipes.ItemResultTargetID)
	if info.Prices.ItemID == 0 {
		info.Prices = collections.InsertPricesDocument(info.Recipes.ItemResultTargetID)
	}
	info.Profits = collections.FindProfitsDocument(recipeID)

	return &info
}
