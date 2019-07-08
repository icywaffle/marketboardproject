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
