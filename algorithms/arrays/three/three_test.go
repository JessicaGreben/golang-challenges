package arrays

import "testing"

var urlCases = []struct {
	name   string
	input  string
	output string
}{
	{name: "one", input: "Hello you", output: "Hello%20you"},
	{name: "two", input: "a", output: "a"},
	{name: "three", input: "", output: ""},
}

func TestUrlify(t *testing.T) {
	for _, c := range urlCases {
		t.Run(c.name, func(t *testing.T) {
			expected := urlify(c.input)
			if expected != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, expected)
			}
		})
	}
}

func TestUrlify2(t *testing.T) {
	for _, c := range urlCases {
		t.Run(c.name, func(t *testing.T) {
			expected := urlify2(c.input)
			if expected != c.output {
				t.Errorf("Expected: %v, Actual: %v\n.", c.output, expected)
			}
		})
	}
}

func getLongStr(n int) string {
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = rune(' ')
	}
	return string(runes)
}

func BenchmarkUrlify(b *testing.B) {
	str := getLongStr(1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		urlify(str)
	}
}

func BenchmarkUrlify2(b *testing.B) {
	str := getLongStr(1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		urlify2(str)
	}
}
