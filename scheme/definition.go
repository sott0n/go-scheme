package scheme

// Definition is a type of define function.
type Definition struct {
	ObjectBase
	environment *Environment
	variable    *Variable
	value       Object
}

func (d *Definition) String() string {
	TopLevel.Bind(d.variable.identifier, d.value)
	return d.variable.identifier
}
