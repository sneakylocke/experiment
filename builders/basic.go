package builders

import (
	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	e "github.com/sneakylocke/experiment"
)

const (
	maximumVariables = 64
)

type BasicBuilder interface {
	Build() (*e.Experiment, error)
	AddFloat(variableName string, weights []uint32, values []float64) error
	AddInt(variableName string, weights []uint32, values []int64) error
	AddBool(variableName string, weights []uint32, values []bool) error
}

type factorialBuilder struct {
	Audience         *e.Audience
	ExperimentName   string
	MaximumVariables int
	IsFactorial      bool
	NumberValues     int
}

func NewSimpleBuilder(experimentName string) BasicBuilder {
	b := &factorialBuilder{}
	b.Audience = e.NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = 1
	b.IsFactorial = false

	return b
}

func NewFactorialBuilder(experimentName string) BasicBuilder {
	b := &factorialBuilder{}
	b.Audience = e.NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = true

	return b
}

func NewAlignedlBuilder(experimentName string) BasicBuilder {
	b := &factorialBuilder{}
	b.Audience = e.NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = false

	return b
}

func (b *factorialBuilder) AddFloat(variableName string, weights []uint32, values []float64) error {
	if err := b.prevalidate(variableName, len(weights), len(values)); err != nil {
		return errors.Annotate(err, "could not add floats")
	}

	b.setupAudience(e.NewFloatValueGroup(variableName, weights, values), variableName)

	return nil
}

func (b *factorialBuilder) AddInt(variableName string, weights []uint32, values []int64) error {
	if err := b.prevalidate(variableName, len(weights), len(values)); err != nil {
		return errors.Annotate(err, "could not add ints")
	}

	b.setupAudience(e.NewIntValueGroup(variableName, weights, values), variableName)

	return nil
}

func (b *factorialBuilder) AddBool(variableName string, weights []uint32, values []bool) error {
	if err := b.prevalidate(variableName, len(weights), len(values)); err != nil {
		return errors.Annotate(err, "could not add bools")
	}

	b.setupAudience(e.NewBoolValueGroup(variableName, weights, values), variableName)

	return nil
}

func (b *factorialBuilder) Build() (*e.Experiment, error) {
	experiment := &e.Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = b.ExperimentName
	experiment.Salt = b.ExperimentName
	experiment.VariableNames = []string{}
	experiment.Enabled = true

	// Add audience
	experiment.Audiences = []e.Audience{*b.Audience}

	// Calculate all available variable names
	for _, valueGroup := range b.Audience.ValueGroups {
		experiment.VariableNames = append(experiment.VariableNames, valueGroup.Name)
	}

	// Non-factorial experiments need the same salt so that the value groups align
	if !b.IsFactorial {
		salt := uuid.NewV4().String()
		for _, valueGroup := range b.Audience.ValueGroups {
			valueGroup.Salt = salt
		}
	}

	if validateErr := experiment.Validate(); validateErr != nil {
		return nil, errors.Annotate(validateErr, "could not validate and build experiment")
	}

	return experiment, nil
}

func (b *factorialBuilder) setupAudience(valueGroup *e.ValueGroup, variableName string) {
	// Set control value to first element
	valueGroup.ControlValue = valueGroup.WeightedValues[0].Value

	// Setup Audience
	b.Audience.ValueGroups[variableName] = *valueGroup
}

func (b *factorialBuilder) prevalidate(variableName string, numberWeights int, numberValues int) error {
	// How many variables will there be if we add this variable
	numberVariables := len(b.Audience.ValueGroups) + 1

	// Prepare logger
	logger := log.WithFields(log.Fields{
		"isFactorial":      b.IsFactorial,
		"maximumVariables": b.MaximumVariables,
		"numberValues":     b.NumberValues,
		"variable":         variableName,
		"numberVariables":  numberVariables,
		"weights length":   numberWeights,
		"values length":    numberValues})

	// Make sure we don't add the same variable twice
	if _, found := b.Audience.ValueGroups[variableName]; found {
		logger.Errorf("cannot set the same variable twice")
		return errors.Errorf("cannot set the same variable twice")
	}

	// Make sure we do not exceed maximum variables
	if numberVariables > b.MaximumVariables {
		logger.Errorf("poorly formed weights or values")
		return errors.Errorf("poorly formed weights or values")
	}

	// Weights and values must be of the same amount
	if numberWeights != numberValues || numberWeights == 0 {
		logger.Errorf("poorly formed weights or values")
		return errors.Errorf("poorly formed weights or values")
	}

	// If we have a non-factorial experiment the number of treatments must be equal
	if !b.IsFactorial && b.NumberValues > 0 && numberValues != b.NumberValues {
		logger.Errorf("non-factorial experiment value groups need the same number of treatments")
		return errors.Errorf("non-factorial experiment value groups need the same number of treatments")
	}

	return nil
}
