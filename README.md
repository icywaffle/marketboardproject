## Marketboard Project
Designed to calculate whether a crafted item will give you profit.

Since it will be impossible to grab information about every single item in the game without putting a stressed load on any API,
it's better to just ask whether or not an item that you see will net you some profit, and what amount of profit.

Then with this, we can fill up a database full of items, eventually, to find
an item that has the highest amount of profit percentage.

## Motivation
Most current APIs only have current marketboard price, and graphs of history,
but they never really show you what items make the most amount of money for 
the price of making them.

The point of this project, is to hopefully show whichs items are worth it to 
just buy materials off of the marketboard and then craft it, and which items
are not worth the effort to craft.

This requires a little more than just marketboard pricing, but also time costs
of obtaining materials, which may be a feature in the futre.

## Tech/framework used
<b>Built with</b>
- [Golang](https://golang.org/)
A simplified programming language, that is great for web development,
and also has very clean syntax.
- [MongoDB](https://www.mongodb.com/)
A No-SQL Database, since there are certain items that do not follow
a strict Schema.
- [Revel](https://revel.github.io/)
A Full-Stack web framework to run the entire project.
- [UIKit](https://getuikit.com/)
A Front-End web framework that minimilistically styles the site.

<b>Additional Dependencies</b>
- [MongoDB-Go-Driver](https://github.com/mongodb/mongo-go-driver)
A MongoDB Driver that allows an easier way to access the mongodb.

## Current Features
Profits / Costs of Items you want to craft.
Sorted List of items with most profits.

## Future Features
Total List of prices and materials that you need for crafting.
Save your searched items into the database so that you can compare which items may net you more profit
A cost of time in how much materials to actually gather.

## Code Example

`xivapi.BaseInformation(collections,recipeID)`
This fills out all the information that can be obtained by the database
for a specific recipeID.

## Development Environment
For current build,

Create a MongoDB server.

Uncomment and paste your Mongo Connection link in `keys/samplekeys.go`

Create an XIVAPI account and obtain your own private key.

Uncomment and past your private key in `keys/samplekeys.go`

Install the MongoDB-Go-Driver, since it's one of the dependencies for the code.

Next, install Revel and create a Revel app in your GOPATH, for example
`revel new marketboardproject`

This allows you to have a "skeleton" app. Then you may be able to copy this project and overwrite over this skeleton, in `$GOPATH/src/marketboardproject.`

Then you just run using the command
`revel run -a marketboardproject`

and finally connect to localhost:9000.

## Testing

Tests can be accessed by going into /@tests in the browser.

Tests can be built by the `marketboard/tests` folder.

## API Reference
- [XIVAPI] (https://xivapi.com/)

## How to use?
Development usage only. This is not yet released on an official server.

## License
MIT © [2019] (Jacob Nguyen)
