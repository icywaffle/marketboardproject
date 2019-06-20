package database

import (
	"context"
	"fmt"
	"log"
	"marketboardproject/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Calls ingredient amounts and item IDs, and returns the results
func Ingredientmaterials(collection *mongo.Collection, recipeID int) *models.Recipes {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Recipes
	collection.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// Call the prices from the database, and return the sold average and the current average
func Ingredientprices(collection *mongo.Collection, itemID int) *models.Prices {
	filter := bson.M{"ItemID": itemID}
	var result models.Prices
	collection.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

func Ingredientprofits(collection *mongo.Collection, recipeID int) *models.Profits {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Profits
	collection.FindOne(context.TODO(), filter).Decode(&result)

	return &result
}

// When calling this function, you should close the cursor after being done with it.
func Profitcomparisons(collection *mongo.Collection) *mongo.Cursor {
	options := options.FindOptions{}

	options.Sort = bson.D{{Key: "ProfitPercentage", Value: -1}}
	limit := int64(20)
	options.Limit = &limit
	cursor, _ := collection.Find(context.Background(), bson.D{}, &options)
	return cursor

}

// Pass information from jsonconv to this to input these values into the database.
func InsertRecipe(collection *mongo.Collection, recipes models.Recipes) {

	insertResult, err := collection.InsertOne(context.TODO(), recipes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted Recipe into Database: ", insertResult.InsertedID)

}

func InsertPrices(collection *mongo.Collection, prices models.Prices, itemID int) {
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

		fmt.Println("Inserted Prices into Database: ", insertResult.InsertedID)
	} else {

		Itemexample := bson.D{
			primitive.E{Key: "ItemID", Value: itemID},
			primitive.E{Key: "Sargatanas", Value: prices.Sargatanas},
		}

		insertResult, err := collection.InsertOne(context.TODO(), Itemexample)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted Prices into Database: ", insertResult.InsertedID)

	}

}

func InsertProfits(collection *mongo.Collection, profits models.Profits) {
	insertResult, err := collection.InsertOne(context.TODO(), profits)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted Profits into Database: ", insertResult.InsertedID)

}

func UpdatePrices(collection *mongo.Collection, prices models.Prices, itemID int) {
	filter := bson.M{"ItemID": itemID}

	var result models.Prices
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Unable to Update to Prices")
	} else {
		collection.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: prices},
		})
		fmt.Println("Updated Item into Prices Collection :", itemID)
	}

}

func UpdateProfits(collection *mongo.Collection, profits models.Profits, recipeID int) {
	filter := bson.M{"RecipeID": recipeID}

	var result models.Profits
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Unable to Update to Profits")
	} else {
		collection.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: profits},
		})
		fmt.Println("Updated Item into Profit Collection :", recipeID)
	}
}

func UpdateRecipes(collection *mongo.Collection, recipes models.Recipes) {
	filter := bson.M{"RecipeID": recipes.ID}

	//recipesdocument, _ := structtobsond(recipes)
	var result models.Recipes
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Unable to Update to Recipes")
	} else {
		collection.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: recipes},
		})
		fmt.Println("Updated Item into Recipe Collection :", recipes.ID)
	}
}
