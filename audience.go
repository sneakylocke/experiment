package experiment

type Audience struct {
	Constraints  []Constraint
	ControlValue Value
	ValueGroups  map[string]ValueGroup
	Exposure     float64
	Enabled      bool
}

func NewAudience() *Audience {
	audience := &Audience{}

	audience.Constraints = make([]Constraint, 1)
	audience.ValueGroups = make(map[string]ValueGroup)
	audience.Exposure = 1.0
	audience.Enabled = true

	return audience
}
