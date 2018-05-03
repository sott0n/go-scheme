package scheme

import (
	"testing"
)

type parserTest struct {
	source  string
	results []string
}

var parserTests = []parserTest{
	makePT("12", "12"),
	makePT("()", "()"),
	makePT("#f", "#f"),
	makePT("#t", "#t"),
	makePT("1234567890", "1234567890"),

	makePT("(+)", "0"),
	makePT("(- 1)", "1"),
	makePT("(*)", "1"),
	makePT("(/ 1)", "1"),

	makePT("(+ 1 2)", "3"),
	makePT("(+ 1 20 300 4000)", "4321"),
	makePT("( + 1 2 3 )", "6"),
	makePT("(+ 1 (+ 2 3) (+ 3 4))", "13"),
	makePT("(- 3 (- 2 3) (+ 3 0))", "1"),
	makePT("(* (* 3 3) 3)", "27"),
	makePT("(/ 100 (/ 4 2))", "50"),
	makePT("(+ (* 100 3) (/ (- 4 2) 2))", "301"),

	makePT("(number? 100)", "#t"),
	makePT("(number? (+ 3 (* 2 8)))", "#t"),
	makePT("(number? #t)", "#f"),
	makePT("(number? ()", "#f"),

	makePT("(define x 1) x", "x", "1"),
	makePT("(define x (+ 1 3)) x", "x", "4"),

	makePT("'12", "12"),
	makePT("'hello", "hello"),
	makePT("'#f", "#f"),
	makePT("'#t", "#t"),
	makePT("'(  1  2  3  )", "(1 2 3)"),
	makePT("'( 1 ( 2 3 ) )", "(1 (2 3))"),
}

func makePT(source string, results ...string) parserTest {
	return parserTest{source: source, results: results}
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		p := NewParser(test.source)
		p.Peek()
		for i := 0; i < len(test.results); i++ {
			result := test.results[i]
			parseObject := p.Parse()
			if parseObject == nil {
				t.Errorf("%s => <nil>; want %s", test.source, result)
				return
			}
			actual := parseObject.String()
			if actual != result {
				t.Errorf("%s => %s; want %s", test.source, actual, result)
			}
		}
	}
}
