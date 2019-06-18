package controllers

import (
	"fmt"
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
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))
	profits, recipes, prices := xivapi.NetItemPrice(recipeID, materialprices, materialingredients)
	fmt.Println(profits, recipes, prices)
	return c.Render(profits, recipes, prices, materialprices, materialingredients)
}

func (c Result) Update() revel.Result {
	// Here we update the database price entries.
	itemID, _ := strconv.Atoi(c.Params.Form.Get("itemID"))
	xivapi.UpdateItemPrices(itemID)
	return c.Redirect("/Result")
}
