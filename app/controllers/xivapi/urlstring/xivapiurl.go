// UrlItemRecipe(userID string) returns a string of the Item Recipe Url.
// UrlSearch(usersearch string) returns a string of the User Search Url
// UrlPrices(useritemid int) 	returns a string of the Item Prices Url.
// XiviapiRecipeConnector(recipeID int) connects to the API, and returns the byteValue of the url page.
package urlstring

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// TODO: ?columns=Attributes,Object.Attribute will significantly lower payload

//Creates the URL for recipes and items
func UrlItemRecipe(userID int) string {
	//Example: https://xivapi.com/Recipe/33180
	basewebsite := []byte("https://xivapi.com/")
	field := []byte("Recipe")
	uniqueID := []byte(strconv.Itoa(userID))
	completefield := append(field[:], '/')
	userinputurl := append(append(basewebsite[:], completefield[:]...), uniqueID[:]...)

	//Finishing the url with the AuthKey
	authkey := []byte(XivAuthKey)
	websiteurl := append(append(userinputurl[:], '?'), authkey[:]...)

	s := string(websiteurl)
	return s
}

func UrlItem(userID int) string {
	//Example: https://xivapi.com/Item/14160
	basewebsite := []byte("https://xivapi.com/")
	field := []byte("Item")
	uniqueID := []byte(strconv.Itoa(userID))
	completefield := append(field[:], '/')
	userinputurl := append(append(basewebsite[:], completefield[:]...), uniqueID[:]...)

	//Finishing the url with the AuthKey
	authkey := []byte(XivAuthKey)
	websiteurl := append(append(userinputurl[:], '?'), authkey[:]...)

	s := string(websiteurl)
	return s
}

//UserInputs some item to search. This appends it to the websiteurl.
func UrlSearch(usersearch string) string {
	//Example: https://xivapi.com/search?string
	basewebsite := []byte("https://xivapi.com/search?string=")

	//Example: https://xivapi.com/search?string=High+Mythrite+Ingot
	var replacer = strings.NewReplacer(" ", "+")
	fixedusersearch := replacer.Replace(usersearch)
	searchfield := []byte(fixedusersearch)
	userinputurl := append(append(basewebsite[:], searchfield[:]...), '&')

	authkey := []byte(XivAuthKey)
	websiteurl := append(userinputurl[:], authkey[:]...)

	s := string(websiteurl)
	return s
}

func UrlPrices(useritemid int) string {
	//Example: https://xivapi.com/market/item/3?servers=Phoenix,Lich,Moogle

	//Produces : https://xivapi.com/market/item/3
	itemwebsitefield := []byte("https://xivapi.com/market/item/")
	itemid := []byte(strconv.Itoa(useritemid))
	basewebsite := append(itemwebsitefield[:], itemid[:]...)

	//Produces :https://xivapi.com/market/item/3?servers=Phoenix,Lich,Moogle&
	//TODO:Let's just use Sargatanas for now for simple structs, then expand later. ?servers=Adamantoise,Cactuar,Faerie,Gilgamesh,Jenova,Midgardsormr,Sargatanas,Siren
	servers := []byte("?servers=Sargatanas")
	userinputurl := append(append(basewebsite[:], servers[:]...), '&')

	//Attaches key to the end.
	authkey := []byte(XivAuthKey)
	websiteurl := append(userinputurl[:], authkey[:]...)

	s := string(websiteurl)
	return s
}

func XiviapiRecipeConnector(websiteurl string) []byte {

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
