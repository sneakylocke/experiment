package builders

import (
	e "github.com/sneakylocke/experiment/src/experiment"
)

func NewBasicFloatExperiment(experimentName string, variableName string, weights []uint32, values []float64) *e.Experiment {
	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewFloatValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	return experiment
}

func NewBasicIntExperiment(experimentName string, variableName string, weights []uint32, values []int64) *e.Experiment {
	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewIntValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	return experiment
}

func NewBasicBoolExperiment(experimentName string, variableName string, weights []uint32, values []bool) *e.Experiment {
	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewBoolValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	return experiment
}

func NewBasicArbitraryExperiment(experimentName string, variableName string, weights []uint32, values []interface{}) *e.Experiment {
	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewArbitraryValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	return experiment
}

func newExperiment(experimentName string, variableName string) *e.Experiment {
	experiment := &e.Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = experimentName
	experiment.VariableNames = []string{variableName}
	experiment.Salt = experimentName
	experiment.Enabled = true

	// Setup Audience
	audience := e.NewAudience()

	// Add audience
	experiment.Audiences = make([]e.Audience, 1, 1)
	experiment.Audiences[0] = *audience

	return experiment
}
