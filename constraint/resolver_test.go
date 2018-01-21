package constraint

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOperatorsFloat64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float64(3.0)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsFloat32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float32(3.0)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsInt(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int(3)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsInt8(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int8(3)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsInt16(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int16(3)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsInt32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int32(3)
	testOperators(t, NewMapContext(context))
}

func TestOperatorsInt64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int64(3)
	testOperators(t, NewMapContext(context))
}

func TestTypesFloat64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float64(3.0)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesFloat32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float32(3.0)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesInt(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int(3)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesInt8(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int8(3)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesInt16(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int16(3)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesInt32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int32(3)
	testTypes(t, NewMapContext(context), 3)
}

func TestTypesInt64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int64(3)
	testTypes(t, NewMapContext(context), 3)
}

func TestOperatorsStrings(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = "cucumbers"
	mapContext := NewMapContext(context)

	resolver := resolver{}

	// EQ
	constraintEQ := NewConstraint("key", OPERATOR_EQ, "cucumbers")
	okEQ, errEQ := resolver.Resolve(constraintEQ, mapContext)
	assert.Nil(t, errEQ)
	assert.True(t, okEQ)

	// NOT EQ
	constraintNEQ := NewConstraint("key", OPERATOR_NOT_EQ, "apples")
	okNEQ, errNEQ := resolver.Resolve(constraintNEQ, mapContext)
	assert.Nil(t, errNEQ)
	assert.True(t, okNEQ)

	// CONTAINS
	constraintContains := NewConstraint("key", OPERATOR_CONTAINS, []string{"apples", "bananas", "cucumbers"})
	okContains, errContains := resolver.Resolve(constraintContains, mapContext)
	assert.Nil(t, errContains)
	assert.True(t, okContains)

	// NOT CONTAINS
	constraintNotContains := NewConstraint("key", OPERATOR_NOT_CONTAINS, []string{"apples", "bananas", "berries"})
	okNotContains, errNotContains := resolver.Resolve(constraintNotContains, mapContext)
	assert.Nil(t, errNotContains)
	assert.True(t, okNotContains)

	// NOT CONTAINS ERROR
	constraintNotContains2 := NewConstraint("key", OPERATOR_NOT_CONTAINS, []int{1, 2, 3, 4})
	okNotContains2, errNotContains2 := resolver.Resolve(constraintNotContains2, mapContext)
	assert.NotNil(t, errNotContains2)
	assert.False(t, okNotContains2)
}

func TestNoKey(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = "cucumbers"
	mapContext := NewMapContext(context)

	resolver := resolver{}

	constraintEQ := NewConstraint("wrong_key", OPERATOR_EQ, "cucumbers")
	okEQ, errEQ := resolver.Resolve(constraintEQ, mapContext)
	assert.NotNil(t, errEQ)
	assert.False(t, okEQ)
}

func TestWeirdKeyValue(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = []string{"this", "does", "not", "work"}
	mapContext := NewMapContext(context)

	resolver := resolver{}

	constraintEQ := NewConstraint("key", OPERATOR_EQ, "cucumbers")
	okEQ, errEQ := resolver.Resolve(constraintEQ, mapContext)
	assert.NotNil(t, errEQ)
	assert.False(t, okEQ)
}

func testTypes(t *testing.T, context Context, constraintValue float64) {
	resolver := resolver{}

	constraint := NewConstraint("key", OPERATOR_EQ, float64(constraintValue))
	ok1, err1 := resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, float32(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, int8(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, int16(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, int32(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, int64(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)

	constraint = NewConstraint("key", OPERATOR_EQ, int(constraintValue))
	ok1, err1 = resolver.Resolve(constraint, context)
	assert.Nil(t, err1)
	assert.True(t, ok1)
}

func testOperators(t *testing.T, context Context) {
	resolver := resolver{}

	// LT
	constraintLT := NewConstraint("key", OPERATOR_LT, 4)
	okLT, errLT := resolver.Resolve(constraintLT, context)
	assert.Nil(t, errLT)
	assert.True(t, okLT)

	// LTE
	constraintLTE := NewConstraint("key", OPERATOR_LTE, 3.0)
	okLTE, errLTE := resolver.Resolve(constraintLTE, context)
	assert.Nil(t, errLTE)
	assert.True(t, okLTE)

	// LTE
	constraintLTE2 := NewConstraint("key", OPERATOR_LTE, 3.0)
	okLTE2, errLTE2 := resolver.Resolve(constraintLTE2, context)
	assert.Nil(t, errLTE2)
	assert.True(t, okLTE2)

	// GT
	constraintGT := NewConstraint("key", OPERATOR_GT, 2)
	okGT, errGT := resolver.Resolve(constraintGT, context)
	assert.Nil(t, errGT)
	assert.True(t, okGT)

	// GTE
	constraintGTE := NewConstraint("key", OPERATOR_GTE, 3.0)
	okGTE, errGTE := resolver.Resolve(constraintGTE, context)
	assert.Nil(t, errGTE)
	assert.True(t, okGTE)

	// GTE
	constraintGTE2 := NewConstraint("key", OPERATOR_GTE, 2.0)
	okGTE2, errGTE2 := resolver.Resolve(constraintGTE2, context)
	assert.Nil(t, errGTE2)
	assert.True(t, okGTE2)

	// NOT EQUAL
	constraintNEQ := NewConstraint("key", OPERATOR_NOT_EQ, 2)
	okNEQ, errNEQ := resolver.Resolve(constraintNEQ, context)
	assert.Nil(t, errNEQ)
	assert.True(t, okNEQ)

	// EQUAL FLOAT
	constraintEQ1 := NewConstraint("key", OPERATOR_EQ, 3.0)
	okEQ1, errEQ1 := resolver.Resolve(constraintEQ1, context)
	assert.Nil(t, errEQ1)
	assert.True(t, okEQ1)

	// EQUAL INT
	constraintEQ2 := NewConstraint("key", OPERATOR_EQ, 3)
	okEQ2, errEQ2 := resolver.Resolve(constraintEQ2, context)
	assert.Nil(t, errEQ2)
	assert.True(t, okEQ2)

	// WRONG CONSTRAINT TYPE
	constraintBad := NewConstraint("key", OPERATOR_EQ, make(map[string]int))
	okBad, errBad := resolver.Resolve(constraintBad, context)
	assert.NotNil(t, errBad)
	assert.False(t, okBad)
}
