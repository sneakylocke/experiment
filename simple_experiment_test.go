package experiment

import (
	"testing"
)

func Test(t *testing.T) {
	// User info
	userID := "knapikt1131"

	// Create an experiment
	name := "simple_experiment"
	variableName := "velocity"
	weights := []uint32{1, 1, 1}
	values := []float64{1, 2, 3}
	experiment := NewSimpleExperiment(name, variableName, weights, values)

	experiments := make([]Experiment, 1)
	experiments = append(experiments, *experiment)

	service := NewExperimentService()
	service.Reload(experiments)
	service.GetVariable(variableName, userID, nil)
}
