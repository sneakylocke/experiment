package builders

import (
	"fmt"
	e "github.com/sneakylocke/experiment/src/experiment"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	MAX_ITERATIONS = 100
)

func TestGetVariable(t *testing.T) {
	// User info
	userID := "some_user_id"

	// Simple experiment 1, evenly distributed
	experiment1 := NewSimpleExperiment(
		"experiment_1",
		"velocity",
		[]uint32{1, 1, 1},
		[]float64{1.0, 2.0, 3.0})

	// Simple experiment 2, forces the value of 6
	experiment2 := NewSimpleExperiment(
		"experiment_2",
		"position",
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

	// Check velocity
	e1, _, _, err1 := service.GetVariable("velocity", userID, nil)
	assert.Equal(t, "experiment_1", e1.Name)
	assert.Nil(t, err1)

	// Check position
	e2, _, v2, err2 := service.GetVariable("position", userID, nil)
	assert.Equal(t, "experiment_2", e2.Name)
	assert.Equal(t, 6.0, v2.FloatValue)
	assert.Nil(t, err2)

	// Check fake variable
	e3, a3, v3, err3 := service.GetVariable("fake_variable", userID, nil)
	assert.Nil(t, e3)
	assert.Nil(t, a3)
	assert.Nil(t, v3)
	assert.NotNil(t, err3)
}

func TestGetVariableDistribution(t *testing.T) {
	floatValues := []float64{1.0, 2.0, 3.0}

	// Simple experiment 1, evenly distributed
	experiment := NewSimpleExperiment(
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

		_, _, v, _ := service.GetVariable("velocity", userID, nil)

		if _, ok := valueMap[v.FloatValue]; !ok {
			valueMap[v.FloatValue] = 0
		}

		valueMap[v.FloatValue]++
	}

	// Make sure all possible outcomes can occure
	assert.True(t, valueMap[1.0] > 0)
	assert.True(t, valueMap[2.0] > 0)
	assert.True(t, valueMap[3.0] > 0)
}

func makeUserID(index int) string {
	return fmt.Sprintf("user_id_%d", index)
}
