package main

import (
	"os"
	"log"
	"bufio"
	"io"
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



var CPU_LOAD_FACTOR = 10
var INPUT_PROFILE string = ""
var MOLTEN_FILE string = ""
var OUTPUT_FILE string = ""
var BUFFER_SIZE = 16384 // 3 times bigger then 4096

func main() {
	cli()

	if MOLTEN_FILE != "" && OUTPUT_FILE != "" {
		pariwise_to_matrix(MOLTEN_FILE, OUTPUT_FILE)
		os.Exit(0)
	}

	data_ := load_profile(INPUT_PROFILE)
	data := *data_

	var f *bufio.Writer
	if OUTPUT_FILE != "" {
		file := open_file(OUTPUT_FILE, os.O_WRONLY)
		defer file.Close()
		f = bufio.NewWriterSize(io.Writer(file), BUFFER_SIZE)
	}else{
		f = bufio.NewWriterSize(os.Stdout, BUFFER_SIZE)
	}

	
	defer f.Flush()
	
	run_data(&data, f)

	log.Println("All threads depleted.")
	log.Println("Done")

}