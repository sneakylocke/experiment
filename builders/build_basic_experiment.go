package builders

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	e "github.com/sneakylocke/experiment"
)

func NewBasicFloatExperiment(experimentName string, variableName string, weights []uint32, values []float64) (*e.Experiment, error) {
	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"experiment": experimentName, "variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return nil, errors.Errorf("poorly formed weights or values")
	}

	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewFloatValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	// Make sure experiment is valid before returning
	if _, err := experiment.Validate(); err != nil {
		return nil, errors.Annotate(err, "could not build basic float experiment")
	}

	return experiment, nil
}

func NewBasicIntExperiment(experimentName string, variableName string, weights []uint32, values []int64) (*e.Experiment, error) {
	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"experiment": experimentName, "variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return nil, errors.Errorf("poorly formed weights or values")
	}

	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewIntValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	// Make sure experiment is valid before returning
	if _, err := experiment.Validate(); err != nil {
		return nil, errors.Annotate(err, "could not build basic float experiment")
	}

	return experiment, nil
}

func NewBasicBoolExperiment(experimentName string, variableName string, weights []uint32, values []bool) (*e.Experiment, error) {
	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"experiment": experimentName, "variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return nil, errors.Errorf("poorly formed weights or values")
	}

	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewBoolValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	// Make sure experiment is valid before returning
	if _, err := experiment.Validate(); err != nil {
		return nil, errors.Annotate(err, "could not build basic float experiment")
	}

	return experiment, nil
}

func NewBasicArbitraryExperiment(experimentName string, variableName string, weights []uint32, values []interface{}) (*e.Experiment, error) {
	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"experiment": experimentName, "variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return nil, errors.Errorf("poorly formed weights or values")
	}

	experiment := newExperiment(experimentName, variableName)

	// Setup ValueGroup
	valueGroup := e.NewArbitraryValueGroup(variableName, weights, values)

	// Setup Audience
	experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	// Make sure experiment is valid before returning
	if _, err := experiment.Validate(); err != nil {
		return nil, errors.Annotate(err, "could not build basic float experiment")
	}

	return experiment, nil
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
