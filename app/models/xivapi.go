package models

type Profits struct {
	ItemID           int     `bson:"ItemID"`
	RecipeID         int     `bson:"RecipeID"`
	MarketboardPrice int     `bson:"MarketboardPrice"`
	MaterialCosts    int     `bson:"MaterialCosts"`
	Profits          int     `bson:"Profits"`
	ProfitPercentage float32 `bson:"ProfitPercentage"`
	Added            int64   `bson:"Added"`
}

type Matprofitmaps struct {
	// These maps can show which materials are different tiers of which crafted items.
	Costs       map[int][10]int
	Ingredients map[int][]int
	Total       map[int]int
}

type Recipes struct {
	Name               string  `bson:"Name" json:"Name"`
	ItemResultTargetID int     `bson:"ItemID" json:"ItemResultTargetID"`
	ID                 int     `bson:"RecipeID" json:"ID"`
	CraftTypeTargetID  int     `bson:"CraftTypeTargetID" json:"CraftTypeTargetID"`
	AmountResult       int     `bson:"AmountResult" json:"AmountResult"`
	IngredientNames    []int   `bson:"IngredientName"`
	IngredientAmounts  []int   `bson:"IngredientAmount"`
	IngredientRecipes  [][]int `bson:"IngredientRecipes"`
}

type Prices struct {
	ItemID     int `bson:"ItemID"`
	Sargatanas struct {
		History []struct {
			Added        int  `json:"Added" bson:"Added"` // Time is in Unix epoch time
			IsHQ         bool `json:"IsHQ" bson:"IsHQ"`
			PricePerUnit int  `json:"PricePerUnit" bson:"PricePerUnit"`
			PriceTotal   int  `json:"PriceTotal" bson:"PriceTotal"`
			PurchaseDate int  `json:"PurchaseDate" bson:"PurchaseDate"`
			Quantity     int  `json:"Quantity" bson:"Quantity"`
		} `json:"History" bson:"History"`
		Prices []struct {
			Added        int  `json:"Added" bson:"Added"`
			IsHQ         bool `json:"IsHQ" bson:"IsHQ"`
			PricePerUnit int  `json:"PricePerUnit" bson:"PricePerUnit"`
			PriceTotal   int  `json:"PriceTotal" bson:"PriceTotal"`
			Quantity     int  `json:"Quantity" bson:"Quantity"`
		} `json:"Prices" bson:"Prices"`
	} `json:"Sargatanas" bson:"Sargatanas"`
	VendorPrice int `json:"PriceMid"`
}

// Lets also learn some dependency injections

// What is dependency injection?
// You basically have some struct that initializes another struct.
// Or it really requires another struct to implement the behavior.

// Then you have a function that initializes (constructor)
// Your function input should be that inner struct.
// We can pass anything into that function, as long as it meets the struct requirements.

// To provide layers, you basically have a struct
// That contains the structs for config, and repository

// Basically, you want structs to be built by some function.
// You don't want a single function to both, build dependencies and the resultant structs.
// You want to separate the dependency from the function if possible.
// So you have an input as a dependency instead!

// The separation of these into structs allows us to build config alone,
// then use this config to build a server connection.
// This also allows us to create mock responses as well.
