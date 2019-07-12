package ffdiscord

import (
	"context"
	"fmt"
	"marketboardproject/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	User *mongo.Collection
}

// Inserts a discord user into the database in order to keep track of tables.
func (coll Collections) InsertUserDocument(DiscordUser *models.DiscordUser) {
	filter := bson.M{"UniqueID": DiscordUser.UniqueID}
	// Testing if the user already exists in the DB.
	var isinDB models.DiscordUser
	err := coll.User.FindOne(context.TODO(), filter).Decode(&isinDB)
	// If it's already in the DB, there's nothing to do here.
	if err != nil {
		coll.User.InsertOne(context.TODO(), DiscordUser)
		fmt.Println("Inserted User into Database: ", DiscordUser.UniqueID)
	}
	fmt.Println("User is already in Database.")
}

// This allows the user to be persistent across pages.
func (coll Collections) FindUserDocument(DiscordUser *models.DiscordUser) {

}
