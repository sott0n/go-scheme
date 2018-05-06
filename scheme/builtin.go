// This file defines built-in procedures for TopLevel environment.

package scheme

import (
	"fmt"
	"strings"
)

// BuiltinProcedures is lists of procedure binded.
func BuiltinProcedures() Binding {
	return Binding{
		"+":              builtinProcedure(plus),
		"-":              builtinProcedure(minus),
		"*":              builtinProcedure(multiply),
		"/":              builtinProcedure(divide),
		"=":              builtinProcedure(equal),
		"<":              builtinProcedure(lessThan),
		"<=":             builtinProcedure(lessEqual),
		">":              builtinProcedure(greaterThan),
		">=":             builtinProcedure(greaterEqual),
		"number?":        builtinProcedure(isNumber),
		"null?":          builtinProcedure(isNull),
		"procedure?":     builtinProcedure(isProcedure),
		"boolean?":       builtinProcedure(isBoolean),
		"pair?":          builtinProcedure(isPair),
		"list?":          builtinProcedure(isList),
		"symbol?":        builtinProcedure(isSymbol),
		"string?":        builtinProcedure(isString),
		"not":            builtinProcedure(not),
		"car":            builtinProcedure(car),
		"cdr":            builtinProcedure(cdr),
		"list":           builtinProcedure(list),
		"string-append":  builtinProcedure(stringAppend),
		"symbol->string": builtinProcedure(symbolToString),
		"string->symbol": builtinProcedure(stringToSymbol),
		"eq?":            builtinProcedure(isEq),
		"neq?":           builtinProcedure(isNeq),
		"equal?":         builtinProcedure(isEqual),
	}
}

func builtinProcedure(function func(Object) Object) *Procedure {
	return &Procedure{function: function}
}

func assertListMinimum(arguments Object, minimum int) {
	if !arguments.isList() {
		compileError("proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() < minimum {
		compileError("procedure requires at least %d argument", minimum)
	}
}

func assertListEqual(arguments Object, length int) {
	if !arguments.isList() {
		compileError("proper list required for function application or macro use")
	} else if arguments.(*Pair).ListLength() != length {
		compileError("Wrong number of arguments: number? requires %d, but got %d",
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
	case *Pair:
		if object.isNull() {
			return "null"
		}
		return "pair"
	default:
		rawTypeName := fmt.Sprintf("%T", object)
		typeName := strings.Replace(rawTypeName, "*scheme.", "", 1)
		return strings.ToLower(typeName)
	}
}

func assertObjectType(object Object, assertType string) {
	if assertType != typeName(object) {
		compileError("%s required, but got %s", assertType, object)
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

func compareNumber(arguments Object, compareFunc func(int, int) bool) Object {
	assertListMinimum(arguments, 2)

	numbers := evaledObjects(arguments.(*Pair).Elements())
	assertObjectsType(numbers, "number")

	oldValue := numbers[0].(*Number).value
	for _, number := range numbers[1:] {
		if !compareFunc(oldValue, number.(*Number).value) {
			return NewBoolean(false)
		}
		oldValue = number.(*Number).value
	}
	return NewBoolean(true)
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
	return compareNumber(arguments, func(a, b int) bool { return a == b })
}

func lessThan(arguments Object) Object {
	return compareNumber(arguments, func(a, b int) bool { return a < b })
}

func lessEqual(arguments Object) Object {
	return compareNumber(arguments, func(a, b int) bool { return a <= b })
}

func greaterThan(arguments Object) Object {
	return compareNumber(arguments, func(a, b int) bool { return a > b })
}

func greaterEqual(arguments Object) Object {
	return compareNumber(arguments, func(a, b int) bool { return a >= b })
}

func isNumber(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isNumber()
	})
}

func isNull(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isNull()
	})
}

func isProcedure(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isProcedure()
	})
}

func isBoolean(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isBoolean()
	})
}

func isPair(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isPair()
	})
}

func isList(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isList()
	})
}

func isSymbol(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isSymbol()
	})
}

func isString(arguments Object) Object {
	return booleanByFunc(arguments, func(object Object) bool {
		return object.isString()
	})
}

func not(arguments Object) Object {
	return booleanByFunc(arguments,
		func(object Object) bool {
			return object.isBoolean() && !object.(*Boolean).value
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

func areIdentical(a Object, b Object) bool {
	if typeName(a) != typeName(b) {
		return false
	}

	switch a.(type) {
	case *Number:
		return a.(*Number).value == b.(*Number).value
	case *Symbol:
		return a.(*Symbol).identifier == b.(*Symbol).identifier
	case *Boolean:
		return a.(*Boolean).value == b.(*Boolean).value
	default:
		return false
	}
}

func areEqual(a Object, b Object) bool {
	if a == nil {
		return true
	}
	if typeName(a) != typeName(b) {
		return false
	} else if areIdentical(a, b) {
		return true
	}

	switch a.(type) {
	case *Pair:
		return areEqual(a.(*Pair).Car, b.(*Pair).Car) && areEqual(a.(*Pair).Cdr, b.(*Pair).Cdr)
	default:
		return false
	}
}

func areSameList(a Object, b Object) bool {
	if typeName(a) != typeName(b) {
		return false
	}

	switch a.(type) {
	case *Pair:
		return areSameList(a.(*Pair).Car, b.(*Pair).Car) && areSameList(a.(*Pair).Cdr, b.(*Pair).Cdr)
	default:
		return areIdentical(a, b)
	}
}

func isEq(arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areIdentical(objects[0], objects[1]))
}

func isNeq(arguments Object) Object {
	return NewBoolean(!isEq(arguments).(*Boolean).value)
}

func isEqual(arguments Object) Object {
	assertListEqual(arguments, 2)

	objects := evaledObjects(arguments.(*Pair).Elements())
	return NewBoolean(areEqual(objects[0], objects[1]))
}
