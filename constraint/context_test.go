package constraint

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float64(3.0)
	testOperators(t, NewMapContext(context))
}

func TestFloat32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = float32(3.0)
	testOperators(t, NewMapContext(context))
}
func TestInt(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int(3)
	testOperators(t, NewMapContext(context))
}
func TestInt8(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int8(3)
	testOperators(t, NewMapContext(context))
}
func TestInt16(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int16(3)
	testOperators(t, NewMapContext(context))
}
func TestInt32(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int32(3)
	testOperators(t, NewMapContext(context))
}
func TestInt64(t *testing.T) {
	context := make(map[string]interface{})
	context["key"] = int64(3)
	testOperators(t, NewMapContext(context))
}

func testOperators(t *testing.T, context Context) {
	resolver := resolver{}

	// LT
	constraintLT := NewConstraint(OPERATOR_LT, 4)
	okLT, errLT := resolver.resolve("key", constraintLT, context)
	assert.Nil(t, errLT)
	assert.True(t, okLT)

	// LTE
	constraintLTE := NewConstraint(OPERATOR_LTE, 3.0)
	okLTE, errLTE := resolver.resolve("key", constraintLTE, context)
	assert.Nil(t, errLTE)
	assert.True(t, okLTE)

	// LTE
	constraintLTE2 := NewConstraint(OPERATOR_LTE, 3.0)
	okLTE2, errLTE2 := resolver.resolve("key", constraintLTE2, context)
	assert.Nil(t, errLTE2)
	assert.True(t, okLTE2)

	// GT
	constraintGT := NewConstraint(OPERATOR_GT, 2)
	okGT, errGT := resolver.resolve("key", constraintGT, context)
	assert.Nil(t, errGT)
	assert.True(t, okGT)

	// GTE
	constraintGTE := NewConstraint(OPERATOR_GTE, 3.0)
	okGTE, errGTE := resolver.resolve("key", constraintGTE, context)
	assert.Nil(t, errGTE)
	assert.True(t, okGTE)

	// GTE
	constraintGTE2 := NewConstraint(OPERATOR_GTE, 2.0)
	okGTE2, errGTE2 := resolver.resolve("key", constraintGTE2, context)
	assert.Nil(t, errGTE2)
	assert.True(t, okGTE2)

	// NOT EQUAL
	constraintNEQ := NewConstraint(OPERATOR_NOT_EQ, 2)
	okNEQ, errNEQ := resolver.resolve("key", constraintNEQ, context)
	assert.Nil(t, errNEQ)
	assert.True(t, okNEQ)

	// EQUAL FLOAT
	constraintEQ1 := NewConstraint(OPERATOR_EQ, 3.0)
	okEQ1, errEQ1 := resolver.resolve("key", constraintEQ1, context)
	assert.Nil(t, errEQ1)
	assert.True(t, okEQ1)

	// EQUAL INT
	constraintEQ2 := NewConstraint(OPERATOR_EQ, 3)
	okEQ2, errEQ2 := resolver.resolve("key", constraintEQ2, context)
	assert.Nil(t, errEQ2)
	assert.True(t, okEQ2)
}
