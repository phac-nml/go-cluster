package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
)

type bucket_tests struct {
	data_length, bucket_size, cpu_modifier, expected int
}

var bucket_size_tests = []bucket_tests{
	bucket_tests{10, 1, 1, 10},
	bucket_tests{10, 2, 1, 5},
}

//func TestCalculateBucketSize(t *testing.T) {
//	for _, test := range bucket_size_tests {
//		if output := CalculateBucketSize(test.data_length, test.bucket_size, test.cpu_modifier); output != test.expected {
//			t.Errorf("Output %d not equal to expected %d", output, test.expected)
//			t.Errorf("Output %+v", output)
//		}
//	}
//}

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
	//f1, _ := ioutil.ReadFile(test_expected_output)
	f1, _ := os.ReadFile(test_expected_output)
	f2, _ := os.ReadFile(test_output_file)

	if !bytes.Equal(f1, f2) {
		t.Fatal("Input and output files to not match.")
	}
}

// / Example calculation of the BucketIndices function
// / Input values are integers of data_length, and then the bucket_size returnting a Bucket struct
//func TestBucketIndices(t *testing.T) {
//	CPU_LOAD_FACTOR = 1
//	buckets := BucketsIndices(100, 50)
//	fmt.Printf("%v\n", buckets)
//
//	buckets = BucketsIndices(10, 1)
//	fmt.Printf("%v\n", buckets)
//}

// / Test redistribution of data across threads
// / This test simulates the what occurs during the RunData function
//func TestBucketResizing(t *testing.T) {
//	CPU_LOAD_FACTOR = 1
//	data_range := 100
//	number_of_cpus := 10
//	bucket_size := CalculateBucketSize(data_range, number_of_cpus, CPU_LOAD_FACTOR)
//	buckets := BucketsIndices(data_range, bucket_size)
//	fmt.Println(buckets)
//	expected_buckets := []Bucket{{0, 10}, {10, 20}, {20, 30}, {30, 40}, {40, 50}, {50, 60}, {60, 70}, {70, 80}, {80, 90}, {90, 100}}
//	for i := range expected_buckets {
//		if expected_buckets[i].start != buckets[i].start || expected_buckets[i].end != buckets[i].end {
//			t.Fatal("Mismatched value in expected bucket outputs.")
//		}
//	}
//
//	// Create a test array of values
//	var slice = make([]int, data_range)
//	for i := range slice {
//		slice[i] = i
//	}
//
//	bucket_index := 0
//	arr_pos := 1
//
//	for range slice {
//		// inner loop using the bucket of values here
//		// Need to deplete buckets...
//
//		fmt.Println(buckets[bucket_index:])
//		buckets[bucket_index].start++
//		if len(buckets) > 1 && arr_pos%bucket_size == 0 {
//			fmt.Println(bucket_index)
//			bucket_index++
//		}
//		arr_pos++
//	}
//
//}

// Testing an alternate method for generating compute indices
//func TestBucketsGeneration(t *testing.T) {
//	var cpus int = 6
//	var profile_sizes int = 100
//	CPU_LOAD_FACTOR = 1
//	minimum_bucket_size := 10
//	buckets := CalculateBucketSize(profile_sizes, cpus, CPU_LOAD_FACTOR)
//	bucket_indices := CreateBucketIndices(profile_sizes, buckets, 0)
//
//	fmt.Println(buckets)
//	fmt.Println(bucket_indices)
//	bucket_index := 0
//
//	for val := range profile_sizes {
//		bucket_indices[bucket_index].start++
//		if len(bucket_indices) > 1 && bucket_indices[bucket_index].Diff() < minimum_bucket_size {
//			buckets = CalculateBucketSize(profile_sizes-val, cpus, CPU_LOAD_FACTOR)
//			fmt.Println(buckets)
//			bucket_indices = CreateBucketIndices(profile_sizes-val, buckets, val)
//			fmt.Println(bucket_indices)
//		}
//		fmt.Println(bucket_indices)
//	}
//}

// Testing the redistribution of bucket indices at runtime
func TestRedistributeBuckets(t *testing.T) {
	var profile_size int = 10000
	var cpus int = 6
	CPU_LOAD_FACTOR = 2
	minimum_bucket_size := cpus * CPU_LOAD_FACTOR
	var buckets int
	buckets, minimum_bucket_size = CalculateBucketSize(profile_size-1, minimum_bucket_size, CPU_LOAD_FACTOR)
	bucket_indices := CreateBucketIndices(profile_size-1, buckets, 0)

	comparisons := make([][]int, profile_size)
	for idx := range comparisons {
		comparisons[idx] = make([]int, 0)
	}

	resize_counter := 0
	for val := range profile_size {
		for _, b := range bucket_indices {
			for i := b.start; i < b.end; i++ {
				comparisons[val] = append(comparisons[val], i)
			}
		}

		if bucket_indices[0].Diff() < minimum_bucket_size {
			buckets, minimum_bucket_size = CalculateBucketSize(profile_size-val, minimum_bucket_size, CPU_LOAD_FACTOR)
			bucket_indices = CreateBucketIndices(profile_size-1, buckets, val)
			fmt.Println(val, len(bucket_indices))
			resize_counter++
		}
		bucket_indices[0].start++
	}
	fmt.Println("Time Resizing bins", resize_counter)

	corrected_profile_size := profile_size - 1
	for idx := range profile_size {
		if len(comparisons[idx]) != corrected_profile_size {
			t.Fatalf("Mismatched number of outputs for entry for index %d: %d != %d", idx, corrected_profile_size, len(comparisons[idx]))
		}
		corrected_profile_size--
	}

}
