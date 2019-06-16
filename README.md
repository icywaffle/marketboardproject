## Marketboard Project
Designed to calculate whether a crafted item will give you profit.

## Motivation
New Motivation: Since it will be impossible to grab information about every single item in the game without putting a stressed load on any API,
it's better to just ask whether or not an item that you see will net you some profit, and what amount of profit.

Then we can store this information into the database to help see which items will actually net you a better amount of profit.

This database begins to become better over time, the more searches that are used, since they will give the database more information with each request.



## Tech/framework used
<b>Built with</b>
- [Golang](https://golang.org/)
- [MongoDB](https://www.mongodb.com/)
- [MongoDB-Go-Driver](https://github.com/mongodb/mongo-go-driver)
- [Revel](https://revel.github.io/)

## Current Features
Profits / Costs of Items you want to craft.

## Future Features
Total List of prices and materials that you need for crafting.
Force update the prices of the materials
Save your searched items into the database so that you can compare which items may net you more profit
A full database with sorted percentages of profits. Then you can use this to find the highest profited items eventually.
   For this however, it needs to be limited per user, so that a single user doesn't crawl the entire xivapi.

## Code Example

`xivapi.NetItemPrice(recipeID int, results *models.Result)`
This Net Item Price function allows you to access the database, to search for the current marketboard prices and determine the profit. If it's not inside the database currently, then it will access xivapi, to be able to find the current prices.


## Installation
For current build,

Install MongoDB, and create a server that uses the default port 27017.

Create an XIVAPI account and obtain your own private key.

Create a file inside `marketboardproject/app/controllers/xivapi/urlstring/`

Then create a go file that contains

`package keys
var XivAuthKey string = "private_key=#######"
`

Next, install revel and create a revel app, for example
`revel new marketboardproject`

Then you may be able to copy and overwrite the folder that it creates in $GOPATH/src/marketboardproject.

Then you just run using the command
`revel run -a marketboardproject`

## API Reference
- [XIVAPI] (https://xivapi.com/)

## How to use?
Development usage only. This is not yet released on an official server.

## License
MIT © [2019] (Jacob Nguyen)
