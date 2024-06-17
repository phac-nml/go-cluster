package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"
)

func TestIdentifyMatches(t *testing.T) {
	test_input_profile := "TestInputs/DistanceMatrix/Random100_input.txt"
	expected_output := "TestInputs/DistanceMatrix/Random100xRandom100.fast-match.txt"

	tmpdir := t.TempDir()
	output_fm := path.Join(tmpdir, "output.txt")
	output_buffer, out_file := CreateOutputBuffer(output_fm)

	t.Logf("Test Input: %s", test_input_profile)
	t.Logf("Test Expected Output: %s", expected_output)
	t.Logf("Test Output: %s", output_fm)

	IdentifyMatches(test_input_profile, test_input_profile, 1, output_buffer)
	output_buffer.Flush()
	out_file.Close()

	f1, _ := ioutil.ReadFile(output_fm)
	f2, _ := ioutil.ReadFile(expected_output)
	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}
