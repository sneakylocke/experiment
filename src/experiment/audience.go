package experiment

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
	return true, nil
}
