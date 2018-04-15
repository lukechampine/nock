// Package nock implements a simple Nock interpreter, according the spec
// available at https://urbit.org/docs/nock/definition
//
// This interpreter assumes that its input is well-formed, and does not
// support atoms larger than a machine int.
package nock

import (
	"strconv"
	"strings"
)

// A Noun is an atom or a cell. An atom is any natural number. A cell is any
// ordered pair of nouns.
type Noun struct {
	atom *int
	cell *[2]Noun
}

// IsAtom returns true if n is an atom.
func (n Noun) IsAtom() bool { return n.atom != nil }

// IsCell returns true if n is a cell.
func (n Noun) IsCell() bool { return n.cell != nil }

// Num returns the integer value of n, which must be an atom.
func (n Noun) Num() int { return *n.atom }

// Head returns the head of n, which must be a cell.
func (n Noun) Head() Noun { return n.cell[0] }

// Tail returns the tail of n, which must be a cell.
func (n Noun) Tail() Noun { return n.cell[1] }

// String implements the fmt.Stringer interface.
func (n Noun) String() string {
	if n.IsAtom() {
		return strconv.Itoa(n.Num())
	}
	return "[" + n.Head().String() + " " + n.Tail().String() + "]"
}

// Atom returns an atom with value i.
func Atom(i int) Noun { return Noun{atom: &i} }

// Cell returns a cell that pairs head with tail.
func Cell(head, tail Noun) Noun { return Noun{cell: &[2]Noun{head, tail}} }

// Loobean returns the atom 0 if b is true, and the atom 1 if b is false.
func Loobean(b bool) Noun { return Atom(map[bool]int{true: 0, false: 1}[b]) }

func wut(n Noun) Noun { return Loobean(n.IsCell()) }
func lus(n Noun) Noun { return Atom(1 + n.Num()) }
func tis(n Noun) Noun { return Loobean(n.Head().String() == n.Tail().String()) }

func fas(i int, n Noun) Noun {
	switch i {
	case 1:
		return n
	case 2:
		return n.Head()
	case 3:
		return n.Tail()
	default:
		return fas(2+i%2, fas(i/2, n))
	}
}

func tar(sub, form Noun) Noun {
	if form.Head().IsCell() {
		return Cell(tar(sub, form.Head()), tar(sub, form.Tail()))
	}
	inst, arg := form.Head(), form.Tail()
	switch inst.Num() {
	case 0:
		// *[a 0 b]         /[b a]
		return fas(arg.Num(), sub)
	case 1:
		// *[a 1 b]         b
		return arg
	case 2:
		// *[a 2 b c]       *[*[a b] *[a c]]
		return tar(tar(sub, arg.Head()), tar(sub, arg.Tail()))
	case 3:
		// *[a 3 b]         ?*[a b]
		return wut(tar(sub, arg))
	case 4:
		// *[a 4 b]         +*[a b]
		return lus(tar(sub, arg))
	case 5:
		// *[a 5 b]         =*[a b]
		return tis(tar(sub, arg))
	case 6:
		// *[a 6 b c d]     *[a 2 [0 1] 2 [1 c d] [1 0] 2 [1 2 3] [1 0] 4 4 b]
		if tar(sub, arg.Head()).Num() == 0 {
			return tar(sub, fas(6, arg))
		}
		return tar(sub, fas(7, arg))
	case 7:
		// *[a 7 b c]       *[a 2 b 1 c]
		return tar(tar(sub, arg.Head()), arg.Tail())
	case 8:
		// *[a 8 b c]       *[a 7 [[7 [0 1] b] 0 1] c]
		return tar(Cell(tar(sub, arg.Head()), sub), arg.Tail())
	case 9:
		// *[a 9 b c]        *[a 7 c [2 [0 1] [0 b]]]
		d := tar(sub, arg.Tail())
		return tar(d, fas(arg.Head().Num(), d))
	case 10:
		// *[a 10 b c]      *[a c]
		if b := arg.Head(); b.IsCell() {
			_ = tar(sub, b.Tail())
		}
		return tar(sub, arg.Tail())
	}
	panic("unreachable")
}

// Nock evaluates the nock function on n.
func Nock(n Noun) Noun {
	if n.IsAtom() || n.Tail().IsAtom() {
		return n
	}
	return tar(n.Head(), n.Tail())
}

// Parse parses a Nock program.
func Parse(s string) Noun {
	s = strings.Replace(s, "[", " [ ", -1)
	s = strings.Replace(s, "]", " ] ", -1)
	n, _ := parseNoun(strings.Fields(strings.TrimSpace(s)))
	return n
}

func parseNoun(s []string) (Noun, []string) {
	if s[0] == "[" {
		return parseCell(s)
	}
	return parseAtom(s)
}

func parseCell(s []string) (Noun, []string) {
	s = s[1:]
	var elems []Noun
	for s[0] != "]" {
		var e Noun
		e, s = parseNoun(s)
		elems = append(elems, e)
	}
	for len(elems) > 1 {
		elems = append(elems[:len(elems)-2], Cell(elems[len(elems)-2], elems[len(elems)-1]))
	}
	return elems[0], s[1:]
}

func parseAtom(s []string) (Noun, []string) {
	i, _ := strconv.Atoi(s[0])
	return Atom(i), s[1:]
}
