package tests

import (
	"fmt"
	"marketboardproject/app/controllers/xivapi"
	"marketboardproject/app/models"

	"github.com/revel/revel/testing"
)

// TestSuite is Revel's Equivalent to Golang's testing.T
type AppTest struct {
	testing.TestSuite
}

type FakeCollections struct {
}

// Mock CollectionHandler Interfaces
// These methods fill out the struct, rather than calling a database or API for information.
func (fake FakeCollections) FindRecipesDocument(recipeID int) *models.Recipes {
	var recipes models.Recipes
	// Using 1, to test if the IDs are actually changing in BaseInformation
	recipes.ID = 1
	// Once it does change, we can pretend 1 is the result that the database doesn't have the info.
	if recipes.ID == 1 {
		recipes = *fake.InsertRecipesDocument(recipeID)
	}

	return &recipes
}

func (fake FakeCollections) FindPricesDocument(itemID int) *models.Prices {
	var prices models.Prices
	prices.ItemID = 1
	if prices.ItemID == 1 {
		prices = *fake.InsertPricesDocument(itemID)
	}
	return &prices
}

func (fake FakeCollections) FindProfitsDocument(info *xivapi.Information, recipeID int) *models.Profits {
	var profits models.Profits
	profits.RecipeID = 1
	if profits.RecipeID == 1 {
		profits = *fake.InsertProfitsDocument(info, recipeID)
	}
	return &profits
}

func (fake FakeCollections) InsertRecipesDocument(recipeID int) *models.Recipes {
	var recipes models.Recipes
	recipes.ID = 33180
	return &recipes
}

func (fake FakeCollections) InsertPricesDocument(itemID int) *models.Prices {
	var prices models.Prices
	prices.ItemID = 24322
	return &prices
}

func (fake FakeCollections) InsertProfitsDocument(info *xivapi.Information, recipeID int) *models.Profits {
	var profits models.Profits
	// Mocks a call to the information, and passes this information to the profits.
	// Profits can now manipulate whatever is called from the information.
	// For the purpose of this test, we can just assign these variables.
	profits.RecipeID = info.Recipes.ID
	profits.ItemID = info.Prices.ItemID
	return &profits
}
func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	// Issues a GET request to the page, and stores this in Response and Response body
	t.Get("/")
	t.AssertOk()
	// There are different types of assert functions given in the documentation.
	// What these Assert functions do is check against your values and return whether a test returns fail or correct
	// Depending on which Assert function that you use.
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) Test_fails_if_BaseInformation_returns_nothing() {
	var testfake FakeCollections
	info := xivapi.BaseInformation(testfake, 33180)
	fmt.Println(info)
	testinfo := xivapi.Information{
		Recipes: &models.Recipes{
			ID: 33180,
		},
		Prices: &models.Prices{
			ItemID: 24322,
		},
		Profits: &models.Profits{
			RecipeID: 33180,
		},
	}
	// For the test, we need to take the struct and put it into arrays.
	expectedarray := [3]int{testinfo.Recipes.ID, testinfo.Prices.ItemID, testinfo.Profits.RecipeID}
	// BaseInformation is broken if it doesn't fill this array with the right info.
	resultarray := [3]int{info.Recipes.ID, info.Prices.ItemID, info.Profits.RecipeID}
	t.AssertEqual(expectedarray, resultarray)
}

func (t *AppTest) After() {
	println("Tear down")
}
