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

// collection := DB.Collection(collectionname)
func (c Result) Obtain() revel.Result {

	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))

	baseinfo := xivapi.BaseInformation(DB, recipeID)

	fmt.Println(baseinfo)
	return c.Render(baseinfo)
}

func (c Result) Profit() revel.Result {
	profitpercentage := xivapi.ProfitInformation(DB)

	return c.Render(profitpercentage)
}

/*
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



func (c Result) UpdateProfits() revel.Result {
	recipeID, _ := strconv.Atoi(c.Params.Form.Get("recipeID"))
	xivapi.ForceUpdateProfits(recipeID)
	return c.Redirect("/Result/Profit")
}
*/
