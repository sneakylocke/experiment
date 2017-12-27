package builders

import e "github.com/sneakylocke/experiment"

type BasicBuilder interface {
	Build() (*e.Experiment, error)
	AddFloat(variableName string, weights []uint32, values []float64) error
	AddInt(variableName string, weights []uint32, values []int64) error
	AddBool(variableName string, weights []uint32, values []bool) error
}

func NewBasicExperiment(experimentName string) *e.Experiment {
	experiment := &e.Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = experimentName
	experiment.VariableNames = []string{}
	experiment.Salt = experimentName
	experiment.Enabled = true

	// Setup Audience
	audience := e.NewAudience()

	// Add audience
	experiment.Audiences = make([]e.Audience, 1, 1)
	experiment.Audiences[0] = *audience

	return experiment
}
