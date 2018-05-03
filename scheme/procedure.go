// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// when procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environmnet which is empty).

package scheme

import "log"

// Procedure is a struction for scheme procedure.
type Procedure struct {
	ObjectBase
	environment *Environment
	function    func(Object) Object
}

var builtinProcedures = Binding{
	"+":       NewProcedure(plus),
	"-":       NewProcedure(minus),
	"*":       NewProcedure(multiply),
	"/":       NewProcedure(divide),
	"number?": NewProcedure(isNumber),
}

// NewProcedure is a function for definition a new procedure.
func NewProcedure(function func(Object) Object) *Procedure {
	return &Procedure{
		environment: nil,
		function:    function,
	}
}

// Eval is Procedure's eval IF.
func (p *Procedure) Eval() Object {
	return p
}

func (p *Procedure) invoke(argument Object) Object {
	return p.function(argument)
}

//
// *** Builtin Procedures ***
//

func plus(arguments Object) Object {
	sum := 0
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil || pair.Car == nil {
			break
		}
		if car := pair.Car.Eval(); car != nil {
			number := car.(*Number)
			sum += number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(sum)
}

func minus(arguments Object) Object {
	if !arguments.IsList() {
		log.Print("Error: proper list required for function application or macro use.")
		return nil
	} else if arguments.(*Pair).ListLength() < 1 {
		log.Print("Error: procedure requires at least one argument.")
		return nil
	}

	pair := arguments.(*Pair)
	difference := pair.Car.Eval().(*Number).value
	list := pair.Cdr
	for {
		if list == nil || list.Car == nil {
			break
		}
		if car := list.Car.Eval(); car != nil {
			number := car.(*Number)
			difference -= number.value
		}
		list = list.Cdr
	}
	return NewNumber(difference)
}

func multiply(arguments Object) Object {
	product := 1
	for arguments != nil {
		pair := arguments.(*Pair)
		if pair == nil || pair.Car == nil {
			break
		}
		if car := pair.Car.Eval(); car != nil {
			number := car.(*Number)
			product *= number.value
		}
		arguments = pair.Cdr
	}
	return NewNumber(product)
}

func divide(arguments Object) Object {
	if !arguments.IsList() {
		log.Print("Error: proper list required for function application or macro use.")
		return nil
	} else if arguments.(*Pair).ListLength() < 1 {
		log.Print("Error: procedure requires at least one argument.")
		return nil
	}
	pair := arguments.(*Pair)
	quotient := pair.Car.Eval().(*Number).value
	list := pair.Cdr
	for {
		if list == nil || list.Car == nil {
			break
		}
		if car := list.Car.Eval(); car != nil {
			number := car.(*Number)
			quotient /= number.value
		}
		list = list.Cdr
	}
	return NewNumber(quotient)
}

func isNumber(object Object) Object {
	if object.IsApplication() {
		object = object.(*Application).applyProcedure()
	}
	if object.IsList() {
		list := object.(*Pair)
		if list.ListLength() == 1 {
			object = list.Car.Eval()
		} else {
			log.Printf("Wrong number of arguments: number? requires 1, but got %d", list.ListLength())
			return nil
		}
	}
	return NewBoolean(object.IsNumber())
}
