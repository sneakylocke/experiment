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
	userID := "some_user_id"

	builder := NewSimpleBuilder("experiment_1")
	builder.AddFloat("variable_1", []uint32{0, 0, 0}, []float64{1.0, 2.0, 3.0})
	experiment1, _ := builder.Build()

	// Assert that the created experiments are valid
	valid1 := experiment1.Validate()

	assert.Nil(t, valid1)

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
	builder1 := NewSimpleBuilder("experiment_1")
	builder1.AddFloat("variable_1", []uint32{1}, []float64{1.0, 2.0, 3.0})
	experiment1, err1 := builder1.Build()

	assert.Nil(t, experiment1)
	assert.NotNil(t, err1)

	builder2 := NewSimpleBuilder("experiment_1")
	builder2.AddInt("variable_1", []uint32{}, []int64{1, 2, 3})
	experiment2, err2 := builder2.Build()

	assert.Nil(t, experiment2)
	assert.NotNil(t, err2)

	builder3 := NewSimpleBuilder("experiment_1")
	builder3.AddBool("variable_1", []uint32{1, 1, 1}, []bool{false, true})
	experiment3, err3 := builder3.Build()

	assert.Nil(t, experiment3)
	assert.NotNil(t, err3)
}

func TestGetFloatVariable(t *testing.T) {
	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	builder1 := NewSimpleBuilder("experiment_1")
	builder1.AddFloat("variable_1", []uint32{1, 1, 1}, []float64{1.0, 2.0, 3.0})
	experiment1, eErr1 := builder1.Build()

	// Simple experiment 2, forces the value of 6
	builder2 := NewSimpleBuilder("experiment_2")
	builder2.AddFloat("variable_2", []uint32{0, 0, 1}, []float64{4.0, 5.0, 6.0})
	experiment2, eErr2 := builder2.Build()

	assert.Nil(t, eErr1)
	assert.Nil(t, eErr2)

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
	builder1 := NewSimpleBuilder("experiment_1")
	builder1.AddBool("variable_1", []uint32{1, 1}, []bool{true, false})
	experiment1, eErr1 := builder1.Build()

	// Simple experiment 2, forces the value of false
	builder2 := NewSimpleBuilder("experiment_2")
	builder2.AddBool("variable_2", []uint32{0, 1}, []bool{true, false})
	experiment2, eErr2 := builder2.Build()

	assert.Nil(t, eErr1)
	assert.Nil(t, eErr2)

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
	builder1 := NewSimpleBuilder("experiment_1")
	builder1.AddInt("variable_1", []uint32{1, 1, 1}, []int64{1, 2, 3})
	experiment1, eErr1 := builder1.Build()

	// Simple experiment 2, forces the value of 6
	builder2 := NewSimpleBuilder("experiment_2")
	builder2.AddInt("variable_2", []uint32{0, 0, 1}, []int64{4, 5, 6})
	experiment2, eErr2 := builder2.Build()

	assert.Nil(t, eErr1)
	assert.Nil(t, eErr2)

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
	builder := NewSimpleBuilder("experiment_1")
	builder.AddFloat("variable_1", []uint32{1, 1, 1}, floatValues)
	experiment, _ := builder.Build()

	// Load the experiment service
	service := e.NewExperimentService()
	service.Reload([]e.Experiment{*experiment})

	// Create map for counting values
	valueMap := make(map[float64]int)

	// Count how many times each outcome occurs
	for i := 0; i < MAX_ITERATIONS; i++ {
		userID := makeUserID(i)

		result, _ := service.GetVariable("variable_1", userID, nil)

		if _, ok := valueMap[result.Value.FloatValue]; !ok {
			valueMap[result.Value.FloatValue] = 0
		}

		valueMap[result.Value.FloatValue]++
	}

	// Make sure all possible outcomes can occur
	assert.True(t, valueMap[1.0] > 0)
	assert.True(t, valueMap[2.0] > 0)
	assert.True(t, valueMap[3.0] > 0)
}

func makeUserID(index int) string {
	return fmt.Sprintf("user_id_%d", index)
}
