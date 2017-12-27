package builders

import (
	//"github.com/juju/errors"
	//log "github.com/sirupsen/logrus"
	e "github.com/sneakylocke/experiment"
)

type FactorialBasicBuilder struct {
}

func NewFactorialBasicBuilder(experimentName string) BasicBuilder {
	b := &FactorialBasicBuilder{}

	return b
}

func (b *FactorialBasicBuilder) AddFloat(variableName string, weights []uint32, values []float64) error {
	return nil
}

func (b *FactorialBasicBuilder) AddInt(variableName string, weights []uint32, values []int64) error {
	return nil
}

func (b *FactorialBasicBuilder) AddBool(variableName string, weights []uint32, values []bool) error {
	return nil
}

func (b *FactorialBasicBuilder) Build() (*e.Experiment, error) {
	return nil, nil
}
