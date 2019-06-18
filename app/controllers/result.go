package controllers

import (
	"marketboardproject/app/controllers/xivapi"
	"marketboardproject/app/models"
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
	var results models.Result
	var recipes models.Recipes
	var prices models.Prices
	// These maps can show which materials are different tiers of which crafted items.
	materialprices := make(map[int][10]int)
	materialingredients := make(map[int][]int)
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))
	xivapi.NetItemPrice(recipeID, &results, &recipes, &prices, materialprices, materialingredients)
	return c.Render(results, recipes, prices, materialprices, materialingredients)
}

func (c Result) Update() revel.Result {
	// Here we update the database price entries.
	itemID, _ := strconv.Atoi(c.Params.Form.Get("itemID"))
	xivapi.UpdateItemPrices(itemID)
	return c.Redirect("/Result")
}
