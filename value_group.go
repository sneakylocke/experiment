package experiment

type ValueGroup struct {
	Name           string
	Salt           string
	WeightedValues []WeightedValue
}

func NewFloatValueGroup(name string, weights []uint32, values []float64) *ValueGroup {
	valueGroup := &ValueGroup{}
	valueGroup.Name = name
	valueGroup.Salt = name
	valueGroup.WeightedValues = make([]WeightedValue, len(weights))

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Weight = weights[i]
		valueGroup.WeightedValues[i].Value = *NewFloatValue(values[i])
	}

	return valueGroup
}
