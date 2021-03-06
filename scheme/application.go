// Application is a type to express scheme procedure application.
// Application has a procedure and its argument as list which consists.
// of Pair.

package scheme

// Application is a struction for application.
type Application struct {
	ObjectBase
	procedure Object
	arguments Object
}

type Invoker interface {
	Invoke(Object) Object
}

// NewApplication is creating new object for application.
func NewApplication(parent Object) *Application {
	return &Application{
		ObjectBase: ObjectBase{parent: parent},
	}
}

// Eval is eval function IF that returns applyProcedure.
func (a *Application) Eval() Object {
	return a.applyProcedure()
}

func (a *Application) String() string {
	// Exceptional handling for special form: quote
	if a.procedure.isVariable() {
		variable := a.procedure.(*Variable)
		if variable.boundedObject(variable.identifier) == builtinSyntaxes["quote"] {
			if a.arguments.isNull() {
				return "(quote)"
			} else {
				return "'" + a.arguments.(*Pair).ElementAt(0).String()
			}
		}
	}
	pair := NewPair(nil)
	pair.Car = a.procedure
	pair.Cdr = a.arguments
	return pair.String()
}

func (a *Application) toList() *Pair {
	list := NewPair(a.Parent())
	list.Car = a.procedure
	list.Car.setParent(list)
	list.Cdr = a.arguments
	list.Cdr.setParent(list)
	return list
}

func (a *Application) applyProcedure() Object {
	evaledObject := a.procedure.Eval()

	switch evaledObject.(type) {
	case Invoker:
		return evaledObject.(Invoker).Invoke(a.arguments)
	default:
		runtimeError("invalid application")
		return nil
	}
}

func (a *Application) isApplication() bool {
	return true
}
