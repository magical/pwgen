package main

import (
	"testing"
)

func TestConstantTimeMinMax(t *testing.T) {
	const (
		maxInt = 1<<31 - 1
	)
	var values = []int{
		0, 1, 99, 100,

		maxInt - 1,
		maxInt,

		maxInt/2 + 1,
		maxInt/2 - 1,
	}
	for _, a := range values {
		for _, b := range values {
			min, max := constantTimeMinMax(a, b)
			if min != a && min != b {
				t.Errorf("min is not one of the input values. a=%d, b=%d, min=%d, max=%d", a, b, min, max)
			}
			if max != a && max != b {
				t.Errorf("max is not one of the input values. a=%d, b=%d, min=%d, max=%d", a, b, min, max)
			}
			if a != b && min == max {
				t.Errorf("lost an input value. a=%d, b=%d, min=%d, max=%d", a, b, min, max)
			}
			if min > max {
				t.Errorf("want min <= max. a=%d, b=%d, min=%d, max=%d", a, b, min, max)
			}
		}
	}
}

func TestIndexWithoutReplacement(t *testing.T) {
	tests := []struct {
		seen, test []int
	}{
		{[]int{0, 1, 2, 10, 12}, []int{3, 4, 5, 6, 7, 8, 9, 11, 13, 14}},
		{[]int{}, []int{0, 1, 2, 3, 4}},

		// invalid
		//{[]int{2, 2}, []int{0, 1, 4, 5}},
		//{[]int{2, 1}, []int{0, 2, 4, 5}},
	}
	for ti, tt := range tests {
		for idx, want := range tt.test {
			got := indexWithoutReplacement(idx, tt.seen)
			if got != want {
				t.Errorf("test#%d, idx=%d: got %d want %d", ti, idx, got, want)
			}
		}
	}
}
