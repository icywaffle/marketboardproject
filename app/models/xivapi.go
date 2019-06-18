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

type Recipes struct {
	Name               string  `json:"Name" bson:"Name"`
	ItemResultTargetID int     `json:"ItemResultTargetID" bson:"ItemID"`
	ID                 int     `json:"ID" bson:"RecipeID"`
	CraftTypeTargetID  int     `json:"CraftTypeTargetID" bson:"CraftTypeTargetID"`
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
