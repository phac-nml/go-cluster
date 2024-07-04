package main

import (
	"testing"
)

// Tests for cli
func Test_cli(t *testing.T) {

	cli()
	if distance_matrix.Used {
		t.Error("Distance matrix should not have been used.")
	}

	if convert_matrix.Used {
		t.Error("Convert matrix should not have been used.")
	}

	if fast_match.Used {
		t.Error("Fast match should not have been used.")
	}

	if tree.Used {
		t.Error("tree should not have been used.")
	}
}
