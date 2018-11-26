package arrays

import "testing"

var palinCases = []struct {
	name   string
	input  string
	output bool
}{
	{name: "1", input: "abc", output: false},
	{name: "2", input: "abcabc", output: true},
	{name: "3", input: "", output: true},
	{name: "4", input: "ЖЖ", output: true},
	{name: "5", input: "a", output: true},
	{name: "6", input: "aabbccd", output: true},
	{name: "7", input: "aad", output: true},
	{name: "8", input: "aadce", output: false},
}

func TestPermPalin(t *testing.T) {
	for _, c := range palinCases {
		t.Run(c.name, func(t *testing.T) {
			b := permPalin(c.input)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}
