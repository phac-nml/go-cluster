// Test of matrix conversion
package main

import (
	"testing"
	"io/ioutil"
	"bytes"
	"path"
)


// Tests of matrix conversion
func TestPairwiseToMatrix(t *testing.T){
	tempdir := t.TempDir()
	test_input_file := "TestInputs/DistanceMatrix/Random100_molten.txt"
	test_expected_file := "TestInputs/DistanceMatrix/Random100_matrix.txt"
	test_output_file := path.Join(tempdir, "output.txt")

	t.Logf("Test Input: %s", test_input_file)
	t.Logf("Test Expected Output: %s", test_expected_file)
	t.Logf("Test Output: %s", test_output_file)

	PairwiseToMatrix(test_input_file, test_output_file)

	f1, _ := ioutil.ReadFile(test_expected_file)
	f2, _ := ioutil.ReadFile(test_output_file)
	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and expected distance matrix files do not match.")
	}

}