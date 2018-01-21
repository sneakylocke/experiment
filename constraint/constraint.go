package constraint

// Constraint is a struct that defines an operator and an object to compare to.
type Constraint struct {
	operator OPERATOR
	value    interface{}
}

// NewConstraint creates and returns a pointer to a Constraint.
func NewConstraint(operator OPERATOR, value interface{}) *Constraint {
	constraint := &Constraint{}
	constraint.operator = operator
	constraint.value = value

	return constraint
}
