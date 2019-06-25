package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"marketboardproject/app/models"
	"marketboardproject/keys"
	"net/http"
	"strconv"
	"time"
)

// Converts Recipe Pages of json, to arrays.

/////////////////Recipe Struct Here//////////////////////////

type AmountIngredient struct {
	//The outer values
	AmountIngredient0 int `json:"AmountIngredient0"`
	AmountIngredient1 int `json:"AmountIngredient1"`
	AmountIngredient2 int `json:"AmountIngredient2"`
	AmountIngredient3 int `json:"AmountIngredient3"`
	AmountIngredient4 int `json:"AmountIngredient4"`
	AmountIngredient5 int `json:"AmountIngredient5"`
	AmountIngredient6 int `json:"AmountIngredient6"`
	AmountIngredient7 int `json:"AmountIngredient7"`
	AmountIngredient8 int `json:"AmountIngredient8"`
	AmountIngredient9 int `json:"AmountIngredient9"`
}

type ItemIngredient struct {
	ItemIngredient0TargetID int `json:"ItemIngredient0TargetID"`
	ItemIngredient1TargetID int `json:"ItemIngredient1TargetID"`
	ItemIngredient2TargetID int `json:"ItemIngredient2TargetID"`
	ItemIngredient3TargetID int `json:"ItemIngredient3TargetID"`
	ItemIngredient4TargetID int `json:"ItemIngredient4TargetID"`
	ItemIngredient5TargetID int `json:"ItemIngredient5TargetID"`
	ItemIngredient6TargetID int `json:"ItemIngredient6TargetID"`
	ItemIngredient7TargetID int `json:"ItemIngredient7TargetID"`
	ItemIngredient8TargetID int `json:"ItemIngredient8TargetID"`
	ItemIngredient9TargetID int `json:"ItemIngredient9TargetID"`
}

type IngredientRecipe struct {
	ItemIngredientRecipe0 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe0"`
	ItemIngredientRecipe1 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe1"`
	ItemIngredientRecipe2 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe2"`
	ItemIngredientRecipe3 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe3"`
	ItemIngredientRecipe4 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe4"`
	ItemIngredientRecipe5 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe5"`
	ItemIngredientRecipe6 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe6"`
	ItemIngredientRecipe7 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe7"`
	ItemIngredientRecipe8 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe8"`
	ItemIngredientRecipe9 []struct {
		ID int `json:"ID"`
	} `json:"ItemIngredientRecipe9"`
}

// idtype should be "recipe", "item", or "market/item"
func ApiConnect(inputid int, idtype string) []byte {
	// MAX Rate limit is 30 Req/s -> 0.033s/Req. It's safer to use something
	time.Sleep(45 * time.Millisecond)
	byteValue := xivapiconnector(websiteurl(inputid, idtype))
	fmt.Println("Connected to API")
	return byteValue
}

func websiteurl(userID int, idtype string) string {
	//Example: https://xivapi.com/Item/14160
	basewebsite := []byte("https://xivapi.com/")
	field := []byte(idtype)
	uniqueID := []byte(strconv.Itoa(userID))
	completefield := append(field[:], '/')
	userinputurl := append(append(basewebsite[:], completefield[:]...), uniqueID[:]...)

	//Add Authkey to the URL
	authkey := []byte(keys.XivAuthKey)
	websiteurl := append(append(userinputurl[:], '?'), authkey[:]...)

	//Example: https://xivapi.com/market/item/3?authkey&servers=Phoenix,Lich,Moogle
	if idtype == "market/item" {
		servers := []byte("&servers=Sargatanas")
		websiteurl = append(websiteurl[:], servers[:]...)
	}

	s := string(websiteurl)
	return s
}

func xivapiconnector(websiteurl string) []byte {

	//What this does, is open the file, and read it
	jsonFile, err := http.Get(websiteurl)
	if err != nil {
		log.Fatalln(err)
	}
	// Takes the jsonFile.Body, and put it into memory as byteValue array.
	byteValue, err := ioutil.ReadAll(jsonFile.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Body.Close()

	return byteValue
}

func Jsonitemrecipe(byteValue []byte) models.Recipes {

	// Unmarshal the information into the structs
	var recipes models.Recipes
	json.Unmarshal(byteValue, &recipes)

	var amount AmountIngredient
	json.Unmarshal(byteValue, &amount)

	var matitemID ItemIngredient
	json.Unmarshal(byteValue, &matitemID)

	// Create the slices
	amountslice := []int{amount.AmountIngredient0,
		amount.AmountIngredient1,
		amount.AmountIngredient2,
		amount.AmountIngredient3,
		amount.AmountIngredient4,
		amount.AmountIngredient5,
		amount.AmountIngredient6,
		amount.AmountIngredient7,
		amount.AmountIngredient8,
		amount.AmountIngredient9}

	matitemIDslice := []int{matitemID.ItemIngredient0TargetID,
		matitemID.ItemIngredient1TargetID,
		matitemID.ItemIngredient2TargetID,
		matitemID.ItemIngredient3TargetID,
		matitemID.ItemIngredient4TargetID,
		matitemID.ItemIngredient5TargetID,
		matitemID.ItemIngredient6TargetID,
		matitemID.ItemIngredient7TargetID,
		matitemID.ItemIngredient8TargetID,
		matitemID.ItemIngredient9TargetID}

	// We need to go through every single possible recipe that can make this item.
	var matrecipeID IngredientRecipe
	json.Unmarshal(byteValue, &matrecipeID)
	matrecipeIDslice := make([][]int, 10)

	//No choice but to unravel for each element, the possible Material Ingredient Recipe IDs 10 times.
	// There is variable length for different elements.
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe0); i++ {
		matrecipeIDslice[0] = append(matrecipeIDslice[0], matrecipeID.ItemIngredientRecipe0[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe1); i++ {
		matrecipeIDslice[1] = append(matrecipeIDslice[1], matrecipeID.ItemIngredientRecipe1[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe2); i++ {
		matrecipeIDslice[2] = append(matrecipeIDslice[2], matrecipeID.ItemIngredientRecipe2[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe3); i++ {
		matrecipeIDslice[3] = append(matrecipeIDslice[3], matrecipeID.ItemIngredientRecipe3[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe4); i++ {
		matrecipeIDslice[4] = append(matrecipeIDslice[4], matrecipeID.ItemIngredientRecipe4[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe5); i++ {
		matrecipeIDslice[5] = append(matrecipeIDslice[5], matrecipeID.ItemIngredientRecipe5[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe6); i++ {
		matrecipeIDslice[6] = append(matrecipeIDslice[6], matrecipeID.ItemIngredientRecipe6[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe7); i++ {
		matrecipeIDslice[7] = append(matrecipeIDslice[7], matrecipeID.ItemIngredientRecipe7[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe8); i++ {
		matrecipeIDslice[8] = append(matrecipeIDslice[8], matrecipeID.ItemIngredientRecipe8[i].ID)
	}
	for i := 0; i < len(matrecipeID.ItemIngredientRecipe9); i++ {
		matrecipeIDslice[9] = append(matrecipeIDslice[9], matrecipeID.ItemIngredientRecipe9[i].ID)
	}

	// These are custom things that we can add to the Recipes documents
	recipes.IngredientNames = matitemIDslice
	recipes.IngredientAmounts = amountslice
	recipes.IngredientRecipes = matrecipeIDslice
	return recipes
}

func Jsonprices(byteValue []byte) models.Prices {

	var prices models.Prices
	json.Unmarshal(byteValue, &prices)

	return prices

}
