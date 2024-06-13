/* Data scheduler handing data slices to threads for resource handling.

Only outputting pairwise comparisons, as a seperate program will convert the file into a matrix

Matthew Wells:2024-02-07
*/

package main

import (
	"log"
	"time"
	"runtime"
	"sync"
	"fmt"
	"bufio"
)

type OutputValue struct {
	profile_1 string
	profile_2 string
	distance float64
}

func calculate_bucket_size(data_length int, runtime_cpus int) int {
	//bucket_size := data_length / (runtime_cpus * 2)
	bucket_size := data_length / runtime_cpus
	return bucket_size
}

type Bucket struct {
	start, end int
}

func buckets_indices(data_length int, bucket_size int) []Bucket {
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

	bucks = append(bucks, Bucket{bucks[len(bucks) - 1].end, data_length})
	
	threads_running := fmt.Sprintf("Using %d threads for running.", len(bucks) - 1)
	log.Println(threads_running)
	profiles_to_thread := fmt.Sprintf("Allocating ~%d profiles per a thread.", window)
	log.Println(profiles_to_thread)
	return bucks
}





func thread_execution(data_slice *[]*Profile, profile_compare *Profile, bucket Bucket, dist_fn func(*[]int, *[]int) float64, array_writes *[]*string) {
	/* Compute profile differences.
	
	data_slice: the data range to use for calculation against the profile to be compared too.
	profile_compare: the profile being compared in all threads
	start_idx: The starting range in the profile to be used to initilize comparisons
	end_idx: The end range to calculate comparisons up too.
	dist_fn: The distance function to use for calculation of differences. Takes pointer to two profile to compare and returns a float 64
	array_writes: Array of values to append writes too
	*/


	format_expression := get_format_string()

	// TODO need to pass a slice properly in the future
	for i := bucket.start; i < bucket.end; i++ {
		x := dist_fn((*data_slice)[i].profile, profile_compare.profile);

		output := fmt.Sprintf(format_expression, profile_compare.name, (*data_slice)[i].name, x);
		(*array_writes)[i-bucket.start] = &output;
	}
}


func run_data(profile_data *[]*Profile, f *bufio.Writer) {
	/* Schedule and arrange the calculation of the data in parallel
	TODO redistribute data across threads at run time
	TODO writing to stdout will be the initial method outputting calculated results, but this will likely change in the future
	*/

	start := time.Now()
	data := *profile_data
	
	dist := distance_functions[DIST_FUNC].function // ! Default value is stored in the dists.go file

	bucket_index := 0
	empty_name := ""
	bucket_size := calculate_bucket_size(len(data), runtime.NumCPU())
	buckets := buckets_indices(len(data), bucket_size)
	arr_pos := 1

	// TODO can create a pool of go routines and pass the profile to compare to each channel
	var wg sync.WaitGroup
	for g := range data[0:] {
		profile_comp := data[g] // copy struct for each thread
		values_write := make([]*[]*string, len(buckets) - bucket_index)
		// TODO an incredible optimization here would be to go lockless, or re-use threads
		for i := bucket_index; i < len(buckets); i++ {
			
			array_writes := make([]*string, buckets[i].end - buckets[i].start)
			values_write[i - bucket_index] = &array_writes
			buckets := buckets[i]
			wg.Add(1)
			go func(){
				thread_execution(&data, profile_comp, buckets, dist, &array_writes)
				wg.Done()
			}()
		}
		wg.Wait() // Wait for everyone to catch up
		buckets[bucket_index].start++ // update the current buckets tracker
		
		for _, i := range values_write {
			for _, value := range *i {
				fmt.Fprintf(f, *value);
			}
		}

		if len(buckets) > 1 && arr_pos % bucket_size == 0 {
			for f := buckets[bucket_index].end - bucket_size; f < buckets[bucket_index].end; f++ {
				data[f].profile = nil;
				data[f].name = empty_name;
				
			}
			bucket_index++
			end := time.Now().Sub(start)
			thread_depletion_time := fmt.Sprintf("One thread depleted in: %fs", end.Seconds())
			log.Println(thread_depletion_time)
			start = time.Now()
		}
		arr_pos++;
	}
	wg.Wait()
	f.Flush();
}

