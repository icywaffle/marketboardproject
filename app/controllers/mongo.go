package controllers

import (
	"context"
	"log"
	"marketboardproject/app/controllers/xivapi"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Test *mongo.Client

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

	Test = client

}
