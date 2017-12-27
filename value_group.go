package experiment

type ValueGroup struct {
	Name           string
	Salt           string
	WeightedValues []WeightedValue
}

func NewFloatValueGroup(name string, weights []uint32, values []float64) *ValueGroup {
	valueGroup := newValueGroup(name, weights)

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Value = *NewFloatValue(values[i])
	}

	return valueGroup
}

func NewIntValueGroup(name string, weights []uint32, values []int64) *ValueGroup {
	valueGroup := newValueGroup(name, weights)

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Value = *NewIntValue(values[i])
	}

	return valueGroup
}

func NewBoolValueGroup(name string, weights []uint32, values []bool) *ValueGroup {
	valueGroup := newValueGroup(name, weights)

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Value = *NewBoolValue(values[i])
	}

	return valueGroup
}

func NewArbitraryValueGroup(name string, weights []uint32, values []interface{}) *ValueGroup {
	valueGroup := newValueGroup(name, weights)

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Value = *NewArbitraryValue(values[i])
	}

	return valueGroup
}

func newValueGroup(name string, weights []uint32) *ValueGroup {
	valueGroup := &ValueGroup{}
	valueGroup.Name = name
	valueGroup.Salt = name
	valueGroup.WeightedValues = make([]WeightedValue, len(weights))

	for i := range valueGroup.WeightedValues {
		valueGroup.WeightedValues[i].Weight = weights[i]
	}

	return valueGroup
}
