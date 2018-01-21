package constraint

import "github.com/juju/errors"

// Resolver is an interface that defines methods needed to resolve whether constraints are satisfied by some Context.
type Resolver interface {
	Resolve(key string, constraint Constraint, context Context) (bool, error)

	resolveFloat64(constraint Constraint, value float64) (bool, error)
	resolveInt64(constraint Constraint, value int64) (bool, error)
	resolveString(constraint Constraint, value string) (bool, error)
}

// resolver is a default implementation of Resolver
type resolver struct {
}

// Resolve returns true is the Constraint is satisfied via the provided Context for a given key.
func (r *resolver) Resolve(key string, constraint *Constraint, context Context) (bool, error) {

	// Attempt to retrieve the value at the key
	value, contextErr := context.value(key)

	if contextErr != nil {
		return false, errors.Annotatef(contextErr, "key not found in context: %s", key)
	}

	// Inspect the type of the value and resolve the constraint appropriately.
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
	// Attempt to force the constraint's value to a float64 for comparison
	floatValue, forceError := r.forceFloat64(constraint.value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.value)
	}

	return r.compareFloat64(constraint.operator, value, floatValue)
}

func (r *resolver) resolveInt64(constraint *Constraint, value int64) (bool, error) {
	// Attempt to force the constraint's value to a int64 for comparison
	intValue, forceError := r.forceInt64(constraint.value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.value)
	}

	return r.compareInt64(constraint.operator, value, intValue)
}

func (r *resolver) resolveString(constraint *Constraint, value string) (bool, error) {
	// Attempt direct string comparison first
	stringValue, stringOk := constraint.value.(string)

	if stringOk {
		switch constraint.operator {
		case OPERATOR_EQ:
			return value == stringValue, nil
		case OPERATOR_NOT_EQ:
			return value != stringValue, nil
		default:
			return false, errors.Errorf("could not compare strings with operator: %d", constraint.operator)
		}
	}

	// Attempt set comparison (contains, not contains)
	strings, stringsOk := constraint.value.([]string)
	if stringsOk {
		return r.arrayCompareString(constraint.operator, value, strings)
	}

	return false, errors.Errorf("could not compare input %s with constraint %+v", value, constraint.value)
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
		return false, errors.Errorf("operator not available for float comparison: %d", operator)
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
		return false, errors.Errorf("operator not available for int comparison: %d", operator)
	}
}

func (r *resolver) arrayCompareString(operator OPERATOR, value string, values []string) (bool, error) {
	// Try to find the value in the array of strings
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

func (r *resolver) forceFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
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
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, errors.Errorf("could not force %+v to float64", value)
	}
}
