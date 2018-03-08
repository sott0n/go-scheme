// Object and ObjectBase is an abstruct class for all scheme expression.
// A return value of a method which returns scheme object is Object.
// And ObjectBase has Object's implementation of String().

package scheme

type Object interface {
	String() string
	IsProcedure() bool
}

type ObjectBase struct {
}

func (o *ObjectBase) String() string {
	return "This type's String() is not implemented yet."
}

func (o *ObjectBase) IsProcedure() bool {
	return false
}
