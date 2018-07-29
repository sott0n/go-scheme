// Parser is a type to analyze scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

// Parser is a struction for analyze scheme source's syntax.
type Parser struct {
	*Lexer
}

// NewParser is a function for definition a new Parser.
func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

// Parse is a function to deal parser.
func (p Parser) Parse(parent Object) Object {
	p.ensureAvailability()
	return p.parseObject(parent)
}

func (p *Parser) parseObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseBlock(parent)
	case '\'':
		return p.parseSingleQuote(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewVariable(token, parent)
	case BooleanToken:
		return NewBoolean(token, parent)
	case StringToken:
		return NewString(token[1:len(token)-1], parent)
	default:
		return nil
	}
}

func (p *Parser) parseBlock(parent Object) Object {
	switch p.PeekToken() {
	case ")":
		p.NextToken()
		return Null
	case "let", "let*", "letrec":
		p.NextToken()
		return p.parseLet(parent)
	}
	return p.parseApplication(parent)
}

// This is for parsing syntax sugar '*** => (quote ***)
func (p *Parser) parseSingleQuote(parent Object) Object {
	if len(p.PeekToken()) == 0 {
		runtimeError("unterminated quote")
	}
	applicaton := NewApplication(parent)
	applicaton.procedure = NewVariable("quote", applicaton)
	applicaton.arguments = NewList(applicaton, p.parseObject(applicaton))
	return applicaton
}

// This function returns *Pair of first object and list from second.
// Scanner position ends with the next of close parentheses.
func (p *Parser) parseList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseList(pair).(*Pair)
	return pair
}

func (p *Parser) parseApplication(parent Object) Object {
	application := NewApplication(parent)
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseLet(parent Object) Object {
	if p.TokenType() == '(' {
		p.NextToken()
	} else {
		syntaxError("malformed let")
	}

	application := NewApplication(parent)
	procedure := new(Procedure)

	procedureArguments := NewPair(procedure)
	applicationArguments := NewPair(application)

	argumentSets := p.parseList(application)
	for _, set := range argumentSets.(*Pair).Elements() {
		if !set.isApplication() || set.(*Application).arguments.(*Pair).ListLength() != 1 {
			syntaxError("malformed let")
		}

		procedureArguments.Append(set.(*Application).procedure)
		applicationArguments.Append(set.(*Application).arguments.(*Pair).ElementAt(0))
	}

	procedure.arguments = procedureArguments
	procedure.body = p.parseList(application)
	procedure.generateFunction(parent)

	application.arguments = applicationArguments
	application.procedure = procedure
	return application
}

func (p *Parser) parseQuotedObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(parent)
	case '\'':
		return p.parseSingleQuote(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewSymbol(token)
	case BooleanToken:
		return NewBoolean(token, parent)
	case ')':
		return nil
	default:
		runtimeError("unterminated quote")
		return nil
	}
}

func (p *Parser) parseQuotedList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseQuotedObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseQuotedList(pair).(*Pair)
	return pair
}

func (p *Parser) ensureAvailability() {
	// Error message will be printed by interpreter.
	recover()
}
