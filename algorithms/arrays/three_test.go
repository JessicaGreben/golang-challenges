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

func BenchmarkUrlify(b *testing.B) {
	for n := 0; n < b.N; n++ {
		urlify("Hello world, it is a wonderful day and great week.           ")
	}
}

func BenchmarkUrlify2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		urlify2("Hello world, it is a wonderful day and great week.          ")
	}
}
