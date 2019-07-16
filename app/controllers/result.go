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
	c.renderdiscorduser()
	return c.RenderTemplate("Result/Index.html")
}

func (c Result) Obtain() revel.Result {

	recipeID, _ := strconv.Atoi(c.Params.Form.Get("updatespecificrecipe"))
	var baseinfo xivapi.Information
	// We have to initialize the maps here, to be able to allow recursive calls.
	var innerinfo xivapi.InnerInformation
	innerrecipes := make(map[int]*models.Recipes)           // Contains the inner recipes for some key = Recipe.ID
	innersimpleprices := make(map[int]*models.SimplePrices) // Contains the inner prices for some key =  Item ID
	innerprofits := make(map[int]*models.Profits)
	innerinfo.Recipes = innerrecipes
	innerinfo.Prices = innersimpleprices
	innerinfo.Profits = innerprofits
	xivapi.BaseInformation(DB, recipeID, innerinfo)

	// The baseinfo should also be in the maps themselves.
	baseinfo.Recipes = innerinfo.Recipes[recipeID]
	baseinfo.Prices = innerinfo.Prices[baseinfo.Recipes.ItemResultTargetID]
	baseinfo.Profits = innerinfo.Profits[recipeID]

	c.ViewArgs["baseinfo"] = baseinfo
	c.ViewArgs["innerinfo"] = innerinfo
	c.renderdiscorduser()
	return c.RenderTemplate("Result/Obtain.html")
}

func (c Result) Profit() revel.Result {
	profitpercentage := xivapi.ProfitInformation(DB)

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
