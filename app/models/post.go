package models

type Post struct {
	Name string `bson:"Name" json:"Name"`
	ID   int    `bson:"ID" json:"ID"`
	Food string `bson:"Food" json:"Food"`
}
