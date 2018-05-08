// Pair is a type which is generated by cons procedure.
// Pair has two pointers, which are named car and cdr.
//
// List is expressed by linked of Pair.
// And procedure application has list which consists of Pair.
// as its arguments.

package scheme

import (
	"fmt"
	"strings"
)

// Pair is a struction for car and cdr.
type Pair struct {
	ObjectBase
	Car Object
	Cdr Object
}

// NewNull is creating clear pair.
func NewNull(parent Object) *Pair {
	return &Pair{ObjectBase: ObjectBase{parent: parent}, Car: nil, Cdr: nil}
}

// Eval is Pair's eval IF.
func (p *Pair) Eval() Object {
	return p
}

// String is a string function with accessing Pair.
func (p *Pair) String() string {
	if p.isNull() {
		return "()"
	} else if p.isList() {
		length := p.ListLength()
		tokens := []string{}
		for i := 0; i < length; i++ {
			tokens = append(tokens, p.ElementAt(i).Eval().String())
		}
		return fmt.Sprintf("(%s)", strings.Join(tokens, " "))
	} else {
		return fmt.Sprintf("(%s . %s)", p.Car, p.Cdr)
	}
}

func (p *Pair) isNull() bool {
	return p.Car == nil && p.Cdr == nil
}

func (p *Pair) isPair() bool {
	return !p.isNull()
}

func (p *Pair) isList() bool {
	pair := p
	for {
		if pair.isNull() {
			return true
		}
		switch pair.Cdr.(type) {
		case *Pair:
			pair = pair.Cdr.(*Pair)
		default:
			return false
		}
	}
}

// Elements is returns each elements.
func (p *Pair) Elements() []Object {
	elements := []Object{}
	pair := p
	for {
		if pair.Car == nil {
			break
		} else {
			elements = append(elements, pair.Car)
		}
		pair = pair.Cdr.(*Pair)
	}
	return elements
}

// ElementAt is return value with specified index.
func (p *Pair) ElementAt(index int) Object {
	return p.Elements()[index]
}

// ListLength returns length of list.
func (p *Pair) ListLength() int {
	if p.isNull() {
		return 0
	}
	return p.Cdr.(*Pair).ListLength() + 1
}

// Append returns pair appended object.
func (p *Pair) Append(object Object) *Pair {
	assertListMinimum(p, 0)

	listTail := p
	for {
		if listTail.isNull() {
			break
		} else {
			listTail = listTail.Cdr.(*Pair)
		}
	}
	listTail.Car = object
	listTail.Cdr = new(Pair)
	return p
}

// AppendList returns pair appended list.
func (p *Pair) AppendList(list Object) *Pair {
	assertListMinimum(p, 0)
	assertListMinimum(list, 0)

	listTail := p
	for {
		if listTail.isNull() {
			break
		} else {
			listTail = listTail.Cdr.(*Pair)
		}
	}
	listTail.Car = list.(*Pair).Car
	listTail.Cdr = list.(*Pair).Cdr
	return p
}
