package arrays

import (
	"reflect"
	"testing"
)

func TestCompress(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zero(tt.input)
			for i, v := range tt.input {
				if ok := reflect.DeepEqual(v, tt.output[i]); !ok {
					t.Errorf("Expected: %v, Actual: %v\n.", tt.output, tt.input)
				}
			}
		})
	}
}

var tests = []struct {
	name   string
	input  [][]int
	output [][]int
}{
	{
		name:   "1",
		input:  [][]int{[]int{1, 1, 0, 1}, []int{1, 1, 1, 1}},
		output: [][]int{[]int{0, 0, 0, 0}, []int{1, 1, 0, 1}},
	},
}
