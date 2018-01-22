package constraint

import "github.com/juju/errors"

// OPERATOR is intended to act as an enum for different types of comparisons.
type OPERATOR = string

const (
	OPERATOR_EQ           = "EQ"
	OPERATOR_NOT_EQ       = "NEQ"
	OPERATOR_LT           = "LT"
	OPERATOR_LTE          = "LTE"
	OPERATOR_GT           = "GT"
	OPERATOR_GTE          = "GTE"
	OPERATOR_CONTAINS     = "CONTAINS"
	OPERATOR_NOT_CONTAINS = "NCONTAINS"
)

func ValidateOperator(operator OPERATOR) error {
	if operator == OPERATOR_EQ ||
		operator == OPERATOR_NOT_EQ ||
		operator == OPERATOR_LT ||
		operator == OPERATOR_LTE ||
		operator == OPERATOR_GT ||
		operator == OPERATOR_GTE ||
		operator == OPERATOR_CONTAINS ||
		operator == OPERATOR_NOT_CONTAINS {
		return nil
	}
	return errors.Errorf("invalid operator: %s", operator)
}
