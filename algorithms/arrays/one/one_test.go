package arrays

import "testing"

var cases = []struct {
	name   string
	input  string
	output bool
}{
	{name: "1", input: "abc", output: true},
	{name: "2", input: "abcabc", output: false},
	{name: "3", input: "", output: true},
	{name: "4", input: "ЖЖ", output: false},
	{name: "5", input: "♘", output: true},

	// case 6 is an example of a character that is more than one code point.
	// {name: "6", input: " ̀àa", output: true}, // "\u0300\u0061\u0300\u0061"
}

func TestSolutionOneA(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := solutionOneA(c.input)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}

func TestSolutionOneB(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := solutionOneB(c.input)
			if b != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, b)
			}
		})
	}
}

func getLongStr(n int) string {
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = rune(i)
	}
	return string(runes)
}

var benchCases = []struct {
	name   string
	input  string
	output bool
}{
	{name: "1", input: getLongStr(1e2)},
	{name: "2", input: getLongStr(1e3)},
	{name: "3", input: getLongStr(1e4)},
	{name: "4", input: getLongStr(1e5)},
}

func BenchmarkSolutionOneA(b *testing.B) {
	for _, c := range benchCases {
		b.Run(c.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				solutionOneA(c.input)
			}
		})
	}
}

func BenchmarkSolutionOneB(b *testing.B) {
	for _, c := range benchCases {
		b.Run(c.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				solutionOneB(c.input)
			}
		})
	}
}
