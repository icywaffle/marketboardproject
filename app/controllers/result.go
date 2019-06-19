package controllers

import (
	"marketboardproject/app/controllers/xivapi"
	"strconv"

	"github.com/revel/revel"
)

type Result struct {
	*revel.Controller
}

func (c Result) Index() revel.Result {
	greetings := "Greetings Earthling"
	return c.Render(greetings)
}

func (c Result) Obtain() revel.Result {

	// These maps can show which materials are different tiers of which crafted items.
	materialprices := make(map[int][10]int)
	materialingredients := make(map[int][]int)
	materialtotal := make(map[int]int)
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))
	profits, recipes, prices := xivapi.NetItemPrice(recipeID, materialprices, materialingredients, materialtotal)
	return c.Render(profits, recipes, prices, materialprices, materialingredients)
}

// Allows user to either update recipes or prices on one page. Also allows a user to do both at the same time.
func (c Result) UpdatePrices() revel.Result {
	itemID, _ := strconv.Atoi(c.Params.Form.Get("itemID"))
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))

	if recipeID != 0 {
		xivapi.ForceUpdateRecipes(recipeID)
	}
	if itemID != 0 {
		xivapi.ForceUpdateItemPrices(itemID)
	}
	return c.Redirect("/Result")
}

func (c Result) Profit() revel.Result {
	profitpercentage := xivapi.CompareProfits()

	return c.Render(profitpercentage)
}

func (c Result) UpdateProfits() revel.Result {
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))
	xivapi.ForceUpdateProfits(recipeID)
	return c.Redirect("/Result/Profit")
}
