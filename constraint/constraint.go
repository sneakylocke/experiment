package constraint

import "github.com/juju/errors"

// Constraint is a struct that defines an operator and an object to compare to.
type Constraint struct {
	key      string
	operator OPERATOR
	value    interface{}
}

// NewConstraint creates and returns a pointer to a Constraint.
func NewConstraint(key string, operator OPERATOR, value interface{}) *Constraint {
	constraint := &Constraint{}
	constraint.key = key
	constraint.operator = operator
	constraint.value = value

	return constraint
}

func (c *Constraint) Validate() error {
	if c.key == "" {
		return errors.Errorf("constraint key must be specified")
	}

	if c.value == nil {
		return errors.Errorf("constraint value must not be nil")
	}

	return nil
}
