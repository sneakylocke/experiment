package experiment

import (
	"github.com/juju/errors"
	"github.com/sneakylocke/experiment/constraint"
)

type Audience struct {
	Name        string                  `json:"name"`
	Constraints []constraint.Constraint `json:"constraints"`
	ValueGroups map[string]*ValueGroup  `json:"valueGroups"`
	Exposure    float64                 `json:"exposure"`
	Enabled     bool                    `json:"enabled"`
}

func NewAudience() *Audience {
	audience := &Audience{}

	audience.Name = ""
	audience.Constraints = make([]constraint.Constraint, 0)
	audience.ValueGroups = make(map[string]*ValueGroup)
	audience.Exposure = 1.0
	audience.Enabled = true

	return audience
}

func (a *Audience) Validate() (bool, error) {
	if a.Exposure < 0 || a.Exposure > 1 {
		return false, errors.Errorf("invalid exposure: %f", a.Exposure)
	}

	return true, nil
}
