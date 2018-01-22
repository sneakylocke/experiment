package constraint

import "github.com/juju/errors"

// Constraint is a struct that defines an Operator, an object to compare to, and the Key/name of what type of thing
// Value is (country, height).
type Constraint struct {
	Key      string      `json:"key"`
	Operator OPERATOR    `json:"operator"`
	Value    interface{} `json:"value"`
}

// NewConstraint creates and returns a pointer to a Constraint.
func NewConstraint(key string, operator OPERATOR, value interface{}) *Constraint {
	constraint := &Constraint{}
	constraint.Key = key
	constraint.Operator = operator
	constraint.Value = value

	return constraint
}

func (c *Constraint) Validate() error {
	if c.Key == "" {
		return errors.Errorf("constraint Key must be specified: %+v", c)
	}

	if c.Value == nil {
		return errors.Errorf("constraint Value must not be nil: %+v", c)
	}

	if err := ValidateOperator(c.Operator); err != nil {
		return err
	}

	return nil
}
