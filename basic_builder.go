package experiment

const (
	audienceName = "default_audience"
)

type BasicBuilder interface {
	Build() (*Experiment, error)
	AddFloats(variableName string, weights []uint32, values []float64) error
	AddInts(variableName string, weights []uint32, values []int64) error
	AddBools(variableName string, weights []uint32, values []bool) error
}

type basicBuilder struct {
	AdvancedBuilder AdvancedBuilder
}

func NewSimpleBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.AdvancedBuilder = NewAdvancedBuilder(experimentName)

	return b
}

func NewFactorialBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.AdvancedBuilder = NewAdvancedFactorialBuilder(experimentName)

	return b
}

func NewAlignedBuilder(experimentName string) BasicBuilder {
	b := &basicBuilder{}
	b.AdvancedBuilder = NewAdvancedAlignedBuilder(experimentName)

	return b
}

func (b *basicBuilder) AddFloats(variableName string, weights []uint32, values []float64) error {
	return b.AdvancedBuilder.AddFloats(variableName, audienceName, weights, values)
}

func (b *basicBuilder) AddInts(variableName string, weights []uint32, values []int64) error {
	return b.AdvancedBuilder.AddInts(variableName, audienceName, weights, values)
}

func (b *basicBuilder) AddBools(variableName string, weights []uint32, values []bool) error {
	return b.AdvancedBuilder.AddBools(variableName, audienceName, weights, values)
}

func (b *basicBuilder) Build() (*Experiment, error) {
	return b.AdvancedBuilder.Build()
}
