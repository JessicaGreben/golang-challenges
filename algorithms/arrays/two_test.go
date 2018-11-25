package arrays

import "testing"

var permCases = []struct {
	name    string
	string1 string
	string2 string
	output  bool
}{
	{name: "one", string1: "abc", string2: "cba", output: true},
	{name: "two", string1: "a", string2: "", output: false},
	{name: "three", string1: "", string2: "", output: true},
	{name: "four", string1: "Ж♘Ж", string2: "ЖЖ♘", output: true},
}

func TestIsPermutation(t *testing.T) {
	for _, c := range permCases {
		t.Run(c.name, func(t *testing.T) {
			b := isPermutation(c.string1, c.string2)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}
