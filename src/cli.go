/* Initial CLI for the the program, this CLI would not be posix complient yet. A library will
need to be installed to do so, as the dev time would be costly.

Matthew Wells: 2024-02-07
*/

package main


import (
	flag "github.com/spf13/pflag"
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

	flag.StringVarP(&INPUT_PROFILE, "profile", "p", INPUT_PROFILE, "File path to your alleles profiles.")
	flag.IntVarP(&CPU_LOAD_FACTOR, "load-factor", "l", CPU_LOAD_FACTOR, 
	`Used to set the minimum number of values needed to use 
multi-threading. e.g. if (number of cpus * load factor) > number of table rows. Only a single thread will be used. `)
	flag.IntVarP(&DIST_FUNC, "distance", "f", 0, distance_func_help)
	flag.StringVarP(&MOLTEN_FILE, "molten", "m", MOLTEN_FILE, "File path to a previously generated output for conversion into a distance matrix.")
	flag.StringVarP(&OUTPUT_FILE, "output", "o", OUTPUT_FILE, "Name of output file.")
	flag.IntVarP(&BUFFER_SIZE, "buffer-size", "b", BUFFER_SIZE, "Larger buffers may increase performance.")
	flag.StringVarP(&COLUMN_DELIMITER, "column-delimiter", "d", COLUMN_DELIMITER, "Column delimiter")
	flag.StringVarP(&MISSING_ALLELE_STRING, "missing-allele-character", "a", MISSING_ALLELE_STRING, "String denoting missing alleles.")
	
	flag.Parse()
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
