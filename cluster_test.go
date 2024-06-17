package main


import (
	"testing"
	"io/ioutil"
	"bytes"
	"path"
)


func TestCluster(t *testing.T) {
	linkage_methods := 0
	test_input_matrix := "TestInputs/DistanceMatrix/Random100_matrix.txt"
	expected_output_tree := "TestInputs/DistanceMatrix/Random100_tree.nwk"

	tmpdir := t.TempDir()
	output_tree := path.Join(tmpdir, "output.nwk")
	output_buffer, out_file := CreateOutputBuffer(output_tree)

	t.Logf("Test Input: %s", test_input_matrix)
	t.Logf("Test Expected Output: %s", expected_output_tree)
	t.Logf("Test Output: %s", output_tree)
	
	Cluster(test_input_matrix, linkage_methods, output_buffer)
	output_buffer.Flush()
	out_file.Close()

	f1, _ := ioutil.ReadFile(output_tree)
	f2, _ := ioutil.ReadFile(expected_output_tree)
	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}