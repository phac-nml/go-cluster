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
func TestBucketIndices(t *testing.T){
	CPU_LOAD_FACTOR = 1
	buckets := BucketsIndices(100, 50);
	fmt.Printf("%v\n", buckets)

	buckets = BucketsIndices(10, 1);
	fmt.Printf("%v\n", buckets)
}

/// Test redistribution of data across threads
/// This test simulates the what occurs during the RunData function
func TestBucketResizing(t *testing.T){
	CPU_LOAD_FACTOR = 1
	data_range := 100
	number_of_cpus := 10
	bucket_size := CalculateBucketSize(data_range, number_of_cpus, CPU_LOAD_FACTOR)
	buckets := BucketsIndices(data_range, bucket_size)
	fmt.Println(buckets)
	expected_buckets := []Bucket{Bucket{0, 10}, Bucket{10, 20}, Bucket{20, 30}, Bucket{30, 40}, Bucket{40, 50}, Bucket{50, 60}, Bucket{60, 70}, Bucket{70, 80}, Bucket{80, 90}, Bucket{90, 100}}
	for i := range expected_buckets {
		if expected_buckets[i].start != buckets[i].start || expected_buckets[i].end != buckets[i].end {
			t.Fatal("Mismatched value in expected bucket outputs.")
		}
	}
	
	// Create a test array of values
	var slice = make([]int, data_range)
	for i := range slice {
		slice[i] = i
	}

	bucket_index := 0
	arr_pos := 1

	for  range slice {
		// inner loop using the bucket of values here
		// Need to deplete buckets...

		
		fmt.Println(buckets[bucket_index:])
		buckets[bucket_index].start++
		if len(buckets) > 1 && arr_pos%bucket_size == 0 {
			fmt.Println(bucket_index)
			bucket_index++
		}
		arr_pos++
	}


}
