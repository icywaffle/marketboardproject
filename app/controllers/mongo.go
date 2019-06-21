package controllers

import (
	"context"
	"log"
	"marketboardproject/app/controllers/xivapi"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// Initializes DB so it would give the Clients so that we can access the database
func InitDB() {
	clientOptions := options.Client().ApplyURI(xivapi.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("Marketboard")

}
