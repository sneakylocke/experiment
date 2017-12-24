package experiment

type Experiment struct {
	Name          string     // Name of the experiment
	VariableNames []string   // Names of variables allowed to be randomized in audiences
	Audiences     []Audience // Details of variables are captured in a single audience. Users may belong to one audience
	Salt          string
	Enabled       bool
}

type WeightedValue struct {
	Value  Value
	Weight uint32
}

type Constraint struct {
}
