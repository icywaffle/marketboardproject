package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Calls ingredient amounts and item IDs, and returns the results
func Ingredientmaterials(collection *mongo.Collection, recipeID int) *Recipes {
	filter := bson.M{"RecipeID": recipeID}
	var result Recipes
	collection.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// Call the prices from the database, and return the sold average and the current average
func Ingredientprices(collection *mongo.Collection, itemID int) *Prices {
	filter := bson.M{"ItemID": itemID}
	var result Prices
	collection.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// Pass information from jsonconv to this to input these values into the database.
func InsertRecipe(collection *mongo.Collection, recipes Recipes, ingredientid []int, ingredientamount []int, ingredientrecipes [][]int) {

	Itemexample := bson.D{
		primitive.E{Key: "Name", Value: recipes.Name},
		primitive.E{Key: "ItemID", Value: recipes.ItemResultTargetID},
		primitive.E{Key: "RecipeID", Value: recipes.ID},
		primitive.E{Key: "CraftTypeTargetID", Value: recipes.CraftTypeTargetID},
		primitive.E{Key: "IngredientName", Value: ingredientid},
		primitive.E{Key: "IngredientAmount", Value: ingredientamount},
		primitive.E{Key: "IngredientRecipes", Value: ingredientrecipes},
	}

	insertResult, err := collection.InsertOne(context.TODO(), Itemexample)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted Item into Database: ", insertResult.InsertedID)

}

func InsertPrices(collection *mongo.Collection, prices Prices, itemID int) {
	// In cases that we don't have any market ready price, we need to grab the gil vendor price instead
	if prices.VendorPrice != 0 {
		Itemexample := bson.D{
			primitive.E{Key: "ItemID", Value: itemID},
			{Key: "Sargatanas", Value: bson.D{
				{Key: "Prices", Value: bson.A{bson.D{
					{Key: "Added", Value: 0},
					{Key: "PricePerUnit", Value: prices.VendorPrice}}}},
			}},
		}
		insertResult, err := collection.InsertOne(context.TODO(), Itemexample)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted Item into Database: ", insertResult.InsertedID)
	} else {

		Itemexample := bson.D{
			primitive.E{Key: "ItemID", Value: itemID},
			primitive.E{Key: "Sargatanas", Value: prices.Sargatanas},
		}

		insertResult, err := collection.InsertOne(context.TODO(), Itemexample)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted Item into Database: ", insertResult.InsertedID)

	}

}

func UpdatePrices(collection *mongo.Collection, prices Prices, itemID int) {
	filter := bson.M{"ItemID": itemID}

	update := bson.D{
		primitive.E{Key: "Sargatanas", Value: prices.Sargatanas},
	}

	var result Prices
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		InsertPrices(collection, prices, itemID)
	} else {
		collection.UpdateOne(context.TODO(), filter, update)
		fmt.Println("Updated Item into Database")
	}

}
