package constraint

import (
	"github.com/juju/errors"
)

type VARIANT = int
type OPERATOR = int

const (
	VARIANT_FLOAT_64 = iota
	VARIANT_INT_64
	VARIANT_STRING
)

const (
	OPERATOR_EQ = iota
	OPERATOR_NOT_EQ
	OPERATOR_LT
	OPERATOR_LTE
	OPERATOR_GT
	OPERATOR_GTE
	OPERATOR_CONTAINS
	OPERATOR_NOT_CONTAINS
)

type Context interface {
	value(key string) (interface{}, error)
}

type MapContext struct {
	context map[string]interface{}
}

func (context *MapContext) value(key string) (interface{}, error) {
	if value, ok := context.context[key]; ok {
		return value, nil
	}

	return nil, errors.Errorf("key '%s' does not exists in context", key)
}

func NewMapContext(context map[string]interface{}) *MapContext {
	c := MapContext{}
	c.context = context
	return &c
}

type Resolver interface {
	resolve(key string, constraint Constraint, context Context) (bool, error)
	resolveFloat64(constraint Constraint, value float64) (bool, error)
	resolveInt64(constraint Constraint, value int64) (bool, error)
	resolveString(constraint Constraint, value string) (bool, error)
}

type Constraint struct {
	variant  VARIANT
	operator OPERATOR
	value    interface{}
}

func NewConstraint(variant VARIANT, operator OPERATOR, value interface{}) *Constraint {
	constraint := &Constraint{}
	constraint.variant = variant
	constraint.operator = operator
	constraint.value = value

	return constraint
}

type resolver struct {
}

func (r *resolver) resolve(key string, constraint *Constraint, context Context) (bool, error) {

	value, contextErr := context.value(key)

	if contextErr != nil {
		return false, errors.Annotatef(contextErr, "key not found in context: %s", key)
	}

	switch valueType := value.(type) {
	case float64:
		return r.resolveFloat64(constraint, float64(valueType))
	case float32:
		return r.resolveFloat64(constraint, float64(valueType))
	case int64:
		return r.resolveFloat64(constraint, float64(int64(valueType)))
	case int:
		return r.resolveFloat64(constraint, float64(int(valueType)))
	case string:
		return r.resolveString(constraint, string(valueType))
	default:
		return false, errors.New("unknown type found")
	}

	return true, nil
}

func (r *resolver) resolveFloat64(constraint *Constraint, value float64) (bool, error) {
	println("Resolve float 64")
	return true, nil
}

func (r *resolver) resolveString(constraint *Constraint, value string) (bool, error) {
	println("Resolve string")
	return true, nil
}
