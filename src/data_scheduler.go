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
	bucket_size := data_length / (runtime_cpus * 2)
	return bucket_size
}

func buckets_indices(data_length int, bucket_size int) [][]int {
	var bucks [][]int
	cpu_load_factor := CPU_LOAD_FACTOR // Need to add description to global options
	window := bucket_size

	cpu_load_string := fmt.Sprintf("CPU load factor x%d", cpu_load_factor)
	log.Println(cpu_load_string)
	
	if window > data_length {
		x := make([]int, 2)
		x[0] = 0
		x[1] = data_length
		bucks = append(bucks, x)
		log.Println("Running single threaded as there are too few entries to justify multithreading.")
		return bucks
	}
	
	if data_length < (runtime.NumCPU() * cpu_load_factor) {
		x := make([]int, 2)
		x[0] = 0
		x[1] = data_length
		bucks = append(bucks, x)
		log.Println("Running single threaded as there are too few entries to justify multithreading.")
		return bucks
	}

	for i := window; i < data_length; i = i + window {
		x := make([]int, 2)
		x[0] = i - window
		x[1] = i
		bucks = append(bucks, x)
	}
	
	x := make([]int, 2)
	x[0] = bucks[len(bucks) - 1][1] 
	x[1] = data_length
	bucks = append(bucks, x)
	
	threads_running := fmt.Sprintf("Using %d threads for running.", len(bucks))
	log.Println(threads_running)
	profiles_to_thread := fmt.Sprintf("Allocating ~%d profiles per a thread.", window)
	log.Println(profiles_to_thread)
	//log.Println(bucks)
	return bucks
}





func thread_execution(data_slice *[]*Profile, waitgroup * sync.WaitGroup, profile_compare *Profile, start_idx int, end_idx int, dist_fn func(*[]int, *[]int) float64, f *bufio.Writer){
	/* Compute profile differences.
	
	data_slice: the data range to use for calculation against the profile to be compared too.
	waitgroup: waitgroup for the go routine to be a part of
	profile_compare: the profile being compared in all threads
	start_idx: The starting range in the profile to be used to initilize comparisons
	end_idx: The end range to calculate comparisons up too.
	dist_fn: The distance function to use for calculation of differences. Takes pointer to two profile to compare and returns a float 64
	*/

	truncate := distance_functions[DIST_FUNC].truncate
	defer f.Flush()

	// TODO need to pass a slice properly in the future
	for i := start_idx; i < end_idx; i++ {
		x := dist_fn((*data_slice)[i].profile, profile_compare.profile);

		if truncate {
			fmt.Fprintf(f, "%s %s %.0f\n", profile_compare.name, (*data_slice)[i].name, x)
		}else {
			fmt.Fprintf(f, "%s %s %.2f\n", profile_compare.name, (*data_slice)[i].name, x)
		}
	}

	defer waitgroup.Done()

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
	empty_array := make([]int, 1)
	empty_name := ""
	bucket_size := calculate_bucket_size(len(data), runtime.NumCPU())
	buckets := buckets_indices(len(data), bucket_size)
	arr_pos := 1


	for g := range data[0:] {
		var wg sync.WaitGroup
		profile_comp := data[g] // copy struct for each thread
		// TODO an incredible optimization here would be to go lockless, or re-use threads
		for i := bucket_index; i < len(buckets); i++ {
			wg.Add(1)
			go thread_execution(&data, &wg, profile_comp, buckets[i][0], buckets[i][1], dist, f)
		}
		wg.Wait() // Wait for everyone to catch up
		buckets[bucket_index][0]++ // update the current buckets tracker

		if len(buckets) > 1 && arr_pos % bucket_size == 0 {
			for f := buckets[bucket_index][1] - bucket_size; f < buckets[bucket_index][1]; f++ {
				data[f].profile = &empty_array;
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
}

