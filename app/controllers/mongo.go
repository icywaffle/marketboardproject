package controllers

import (
	"context"
	"log"
	"marketboardproject/app/controllers/ffdiscord"
	"marketboardproject/app/controllers/xivapi"
	"marketboardproject/keys"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB xivapi.Collections
var UserDB ffdiscord.Collections

// Initializes DB so it would give the Clients so that we can access the database
func InitDB() {
	clientOptions := options.Client().ApplyURI(keys.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("Marketboard")

	DB = xivapi.Collections{
		Prices:  database.Collection("Prices"),
		Recipes: database.Collection("Recipes"),
		Profits: database.Collection("Profits")}

	UserDB = ffdiscord.Collections{
		User: database.Collection("User")}
}
