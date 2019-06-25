package xivapi

import (
	"context"
	"fmt"
	"time"

	"marketboardproject/app/controllers/xivapi/database"
	"marketboardproject/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//6/25/19 - 1PM : Added new info to the struct.
var UpdatedStructTime = int64(1561493761)

type Collections struct {
	Prices  *mongo.Collection
	Recipes *mongo.Collection
	Profits *mongo.Collection
}

type Information struct {
	Recipes       *models.Recipes
	Prices        *models.Prices
	Profits       *models.Profits
	Matprofitmaps *models.Matprofitmaps
}

type CollectionHandler interface {
	FindRecipesDocument(recipeID int) *models.Recipes
	FindPricesDocument(itemID int) *models.Prices
	FindProfitsDocument(info *Information, recipeID int) *models.Profits
	InsertRecipesDocument(recipeID int) *models.Recipes
	InsertPricesDocument(itemID int) *models.Prices
	InsertProfitsDocument(info *Information, recipeID int) *models.Profits
	FillProfitMaps(info *Information, matprofitmaps *models.Matprofitmaps)
}

type ProfitHandler interface {
	ProfitDescCursor() []*models.Profits
}

func (coll Collections) FindRecipesDocument(recipeID int) *models.Recipes {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Recipes
	coll.Recipes.FindOne(context.TODO(), filter).Decode(&result)
	// If the ID returns zero, then it's not in the database. We need to insert one.
	// Also, we need to force update when we update the struct with more info.
	if result.ID == 0 || result.Added < UpdatedStructTime {
		result = *coll.InsertRecipesDocument(recipeID)
	}
	return &result
}
func (coll Collections) FindPricesDocument(itemID int) *models.Prices {
	filter := bson.M{"ItemID": itemID}
	var result models.Prices
	coll.Prices.FindOne(context.TODO(), filter).Decode(&result)
	if result.ItemID == 0 || result.Added < UpdatedStructTime {
		result = *coll.InsertPricesDocument(itemID)
	}

	return &result
}

func (coll Collections) FindProfitsDocument(info *Information, recipeID int) *models.Profits {
	filter := bson.M{"RecipeID": recipeID}
	var result models.Profits
	coll.Profits.FindOne(context.TODO(), filter).Decode(&result)
	if result.RecipeID == 0 || result.Added < UpdatedStructTime {
		result = *coll.InsertProfitsDocument(info, recipeID)
	}
	return &result
}

// Will insert a document, or update it if it's already in the collection.
func (coll Collections) InsertRecipesDocument(recipeID int) *models.Recipes {
	byteValue := database.ApiConnect(recipeID, "recipe")
	result := database.Jsonitemrecipe(byteValue)
	// These variables are not in the json file.
	now := time.Now()
	result.Added = now.Unix()
	// Testing if there's an entry in the DB
	var isinDB models.Recipes
	filter := bson.M{"RecipeID": recipeID}
	err := coll.Recipes.FindOne(context.TODO(), filter).Decode(&isinDB)
	if err != nil {
		coll.Recipes.InsertOne(context.TODO(), result)
		fmt.Println("Inserted Recipe into Database: ", result.ID)
	} else {
		coll.Recipes.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: result},
		})
		fmt.Println("Updated Item into Recipe Collection :", result.ID)
	}

	return &result
}

func (coll Collections) InsertPricesDocument(itemID int) *models.Prices {
	byteValue := database.ApiConnect(itemID, "market/item")
	result := database.Jsonprices(byteValue)
	// ItemID is not part of the Json file.
	result.ItemID = itemID
	// If we do have an empty result, it means that we need to search for the vendor prices.
	if result.ItemID != 0 && len(result.Sargatanas.Prices) == 0 {
		byteValue = database.ApiConnect(result.ItemID, "item")
		result = database.Jsonprices(byteValue)
		// We have to rewrite the Item ID into this
		result.ItemID = itemID
	}
	//These variables are not in the json file.
	now := time.Now()
	result.Added = now.Unix()

	filter := bson.M{"ItemID": itemID}
	var isinDB models.Prices
	err := coll.Prices.FindOne(context.TODO(), filter).Decode(&isinDB)
	if err != nil {
		coll.Prices.InsertOne(context.TODO(), result)
		fmt.Println("Inserted Prices into Database: ", result.ItemID)
	} else {
		coll.Prices.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: result},
		})
		fmt.Println("Updated Item into Prices Collection :", result.ItemID)
	}

	return &result
}

// Creates and then inserts the profits document
func (coll Collections) InsertProfitsDocument(info *Information, recipeID int) *models.Profits {
	var profits models.Profits
	profits.ItemID = info.Recipes.ItemResultTargetID
	profits.RecipeID = info.Recipes.ID

	// Some items won't have a market history, because they're from vendors.
	var itempriceperunit int
	if info.Prices.VendorPrice != 0 {
		itempriceperunit = info.Prices.VendorPrice
	} else {
		itempriceperunit = info.Prices.Sargatanas.Prices[0].PricePerUnit
	}
	profits.MarketboardPrice = itempriceperunit

	materialcosts := coll.findsum(info, info.Matprofitmaps)
	profits.MaterialCosts = materialcosts
	// We may get multiple items per craft.
	profits.Profits = itempriceperunit*info.Recipes.AmountResult - materialcosts
	profits.ProfitPercentage = (float32(profits.Profits)) / float32(materialcosts) * 100

	now := time.Now()
	profits.Added = now.Unix()

	filter := bson.M{"RecipeID": recipeID}

	var isinDB models.Profits
	err := coll.Profits.FindOne(context.TODO(), filter).Decode(&isinDB)
	if err != nil {
		coll.Profits.InsertOne(context.TODO(), profits)
		fmt.Println("Inserted Profits into Database: ", profits.RecipeID)
	} else {
		coll.Profits.UpdateOne(context.TODO(), filter, bson.D{
			{Key: "$set", Value: profits},
		})
		fmt.Println("Updated Item into Profit Collection :", profits.RecipeID)
	}
	return &profits

}

// This recursive function, calls through the materials of materials etc, and fills a map up.
func (coll Collections) FillProfitMaps(info *Information, matprofitmaps *models.Matprofitmaps) {

	// Price array, will allow us to take the price of crafting a material,
	// only if the material is craftable.
	var pricearray [10]int
	// We need to search through the base materials
	for i := 0; i < len(info.Recipes.IngredientNames); i++ {
		// Zero is an invalid material ID
		if info.Recipes.IngredientNames[i] != 0 {
			var matpriceinfo models.Prices
			matpriceinfo = *coll.FindPricesDocument(info.Recipes.IngredientNames[i])
			// We need to also deal with vendor prices, since they won't have market prices.
			if matpriceinfo.VendorPrice != 0 {
				pricearray[i] = matpriceinfo.VendorPrice * info.Recipes.IngredientAmounts[i]
			} else {
				if len(matpriceinfo.Sargatanas.Prices) > 0 {
					pricearray[i] = matpriceinfo.Sargatanas.Prices[0].PricePerUnit * info.Recipes.IngredientAmounts[i]
				} else {
					// If the market's empty, there's no available price for it.
					// This is something we have to integrate in the future.
					pricearray[i] = 0
				}
			}
		} else {
			// Zero should be skipped, since it's not a valid item.
			continue
		}
		// After receiving the price information for specific materials, we fill the maps here.
		matprofitmaps.Costs[info.Recipes.ItemResultTargetID] = pricearray
		matprofitmaps.Ingredients[info.Recipes.ItemResultTargetID] = info.Recipes.IngredientNames
	}

	// If there's a recipe for a material, we want to go in one more materialprices, and keep appending to it.
	for i := 0; i < len(info.Recipes.IngredientRecipes); i++ {
		if len(info.Recipes.IngredientRecipes[i]) != 0 {
			// Creating this matinfo allows us to recursively call and create new instances
			// And this will allow us to get information about mats of mats etc.
			var matinfo Information
			matinfo.Recipes = coll.FindRecipesDocument(info.Recipes.IngredientRecipes[i][0])
			// We can then use this new material information, to fill the maps some more.
			coll.FillProfitMaps(&matinfo, matprofitmaps)
		}
	}

	// Then we can finally fill the maps, once we've finished looping through the
	// materials with recipes.
	for itemID, pricearray := range matprofitmaps.Costs {
		var pricesum int
		for i := 0; i < len(pricearray); i++ {
			pricesum += pricearray[i]
		}
		matprofitmaps.Total[itemID] = pricesum
	}

}
func (coll Collections) findsum(info *Information, matprofitmaps *models.Matprofitmaps) int {

	var tiersum int

	// Some materials are base items, so these base items won't have a map key for prices.
	temppricearray := matprofitmaps.Costs[info.Recipes.ItemResultTargetID]
	for i := 0; i < len(info.Recipes.IngredientNames); i++ {
		materialtotalprice, ok := matprofitmaps.Total[info.Recipes.IngredientNames[i]]
		if ok {
			// If a material also has a recipe, then we want to recursively call for it's material prices.
			_, innerrecipe := matprofitmaps.Ingredients[info.Recipes.IngredientNames[i]]
			if innerrecipe {
				var materialinfo Information
				// We're going to need to find the information about the inner materials.
				// For now, we will deal with just the first recipe.
				materialinfo.Recipes = coll.FindRecipesDocument(info.Recipes.IngredientRecipes[i][0])
				materialinfo.Prices = coll.FindPricesDocument(info.Recipes.IngredientNames[i])
				// We we need to redefine the materialtotalprice with the one that is found by looking at the prices of the materials within the materials.
				// We also need to pass the main maps, to fill it up.
				materialtotalprice = coll.findsum(&materialinfo, matprofitmaps)
			}
			temppricearray[i] = materialtotalprice
		}

		tiersum += temppricearray[i]
	}

	return tiersum
}

// Gives a Descending Sorted Array, of 20 items with the most profit from the DB
func (coll Collections) ProfitDescCursor() []*models.Profits {
	options := options.FindOptions{}
	options.Sort = bson.D{{Key: "ProfitPercentage", Value: -1}}
	// We can set this to be bigger later in the future
	limit := int64(20)
	options.Limit = &limit
	cursor, _ := coll.Profits.Find(context.Background(), bson.D{}, &options)

	var allprofits []*models.Profits
	for cursor.Next(context.TODO()) {
		var tempprofits models.Profits
		cursor.Decode(&tempprofits)

		allprofits = append(allprofits, &tempprofits)
	}
	defer cursor.Close(context.TODO())

	return allprofits

}

// Fills out information about a crafted recipe.
func BaseInformation(collections CollectionHandler, recipeID int) *Information {
	var info Information
	info.Recipes = collections.FindRecipesDocument(recipeID)

	info.Prices = collections.FindPricesDocument(info.Recipes.ItemResultTargetID)

	// We need to initialize our maps
	var matprofitmaps models.Matprofitmaps
	matprofitmaps.Costs = make(map[int][10]int)
	matprofitmaps.Ingredients = make(map[int][]int)
	matprofitmaps.Total = make(map[int]int)
	collections.FillProfitMaps(&info, &matprofitmaps)
	info.Matprofitmaps = &matprofitmaps

	info.Profits = collections.FindProfitsDocument(&info, recipeID)

	return &info
}

func ProfitInformation(profit ProfitHandler) []*models.Profits {

	return profit.ProfitDescCursor()
}
