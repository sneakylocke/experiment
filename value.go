package experiment

type Value struct {
	FloatValue     float64
	IntValue       int64
	BoolValue      bool
	ArbitraryValue interface{}
}

func NewFloatValue(v float64) *Value {
	value := &Value{}
	value.FloatValue = v
	return value
}

func NewIntValue(v int64) *Value {
	value := &Value{}
	value.IntValue = v
	return value
}

func NewBoolValue(v bool) *Value {
	value := &Value{}
	value.BoolValue = v
	return value
}

func NewArbitraryValue(v interface{}) *Value {
	value := &Value{}
	value.ArbitraryValue = v
	return value
}
