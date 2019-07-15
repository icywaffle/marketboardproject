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
	baseinfo.InnerRecipes = make(map[int]*models.Recipes)           // Contains the inner recipes for some key = Recipe.ID
	baseinfo.InnerSimplePrices = make(map[int]*models.SimplePrices) // Contains the inner prices for some key =  Item ID
	baseinfo.InnerProfits = make(map[int]*models.Profits)           // Contains the profits for the inner recipes for some key = Recipe.Id
	xivapi.BaseInformation(DB, recipeID, baseinfo)

	c.ViewArgs["baseinfo"] = baseinfo
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
