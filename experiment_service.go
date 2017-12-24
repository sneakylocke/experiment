package experiment

import (
	"errors"
	"hash/fnv"
)

const (
	denominator = 10000
)

type IExperimentService interface {
	Reload(experiments []Experiment) error
	GetVariable(name string, userID string, userInfo interface{}) (*Value, *Experiment, error)
}

type experimentService struct {
	constraintResolver ConstraintResolver
	experiments        []Experiment
	variableMap        map[string][]Experiment
}

func NewExperimentService() *experimentService {
	service := &experimentService{}
	service.constraintResolver = &BasicResolver{}
	service.experiments = make([]Experiment, 1)
	service.variableMap = make(map[string][]Experiment)

	return service
}

func (service *experimentService) Reload(experiments []Experiment) error {
	service.experiments = experiments
	service.variableMap = make(map[string][]Experiment)

	for _, experiment := range experiments {
		for _, variableName := range experiment.VariableNames {
			if service.variableMap[variableName] == nil {
				service.variableMap[variableName] = make([]Experiment, 1)
			}

			service.variableMap[variableName] = append(service.variableMap[variableName], experiment)
		}
	}

	return nil
}

func (service *experimentService) GetVariable(variableName string, userID string, userInfo interface{}) (*Value, *Experiment, error) {
	experiments, experimentsOk := service.variableMap[variableName]

	if !experimentsOk {
		return nil, nil, nil
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

			if service.constraintResolver.PassesConstraints(userInfo, audience.Constraints) {
				value, err := service.getVariable(&experiment, &audience, variableName, userID)

				if err == nil {
					return value, &experiment, err
				} else {
					return nil, nil, err
				}
			}
		}
	}

	return nil, nil, errors.New("failed to find variable")
}

func (service *experimentService) getVariable(experiment *Experiment, audience *Audience, variableName string, userID string) (*Value, error) {
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
		return &audience.ControlValue, nil
	}

	// Build a distribution in order to randomize which value is returned
	var weightSum uint32 = 0
	weights := make([]uint32, len(valueGroup.WeightedValues))

	for i, value := range valueGroup.WeightedValues {

		weightSum += value.Weight
		weights[i] = weightSum
	}

	// Create a valueGroup index based on experiment, variable, and user
	hash := hashNumber % uint32(weightSum)

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
