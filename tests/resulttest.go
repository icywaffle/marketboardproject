package tests

import (
	"marketboardproject/app/controllers"
	"marketboardproject/app/controllers/xivapi"
	"marketboardproject/app/models"

	"github.com/revel/revel/testing"
)

// Tests for the Results Folder
type ResultTest struct {
	testing.TestSuite
}

type FakeCollections struct {
}

// Mock CollectionHandler Interfaces

// The find functions just simply find if it's in the database.
func (fake FakeCollections) FindRecipesDocument(recipeID int) *models.Recipes {
	var recipes models.Recipes
	recipes.ID = 33180
	return &recipes
}

func (fake FakeCollections) FindPricesDocument(itemID int) *models.Prices {
	var prices models.Prices
	prices.ItemID = 24322
	return &prices
}

func (fake FakeCollections) FindProfitsDocument(recipeID int) *models.Profits {
	var profits models.Profits
	profits.RecipeID = 33180
	return &profits
}

// These insert functions simply insert/update into the database, and return you a result.
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
	// Mocks a call to the information, and it should load up the structs with these info.
	// The profits calculations are a bit more complicated, but this mock is a relatively simple idea.
	profits.RecipeID = info.Recipes.ID
	profits.ItemID = info.Prices.ItemID
	return &profits
}
func (fake FakeCollections) FillProfitMaps(info *xivapi.Information, matprofitmaps *models.Matprofitmaps, forceupdate bool) {
	// This mocks the recalls to the collections, and inserting it into the maps.
	matprofitmaps.Costs[24322] = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
}

func (fake FakeCollections) ProfitDescCursor() []*models.Profits {
	var fakearray []*models.Profits

	// Mocks getting a cursor from the database, and appending it to a full array of information.
	var tempprofits models.Profits
	tempprofits.RecipeID = 33180
	var tempprofits2 models.Profits
	tempprofits2.RecipeID = 33181
	fakearray = append(fakearray, &tempprofits)
	fakearray = append(fakearray, &tempprofits2)

	return fakearray
}
func (t *ResultTest) Before() {
	println("Set up")
}

// Unit test for BaseInformation
func (t *ResultTest) Test_fails_if_BaseInformation_returns_nothing() {
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

func (t *ResultTest) Test_fails_if_BaseInformation_maps_are_nil() {
	var testfake FakeCollections
	info := xivapi.BaseInformation(testfake, 33180)
	testmap := make(map[int][10]int)
	testmap[24322] = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	testinfo := xivapi.Information{
		Matprofitmaps: &models.Matprofitmaps{
			Costs: testmap,
		},
	}

	t.AssertEqual(testinfo.Matprofitmaps.Costs, info.Matprofitmaps.Costs)
}

func (t *ResultTest) Test_fails_if_InsertInformation_maps_are_nil() {
	var testfake FakeCollections
	info := xivapi.InsertInformation(testfake, 33180, false)
	testmap := make(map[int][10]int)
	testmap[24322] = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	testinfo := xivapi.Information{
		Matprofitmaps: &models.Matprofitmaps{
			Costs: testmap,
		},
	}

	t.AssertEqual(testinfo.Matprofitmaps.Costs, info.Matprofitmaps.Costs)
}

// Unit test for ProfitInformation
func (t *ResultTest) Test_fails_if_ProfitInformation_returns_nothing() {
	var fakecollection FakeCollections
	resultarray := xivapi.ProfitInformation(fakecollection)
	expectedarray := []*models.Profits{{RecipeID: 33180}, {RecipeID: 33181}}

	t.AssertEqual(expectedarray, resultarray)
}

// Functional test for Database Connection
func (t *ResultTest) Test_fails_if_missing_DB_collection_or_InitDB_failed_to_connect() {
	dbflag := true
	if controllers.DB.Prices == nil || controllers.DB.Recipes == nil || controllers.DB.Profits == nil {
		dbflag = false
	}
	t.Assert(dbflag)
}

func (t *ResultTest) After() {
	println("Tear down")
}
