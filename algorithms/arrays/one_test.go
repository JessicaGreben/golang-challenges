package arrays

import "testing"

var cases = []struct {
	name   string
	input  string
	output bool
}{
	{name: "one", input: "abc", output: true},
	{name: "two", input: "abcabc", output: false},
	{name: "three", input: "", output: true},
	{name: "four", input: "ЖЖ", output: false},
	{name: "four", input: "♘", output: true},
}

func TestSolutionOneA(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, _ := solutionOneA(c.input)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}

func TestSolutionOneB(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, _ := solutionOneB(c.input)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}
