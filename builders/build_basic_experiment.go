package builders

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	e "github.com/sneakylocke/experiment"
)

type OneVariableBasicBuilder struct {
	Experiment *e.Experiment
	Sealed     bool // If true no more changes to the experiment can be made
}

func NewOneVariableBasicBuilder(experimentName string) BasicBuilder {
	b := &OneVariableBasicBuilder{}
	b.Experiment = NewBasicExperiment(experimentName)
	b.Sealed = false

	return b
}

func (b *OneVariableBasicBuilder) AddFloat(variableName string, weights []uint32, values []float64) error {
	if b.Sealed {
		return errors.Errorf("cannot add another variable to a single variable experiment")
	}

	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return errors.Errorf("poorly formed weights or values")
	}

	// Add variable name
	b.Experiment.VariableNames = append(b.Experiment.VariableNames, variableName)

	// Setup ValueGroup
	valueGroup := e.NewFloatValueGroup(variableName, weights, values)

	// Setup Audience
	b.Experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	b.Experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	b.Sealed = true

	return nil
}

func (b *OneVariableBasicBuilder) AddInt(variableName string, weights []uint32, values []int64) error {
	if b.Sealed {
		return errors.Errorf("cannot add another variable to a single variable experiment")
	}

	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return errors.Errorf("poorly formed weights or values")
	}

	// Add variable name
	b.Experiment.VariableNames = append(b.Experiment.VariableNames, variableName)

	// Setup ValueGroup
	valueGroup := e.NewIntValueGroup(variableName, weights, values)

	// Setup Audience
	b.Experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	b.Experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	b.Sealed = true

	return nil
}

func (b *OneVariableBasicBuilder) AddBool(variableName string, weights []uint32, values []bool) error {
	if b.Sealed {
		return errors.Errorf("cannot add another variable to a single variable experiment")
	}

	if len(weights) != len(values) || len(weights) == 0 {
		logger := log.WithFields(log.Fields{"variable": variableName, "weights length": len(weights), "values length": values})
		logger.Errorf("poorly formed weights or values")
		return errors.Errorf("poorly formed weights or values")
	}

	// Add variable name
	b.Experiment.VariableNames = append(b.Experiment.VariableNames, variableName)

	// Setup ValueGroup
	valueGroup := e.NewBoolValueGroup(variableName, weights, values)

	// Setup Audience
	b.Experiment.Audiences[0].ControlValue = valueGroup.WeightedValues[0].Value
	b.Experiment.Audiences[0].ValueGroups[variableName] = *valueGroup

	b.Sealed = true

	return nil
}

func (b *OneVariableBasicBuilder) Build() (*e.Experiment, error) {
	if !b.Sealed {
		return nil, errors.Errorf("A float, bool, or int array should be added before building")
	}

	if _, err := b.Experiment.Validate(); err != nil {
		return nil, errors.Annotate(err, "could not build experiment")
	}

	return b.Experiment, nil
}
