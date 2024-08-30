// Data Scheduler
//
// Data scheduler handing data slices to threads for resource handling.
// Only outputting pairwise comparisons, as a separate program will convert the file into a matrix
//
// ? While the molten data format may be annoying it allows one to potentially add resume functionality
//
// Matthew Wells:2024-02-07
package main

import (
	"bufio"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

/*
Determine how many bins of the input dataset should be processed when running the program.
The bucket size means the x Profiles will be processed by a thread, which will directly
relate to how many go routines are run at a time.
*/
func CalculateBucketSize(data_length int, minimum_bins int, bucket_increase int) (int, int) {

	if minimum_bins == 0 {
		log.Fatal("You must have a CPU modifier value greater than 0")
	}

	if data_length < minimum_bins {
		return data_length, 1
	}

	bucket_size := (data_length / minimum_bins) + bucket_increase

	if data_length < bucket_size {
		return data_length, 1
	}

	return bucket_size, minimum_bins
}

// A pair containing the start and end values for a given range of data to be processed.
type Bucket struct {
	start, end int
}

// Get the difference in indices between the two bucket fields
func (t *Bucket) Diff() int {
	return t.end - t.start
}

// The distance metric for a given comparison
type ComparedProfile struct {
	compared, reference *string
	distance            float64
}

// Calculate the initial bin sizes to use for running profiles in parallel
func CreateBucketIndices(data_length int, bucket_size int, modifier int) []Bucket {
	var buckets []Bucket

	if (data_length - modifier) <= bucket_size {
		// Just return the one set of indices the values are small enough
		buckets = append(buckets, Bucket{modifier, data_length})

	}

	for i := (bucket_size + modifier); i < data_length; i = i + bucket_size {
		new_bucket := Bucket{i - bucket_size, i}
		buckets = append(buckets, new_bucket)
	}
	buckets[len(buckets)-1].end = data_length

	return buckets
}

// Compute profile differences in a given go routine.
//
// data_slice: the data range to use for calculation against the profile to be compared too.
// profile_compare: the profile being compared in all threads
// bucket: The start and end range of the data set to write to
// dist_fn: The distance function to use for calculation of differences. Takes pointer to two profile to compare and returns a float 64
// array_writes: Array of values to append writes too
func ThreadExecution(data_slice *[]*Profile, profile_compare *Profile, bucket Bucket, dist_fn func(*[]int, *[]int) float64, array_writes *[]*ComparedProfile) {

	for i := bucket.start; i < bucket.end; i++ {
		x := dist_fn((*data_slice)[i].profile, profile_compare.profile)

		output := ComparedProfile{&profile_compare.name, &(*data_slice)[i].name, x}
		(*array_writes)[i-bucket.start] = &output
	}
}

/*
Main run loop to create a distance matrix. It create the outputs and will write
them directly to the passed in bufio.Writer.
*/
func RunData(profile_data *[]*Profile, f *bufio.Writer) {
	/* Schedule and arrange the calculation of the data in parallel
	This function is quite large and likely has room for optimization.

	Once day an incredible optimization here would be to go lockless, or re-use threads
	*/

	start := time.Now()
	data := *profile_data

	dist := distance_functions[DIST_FUNC].function
	bucket_index := 0
	cpu_modifier := BUCKET_SCALE
	data_size := len(data)
	minimum_buckets := runtime.NumCPU() * cpu_modifier
	bucket_size, _ := CalculateBucketSize(data_size, minimum_buckets, cpu_modifier)
	buckets := CreateBucketIndices(data_size, bucket_size, 0)
	format_expression := GetFormatString()
	initial_bucket_location := buckets[0].start
	var wg sync.WaitGroup

	for idx := range data {
		profile_comp := data[idx] // copy struct for each thread
		values_write := make([]*[]*ComparedProfile, len(buckets)-bucket_index)
		for b_idx, b := range buckets {
			array_writes := make([]*ComparedProfile, b.Diff())
			values_write[b_idx] = &array_writes
			wg.Add(1)
			go func(output_array *[]*ComparedProfile, bucket_compute Bucket, profile_compare *Profile) {
				ThreadExecution(&data, profile_compare, bucket_compute, dist, output_array)
				wg.Done()
			}(&array_writes, b, profile_comp)
		}

		wg.Wait()          // Wait for everyone to catch up
		buckets[0].start++ // update the current buckets tracker

		for _, i := range values_write {
			for _, value := range *i {
				fmt.Fprintf(f, format_expression, *(*value).compared, *(*value).reference, (*value).distance)
			}
		}

		resize_ratio := buckets[len(buckets)-1].Diff() >> 2
		if len(buckets) != 1 && buckets[0].Diff() < resize_ratio {

			bucket_size, minimum_buckets = CalculateBucketSize(data_size-idx, minimum_buckets, cpu_modifier)
			buckets = CreateBucketIndices(data_size, bucket_size, idx)
			for index := initial_bucket_location; index < buckets[0].start; index++ {
				data[index] = nil
			}
			initial_bucket_location = buckets[0].start
			buckets[0].start++ // start index is reserved so needs to be incremented
			end := time.Since(start)
			thread_depletion_time := fmt.Sprintf("Redistributing data across %d threads processed %d/%d profiles. %fs", len(buckets), idx, data_size, end.Seconds())
			log.Println(thread_depletion_time)
			start = time.Now()
		}
	}

	wg.Wait()
	f.Flush()
}
