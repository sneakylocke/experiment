package experiment

import (
	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	maximumVariables = 64
)

type BasicBuilder interface {
	Build() (*Experiment, error)
	AddFloat(variableName string, weights []uint32, values []float64) error
	AddInt(variableName string, weights []uint32, values []int64) error
	AddBool(variableName string, weights []uint32, values []bool) error
}

type basicBuilder struct {
	Audience         *Audience
	ExperimentName   string
	MaximumVariables int
	IsFactorial      bool
	NumberValues     int
	FirstWeights     []uint32
}

func NewSimpleBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.Audience = NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = 1
	b.IsFactorial = false

	return b
}

func NewFactorialBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.Audience = NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = true

	return b
}

func NewAlignedBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.Audience = NewAudience()
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = false

	return b
}

func (b *basicBuilder) AddFloat(variableName string, weights []uint32, values []float64) error {
	if err := b.prevalidate(variableName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add floats")
	}

	b.setup(NewFloatValueGroup(variableName, weights, values), weights, variableName)

	return nil
}

func (b *basicBuilder) AddInt(variableName string, weights []uint32, values []int64) error {
	if err := b.prevalidate(variableName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add ints")
	}

	b.setup(NewIntValueGroup(variableName, weights, values), weights, variableName)

	return nil
}

func (b *basicBuilder) AddBool(variableName string, weights []uint32, values []bool) error {
	if err := b.prevalidate(variableName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add bools")
	}

	b.setup(NewBoolValueGroup(variableName, weights, values), weights, variableName)

	return nil
}

func (b *basicBuilder) Build() (*Experiment, error) {
	if postValidateErr := b.postValidate(); postValidateErr != nil {
		return nil, errors.Annotate(postValidateErr, "could not build experiment")
	}

	experiment := &Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = b.ExperimentName
	experiment.Salt = b.ExperimentName
	experiment.VariableNames = []string{}
	experiment.Enabled = true

	// Add audience
	experiment.Audiences = []Audience{*b.Audience}

	// Calculate all available variable names
	for _, valueGroup := range b.Audience.ValueGroups {
		experiment.VariableNames = append(experiment.VariableNames, valueGroup.Name)
	}

	// Do extra preparations for an aligned experiment
	if !b.IsFactorial {
		// Aligned experiments need the value group salt to be the same
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

func (b *basicBuilder) setup(valueGroup *ValueGroup, weights []uint32, variableName string) {
	// Set control value to first element
	valueGroup.ControlValue = valueGroup.WeightedValues[0].Value

	// Setup Audience
	b.Audience.ValueGroups[variableName] = valueGroup

	if b.FirstWeights == nil {
		b.FirstWeights = weights
	}
}

func (b *basicBuilder) postValidate() error {
	if len(b.Audience.ValueGroups) == 0 {
		return nil
	}

	if !b.IsFactorial {
		length := 0
		for _, valueGroup := range b.Audience.ValueGroups {
			currentLength := len(valueGroup.WeightedValues)

			if length != 0 && length != currentLength {
				return errors.New("aligned experiments need the same number of values")
			}

			length = currentLength
		}
	}

	return nil
}

func (b *basicBuilder) prevalidate(variableName string, weights []uint32, numberValues int) error {
	numberWeights := len(weights)

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
	if !b.IsFactorial && b.FirstWeights != nil {
		if len(b.FirstWeights) != len(weights) {
			logger.Errorf("aligned experiments should have the same number of weights")
			return errors.Errorf("aligned experiments should have the same number of weights")
		}

		for index, weight := range b.FirstWeights {
			if weights[index] != weight {
				logger.Errorf("aligned experiments should have the same weights")
				return errors.Errorf("aligned experiments should have the same weights")
			}
		}

	}

	return nil
}
