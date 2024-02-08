package main

import (
	_ "io"
	_ "bufio"
	"os"
	"fmt"
	_ "log"
	_ "strconv"
	_ "strings"
	_ "sync"
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


var CPU_LOAD_FACTOR = 10
var INPUT_PROFILE string = ""
var MOLTEN_FILE string = ""
var OUTPUT_FILE string = ""
var BUFFER_SIZE = 16384 // 3 times bigger then 4096
//var output_channel_length = 2000 // could probably make this a bucket length
//var output_channel_length = 10 // Making smaller for testing

func main() {
	cli()

	if MOLTEN_FILE != "" && OUTPUT_FILE != "" {
		pariwise_to_matrix(MOLTEN_FILE, OUTPUT_FILE)
		os.Exit(0)
	}

	data_ := load_profile(INPUT_PROFILE)
	data := *data_

	run_data(&data)

	fmt.Fprintf(os.Stderr, "All threads depleted.\n")
	fmt.Fprintf(os.Stderr, "Done\n")

}