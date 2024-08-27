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

	bucket_size := data_length / minimum_bins
	if data_length < bucket_size {
		return data_length, 1
	}

	if bucket_size < minimum_bins {
		bucket_size *= bucket_increase
		minimum_bins = data_length / bucket_size
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

	if (data_length-modifier) < bucket_size || bucket_size > data_length {
		// Just return the one set of indices the values are small enough
		buckets = append(buckets, Bucket{modifier, data_length})
		return buckets
	}

	for i := (bucket_size + modifier); i < data_length; i = i + bucket_size {
		new_bucket := Bucket{i - bucket_size, i}
		buckets = append(buckets, new_bucket)
	}

	final_start := buckets[len(buckets)-1].end
	final_end := data_length

	if final_end-final_start < bucket_size {
		// Extend the last index if required if it is very small
		buckets[len(buckets)-1].end = data_length
	} else {
		buckets = append(buckets, Bucket{final_start, final_end})
	}
	return buckets
}

/*
For a given data set determine the the start and end range of each of the bins to be used.
e.g. if a dataset has 1000 profiles, and our bucket size is 500 we will create bins with
an of [0, 500], [500, 1000]
*/
func BucketsIndices(data_length int, bucket_size int) []Bucket {
	var bucks []Bucket
	cpu_load_factor := CPU_LOAD_FACTOR // Need to add description to global options
	window := bucket_size

	cpu_load_string := fmt.Sprintf("CPU load factor x%d", cpu_load_factor)
	log.Println(cpu_load_string)

	if window > data_length {
		bucks = append(bucks, Bucket{0, data_length})
		log.Println("Running single threaded as there are too few entries to justify multithreading.")
		return bucks
	}

	if data_length < (runtime.NumCPU() * cpu_load_factor) {
		bucks = append(bucks, Bucket{0, data_length})
		log.Println("Running single threaded as there are too few entries to justify multithreading.")
		return bucks
	}

	for i := window; i < data_length; i = i + window {
		bucks = append(bucks, Bucket{i - window, i})
	}

	bucks = append(bucks, Bucket{bucks[len(bucks)-1].end, data_length})

	threads_running := fmt.Sprintf("Using %d threads for running.", len(bucks)-1)
	log.Println(threads_running)
	profiles_to_thread := fmt.Sprintf("Allocating ~%d profiles per a thread.", window)
	log.Println(profiles_to_thread)
	return bucks
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
	TODO redistribute data across threads at run time
	*/

	start := time.Now()
	data := *profile_data

	dist := distance_functions[DIST_FUNC].function

	bucket_index := 0
	empty_name := ""
	const cpu_modifier = 2
	minimum_buckets := runtime.NumCPU() * cpu_modifier
	bucket_size, _ := CalculateBucketSize(len(data), minimum_buckets, cpu_modifier)
	buckets := BucketsIndices(len(data), bucket_size)
	arr_pos := 1
	format_expression := GetFormatString()

	// TODO redistribute threads at run time
	var wg sync.WaitGroup
	for g := range data {
		profile_comp := data[g] // copy struct for each thread
		values_write := make([]*[]*ComparedProfile, len(buckets)-bucket_index)
		// TODO an incredible optimization here would be to go lockless, or re-use threads
		for i := bucket_index; i < len(buckets); i++ {
			array_writes := make([]*ComparedProfile, buckets[i].end-buckets[i].start)
			values_write[i-bucket_index] = &array_writes
			wg.Add(1)
			go func(output_array *[]*ComparedProfile, bucket_compute Bucket, profile_compare *Profile) {
				ThreadExecution(&data, profile_compare, bucket_compute, dist, output_array)
				wg.Done()
			}(&array_writes, buckets[i], profile_comp)
		}
		wg.Wait()                     // Wait for everyone to catch up
		buckets[bucket_index].start++ // update the current buckets tracker

		for _, i := range values_write {
			for _, value := range *i {
				fmt.Fprintf(f, format_expression, *(*value).compared, *(*value).reference, (*value).distance)
			}
		}

		if len(buckets) > 1 && arr_pos%bucket_size == 0 {
			for f := buckets[bucket_index].end - bucket_size; f < buckets[bucket_index].end; f++ {
				data[f].profile = nil
				data[f].name = empty_name

			}
			bucket_index++
			// TODO re-distribute across all cores here, no need to deplete a thread thats not using all resources
			end := time.Since(start)
			thread_depletion_time := fmt.Sprintf("One thread depleted in: %fs", end.Seconds())
			log.Println(thread_depletion_time)
			start = time.Now()
		}
		arr_pos++
	}
	wg.Wait()
	f.Flush()
}
