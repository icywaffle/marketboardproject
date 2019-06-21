package tests

import (
	"github.com/revel/revel/testing"
)

// TestSuite is Revel's Equivalent to Golang's testing.T
type AppTest struct {
	testing.TestSuite
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

// Let's try to write a unit test for Ingredient Materials.
// A unit test should test only one method.
// A unit test should provide some specific arguments into that method.
// A unit test tests the result that is as expected
func (t *AppTest) TestIngredientMaterials() {
	// What specific behavior are you testing?
	// You should be careful of redundancies.
		// One logical assertion per test.
	// Test only one code unit at a time.
	// No need to look inside the method to know what it's doing.
	// Changing the internals should not cause the test to fail.

	// You should not directly test that private methods are being called.
		// If you do, then use a code coverage tool.

	// If the method calls other public methods in other packages, and these calls
	// are guaranteed by the interface, then you can test that these calls are made by a mocking framework
		 // What does this really mean?
		
	// Do not use the method, or internal code, to generate the expected result dynamically.
	// The expected result should be hard-coded into the test case.
		// This is so that it doesn't change when implementation changes.

	// Add more test cases, to see if you have missed any interesting paths.

	// You don't need ot know HOW things are done, just that they correctly get stuff done.

	// Test Driven - Thinking about writing tests and code at the same time.

	// Summary : You want to test the final result, and see if the function behaves the way it does,
	// Given that you've provided the correct inputs that you're going to provide.
	// Unit tests are meant for also purposes that our code could fail.
	// For example : We want to simulate some requests in the case that we would receive some error page.
		// Like getting an error page when we're requesting some recipe ID 3900.
		// We don't want to spam the API for these because we can expect a response.
		// So unit testing allows us to simulate these bad responses.

	// Test should not depend on ANY OUTSIDE RESOURCES (DATABASES, APIs, ETC.)
	// If you need them, you should simulate them
		// This is because you want to control what the return should be given a specific input.

		// We need Inversion of Control.
		// This means we create a method for the input.
		// Like https://blog.drewolson.org/dependency-injection-in-go

	Ingredientmaterials(collection *mongo.Collection, recipeID int) *models.Recipes
}

func (t *AppTest) After() {
	println("Tear down")
}
