/*
TODO golang allows benchmarks to be added
*/
package main

import (
	"fmt"
	"testing"
)

type addTest struct {
	profile1                   []int
	profile2                   []int
	expected, expected_missing float64
}

var hamming_dist_tests = []addTest{
	addTest{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, float64(0), float64(0)},
	addTest{[]int{2, 1, 1, 2}, []int{2, 2, 2, 2}, float64(2), float64(2)},
	addTest{[]int{2, 0, 0, 2}, []int{2, 0, 0, 2}, float64(0), float64(0)},
	addTest{[]int{1, 1, 1}, []int{2, 2, 2}, float64(3), float64(3)},
	addTest{[]int{1}, []int{2}, float64(1), float64(1)},
}

var scaled_dist_tests = []addTest{
	addTest{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, float64(100.0 * ((4.0 - 4.0) / 4.0)), float64(100.0 * ((4.0 - 4.0) / 4.0))},
	addTest{[]int{2, 1, 1, 2}, []int{2, 2, 2, 2}, float64(100.0 * ((4.0 - 2.0) / 4.0)), float64(100.0 * ((4.0 - 2.0) / 4.0))},
	addTest{[]int{2, 0, 0, 2}, []int{2, 0, 0, 2}, float64(100.0 * ((2.0 - 2.0) / 2.0)), float64(100.0 * ((4.0 - 4.0) / 4.0))},
	addTest{[]int{1, 1, 1}, []int{2, 2, 2}, float64(100.0 * ((3.0 - 0.0) / 3.0)), float64(100.0 * ((3.0 - 0.0) / 3.0))},
	addTest{[]int{1}, []int{2}, float64(100.0), float64(100.0)},
}

func Test_hamming_distance(t *testing.T) {
	for _, test := range hamming_dist_tests {
		if output := hamming_distance(&test.profile1, &test.profile2); output != test.expected {
			t.Errorf("Output %f not equal to expected %f", output, test.expected)
		}
	}
}

func Test_hamming_distance_missing(t *testing.T) {
	for _, test := range hamming_dist_tests {
		if output := hamming_distance(&test.profile1, &test.profile2); output != test.expected_missing {
			t.Errorf("Output %f not equal to expected %f", output, test.expected_missing)
		}
	}
}

func Test_scaled_distance(t *testing.T) {
	for _, test := range scaled_dist_tests {
		if output := scaled_distance(&test.profile1, &test.profile2); output != test.expected {
			t.Errorf("Output %f not equal to expected %f", output, test.expected)
		}
	}
}

func Test_scaled_distance_missing(t *testing.T) {
	for _, test := range scaled_dist_tests {
		if output := scaled_distance(&test.profile1, &test.profile2); output != test.expected_missing {
			t.Errorf("Output %f not equal to expected %f", output, test.expected_missing)
		}
	}
}

func Example_hamming_distance() {
	fmt.Println(hamming_distance(&[]int{1, 2, 3, 4}, &[]int{1, 2, 3, 4}))
	// Output: 0
}
