package nock

import (
	"testing"
)

func TestNockCommon(t *testing.T) {
	tests := map[string]string{
		`[0 1 2]`:      `2`,
		`[57 [0 1]]`:   `57`,
		`[42 [4 0 1]]`: `43`,
		`[57 [4 0 1]]`: `58`,

		`[[132 19] [0 3]]`:   `19`,
		`[42 [1 153 218]]`:   `[153 218]`,
		`[[132 19] [4 0 3]]`: `20`,

		`[42 [8 [4 0 1] [0 1]]]`:   `[43 42]`,
		`[42 [[4 0 1] [3 0 1]]]`:   `[43 1]`,
		`[42 [8 [4 0 1] [4 0 3]]]`: `43`,
		`[42 [7 [4 0 1] [4 0 1]]]`: `44`,

		`[[[4 5] [6 14 15]] [0 7]]`:      `[14 15]`,
		`[77 [2 [1 42] [1 1 153 218]]]`:  `[153 218]`,
		`[42 [6 [1 0] [4 0 1] [1 233]]]`: `43`,
		`[42 [6 [1 1] [4 0 1] [1 233]]]`: `233`,

		// Decrement.
		`[42 [8 [1 0] 8 [1 6 [5 [0 7] 4 0 6] [0 6] 9 2 [0 2] [4 0 6] 0 7] 9 2 0 1]]`: `41`,
	}
	for in, exp := range tests {
		out := Nock5(Parse(in)).String()
		if out != exp {
			t.Errorf("For Nock 5 expected %q, got %q\n", exp, out)
		}
		out = Nock4(Parse(in)).String()
		if out != exp {
			t.Errorf("For Nock 4 expected %q, got %q\n", exp, out)
		}
	}
}

func TestNock5(t *testing.T) {
	tests := map[string]string{
		`[[132 19] [10 37 [4 0 3]]]`: `20`,
	}
	for in, exp := range tests {
		out := Nock5(Parse(in)).String()
		if out != exp {
			t.Errorf("Expected %q, got %q\n", exp, out)
		}
	}
}

func TestNock4(t *testing.T) {
	tests := map[string]string{
		`[[132 19] [11 37 [4 0 3]]]`:       `20`,
		`[42 [10 [6 [0 1]] [1 20 60 70]]]`: `[20 [42 70]]`,
	}
	for in, exp := range tests {
		out := Nock4(Parse(in)).String()
		if out != exp {
			t.Errorf("Expected %q, got %q\n", exp, out)
		}
	}
}
