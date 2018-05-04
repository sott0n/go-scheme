// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"strings"
)

var builtinProcedures = Binding{
	"+":              BuiltinProcedure(plus),
	"-":              BuiltinProcedure(minus),
	"*":              BuiltinProcedure(multiply),
	"/":              BuiltinProcedure(divide),
	"=":              BuiltinProcedure(equal),
	"number?":        BuiltinProcedure(isNumber),
	"null?":          BuiltinProcedure(isNull),
	"procedure?":     BuiltinProcedure(isProcedure),
	"boolean?":       BuiltinProcedure(isBoolean),
	"pair?":          BuiltinProcedure(isPair),
	"list?":          BuiltinProcedure(isList),
	"symbol?":        BuiltinProcedure(isSymbol),
	"string?":        BuiltinProcedure(isString),
	"not":            BuiltinProcedure(not),
	"car":            BuiltinProcedure(car),
	"cdr":            BuiltinProcedure(cdr),
	"list":           BuiltinProcedure(list),
	"string-append":  BuiltinProcedure(stringAppend),
	"symbol->string": BuiltinProcedure(symbolToString),
	"string->symbol": BuiltinProcedure(stringToSymbol),
}

// BuiltinProcedure is definition of procedure.
func BuiltinProcedure(function func(Object) Object) *Procedure {
	return &Procedure{
		environment: nil,
		function:    function,
		arguments:   nil,
		body:        nil,
	}
}

func assertListMinimum(arguments Object, minimum int) {
	if !arguments.IsList() {
		compileError("proper list required for function application or macro use.")
	} else if arguments.(*Pair).ListLength() < minimum {
		compileError("procedure requires at least %d arguments.", minimum)
	}
}

func assertListEqual(arguments Object, length int) {
	if !arguments.IsList() {
		compileError("proper list required for function application or macro use.")
	} else if arguments.(*Pair).ListLength() != length {
		compileError("Wrong number of arguments: number? requires %d, but got %d.",
			length, arguments.(*Pair).ListLength())
	}
}

func assertObjectsType(objects []Object, typeName string) {
	for _, object := range objects {
		assertObjectType(object, typeName)
	}
}

func typeName(object Object) string {
	switch object.(type) {
	case *Number:
		return "number"
	case *String:
		return "string"
	case *Symbol:
		return "symbol"
	case *Procedure:
		return "procedure"
	case *Boolean:
		return "boolean"
	case *Pair:
		if object.IsNull() {
			return "null"
		}
		return "pair"
	default:
		return "Not Implemented typeName."
	}
}

func assertObjectType(object Object, assertType string) {
	if assertType != typeName(object) {
		compileError("%s required, but got %s.", assertType, object)
	}
}

func evaledObjects(objects []Object) []Object {
	evaledObjects := []Object{}

	for _, object := range objects {
		evaledObjects = append(evaledObjects, object.Eval())
	}
	return evaledObjects
}

func booleanByFunc(arguments Object, typeCheckFunc func(Object) bool) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	return NewBoolean(typeCheckFunc(object))
}

func plus(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	sum := 0
	for _, number := range numbers {
		sum += number.(*Number).value
	}
	return NewNumber(sum)
}

func minus(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	difference := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		difference -= number.(*Number).value
	}
	return NewNumber(difference)
}

func multiply(arguments Object) Object {
	assertListMinimum(arguments, 0)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	product := 1
	for _, number := range numbers {
		product *= number.(*Number).value
	}
	return NewNumber(product)
}

func divide(arguments Object) Object {
	assertListMinimum(arguments, 1)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	quotient := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		quotient /= number.(*Number).value
	}
	return NewNumber(quotient)
}

func equal(arguments Object) Object {
	assertListMinimum(arguments, 2)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	firstValue := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		if firstValue != number.(*Number).value {
			return NewBoolean(false)
		}
	}
	return NewBoolean(true)
}

func isNumber(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsNumber()
	})
}

func isNull(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsNull()
	})
}

func isProcedure(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsProcedure()
	})
}

func isBoolean(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsBoolean()
	})
}

func isPair(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsPair()
	})
}

func isList(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsList()
	})
}

func isSymbol(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsSymbol()
	})
}

func isString(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.IsString()
	})
}

func not(arguments Object) Object {
	return booleanByFunc(arguments,
		func(object Object) bool {
			return object.IsBoolean() && !object.(*Boolean).value
		})
}

func car(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Car
}

func cdr(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "pair")
	return object.(*Pair).Cdr
}

func list(arguments Object) Object {
	return arguments
}

func stringAppend(arguments Object) Object {
	assertListMinimum(arguments, 0)

	stringObjects := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(stringObjects, "string")

	texts := []string{}
	for _, stringObject := range stringObjects {
		texts = append(texts, stringObject.(*String).text)
	}
	return NewString(strings.Join(texts, ""))
}

func symbolToString(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "symbol")
	return NewString(object.(*Symbol).identifier)
}

func stringToSymbol(arguments Object) Object {
	assertListEqual(arguments, 1)

	object := arguments.(*Pair).ElementAt(0).Eval()
	assertObjectType(object, "string")
	return NewSymbol(object.(*String).text)
}
