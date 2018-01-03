package experiment

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Constants for evaluating tests
const (
	maxIterations = 100
)

// Testing ---------------------------------------------------------------------------------
func testSimpleGetVariable(t *testing.T, config basicTestConfig) {
	userID := "some_user_id"

	// Create evenly distributed weights and one weights array that forces a value
	weights1 := make([]uint32, config.numberWeights(), config.numberWeights())
	weights2 := make([]uint32, config.numberWeights(), config.numberWeights())

	for i := range weights1 {
		weights1[i] = 1
	}

	weights2[config.numberWeights()-1] = 1

	// Create the experiments
	builder1 := NewSimpleBuilder("experiment_1")
	builder2 := NewSimpleBuilder("experiment_2")

	switch config.variableType {
	case Float:
		builder1.AddFloat("variable_1", weights1, config.floats)
		builder2.AddFloat("variable_2", weights2, config.floats)
	case Int:
		builder1.AddInt("variable_1", weights1, config.ints)
		builder2.AddInt("variable_2", weights2, config.ints)
	case Bool:
		builder1.AddBool("variable_1", weights1, config.bools)
		builder2.AddBool("variable_2", weights2, config.bools)
	}

	experiment1, eErr1 := builder1.Build()
	experiment2, eErr2 := builder2.Build()

	// An experiment should have been built
	assert.Nil(t, eErr1)
	assert.Nil(t, eErr2)

	// Load the experiment service
	service := NewExperimentService()
	service.Reload([]Experiment{*experiment1, *experiment2})

	// We should find 'variable_1' and it should be associated with 'experiment_1'
	result1, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result1.Experiment.Name)
	assert.Nil(t, err1)

	// Test the second experiment for many different user ids
	for i := 0; i < maxIterations; i++ {
		// We should find 'variable_2' and it should be associated with 'experiment_2'
		result2, err2 := service.GetVariable("variable_2", makeUserID(i), nil)
		assert.Equal(t, "experiment_2", result2.Experiment.Name)
		assert.Nil(t, err2)

		// We put on the weight on the last variant. Ensure the last variant is returned in the result
		switch config.variableType {
		case Float:
			assert.Equal(t, config.floats[config.numberWeights()-1], result2.Value.FloatValue)
		case Int:
			assert.Equal(t, config.ints[config.numberWeights()-1], result2.Value.IntValue)
		case Bool:
			assert.Equal(t, config.bools[config.numberWeights()-1], result2.Value.BoolValue)
		}
	}
	// A fake variable should return an error and nils for the result
	result, err3 := service.GetVariable("fake_variable", userID, nil)
	assert.Nil(t, result.Experiment)
	assert.Nil(t, result.Audience)
	assert.Nil(t, result.Value)
	assert.NotNil(t, err3)

	// Make sure experiments that are reloaded have the same randomization
	testReload(t, builder1, "variable_1")
}

func testReload(t *testing.T, builder BasicBuilder, variableName string) {

	var firstValue *Value = nil

	// Repeatedly reload the service with the same experiment. The random value should be the same each time
	for i := 0; i < maxIterations; i++ {
		experiment, _ := builder.Build()
		service := NewExperimentService()
		service.Reload([]Experiment{*experiment})

		result, _ := service.GetVariable(variableName, "", nil)

		if firstValue == nil {
			firstValue = result.Value
		} else {
			assert.Equal(t, firstValue.IntValue, result.Value.IntValue)
			assert.Equal(t, firstValue.FloatValue, result.Value.FloatValue)
			assert.Equal(t, firstValue.BoolValue, result.Value.BoolValue)
		}
	}
}

func TestFactorial(t *testing.T) {
	// Create a factorial experiment with two int variables
	ints1 := []int64{2, 3, 5}
	ints2 := []int64{7, 11}

	builder := NewFactorialBuilder("experiment_1")
	builder.AddInt("int_1", []uint32{1, 1, 1}, ints1)
	builder.AddInt("int_2", []uint32{1, 1}, ints2)
	experiment, err := builder.Build()

	// This should be a valid experiment
	assert.Nil(t, err, "experiment should be built")

	// Prepare the service
	service := NewExperimentService()
	service.Reload([]Experiment{*experiment})

	// Count how each value shows up
	countMap := make(map[int64]int)
	for i := 0; i < maxIterations; i++ {
		uuid := makeUserID(i)
		result1, err1 := service.GetVariable("int_1", uuid, nil)
		result2, err2 := service.GetVariable("int_2", uuid, nil)

		assert.Nil(t, err1, "int_1 variable should exist")
		assert.Nil(t, err2, "int_2 variable should exist")

		intValue1 := result1.Value.IntValue
		intValue2 := result2.Value.IntValue

		combination := intValue1 * intValue2

		if _, ok := countMap[combination]; !ok {
			countMap[combination] = 0
		}

		countMap[combination]++
	}

	assert.Equal(t, len(countMap), len(ints1)*len(ints2))

	// Make sure experiments that are reloaded have the same randomization
	testReload(t, builder, "int_1")
}

func TestNewAlignedBuilderFail(t *testing.T) {
	// Create a factorial experiment with two int variables
	ints1 := []int64{2, 3, 5}
	ints2 := []int64{7, 11}

	builder := NewAlignedBuilder("experiment_1")
	err1 := builder.AddInt("int_1", []uint32{1, 1, 1}, ints1)
	err2 := builder.AddInt("int_2", []uint32{1, 1}, ints2)

	_, err := builder.Build()

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
	assert.Nil(t, err)
}

func TestNewAlignedBuilderSuccess(t *testing.T) {
	// Create a factorial experiment with two int variables
	ints1 := []int64{1, 2, 3}
	ints2 := []int64{1, 2, 3}

	builder := NewAlignedBuilder("experiment_1")
	builder.AddInt("int_1", []uint32{10, 1, 1}, ints1)
	builder.AddInt("int_2", []uint32{10, 1, 1}, ints2)
	experiment, err := builder.Build()

	// This should be a valid experiment
	assert.Nil(t, err, "experiment should built because they are aligned")

	// Prepare the service
	service := NewExperimentService()
	service.Reload([]Experiment{*experiment})

	countMap := make(map[int64]int)
	for i := 0; i < maxIterations; i++ {
		uuid := makeUserID(i)
		result1, err1 := service.GetVariable("int_1", uuid, nil)
		result2, err2 := service.GetVariable("int_2", uuid, nil)

		assert.Nil(t, err1, "int_1 variable should exist")
		assert.Nil(t, err2, "int_2 variable should exist")

		intValue1 := result1.Value.IntValue
		intValue2 := result2.Value.IntValue

		combination := intValue1 * intValue2

		if _, ok := countMap[combination]; !ok {
			countMap[combination] = 0
		}

		assert.Equal(t, intValue1, intValue2)

		countMap[combination]++
	}

	assert.Equal(t, len(countMap), len(ints1))

	// Make sure experiments that are reloaded have the same randomization
	testReload(t, builder, "int_1")
}

func TestSimpleFloats(t *testing.T) {
	config := basicTestConfig{}
	config.variableType = Float
	config.floats = []float64{1.0, 2.0, 3.0}

	testSimpleGetVariable(t, config)
}

func TestSimpleInts(t *testing.T) {
	config := basicTestConfig{}
	config.variableType = Int
	config.ints = []int64{1, 2, 3}

	testSimpleGetVariable(t, config)
}

func TestSimpleBools(t *testing.T) {
	config := basicTestConfig{}
	config.variableType = Bool
	config.bools = []bool{true, false}

	testSimpleGetVariable(t, config)
}

func TestSimpleFullyDistributed(t *testing.T) {
	floatValues := []float64{1.0, 2.0, 3.0}

	// Simple experiment 1, evenly distributed
	builder := NewSimpleBuilder("experiment_1")
	builder.AddFloat("variable_1", []uint32{1, 1, 1}, floatValues)
	experiment, _ := builder.Build()

	// Load the experiment service
	service := NewExperimentService()
	service.Reload([]Experiment{*experiment})

	// Create map for counting values
	valueMap := make(map[float64]int)

	// Count how many times each outcome occurs
	for i := 0; i < maxIterations; i++ {
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

func TestSimpleNoWeights(t *testing.T) {
	userID := "some_user_id"

	builder := NewSimpleBuilder("experiment_1")
	builder.AddFloat("variable_1", []uint32{0, 0, 0}, []float64{1.0, 2.0, 3.0})
	experiment1, _ := builder.Build()

	// Assert that the created experiments are valid
	valid1 := experiment1.Validate()

	assert.Nil(t, valid1)

	// Load the experiment service
	service := NewExperimentService()
	service.Reload([]Experiment{*experiment1})

	// Check variable_1. Should be control
	result, err1 := service.GetVariable("variable_1", userID, nil)
	assert.Equal(t, "experiment_1", result.Experiment.Name)
	assert.Equal(t, 1.0, result.Value.FloatValue)
	assert.Nil(t, err1)
}

func TestSimpleFail(t *testing.T) {
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

// Test Helpers ---------------------------------------------------------------------------------

// A simple type to help build test configs
type variableType int

const (
	Float variableType = iota
	Int
	Bool
)

// A struct to help generate tests
type basicTestConfig struct {
	variableType variableType
	floats       []float64
	ints         []int64
	bools        []bool
}

func (c *basicTestConfig) numberWeights() int {
	maximum := maxInt(len(c.floats), len(c.ints))
	return maxInt(maximum, len(c.bools))
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func makeUserID(index int) string {
	return fmt.Sprintf("user_id_%d", index)
}
