package main

import (
	"bytes"
	"os"
	"path"
	"testing"
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

	f1, _ := os.ReadFile(output_tree)
	f2, _ := os.ReadFile(expected_output_tree)
	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}

type LinkageMethodTest struct {
	linkage_method int
	expected       string
}

var LinkageMethodTests = []LinkageMethodTest{
	{0, "average"},
	{1, "centroid"},
	{2, "complete"},
	{3, "mcquitty"},
	{4, "median"},
	{5, "single"},
	{6, "ward"},
}

func TestGetLinkageMethod(t *testing.T) {
	for _, test := range LinkageMethodTests {
		if output := GetLinkageMethod(test.linkage_method); output != test.expected {
			t.Errorf("Output not equal to expected %s %s", output, test.expected)
			t.Errorf("Output %+v", output)
		}
	}
}
