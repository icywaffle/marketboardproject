package models

type Result struct {
	ItemID           int
	RecipeID         int
	MarketboardPrice int
	MaterialCosts    int
	Profits          int
	ProfitPercentage float32
}
