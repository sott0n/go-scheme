// Parser is a type to analyze scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import "fmt"

// Parser is a struction for analyze scheme source's syntax.
type Parser struct {
	*Lexer
}

// NewParser is a function for definition a new Parser.
func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

// Parse is a function to deal parser.
func (p Parser) Parse(environment *Environment) Object {
	p.ensureAvailability()
	return p.parseObject(environment)
}

func (p *Parser) parseObject(environment *Environment) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		peekToken := p.PeekToken()
		if p.TokenType() == ')' {
			p.NextToken()
			return new(Pair)
		} else if peekToken == "define" {
			p.NextToken()
			return p.parseDefinition(environment)
		} else if peekToken == "quote" {
			p.NextToken()
			object := p.parseQuotedList(environment)
			if !object.IsList() || object.(*Pair).ListLength() != 1 {
				compileError("syntax-error: malformed quote.")
			}
			return object.(*Pair).Car
		} else if peekToken == "lambda" {
			p.NextToken()
			return p.parseProcedure(environment)
		}
		return p.parseApplication(environment)

	case '\'':
		return p.parseQuotedObject(environment)
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewVariable(token, environment)
	case BooleanToken:
		return NewBoolean(token)
	case StringToken:
		return NewString(token[1 : len(token)-1])
	default:
		return nil
	}
}

// This function returns *Pair of first object and list from second.
// Returns value is Object because if a method returns nil which is not
// interface type, the method's result cannot be judged as nil.
func (p *Parser) parseList(environment *Environment) Object {
	car := p.parseObject(environment)
	if car == nil {
		return new(Pair)
	}
	cdr := p.parseList(environment).(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}

func (p *Parser) parseApplication(environment *Environment) Object {
	firstObject := p.parseObject(environment)
	if firstObject == nil {
		runtimeError("Unexpected flow: procedure application car is nil.")
	}
	list := p.parseList(environment)
	if list == nil {
		runtimeError("Unexpected flow: procedure application cdr is nil.")
	}
	return &Application{
		procedureVariable: firstObject,
		arguments:         list,
		environment:       environment,
	}
}

func (p *Parser) parseProcedure(environment *Environment) Object {
	if p.TokenType() == '(' {
		p.NextToken()
		return NewProcedure(
			environment,
			p.parseList(environment),
			p.parseList(environment),
		)
	}
	runtimeError("Not implemented yet.")
	return nil
}

func (p *Parser) parseDefinition(environment *Environment) Object {
	object := p.parseList(environment)
	if !object.IsList() || object.(*Pair).ListLength() != 2 {
		runtimeError("Compile Error: syntax-error: (define).")
	}
	list := object.(*Pair)
	variable := list.ElementAt(0).(*Variable)
	value := list.ElementAt(1)

	return &Definition{
		environment: environment,
		variable:    variable,
		value:       value,
	}
}

func (p *Parser) parseQuotedObject(environment *Environment) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(environment)
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewSymbol(token)
	case BooleanToken:
		return NewBoolean(token)
	default:
		return nil
	}
}

func (p *Parser) parseQuotedList(environment *Environment) Object {
	car := p.parseQuotedObject(environment)
	if car == nil {
		return new(Pair)
	}
	cdr := p.parseQuotedList(environment).(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}

func (p *Parser) ensureAvailability() {
	// Error message will be printed by interpreter.
	recover()
}

func compileError(format string, a ...interface{}) {
	runtimeError("Compile Error: "+format, a...)
}

func runtimeError(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
