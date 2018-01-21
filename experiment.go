package experiment

import "github.com/juju/errors"

type Experiment struct {
	Name          string     `json:"name"`          // Name of the experiment
	VariableNames []string   `json:"variableNames"` // Names of variables allowed to be randomized in audiences
	Audiences     []Audience `json:"audiences"`     // Details of variables are captured in a single audience. Users may belong to one audience
	Salt          string     `json:"salt"`
	Enabled       bool       `json:"enabled"`
}

type WeightedValue struct {
	Value  Value  `json:"value"`
	Weight uint32 `json:"weight"`
}

func (e *Experiment) Validate() error {
	if e.Name == "" {
		return errors.New("no name")
	}

	if e.Salt == "" {
		return errors.New("no salt")
	}

	if e.VariableNames == nil || len(e.VariableNames) == 0 {
		return errors.New("no variable names")
	}

	if e.Audiences == nil || len(e.Audiences) == 0 {
		return errors.New("no audiences")
	}

	// Validate individual audiences
	for _, audience := range e.Audiences {
		audienceValid, audienceErr := audience.Validate()

		if !audienceValid {
			return errors.Annotate(audienceErr, "invalid audience")
		}
	}

	// Validate experiment names correspond to audience variable names
	nameMap := make(map[string]bool)
	for _, audience := range e.Audiences {
		for _, valueGroup := range audience.ValueGroups {
			nameMap[valueGroup.Name] = true
		}
	}

	experimentNameMap := make(map[string]bool)
	for _, name := range e.VariableNames {
		experimentNameMap[name] = true
	}

	// Check experiment variable names
	for name, _ := range nameMap {
		if _, ok := experimentNameMap[name]; !ok {
			return errors.Errorf("experiment variable names missing %s", name)
		}
	}

	// Check audience variable names
	for name, _ := range experimentNameMap {
		if _, ok := nameMap[name]; !ok {
			return errors.Errorf("audience variable names missing %s", name)
		}
	}

	return nil
}
