// Parser is a type to analyze scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import (
	"log"
)

// Parser is a struction for analyze scheme source's syntax.
type Parser struct {
	*Lexer
}

// NewParser is a function for definition a new Parser.
func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

// Parse is a function to deal parser.
func (p Parser) Parse() Object {
	return p.parseObject()
}

func (p *Parser) parseObject() Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		if p.TokenType() == ')' {
			p.NextToken()
			return new(Pair)
		}
		firstObject := p.parseObject()
		if firstObject == nil {
			log.Print("Unexpected flow: procedure application car is nil.")
			return nil
		}
		list := p.parseList()
		if list == nil {
			log.Print("Unexpected flow: procedure application cdr is nil.")
			return nil
		}
		return &Application{
			procedureVariable: firstObject,
			arguments:         list,
		}
	case ')':
		return nil
	case EOF:
		return nil
	case IntToken:
		return NewNumber(token)
	case IdentifierToken:
		return NewVariable(token)
	default:
		return nil
	}
}

// This function returns *Pair of first object and list from second.
// Returns value is Object because if a method returns nil which is not
// interface type, the method's result cannot be judged as nil.
func (p *Parser) parseList() Object {
	car := p.Parse()
	if car == nil {
		return new(Pair)
	}
	cdr := p.parseList().(*Pair)
	return &Pair{Car: car, Cdr: cdr}
}
