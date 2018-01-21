package constraint

import (
	"github.com/juju/errors"
)

type OPERATOR = int

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

func comparableOperator(operator OPERATOR) bool {
	return operator == OPERATOR_EQ ||
		operator == OPERATOR_NOT_EQ ||
		operator == OPERATOR_GT ||
		operator == OPERATOR_GTE ||
		operator == OPERATOR_LT ||
		operator == OPERATOR_LTE
}

func setOperator(operator OPERATOR) bool {
	return operator == OPERATOR_CONTAINS ||
		operator == OPERATOR_NOT_CONTAINS
}

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
	operator OPERATOR
	value    interface{}
}

func NewConstraint(operator OPERATOR, value interface{}) *Constraint {
	constraint := &Constraint{}
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
		return r.resolveFloat64(constraint, valueType)
	case float32:
		return r.resolveFloat64(constraint, float64(valueType))
	case int:
		return r.resolveInt64(constraint, int64(valueType))
	case int8:
		return r.resolveInt64(constraint, int64(valueType))
	case int16:
		return r.resolveInt64(constraint, int64(valueType))
	case int32:
		return r.resolveInt64(constraint, int64(valueType))
	case int64:
		return r.resolveInt64(constraint, int64(valueType))
	case string:
		return r.resolveString(constraint, string(valueType))
	default:
		return false, errors.New("unknown type found")
	}

	return true, nil
}

func (r *resolver) resolveFloat64(constraint *Constraint, value float64) (bool, error) {
	floatValue, forceError := r.forceFloat64(constraint.value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.value)
	}

	if comparableOperator(constraint.operator) {
		return r.compareFloat64(constraint.operator, value, floatValue)
	}

	return false, errors.Errorf("could not compare %f with %+v", value, constraint.value)
}
func (r *resolver) resolveInt64(constraint *Constraint, value int64) (bool, error) {
	intValue, forceError := r.forceInt64(constraint.value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.value)
	}

	if comparableOperator(constraint.operator) {
		return r.compareInt64(constraint.operator, value, intValue)
	}

	return false, errors.Errorf("could not compare %f with %+v", value, constraint.value)
}

func (r *resolver) compareFloat64(operator OPERATOR, left float64, right float64) (bool, error) {
	switch operator {
	case OPERATOR_EQ:
		return left == right, nil
	case OPERATOR_NOT_EQ:
		return left != right, nil
	case OPERATOR_LT:
		return left < right, nil
	case OPERATOR_LTE:
		return left <= right, nil
	case OPERATOR_GT:
		return left > right, nil
	case OPERATOR_GTE:
		return left >= right, nil
	default:
		return false, errors.Errorf("operator not available for comparison: %d", operator)
	}
}
func (r *resolver) compareInt64(operator OPERATOR, left int64, right int64) (bool, error) {
	switch operator {
	case OPERATOR_EQ:
		return left == right, nil
	case OPERATOR_NOT_EQ:
		return left != right, nil
	case OPERATOR_LT:
		return left < right, nil
	case OPERATOR_LTE:
		return left <= right, nil
	case OPERATOR_GT:
		return left > right, nil
	case OPERATOR_GTE:
		return left >= right, nil
	default:
		return false, errors.Errorf("operator not available for comparison: %d", operator)
	}
}

func (r *resolver) arrayCompareString(operator OPERATOR, value string, values []string) (bool, error) {

	found := false
	for _, v := range values {
		if v == value {
			found = true
			break
		}
	}

	switch operator {
	case OPERATOR_CONTAINS:
		return found, nil
	case OPERATOR_NOT_CONTAINS:
		return !found, nil
	default:
		return false, errors.Errorf("operator not available for comparison: %d", operator)
	}
}

func (r *resolver) resolveString(constraint *Constraint, value string) (bool, error) {
	println("Resolve string")
	return true, nil
}

func (r *resolver) forceFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, errors.Errorf("could not force %+v to float64", value)
	}
}
func (r *resolver) forceInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, errors.Errorf("could not force %+v to float64", value)
	}
}
