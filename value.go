package experiment

type Value struct {
	FloatValue float64 `json:"float"`
	IntValue   int64   `json:"int"`
	BoolValue  bool    `json:"bool"`
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
