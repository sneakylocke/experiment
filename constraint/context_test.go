package constraint

import "testing"

func TestMatch(t *testing.T) {
	context := make(map[string]interface{})
	context["books"] = 3
	context["height"] = 5.8
	context["country"] = "USA"

	c := NewMapContext(context)

	resolver := resolver{}

	bookConstraint := NewConstraint(VARIANT_INT_64, OPERATOR_EQ, 3)
	countryConstraint := NewConstraint(VARIANT_INT_64, OPERATOR_EQ, "USA")

	ok1, err1 := resolver.resolve("books", bookConstraint, c)
	ok2, err2 := resolver.resolve("country", countryConstraint, c)

	println("First one")
	println(ok1)
	println(err1)
	println("Second")
	println(ok2)
	println(err2)

}
