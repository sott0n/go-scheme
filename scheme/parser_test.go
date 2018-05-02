package scheme

import (
	"testing"
	"text/scanner"
)

type parserTest struct {
	source string
	result string
}

var parserTests = []parserTest{
	{"()", "()"},
	{"10", "10"},
	{"123456789", "123456789"},
	{"(+)", "0"},
	{"(- 1)", "1"},
	{"(*)", "1"},
	{"(+ 1 2)", "3"},
	{"(+ 1 20 300 4000)", "4321"},
	{"( + 1 2 3 )", "6"},
	{"(+ 1 (+ 2 3) (+ 3 4))", "13"},
	{"(- 3 (- 2 3) (+ 3 0))", "1"},
	{"(* (* 3 3) 3)", "27"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		p := NewParser(test.source)
		if p.Peek() != scanner.EOF {
			parseObject := p.Parse()
			if parseObject == nil {
				t.Errorf("%s => <nil>; want %s", test.source, test.result)
				return
			}
			actual := parseObject.String()
			if actual != test.result {
				t.Errorf("%s => %s; want %s", test.source, actual, test.result)
			}
		}
	}
}
