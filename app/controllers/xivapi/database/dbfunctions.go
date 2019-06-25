package database

import (
	"context"
	"fmt"
	"marketboardproject/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
