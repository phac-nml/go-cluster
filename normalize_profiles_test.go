/*
	Tests for
*/

package main

import "testing"

type value_insertion struct {
	value    string
	expected int
}

var value_insertion_tests = []value_insertion{
	value_insertion{"test", 1},
	value_insertion{"test", 1},
	value_insertion{"test2", 2},
}

func Test_InsertValue_NewProfile(t *testing.T) {
	pLookup := NewProfileLookup()
	for _, test := range value_insertion_tests {
		if output := pLookup.InsertValue(&test.value); output != test.expected {
			t.Errorf("Output: %d does not match expected %d", output, test.expected)
		}
	}
}
