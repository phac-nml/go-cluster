package main

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"
	"fmt"
)

type bucket_tests struct {
	data_length, bucket_size, cpu_modifier, expected int
}

var bucket_size_tests = []bucket_tests{
	bucket_tests{10, 1, 1, 10},
	bucket_tests{10, 2, 1, 5},
}

func TestCalculateBucketSize(t *testing.T) {
	for _, test := range bucket_size_tests {
		if output := CalculateBucketSize(test.data_length, test.bucket_size, test.cpu_modifier); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
			t.Errorf("Output %+v", output)
		}
	}
}

// Test that molten file output is the same.
func TestRunData(t *testing.T) {
	tempdir := t.TempDir()

	t.Log("Starting end to end test for distance calculations.")
	test_input := "TestInputs/DistanceMatrix/Random100_input.txt"
	test_expected_output := "TestInputs/DistanceMatrix/Random100_molten.txt"
	test_output_file := path.Join(tempdir, "output.txt")

	t.Logf("Test Input: %s", test_input)
	t.Logf("Test Expected Output: %s", test_expected_output)
	t.Logf("Test Output Temp File: %s", test_output_file)
	t.Log("Creating output buffer.")
	out_buffer, out_file := CreateOutputBuffer(test_output_file)

	t.Log("Loading test allele profiles.")
	test_data := LoadProfile(test_input)
	RunData(test_data, out_buffer)
	out_buffer.Flush()
	out_file.Close()

	// Compare outputs line by line
	f1, _ := ioutil.ReadFile(test_expected_output)
	f2, _ := ioutil.ReadFile(test_output_file)

	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}



/// Example calculation of the BucketIndices function
/// Input values are integers of data_length, and then the bucket_size returnting a Bucket struct
func TestBucketIndicesF(t *testing.T){
	CPU_LOAD_FACTOR = 1
	buckets := BucketsIndices(100, 50);
	fmt.Printf("%v\n", buckets)

	buckets = BucketsIndices(10, 1);
	fmt.Printf("%v\n", buckets)
}
