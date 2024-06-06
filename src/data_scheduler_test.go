
package main

import "testing"


type bucket_tests struct {
	data_length, bucket_size, expected int
}

var bucket_size_tests = []bucket_tests {
	bucket_tests{10, 1, 5},
}

func Test_calculate_bucket_size(t *testing.T){
	for _, test := range bucket_size_tests {
		if output := calculate_bucket_size(test.data_length, test.bucket_size); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
		}
	}
}