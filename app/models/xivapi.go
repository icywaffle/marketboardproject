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
	// We can add name maps, and icon ID maps here.
}

type Recipes struct {
	Name               string  `bson:"Name" json:"Name"`
	IconID             int     `bson:"IconID" json:"IconID"`
	ItemResultTargetID int     `bson:"ItemID" json:"ItemResultTargetID"`
	ID                 int     `bson:"RecipeID" json:"ID"`
	CraftTypeTargetID  int     `bson:"CraftTypeTargetID" json:"CraftTypeTargetID"`
	AmountResult       int     `bson:"AmountResult" json:"AmountResult"`
	IngredientNames    []int   `bson:"IngredientName"`
	IngredientAmounts  []int   `bson:"IngredientAmount"`
	IngredientRecipes  [][]int `bson:"IngredientRecipes"`
	Added              int64   `bson:"Added"`
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
	VendorPrice int   `json:"PriceMid" bson:"VendorPrice"`
	Added       int64 `bson:"Added"`
}
