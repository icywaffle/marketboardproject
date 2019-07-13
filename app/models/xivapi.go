package models

type Profits struct {
	Name             string  `bson:"Name"`
	ItemID           int     `bson:"ItemID"`
	RecipeID         int     `bson:"RecipeID"`
	IconID           int     `bson:"IconID"`
	MarketboardPrice int     `bson:"MarketboardPrice"`
	MaterialCosts    int     `bson:"MaterialCosts"`
	Profits          int     `bson:"Profits"`
	ProfitPercentage float32 `bson:"ProfitPercentage"`
	Added            int64   `bson:"Added"`
}

type Matprofitmaps struct {
	Costs       map[int][10]int  //[itemID][prices per material]
	Ingredients map[int][]int    //[itemID][recipeIDs per material]
	Total       map[int]int      //[itemID]totalcost of one item
	Names       map[int][]string //[itemID][names per material]
	IconID      map[int][]int    //[itemID][icondid per material]

}

type Recipes struct {
	Name               string   `bson:"Name" json:"Name"`
	IconID             int      `bson:"IconID" json:"IconID"`
	ItemResultTargetID int      `bson:"ItemID" json:"ItemResultTargetID"`
	ID                 int      `bson:"RecipeID" json:"ID"`
	CraftTypeTargetID  int      `bson:"CraftTypeTargetID" json:"CraftTypeTargetID"`
	AmountResult       int      `bson:"AmountResult" json:"AmountResult"`
	IngredientNames    []string `bson:"IngredientNames"`
	IngredientID       []int    `bson:"IngredientID"`
	IngredientIconID   []int    `bson:"IngredientIconID"`
	IngredientAmounts  []int    `bson:"IngredientAmount"`
	IngredientRecipes  [][]int  `bson:"IngredientRecipes"`
	Added              int64    `bson:"Added"`
}

type Prices struct {
	ItemID     int `bson:"ItemID"`
	Sargatanas struct {
		History []struct {
			Added        int  `json:"Added" bson:"Added"` // XIVAPI added time
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
	Added       int64 `bson:"Added"` // Database added time.
}
