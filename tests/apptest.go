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
	var test models.Recipes
	test.ID = 33180
	return &test
}

func (fake FakeCollections) FindPricesDocument(itemID int) *models.Prices {
	var test models.Prices
	test.ItemID = 24322

	return &test
}

func (fake FakeCollections) FindProfitsDocument(recipeID int) *models.Profits {
	var test models.Profits
	test.RecipeID = 33180

	return &test
}

func (fake FakeCollections) InsertRecipesDocument(recipeID int) *models.Recipes {
	var test models.Recipes
	// If a database doesn't have the ID, we can insert one.
	// The real method will call the API.
	test.ID = 33180
	fmt.Println("Inserted Test Recipes")
	return &test
}

func (fake FakeCollections) InsertPricesDocument(itemID int) *models.Prices {
	var test models.Prices
	test.ItemID = 24322
	fmt.Println("Inserted Test Prices")
	return &test
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
