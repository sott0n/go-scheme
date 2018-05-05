// Object and ObjectBase is an abstruct class for all scheme expression.
// A return value of a method which returns scheme object is Object.
// And ObjectBase has Object's implementation of String().

package scheme

// Object is an abstruct class for scheme object.
type Object interface {
	Eval() Object
	String() string
	IsNumber() bool
	IsBoolean() bool
	IsProcedure() bool
	IsNull() bool
	IsPair() bool
	IsList() bool
	IsSymbol() bool
	IsString() bool
	IsVariable() bool
	IsApplication() bool
}

// ObjectBase is an abstruct class for base scheme object.
type ObjectBase struct {
}

// Eval is object's eval IF.
func (o *ObjectBase) Eval() Object {
	runtimeError("This type's String() is not implemented yet.")
	return nil
}

func (o *ObjectBase) String() string {
	runtimeError("This object's String() is not implemented yet.")
	return ""
}

// IsNumber is an interface function of number boolean.
func (o *ObjectBase) IsNumber() bool {
	return false
}

// IsBoolean is an interface function of boolean.
func (o *ObjectBase) IsBoolean() bool {
	return false
}

// IsProcedure is an interface function of procedure boolean.
func (o *ObjectBase) IsProcedure() bool {
	return false
}

// IsNull is an interface function of procedure null boolean.
func (o *ObjectBase) IsNull() bool {
	return false
}

// IsPair is an interface function of pair boolean.
func (o *ObjectBase) IsPair() bool {
	return false
}

// IsList is an interface function of list boolean.
func (o *ObjectBase) IsList() bool {
	return false
}

// IsVariable is an interface function of variable boolean.
func (o *ObjectBase) IsVariable() bool {
	return false
}

// IsSymbol is an interface function of symbol boolean.
func (o *ObjectBase) IsSymbol() bool {
	return false
}

// IsString is an interface function of string boolean.
func (o *ObjectBase) IsString() bool {
	return false
}

// IsApplication is an interface function of application boolean.
func (o *ObjectBase) IsApplication() bool {
	return false
}
