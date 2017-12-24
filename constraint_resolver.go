package experiment

type ConstraintResolver interface {
	PassesConstraints(userInfo interface{}, constraints []Constraint) bool
}

type BasicResolver struct{}

func (resolver *BasicResolver) PassesConstraints(userInfo interface{}, constraints []Constraint) bool {
	return true
}
