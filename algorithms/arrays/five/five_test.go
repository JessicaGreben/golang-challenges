package arrays

import "testing"

var oneCases = []struct {
	name   string
	str1   string
	str2   string
	output bool
}{
	{name: "1", str1: "abcd", str2: "abcd", output: true},
	{name: "2", str1: "a", str2: "abcd", output: false},
}

func TestOneAway(t *testing.T) {
	for _, c := range oneCases {
		t.Run(c.name, func(t *testing.T) {
			b := oneAway(c.str1, c.str2)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}

var insertCases = []struct {
	name   string
	str1   string
	str2   string
	output bool
}{
	{name: "1", str1: "pale", str2: "ple", output: true},
	{name: "2", str1: "pales", str2: "pale", output: true},
	{name: "3", str1: "pale", str2: "bale", output: true},
	{name: "4", str1: "pale", str2: "bake", output: false},
	{name: "5", str1: "pale", str2: "pakes", output: false},
}

func TestCheck(t *testing.T) {
	for _, c := range insertCases {
		t.Run(c.name, func(t *testing.T) {
			b := check(c.str1, c.str2)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}
