package experiment

import "github.com/juju/errors"

type Audience struct {
	Name         string
	Constraints  []Constraint
	ControlValue Value
	ValueGroups  map[string]ValueGroup
	Exposure     float64
	Enabled      bool
}

func NewAudience() *Audience {
	audience := &Audience{}

	audience.Name = ""
	audience.Constraints = make([]Constraint, 1)
	audience.ValueGroups = make(map[string]ValueGroup)
	audience.Exposure = 1.0
	audience.Enabled = true

	return audience
}

func (a *Audience) validate() (bool, error) {
	if a.Exposure < 0 || a.Exposure > 1 {
		return false, errors.Errorf("invalid exposure: %f", a.Exposure)
	}

	return true, nil
}
