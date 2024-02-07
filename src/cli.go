/* Initial CLI for the the program, this CLI would not be posix complient yet. A library will
need to be installed to do so, as the dev time would be costly.

Matthew Wells: 2024-02-07
*/

package main


import (
	"flag"
	"fmt"
	"os"
)

func cli() {

	distance_func_help:= fmt.Sprintf(`Enter an integer denoting the distance function you would like to use:
	%s: %d
	%s: %d
	%s: %d
	%s: %d`, 
	ham.help, ham.assignment,
	ham_missing.help, ham_missing.assignment,
	scaled.help, scaled.assignment,
	scaled_missing.help, scaled_missing.assignment)


	flag.StringVar(&INPUT_PROFILE, "profile", "", "File path to your alleles profiles.")
	flag.IntVar(&CPU_LOAD_FACTOR, "load-factor", CPU_LOAD_FACTOR, 
	`Used to set the minimum number of values needed to use 
multi-threading. e.g. load-factor * available cpus = minimum number of profiles required for multithreading.
It is best to use the default value which indicates that you have more rows of data then CPUs on your computer.`)
	flag.IntVar(&DIST_FUNC, "distance", 0, distance_func_help)

	flag.Parse()
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Printf("File: %s", INPUT_PROFILE)
}
