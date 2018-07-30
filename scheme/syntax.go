// This file is for statements by syntax form, such as set!

package scheme

import "fmt"

var (
	builtinSyntaxes = Binding{
		"set!":   NewSyntax(setSyntax),
		"if":     NewSyntax(ifSyntax),
		"lambda": NewSyntax(lambdaSyntax),
		"let":    NewSyntax(letSyntax),
		"let*":   NewSyntax(letSyntax),
		"letrec": NewSyntax(letSyntax),
		"and":    NewSyntax(andSyntax),
		"or":     NewSyntax(orSyntax),
		"quote":  NewSyntax(quoteSyntax),
		"begin":  NewSyntax(beginSyntax),
		"define": NewSyntax(defineSyntax),
		"cond":   NewSyntax(condSyntax),
		"do":     NewSyntax(doSyntax),
	}
)

// Syntax is a type for updating value.
type Syntax struct {
	ObjectBase
	function func(*Syntax, Object) Object
}

// NewSyntax is definition of creating new set pair.
func NewSyntax(function func(*Syntax, Object) Object) *Syntax {
	return &Syntax{ObjectBase: ObjectBase{parent: nil}, function: function}
}

// Invoke is for evaluating set object.
func (s *Syntax) Invoke(arguments Object) Object {
	return s.function(s, arguments)
}

func (s *Syntax) String() string {
	return fmt.Sprintf("#<syntax %s>", s.Bounder())
}

func (s *Syntax) isSyntax() bool {
	return true
}

func (s *Syntax) malformedError() {
	syntaxError("malformed %s: %s", s.Bounder(), s.Bounder().Parent())
}

func (s *Syntax) assertListEqual(arguments Object, length int) {
	if !arguments.isList() || arguments.(*Pair).ListLength() != length {
		s.malformedError()
	}
}

func (s *Syntax) assertListMinimum(arguments Object, minimum int) {
	if !arguments.isList() || arguments.(*Pair).ListLength() < minimum {
		s.malformedError()
	}
}

func (s *Syntax) assertListRange(arguments Object, lengthRange []int) {
	if !arguments.isList() {
		s.malformedError()
	}

	for _, length := range lengthRange {
		if length == arguments.(*Pair).ListLength() {
			return
		}
	}
	s.malformedError()
}

func setSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 2)
	elements := arguments.(*Pair).Elements()

	variable := elements[0]
	if !variable.isVariable() {
		s.malformedError()
	}
	value := elements[1].Eval()
	s.Bounder().set(variable.(*Variable).identifier, value)
	return value
}

func ifSyntax(s *Syntax, arguments Object) Object {
	s.assertListRange(arguments, []int{2, 3})
	elements := arguments.(*Pair).Elements()

	result := elements[0].Eval()
	if result.isBoolean() && !result.(*Boolean).value {
		if len(elements) == 3 {
			return elements[2].Eval()
		} else {
			return undef
		}
	} else {
		return elements[1].Eval()
	}
}

// Return elements in list object.
// With type assertion (syntax form specific error message)
// and auto cast to list from application.
func (s *Syntax) elementsMinimum(list Object, minimum int) []Object {
	if list.isApplication() {
		list = list.(*Application).toList()
	}
	s.assertListMinimum(list, minimum)
	return list.(*Pair).Elements()
}

func applicationToList(object Object) Object {
	if object.isApplication() {
		return object.(*Application).toList()
	} else {
		return object
	}
}

func andSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := Object(NewBoolean(true))
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
		if lastResult.isBoolean() && lastResult.(*Boolean).value == false {
			return NewBoolean(false)
		}
	}
	return lastResult
}

func orSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := Object(NewBoolean(false))
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
		if !lastResult.isBoolean() || lastResult.(*Boolean).value != false {
			return lastResult
		}
	}
	return lastResult
}

func beginSyntax(s *Syntax, arguments Object) Object {
	s.assertListMinimum(arguments, 0)

	lastResult := undef
	for _, object := range arguments.(*Pair).Elements() {
		lastResult = object.Eval()
	}
	return lastResult
}

func defineSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 2)
	elements := arguments.(*Pair).Elements()

	if !elements[0].isVariable() {
		syntaxError("%s", s.Bounder().Parent())
	}
	variable := elements[0].(*Variable)
	s.Bounder().define(variable.identifier, elements[1].Eval())

	return NewSymbol(variable.identifier)
}

func quoteSyntax(s *Syntax, arguments Object) Object {
	s.assertListEqual(arguments, 1)
	object := arguments.(*Pair).ElementAt(0)

	p := NewParser(object.String())
	p.Peek()
	return p.parseQuotedObject(s.Bounder())
}

func condSyntax(s *Syntax, arguments Object) Object {
	elements := s.elementsMinimum(arguments, 0)
	if len(elements) == 0 {
		syntaxError("at least one clause is required for cond")
	}

	// First: syntax check
	elseExists := false
	for _, element := range elements {
		if elseExists {
			syntaxError("'else' clause followed by more clauses")
		} else if element.isApplication() && element.(*Application).procedure.isVariable() &&
			element.(*Application).procedure.(*Variable).identifier == "else" {
			elseExists = true
		}

		if element.isNull() || !element.isApplication() {
			syntaxError("bad clause in cond")
		}
	}

	// Second: eval cases
	for _, element := range elements {
		lastResult := undef
		application := element.(*Application)

		isElse := application.procedure.isVariable() && application.procedure.(*Variable).identifier == "else"
		if !isElse {
			lastResult = application.procedure.Eval()
		}

		// first element is 'else' or not '#f'
		if isElse || !lastResult.isBoolean() || lastResult.(*Boolean).value == true {
			for _, object := range application.arguments.(*Pair).Elements() {
				lastResult = object.Eval()
			}
			return lastResult
		}
	}
	return undef
}

func lambdaSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())

	elements := s.elementsMinimum(arguments, 1)
	variables := s.elementsMinimum(elements[0], 0)

	// generate function
	closure.function = func(givenArguments Object) Object {
		// assert given arguments
		givenElements := s.elementsMinimum(givenArguments, 0)
		if len(variables) != len(givenElements) {
			compileError("wrong number of arguments: requires %d, but got %d", len(variables), len(givenElements))
		}

		// define arguments to local scope
		for index, variable := range variables {
			object := givenElements[index].Eval()
			if variable.isVariable() {
				closure.localBinding[variable.(*Variable).identifier] = object
			}
		}

		// returns last eval result
		lastResult := undef
		for _, element := range elements[1:] {
			lastResult = element.Eval()
		}
		return lastResult
	}
	return closure
}

func doSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())

	// Parse iterator list and define first variable
	elements := s.elementsMinimum(arguments, 2)
	iteratorBodies := s.elementsMinimum(elements[0], 0)

	for _, iteratorBody := range iteratorBodies {
		iteratorElements := s.elementsMinimum(iteratorBody, 2)
		if len(iteratorElements) > 3 {
			compileError("bad update expr in %s: %s", s.Bounder(), s.Bounder().Parent())
		}

		variable := iteratorElements[0]
		value := iteratorElements[1]
		if variable.isVariable() {
			closure.localBinding[variable.(*Variable).identifier] = value.Eval()
		}
	}

	// eval test ->
	//   true: eval testBody and returns its result
	//  false: eval continueBody, eval iterator's update
	testElements := s.elementsMinimum(elements[1], 1)
	continueElements := elements[2:]

	for {
		testResult := testElements[0].Eval()
		if !testResult.isBoolean() || testResult.(*Boolean).value == true {
			for _, element := range testElements[1:] {
				testResult = element.Eval()
			}
			return testResult
		} else {
			// eval continueBody
			for _, element := range continueElements {
				element.Eval()
			}

			// update iterators
			for _, iteratorBody := range iteratorBodies {
				iteratorElements := s.elementsMinimum(iteratorBody, 2)

				if len(iteratorElements) == 3 {
					variable := iteratorElements[0]
					if variable.isVariable() {
						closure.localBinding[variable.(*Variable).identifier] = iteratorElements[2].Eval()
					}
				}
			}
		}
	}
	return undef
}

func letSyntax(s *Syntax, arguments Object) Object {
	closure := WrapClosure(arguments.Parent())

	elements := s.elementsMinimum(arguments, 1)
	argumentElements := s.elementsMinimum(elements[0], 0)

	// define arguments to local scope
	for _, argumentElement := range argumentElements {
		variableElements := s.elementsMinimum(argumentElement, 2)
		variable := variableElements[0]

		if variable.isVariable() {
			closure.localBinding[variable.(*Variable).identifier] = variableElements[1].Eval()
		}
	}

	// eval body
	lastResult := undef
	for _, element := range elements[1:] {
		lastResult = element.Eval()
	}
	return lastResult
}
