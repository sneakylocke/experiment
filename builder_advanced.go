package experiment

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sneakylocke/experiment/constraint"
)

const (
	maximumVariables = 64
)

type AdvancedBuilder interface {
	Build() (*Experiment, error)
	AddFloats(variableName string, audienceName string, weights []uint32, values []float64) error
	AddInts(variableName string, audienceName string, weights []uint32, values []int64) error
	AddBools(variableName string, audienceName string, weights []uint32, values []bool) error

	AddConstraint(audienceName string, constraint *constraint.Constraint) error
}

type advancedBuilder struct {
	Audiences []*Audience

	ExperimentName   string
	MaximumVariables int
	IsFactorial      bool
	NumberValues     int
	FirstWeights     []uint32
}

func NewAdvancedBuilder(experimentName string) AdvancedBuilder {
	b := &advancedBuilder{}
	b.Audiences = make([]*Audience, 0, 10)
	b.ExperimentName = experimentName
	b.MaximumVariables = 1
	b.IsFactorial = false

	return b
}

func NewAdvancedFactorialBuilder(experimentName string) AdvancedBuilder {
	b := &advancedBuilder{}
	b.Audiences = make([]*Audience, 0, 10)
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = true

	return b
}

func NewAdvancedAlignedBuilder(experimentName string) AdvancedBuilder {
	b := &advancedBuilder{}
	b.Audiences = make([]*Audience, 0, 10)
	b.ExperimentName = experimentName
	b.MaximumVariables = maximumVariables
	b.IsFactorial = false

	return b
}

func (b *advancedBuilder) AddFloats(variableName string, audienceName string, weights []uint32, values []float64) error {
	if err := b.preValidate(variableName, audienceName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add floats")
	}

	b.setup(NewFloatValueGroup(variableName, weights, values), weights, variableName, audienceName)

	return nil
}

func (b *advancedBuilder) AddInts(variableName string, audienceName string, weights []uint32, values []int64) error {
	if err := b.preValidate(variableName, audienceName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add ints")
	}

	b.setup(NewIntValueGroup(variableName, weights, values), weights, variableName, audienceName)

	return nil
}

func (b *advancedBuilder) AddBools(variableName string, audienceName string, weights []uint32, values []bool) error {
	if err := b.preValidate(variableName, audienceName, weights, len(values)); err != nil {
		return errors.Annotate(err, "could not add bools")
	}

	b.setup(NewBoolValueGroup(variableName, weights, values), weights, variableName, audienceName)

	return nil
}

func (b *advancedBuilder) AddConstraint(audienceName string, constraint *constraint.Constraint) error {
	if constraint == nil {
		return errors.Errorf("cannot add a nil constraint")
	}

	for _, audience := range b.Audiences {
		if audience.Name == audienceName {
			audience.Constraints = append(audience.Constraints, *constraint)
			return nil
		}
	}

	return errors.Errorf("could not find existing audience with name: %s", audienceName)
}

func (b *advancedBuilder) Build() (*Experiment, error) {
	if postValidateErr := b.postValidate(); postValidateErr != nil {
		return nil, errors.Annotate(postValidateErr, "could not build experiment")
	}

	experiment := &Experiment{}

	// Setup the simpler aspects of the experiment
	experiment.Name = b.ExperimentName
	experiment.Salt = b.ExperimentName
	experiment.VariableNames = []string{}
	experiment.Enabled = true

	// Add audiences and find names
	names := make(map[string]bool)

	experiment.Audiences = make([]Audience, 0, 10)
	for _, audience := range b.Audiences {
		experiment.Audiences = append(experiment.Audiences, *audience)

		// Collect the variable names
		for _, valueGroup := range audience.ValueGroups {
			names[valueGroup.Name] = true
		}
	}

	// Calculate all available variable names
	for name := range names {
		experiment.VariableNames = append(experiment.VariableNames, name)
	}

	// Do extra preparations for an aligned experiment
	if !b.IsFactorial {
		// Aligned experiments need the value group salt to be the same
		for _, audience := range experiment.Audiences {
			for _, valueGroup := range audience.ValueGroups {
				valueGroup.Salt = experiment.Salt
			}
		}
	}

	if validateErr := experiment.Validate(); validateErr != nil {
		return nil, errors.Annotate(validateErr, "could not Validate and build experiment")
	}

	return experiment, nil
}

func (b *advancedBuilder) setup(valueGroup *ValueGroup, weights []uint32, variableName string, audienceName string) {
	// Load or create an audience
	audience := b.getAudience(audienceName)

	// Set control value to first element
	valueGroup.ControlValue = valueGroup.WeightedValues[0].Value

	// Setup Audience
	audience.ValueGroups[variableName] = valueGroup

	if b.FirstWeights == nil {
		b.FirstWeights = weights
	}
}

func (b *advancedBuilder) preValidate(variableName string, audienceName string, weights []uint32, numberValues int) error {
	// Load or create an audience
	audience := b.getAudience(audienceName)

	numberWeights := len(weights)

	// How many variables will there be if we add this variable
	numberVariables := len(audience.ValueGroups) + 1

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
	if _, found := audience.ValueGroups[variableName]; found {
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

func (b *advancedBuilder) postValidate() error {
	// Validate each audience constructed
	for _, audience := range b.Audiences {
		err := b.postValidateAudience(audience)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *advancedBuilder) postValidateAudience(audience *Audience) error {
	if len(audience.ValueGroups) == 0 {
		return nil
	}

	if !b.IsFactorial {
		length := 0
		for _, valueGroup := range audience.ValueGroups {
			currentLength := len(valueGroup.WeightedValues)

			if length != 0 && length != currentLength {
				return errors.New("aligned experiments need the same number of values")
			}

			length = currentLength
		}
	}

	return nil
}

func (b *advancedBuilder) getAudience(audienceName string) *Audience {
	// Try to find already existing audience
	for _, audience := range b.Audiences {
		if audience.Name == audienceName {
			return audience
		}
	}

	// If audience was not found, make one and add to the list
	audience := NewAudience()
	audience.Name = audienceName

	b.Audiences = append(b.Audiences, audience)

	return audience
}
