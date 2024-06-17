/*
	Need to add tests for profile ingestion, as a trie, a map or HAMT could be best
	...
*/

package main

import "testing"

var profile_name string = "test"
var profile []int = []int{1, 2, 3}

func Test_NewProfile(t *testing.T) {
	if output := NewProfile(profile_name, &profile); output.name != profile_name || output.profile != &profile {
		t.Errorf("Output %+v not equal to expected.", *output)
	}
}
