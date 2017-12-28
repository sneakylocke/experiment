package experiment

import "github.com/juju/errors"

type Experiment struct {
	Name          string     // Name of the experiment
	VariableNames []string   // Names of variables allowed to be randomized in audiences
	Audiences     []Audience // Details of variables are captured in a single audience. Users may belong to one audience
	Salt          string
	Enabled       bool
}

type WeightedValue struct {
	Value  Value
	Weight uint32
}

type Constraint struct {
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
		return errors.New("no audiences names")
	}

	for _, audience := range e.Audiences {
		audienceValid, audienceErr := audience.validate()

		if !audienceValid {
			return errors.Annotate(audienceErr, "invalid audience")
		}
	}

	return nil
}
