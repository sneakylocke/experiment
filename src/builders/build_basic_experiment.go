package builders

import (
	e "github.com/sneakylocke/experiment/src/experiment"
)

func NewSimpleExperiment(experimentName string, variableName string, weights []uint32, values []float64) *e.Experiment {
	experiment := &e.Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = experimentName
	experiment.VariableNames = []string{variableName}
	experiment.Salt = experimentName
	experiment.Enabled = true

	// Setup ValueGroup
	valueGroup := e.NewFloatValueGroup(variableName, weights, values)

	// Setup Audience
	audience := e.NewAudience()
	audience.ControlValue = valueGroup.WeightedValues[0].Value
	audience.ValueGroups[variableName] = *valueGroup

	// Add audience
	experiment.Audiences = make([]e.Audience, 1, 1)
	experiment.Audiences[0] = *audience

	return experiment
}
