package constraint

import "github.com/juju/errors"

// Resolver is an interface that defines methods needed to resolve whether constraints are satisfied by some Context.
type Resolver interface {
	Resolve(constraint *Constraint, context Context) (bool, error)
}

// NewDefaultResolver returns a basic implementation of a Resolver
func NewDefaultResolver() Resolver {
	return &resolver{}
}

// resolver is a default implementation of Resolver
type resolver struct {
}

// Resolve returns true is the Constraint is satisfied via the provided Context for a given Key.
func (r *resolver) Resolve(constraint *Constraint, context Context) (bool, error) {
	if context == nil {
		return false, errors.Errorf("no context provided")
	}

	// Attempt to retrieve the Value at the Key
	value, contextErr := context.value(constraint.Key)

	if contextErr != nil {
		return false, errors.Annotatef(contextErr, "Key not found in context: %s", constraint.Key)
	}

	// Inspect the type of the Value and resolve the constraint appropriately.
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
	// Attempt to force the constraint's Value to a float64 for comparison
	floatValue, forceError := r.forceFloat64(constraint.Value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.Value)
	}

	return r.compareFloat64(constraint.Operator, value, floatValue)
}

func (r *resolver) resolveInt64(constraint *Constraint, value int64) (bool, error) {
	// Attempt to force the constraint's Value to a int64 for comparison
	intValue, forceError := r.forceInt64(constraint.Value)

	if forceError != nil {
		return false, errors.Annotatef(forceError, "could not compare %f with %+v", value, constraint.Value)
	}

	return r.compareInt64(constraint.Operator, value, intValue)
}

func (r *resolver) resolveString(constraint *Constraint, value string) (bool, error) {
	// Attempt direct string comparison first
	stringValue, stringOk := constraint.Value.(string)

	if stringOk {
		switch constraint.Operator {
		case OPERATOR_EQ:
			return value == stringValue, nil
		case OPERATOR_NOT_EQ:
			return value != stringValue, nil
		default:
			return false, errors.Errorf("could not compare strings with Operator: %d", constraint.Operator)
		}
	}

	// Attempt set comparison (contains, not contains)
	strings, stringsOk := constraint.Value.([]string)
	if stringsOk {
		return r.arrayCompareString(constraint.Operator, value, strings)
	}

	// Having trouble parsing an array of strings from a JSON file
	moreStrings, moreStringsOk := constraint.Value.([]interface{})
	if moreStringsOk {
		strings := make([]string, 0, len(moreStrings))

		for _, interfaceObject := range moreStrings {
			s, sOk := interfaceObject.(string)

			if !sOk {
				return false, errors.Errorf("expected to parse an array of strings, found %+v", interfaceObject)
			}

			strings = append(strings, s)
		}

		return r.arrayCompareString(constraint.Operator, value, strings)
	}

	return false, errors.Errorf("could not compare input %s with constraint %+v", value, constraint.Value)
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
		return false, errors.Errorf("Operator not available for float comparison: %d", operator)
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
		return false, errors.Errorf("Operator not available for int comparison: %d", operator)
	}
}

func (r *resolver) arrayCompareString(operator OPERATOR, value string, values []string) (bool, error) {
	// Try to find the Value in the array of strings
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
		return false, errors.Errorf("Operator not available for comparison: %d", operator)
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
