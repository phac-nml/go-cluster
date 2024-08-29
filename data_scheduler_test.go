package main

import (
	"bytes"
	"os"
	"path"
	"testing"
)


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
	f1, _ := os.ReadFile(test_expected_output)
	f2, _ := os.ReadFile(test_output_file)
	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}

func TestRunDataSmall(t *testing.T) {
	tempdir := t.TempDir()

	t.Log("Starting end to end test for distance calculations.")
	test_input := "TestInputs/DistanceMatrix/Random_2_input.txt"
	test_output_file := path.Join(tempdir, "output.txt")

	t.Logf("Test Input: %s", test_input)
	t.Logf("Test Output Temp File: %s", test_output_file)
	t.Log("Creating output buffer.")
	out_buffer, out_file := CreateOutputBuffer(test_output_file)

	t.Log("Loading test allele profiles.")
	test_data := LoadProfile(test_input)
	RunData(test_data, out_buffer)
	out_buffer.Flush()
	out_file.Close()

	// Compare outputs line by line
	f2, _ := os.ReadFile(test_output_file)
	output := string(f2)
	if output != "1 1 0\n" {
		t.Fatal("Input does not equal output.")
	}
}

// Testing the redistribution of bucket indices at runtime
func TestRedistributeBuckets(t *testing.T) {

	// TODO thread levels are being altered too frequently, need to change som parameters
	var profile_size int = 100
	var cpus int = 6
	BUCKET_SCALE = 3

	minimum_bucket_size := cpus * BUCKET_SCALE
	var buckets int
	buckets, minimum_bucket_size = CalculateBucketSize(profile_size, minimum_bucket_size, BUCKET_SCALE)
	bucket_indices := CreateBucketIndices(profile_size, buckets, 0)

	comparisons := make([][]int, profile_size)
	for idx := range comparisons {
		comparisons[idx] = make([]int, 0)
	}

	corrected_profile_size := profile_size
	for val := range profile_size {
		for _, b := range bucket_indices {
			for i := b.start; i < b.end; i++ {
				comparisons[val] = append(comparisons[val], i)
			}
		}

		resize_ratio := bucket_indices[len(bucket_indices)-1].Diff() >> 2
		if len(bucket_indices) != 1 && bucket_indices[0].Diff() < resize_ratio {

			buckets, minimum_bucket_size = CalculateBucketSize(profile_size-val, minimum_bucket_size, BUCKET_SCALE)
			bucket_indices = CreateBucketIndices(profile_size, buckets, val)
		}
		bucket_indices[0].start++
	}

	profile_sizes := 0
	// Check correct number of values computed
	for idx := profile_sizes; idx != 0; idx-- {
		if len(comparisons[idx]) != corrected_profile_size {
			t.Fatalf("Mismatched number of outputs for entry for index %d: %d != %d", idx, profile_sizes, len(comparisons[idx]))
		}
		profile_sizes++
	}

	// Check content
	for val := range corrected_profile_size {
		idx := val
		for _, i := range comparisons[val] {
			if i != idx {
				t.Fatalf("Index: %d, Incorrect outputs: %d != %d", val, i, idx)
			}
			idx++
		}
	}
}
