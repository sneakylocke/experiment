package constraint

// OPERATOR is intended to act as an enum for different types of comparisons.
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
