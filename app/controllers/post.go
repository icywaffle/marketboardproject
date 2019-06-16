package controllers

import (
	"context"
	"log"
	"marketboardproject/app/models"
	"strconv"

	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	*revel.Controller
}

// The issue is that the stuff from post.go is not being processed by the App.Index for some reason.
// This means I need to start to understand how to route these together.
// I'm not understanding how the routes affect the functions that are here, nor do I understand how the files are linked togehter.

func (c Post) Index() revel.Result {
	posts := []models.Post{}
	collection := Test.Database("Marketboard").Collection("babysteps")

	findOptions := options.Find()
	findOptions.SetLimit(5)

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return c.Render(posts)

}

func (c Post) Create() revel.Result {

	collection := Test.Database("Marketboard").Collection("babysteps")

	var post models.Post
	post.ID, _ = strconv.Atoi(c.Params.Form.Get("ID"))

	// Insert into the document, whatever was inputted as the post.Name from HTML.
	Itemexample := bson.D{
		primitive.E{Key: "ID", Value: post.ID},
	}
	_, err := collection.InsertOne(context.TODO(), Itemexample)
	if err != nil {
		log.Fatal(err)
	}

	// Tells the web to redirect back to some webpage here.
	return c.Redirect("/")
}
