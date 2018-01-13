package experiment

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestValid1(t *testing.T) {
	testValid(t, "testdata/experiments/valid_1.json")
}

func TestInvalidNoAudience(t *testing.T) {
	testInvalid(t, "testdata/experiments/invalid_no_audience.json")
}

func TestInvalidMissingValueGroup(t *testing.T) {
	testInvalid(t, "testdata/experiments/invalid_missing_value_group.json")
}

func TestInvalidMissingVariableName(t *testing.T) {
	testInvalid(t, "testdata/experiments/invalid_missing_variable_name.json")
}

func testValid(t *testing.T, file string) {
	data, _ := ioutil.ReadFile(file)

	experiment := &Experiment{}
	json.Unmarshal(data, experiment)

	assert.Nil(t, experiment.Validate())
}

func testInvalid(t *testing.T, file string) {
	data, _ := ioutil.ReadFile(file)

	experiment := &Experiment{}
	json.Unmarshal(data, experiment)

	assert.NotNil(t, experiment.Validate())
}
