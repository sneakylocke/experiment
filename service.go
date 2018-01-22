package experiment

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/sneakylocke/experiment/constraint"
	"hash/fnv"
)

const (
	denominator = 10000
)

type GetVariableResult struct {
	Experiment *Experiment
	Audience   *Audience
	Value      *Value
}

type Service interface {
	Reload(experiments []Experiment) error
	GetVariable(name string, userID string, context constraint.Context) (*GetVariableResult, error)
}

type service struct {
	resolver    constraint.Resolver
	experiments []Experiment
	variableMap map[string][]Experiment
}

func NewService() *service {
	service := &service{}
	service.resolver = constraint.NewDefaultResolver()
	service.experiments = make([]Experiment, 0)
	service.variableMap = make(map[string][]Experiment)

	return service
}

func (service *service) Reload(experiments []Experiment) error {
	service.experiments = experiments
	service.variableMap = make(map[string][]Experiment)

	for _, experiment := range experiments {
		for _, variableName := range experiment.VariableNames {
			if service.variableMap[variableName] == nil {
				service.variableMap[variableName] = make([]Experiment, 0, 1)
			}

			service.variableMap[variableName] = append(service.variableMap[variableName], experiment)
		}
	}

	return nil
}

func (service *service) GetVariable(variableName string, userID string, context constraint.Context) (*GetVariableResult, error) {
	experiments, experimentsOk := service.variableMap[variableName]

	if !experimentsOk {
		return nil, errors.Errorf("no experiment matching variable '%s'", variableName)
	}

	for _, experiment := range experiments {
		if !experiment.Enabled {
			continue
		}

		if experiment.Audiences == nil {
			continue
		}

		for _, audience := range experiment.Audiences {
			if !audience.Enabled {
				continue
			}

			// By default the constraints are met
			constraintsMet := true

			// If we actually have constraints for this experiment, check them
			if audience.Constraints != nil {

				// All constraints must be passed
				for _, constraint := range audience.Constraints {
					resolveOk, _ := service.resolver.Resolve(&constraint, context)

					// TODO: Log the error. We don't want to stop evaluating so we continue

					if !resolveOk {
						constraintsMet = false
						break
					}
				}
			}

			// If constraints are met extract the value
			if constraintsMet {
				value, err := service.getVariable(&experiment, &audience, variableName, userID)

				if err == nil {
					return &GetVariableResult{Experiment: &experiment, Audience: &audience, Value: value}, nil
				} else {
					return nil, errors.Annotatef(err, "error getting variable")
				}
			}
		}
	}

	return nil, errors.New("failed to find variable or could not meet constraints with given context")
}

func (service *service) getVariable(experiment *Experiment, audience *Audience, variableName string, userID string) (*Value, error) {
	valueGroup, ok := audience.ValueGroups[variableName]

	if !ok {
		return nil, errors.New("failed to find value group for variable name")
	}

	// Create a hash string from salts + userID
	hashString := experiment.Salt + valueGroup.Salt + userID
	hashNumber := getHash(hashString)

	// Check if the exposure indicates we should be in control
	fraction := float64(hashNumber%denominator) / denominator

	// Return the control value if there is not enough exposure for this user
	if fraction > audience.Exposure {
		return &valueGroup.ControlValue, nil
	}

	// Build a distribution in order to randomize which value is returned
	var weightSum uint32 = 0
	weights := make([]uint32, len(valueGroup.WeightedValues))

	for i, value := range valueGroup.WeightedValues {
		weightSum += value.Weight
		weights[i] = weightSum
	}

	// Return control if there are no weights
	if weightSum == 0 {
		log.Warnf("Weight sum is 0, returning the control value")
		return &valueGroup.ControlValue, nil
	}

	// Create a valueGroup index based on experiment, variable, and user
	hash := hashNumber % weightSum

	// Find the appropriate value to return based on the hash
	for i, weight := range weights {
		if hash < weight {
			return &valueGroup.WeightedValues[i].Value, nil
		}
	}

	return nil, errors.New("failed to find value")
}

func getHash(s string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return hash.Sum32()
}
