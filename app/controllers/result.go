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

	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))

	baseinfo := xivapi.BaseInformation(DB, recipeID)

	// We need to lock here, to prevent multiple users from updating/calling from an outdated database
	// When one person inserts a new item, the second person will still have an outdated database.
	// This also allows multiple people to search for different items without being locked behind mutex
	if baseinfo.Recipes.ID == 0 || baseinfo.Profits.ItemID == 0 || baseinfo.Prices.ItemID == 0 {
		Mutex.Lock()
		baseinfo = xivapi.InsertInformation(DB, recipeID)
		Mutex.Unlock()
	}

	return c.Render(baseinfo)
}

func (c Result) Profit() revel.Result {
	profitpercentage := xivapi.ProfitInformation(DB)

	return c.Render(profitpercentage)
}

func (c Result) Search() revel.Result {
	recipename := c.Params.Form.Get("recipename")
	return c.Render(recipename)
}
