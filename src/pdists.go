package main

import (
	_ "io"
	_ "bufio"
	"os"
	"fmt"
	_ "log"
	_ "strconv"
	_ "strings"
	"sync"
	_ "runtime"
	_ "time"
)

func hamming_dist(d1 *[]int, d2 *[]int) int {
	var diffs int = 0
	for g := range *d1 {
		if (*d1)[g] == (*d2)[g] {
			diffs++
		}
	}
	//if diffs == int(len(*d1)) {
	//	fmt.Fprintf(os.Stderr, "--%d\n%d\n--", *d1, *d2)
	//}
	return diffs
}

//func simple_thread(d *[]*[]int,  wg *sync.WaitGroup, c chan int, vals *[]int, start int, end int, arr_pos *int){
func simple_thread(d *[]*Profile,  wg *sync.WaitGroup, vals *Profile, start int, end int, init_spot int){

	for i := start; i < end; i++ { // todo make this a slice in the future
		x := hamming_dist((*d)[i].profile, vals.profile)
		//_ = x
		// todo can use buffer to make this faster
		fmt.Fprintf(os.Stdout, "Postion: %s ID: %s val: %d\n", vals.name, (*d)[i].name, x)
	}
	defer wg.Done()
}

//func buckets_indices(size int, window int) [][]int {
//	var bucks [][]int
//	cpu_load_factor := 10
//	
//	if window > size{
//		x := make([]int, 2)
//		x[0] = int(0)
//		x[1] = int(size)
//		bucks = append(bucks, x)
//		return bucks
//	}
//	
//	if size < (runtime.NumCPU() * cpu_load_factor) {
//		x := make([]int, 2)
//		x[0] = int(0)
//		x[1] = int(size)
//		bucks = append(bucks, x)
//		return bucks
//	}
//
//	for i := window; i < size; i = i + window {
//		x := make([]int, 2)
//		x[0] = int(i)-int(window)
//		x[1] = int(i)
//		bucks = append(bucks, x)
//	}
//	
//	x := make([]int, 2)
//	x[0] = int(size)- int(window) 
//	x[1] = int(size)
//	bucks = append(bucks, x)
//	return bucks
//}

var CPU_LOAD_FACTOR = 10
var INPUT_PROFILE string = ""

func main() {
	cli()

	//file, err := os.Open(os.Args[1])
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()

	//scanner := bufio.NewScanner(file) // ! caps out at 64K a line

	//var data []*[]int
	
	//for scanner.Scan() {
	//	input_text := strings.Split(strings.TrimSuffix(scanner.Text(), "\n"), "\t")
	//	data_in := make([]int, len(input_text))
	//	for f, x := range input_text[0:len(input_text)-1] {
	//		i, err := strconv.ParseInt(x, 10, 64)
	//		if err != nil{
	//			log.Println("Improperly formatted character", err);
	//		}else{
	//			data_in[f] = int(i)
	//		}
	//		
	//	}
	//	data = append(data, &data_in)
	//}

	//if err := scanner.Err(); err != nil {
	//	log.Fatal(err)
	//}

	//data_ := load_profile(os.Args[1])
	data_ := load_profile(INPUT_PROFILE)
	data := *data_

	run_data(&data)

	// Trending towards more being better
	//bucket_size := (len(data) / (runtime.NumCPU() * 2))
//
	//bucks := buckets_indices(len(data), bucket_size)
	//var arr_pos int = 1
//
	//var buck_idx int = 0
	//empty_array := make([]int, 1)
//
	//fmt.Fprintf(os.Stderr, "Threads used: %d\n", len(bucks))
	////fmt.Fprintf(os.Stderr, "Threads used: %d\n", bucks)
	//fmt.Fprintf(os.Stderr, "Allocating ~%d rows to a thread\n", bucket_size)
//
//
	//// todo can redistribute data across threads in the future at runtime
	//start := time.Now()
	//for g := range data[0:] {
	//	var wg sync.WaitGroup
	//	for i := int(buck_idx); i < len(bucks); i++ {
	//		wg.Add(1)
	//		go simple_thread(&data, &wg, data[g], bucks[i][0], bucks[i][1], int(bucks[buck_idx][1]) - int(bucket_size))
	//	}
	//	wg.Wait()
	//	bucks[buck_idx][0]++ // increment value by one to keep same data being ran twice
//
	//	if len(bucks) > 1 && arr_pos % int(bucket_size) == 0{
	//		// mark values no longer needed for garbage collection
	//		for f := int(bucks[buck_idx][1]) - int(bucket_size); f < bucks[buck_idx][1]; f++ {
	//			data[f].profile = &empty_array; // nuke array values for garbage collection
	//		}
	//		buck_idx++
	//		end := time.Now().Sub(start)
	//		fmt.Fprintf(os.Stderr, "One thread depleted in: %fs\n", end.Seconds())
	//		start = time.Now()
	//	}
	//	
	//	arr_pos++
	//}

	fmt.Fprintf(os.Stderr, "All threads depleted.\n")
	fmt.Fprintf(os.Stderr, "Done\n")

}