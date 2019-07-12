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
	c.renderdiscorduser()
	return c.RenderTemplate("Result/Index.html")
}

func (c Result) Obtain() revel.Result {

	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))

	baseinfo := xivapi.BaseInformation(DB, recipeID)

	// We need to lock here, to prevent multiple users from updating/calling from an outdated database
	// When one person inserts a new item, the second person will still have an outdated database.
	// This also allows multiple people to search for different items without being locked behind mutex
	// Added is based off of when the database adds it. If it's zero, then it was never in the database.
	if baseinfo.Recipes.Added < xivapi.UpdatedRecipesStructTime || baseinfo.Profits.Added < xivapi.UpdatedProfitsStructTime || baseinfo.Prices.Added == xivapi.UpdatedPricesStructTime {
		Mutex.Lock()
		baseinfo = xivapi.InsertInformation(DB, recipeID)
		Mutex.Unlock()
	}
	c.ViewArgs["baseinfo"] = baseinfo
	c.renderdiscorduser()
	return c.RenderTemplate("Result/Obtain.html")
}

func (c Result) Profit() revel.Result {
	profitpercentage := xivapi.ProfitInformation(DB)
	// To update profits, we actually need all the previous information.
	// And check through all the items to make sure we're updating them.
	for i := 0; i < len(profitpercentage); i++ {
		if profitpercentage[i].Added < xivapi.UpdatedProfitsStructTime {
			Mutex.Lock()
			baseinfo := xivapi.InsertInformation(DB, profitpercentage[i].RecipeID)
			profitpercentage[i] = baseinfo.Profits
			Mutex.Unlock()
		}
	}

	c.renderdiscorduser()
	c.ViewArgs["profitpercentage"] = profitpercentage
	return c.RenderTemplate("Result/Profit.html")
}

func (c Result) Search() revel.Result {
	recipename := c.Params.Form.Get("recipename")
	c.ViewArgs["recipename"] = recipename

	c.renderdiscorduser()
	return c.RenderTemplate("Result/Search.html")
}

// Adds discordmap to the ViewArgs
func (c Result) renderdiscorduser() {
	discorduser, _ := c.Session.Get("discordinfo")
	if discorduser != nil {
		discordmap, _ := discorduser.(map[string]interface{})
		c.ViewArgs["discordmap"] = discordmap
	} else {
		c.ViewArgs["discordmap"] = nil
	}
}
