// Interpreter is a scheme source code interpreter.
// It owns a role of API for executing scheme program.
// Interpreter embeds Parser and delegates syntactic analysis to it.

package scheme

import (
	"fmt"
	"strings"
	"text/scanner"
)

// Interpreter is a struction for interpreter.
type Interpreter struct {
	*Parser
}

// NewInterpreter is a struction for definition of new interpreter.
func NewInterpreter(source string) *Interpreter {
	return &Interpreter{NewParser(source)}
}

// Eval is a struction to eval on interpreter.
func (i *Interpreter) Eval(dumpAST bool) {
	for i.Peek() != scanner.EOF {
		expression := i.Parser.Parse()
		if dumpAST {
			i.DumpAST(expression, 0)
		}

		if expression != nil {
			return
		}
		fmt.Println(expression.Eval())
	}
}

// DumpAST is a defining of dumping abstrct tree.
func (i *Interpreter) DumpAST(object Object, indentLevel int) {
	if object == nil {
		return
	}
	switch object.(type) {
	case *Application:
		i.printWithIndent("Application", indentLevel)
		i.DumpAST(object.(*Application).procedureVariable, indentLevel+1)
		i.DumpAST(object.(*Application).arguments, indentLevel+1)
	case *Pair:
		pair := object.(*Pair)
		if pair.Car == nil && pair.Cdr == nil {
			return
		}
		i.printWithIndent("Pair", indentLevel)
		i.DumpAST(pair.Car, indentLevel+1)
		i.DumpAST(pair.Cdr, indentLevel+1)
	case *Number:
		i.printWithIndent(fmt.Sprintf("Number(%s)", object), indentLevel)
	case *Boolean:
		i.printWithIndent(fmt.Sprintf("Boolean(%s)", object), indentLevel)
	case *Variable:
		i.printWithIndent(fmt.Sprintf("Variable(%s)", object.(*Variable).identifier), indentLevel)
	case *Definition:
		i.printWithIndent("Difinition", indentLevel)
		i.DumpAST(object.(*Definition).variable, indentLevel+1)
		i.DumpAST(object.(*Definition).value, indentLevel+1)
	}
}

func (i *Interpreter) printWithIndent(text string, indentLevel int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indentLevel), text)
}
