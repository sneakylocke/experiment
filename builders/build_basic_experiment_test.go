package builders

import (
	"fmt"
	e "github.com/sneakylocke/experiment"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	MAX_ITERATIONS = 100
)

func TestNoWeights(t *testing.T) {

	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	experiment1, _ := NewBasicFloatExperiment(
		"experiment_1",
		"variable_1",
		[]uint32{0, 0, 0},
		[]float64{1.0, 2.0, 3.0})

	// Assert that the created experiments are valid
	valid1, _ := experiment1.Validate()

	assert.True(t, valid1)

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment1})

	// Check variable_1. Should be control
	result, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result.Experiment.Name)
	assert.Equal(t, 1.0, result.Value.FloatValue)
	assert.Nil(t, err1)

}

func TestFail(t *testing.T) {
	experiment1, err1 := NewBasicFloatExperiment("experiment_1", "variable_1", []uint32{1}, []float64{1.0, 2.0, 3.0})
	assert.Nil(t, experiment1)
	assert.NotNil(t, err1)

	experiment2, err2 := NewBasicIntExperiment("experiment_1", "variable_1", []uint32{}, []int64{1, 2, 3})
	assert.Nil(t, experiment2)
	assert.NotNil(t, err2)

	experiment3, err3 := NewBasicBoolExperiment("experiment_1", "variable_1", []uint32{1, 1, 1}, []bool{false, true})
	assert.Nil(t, experiment3)
	assert.NotNil(t, err3)
}

func TestGetFloatVariable(t *testing.T) {
	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	experiment1, _ := NewBasicFloatExperiment(
		"experiment_1",
		"variable_1",
		[]uint32{1, 1, 1},
		[]float64{1.0, 2.0, 3.0})

	// Simple experiment 2, forces the value of 6
	experiment2, _ := NewBasicFloatExperiment(
		"experiment_2",
		"variable_2",
		[]uint32{0, 0, 1},
		[]float64{4.0, 5.0, 6.0})

	// Assert that the created experiments are valid
	valid1, _ := experiment1.Validate()
	valid2, _ := experiment2.Validate()

	assert.True(t, valid1)
	assert.True(t, valid2)

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment1, *experiment2})

	// Check variable_1
	result, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result.Experiment.Name)
	assert.Nil(t, err1)

	// Check variable_2
	result, err2 := service.GetVariable("variable_2", userID, nil)
	assert.Equal(t, "experiment_2", result.Experiment.Name)
	assert.Equal(t, 6.0, result.Value.FloatValue)
	assert.Nil(t, err2)

	// Check fake variable
	result, err3 := service.GetVariable("fake_variable", userID, nil)
	assert.Nil(t, result.Experiment)
	assert.Nil(t, result.Audience)
	assert.Nil(t, result.Value)
	assert.NotNil(t, err3)
}

func TestGetBoolVariable(t *testing.T) {
	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	experiment1, _ := NewBasicBoolExperiment(
		"experiment_1",
		"variable_1",
		[]uint32{1, 1},
		[]bool{true, false})

	// Simple experiment 2, forces the value of 6
	experiment2, _ := NewBasicBoolExperiment(
		"experiment_2",
		"variable_2",
		[]uint32{0, 1},
		[]bool{true, false})

	// Assert that the created experiments are valid
	valid1, _ := experiment1.Validate()
	valid2, _ := experiment2.Validate()

	assert.True(t, valid1)
	assert.True(t, valid2)

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment1, *experiment2})

	// Check variable_1
	result, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result.Experiment.Name)
	assert.Nil(t, err1)

	// Check variable_2
	result, err2 := service.GetVariable("variable_2", userID, nil)
	assert.Equal(t, "experiment_2", result.Experiment.Name)
	assert.Equal(t, false, result.Value.BoolValue)
	assert.Nil(t, err2)

	// Check fake variable
	result, err3 := service.GetVariable("fake_variable", userID, nil)
	assert.Nil(t, result.Experiment)
	assert.Nil(t, result.Audience)
	assert.Nil(t, result.Value)
	assert.NotNil(t, err3)
}

func TestGetIntVariable(t *testing.T) {
	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	experiment1, _ := NewBasicIntExperiment(
		"experiment_1",
		"variable_1",
		[]uint32{1, 1, 1},
		[]int64{1, 2, 3})

	// Simple experiment 2, forces the value of 6
	experiment2, _ := NewBasicIntExperiment(
		"experiment_2",
		"variable_2",
		[]uint32{0, 0, 1},
		[]int64{4, 5, 6})

	// Assert that the created experiments are valid
	valid1, _ := experiment1.Validate()
	valid2, _ := experiment2.Validate()

	assert.True(t, valid1)
	assert.True(t, valid2)

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment1, *experiment2})

	// Check variable_1
	result, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result.Experiment.Name)
	assert.Nil(t, err1)

	// Check variable_2
	result, err2 := service.GetVariable("variable_2", userID, nil)
	assert.Equal(t, "experiment_2", result.Experiment.Name)
	assert.Equal(t, int64(6), result.Value.IntValue)
	assert.Nil(t, err2)

	// Check fake variable
	result, err3 := service.GetVariable("fake_variable", userID, nil)
	assert.Nil(t, result.Experiment)
	assert.Nil(t, result.Audience)
	assert.Nil(t, result.Value)
	assert.NotNil(t, err3)
}

func TestGetVariableDistribution(t *testing.T) {
	floatValues := []float64{1.0, 2.0, 3.0}

	// Simple experiment 1, evenly distributed
	experiment, _ := NewBasicFloatExperiment(
		"experiment_1",
		"velocity",
		[]uint32{1, 1, 1},
		floatValues)

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment})

	// Create map for counting values
	valueMap := make(map[float64]int)

	// Count how many times each outcome occurs
	for i := 0; i < MAX_ITERATIONS; i++ {
		userID := makeUserID(i)

		result, _ := service.GetVariable("velocity", userID, nil)

		if _, ok := valueMap[result.Value.FloatValue]; !ok {
			valueMap[result.Value.FloatValue] = 0
		}

		valueMap[result.Value.FloatValue]++
	}

	// Make sure all possible outcomes can occure
	assert.True(t, valueMap[1.0] > 0)
	assert.True(t, valueMap[2.0] > 0)
	assert.True(t, valueMap[3.0] > 0)
}

func makeUserID(index int) string {
	return fmt.Sprintf("user_id_%d", index)
}
