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

func benchmarkOneA(str string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		solutionOneA(str)
	}
}

func BenchmarkOneA1(b *testing.B)  { benchmarkOneA("a", b) }
func BenchmarkOneA2(b *testing.B)  { benchmarkOneA("ab", b) }
func BenchmarkOneA3(b *testing.B)  { benchmarkOneA("abcd", b) }
func BenchmarkOneA10(b *testing.B) { benchmarkOneA("abcdefghij", b) }
func BenchmarkOneA20(b *testing.B) { benchmarkOneA("abcdefghijklmnopqrst", b) }
func BenchmarkOneA40(b *testing.B) { benchmarkOneA("abcdefghijklmnopqrstuvwxyz0123456789àáâä", b) }

func benchmarkOneB(str string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		solutionOneB(str)
	}
}

func BenchmarkOneB1(b *testing.B)  { benchmarkOneB("a", b) }
func BenchmarkOneB2(b *testing.B)  { benchmarkOneB("ab", b) }
func BenchmarkOneB3(b *testing.B)  { benchmarkOneB("abcd", b) }
func BenchmarkOneB10(b *testing.B) { benchmarkOneB("abcdefghij", b) }
func BenchmarkOneB20(b *testing.B) { benchmarkOneB("abcdefghijklmnopqrst", b) }
func BenchmarkOneB40(b *testing.B) { benchmarkOneB("abcdefghijklmnopqrstuvwxyz0123456789àáâä", b) }
