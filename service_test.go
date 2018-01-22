package experiment

import (
	"encoding/json"
	"github.com/sneakylocke/experiment/constraint"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestAudience1(t *testing.T) {
	context := make(map[string]interface{})
	context["country"] = "USA"
	context["temperature"] = 75
	mapContext := constraint.NewMapContext(context)

	testAudience(t, "audience_1", "testdata/experiments/constraints_test_1.json", "a", mapContext)
}

func TestAudience2(t *testing.T) {
	context := make(map[string]interface{})
	context["country"] = "ITALY"
	context["food"] = "banana"
	mapContext := constraint.NewMapContext(context)

	testAudience(t, "audience_2", "testdata/experiments/constraints_test_1.json", "a", mapContext)
}

func TestAudience3(t *testing.T) {
	context := make(map[string]interface{})
	mapContext := constraint.NewMapContext(context)
	context["country"] = "NOT_ITALY"
	context["food"] = "not_banana"

	testAudience(t, "audience_3", "testdata/experiments/constraints_test_1.json", "a", mapContext)
}

func testAudience(t *testing.T, expectedAudienceName string, fileName string, variableName string, context constraint.Context) {
	experiment := loadExperiment(t, fileName)
	service := NewService()
	service.Reload([]Experiment{*experiment})

	result, err := service.GetVariable(variableName, "userID", context)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAudienceName, result.Audience.Name)

}

func loadExperiment(t *testing.T, file string) *Experiment {
	data, err := ioutil.ReadFile(file)

	assert.Nil(t, err)

	experiment := &Experiment{}
	unmarshalErr := json.Unmarshal(data, experiment)

	assert.Nil(t, unmarshalErr)
	assert.Nil(t, experiment.Validate())

	return experiment
}
