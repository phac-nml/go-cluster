// Parallel Distances
//
// A utility for basic calculation of distances on allele profiles (in parallel), fast-matching, clustering 
// and dendrogram generation
//
// This utility is still in development however.
package main

import (
	"os"
	"log"
	"fmt"
	"github.com/integrii/flaggy"
)


var CPU_LOAD_FACTOR int = 100
var FM_THREAD_LIMIT int64 = 100
var COLUMN_DELIMITER = "\t"
var NEWLINE_CHARACTER = "\n"
var MISSING_ALLELE_STRING = "0"
var INPUT_PROFILE string = ""
var OUTPUT_FILE string = ""
var REFERENCE_PROFILES string = ""
var MATCH_THRESHOLD float64 = 10
var BUFFER_SIZE int = 16384 // 3 times bigger then 4096
var LINKAGE_METHOD int = 0
var distance_matrix *flaggy.Subcommand
var convert_matrix *flaggy.Subcommand
var fast_match *flaggy.Subcommand
var tree *flaggy.Subcommand
const version string = "0.0.1"


func cli() {
	flaggy.SetName("Parallel Distances")
	flaggy.SetDescription("A program for getting distances between allelic profiles and creating distance matrices.")
	flaggy.SetVersion(version);
	flaggy.DefaultParser.ShowHelpOnUnexpected = true;
	
	distance_matrix = flaggy.NewSubcommand("distances")
	distance_matrix.Description = "Compute all pairwise distances between the specified input profile."

	distance_func_help:= fmt.Sprintf(`Enter an integer denoting the distance function you would like to use:
	%s: %d
	%s: %d
	%s: %d
	%s: %d`, 
	ham.help, ham.assignment,
	ham_missing.help, ham_missing.assignment,
	scaled.help, scaled.assignment,
	scaled_missing.help, scaled_missing.assignment)

	buffer_help := fmt.Sprintf("The default buffer size is: %d. Larger buffers may increase performance.", BUFFER_SIZE)
	distance_matrix.String(&INPUT_PROFILE, "i", "input", "File path to your alleles profiles.")
	distance_matrix.Int(&CPU_LOAD_FACTOR, "l", "load-factor",
	`Used to set the minimum number of values needed to use 
multi-threading. e.g. if (number of cpus * load factor) > number of table rows. Only a single thread will be used. `)
	distance_matrix.Int(&DIST_FUNC, "d", "distance", distance_func_help)
	distance_matrix.String(&OUTPUT_FILE, "o", "output", "Name of output file. If nothing is specified results will be sent to stdout.")
	distance_matrix.Int(&BUFFER_SIZE, "b", "buffer-size", buffer_help)
	distance_matrix.String(&COLUMN_DELIMITER, "c", "column-delimiter", "Column delimiter")
	distance_matrix.String(&MISSING_ALLELE_STRING, "m", "missing-allele-character", "String denoting missing alleles.")
	

	convert_matrix = flaggy.NewSubcommand("convert")
	convert_matrix.Description = "Convert the pairwise distance generated by the program into a distance matrix."

	convert_matrix.String(&INPUT_PROFILE, "i", "input", "File path to a previously generated output for conversion into a distance matrix.")
	convert_matrix.String(&OUTPUT_FILE, "o", "output", "Name of output file. If nothing is specified results will be sent to stdout.")


	fast_match = flaggy.NewSubcommand("fast-match")
	fast_match.Description = "Tabulate distances between a query profile and reference profiles. Only distances exceeding a threshold will be kept."
	thread_limit_help := fmt.Sprintf("Limit the number of goroutines run at one time. Default: %d", FM_THREAD_LIMIT)
	fast_match.String(&INPUT_PROFILE, "i", "input", "File path to profiles for querying.")
	fast_match.String(&REFERENCE_PROFILES, "r", "reference", "File path to reference profiles to query against.")
	fast_match.String(&COLUMN_DELIMITER, "c", "column-delimiter", "Column delimiter")
	fast_match.String(&MISSING_ALLELE_STRING, "m", "missing-allele-character", "String denoting missing alleles.")
	fast_match.Int(&DIST_FUNC, "d", "distance", distance_func_help)
	fast_match.Float64(&MATCH_THRESHOLD, "t", "threshold", "Threshold for matching alleles.")
	fast_match.String(&OUTPUT_FILE, "o", "output", "Name of output file. If nothing is specified results will be sent to stdout.")
	fast_match.Int64(&FM_THREAD_LIMIT, "l", "goroutine-limit", thread_limit_help)

	tree = flaggy.NewSubcommand("tree")
	tree.String(&INPUT_PROFILE, "i", "input", "File path to previously generate distance matrix.")
	tree.String(&OUTPUT_FILE, "o", "output", "Name of output file. If nothing is specified results will be sent to stdout.")
	tree.Int(&LINKAGE_METHOD, "l", "linkage-method", linkage_methods_help)
	
	flaggy.AttachSubcommand(distance_matrix, 1);
	flaggy.AttachSubcommand(convert_matrix, 1);
	flaggy.AttachSubcommand(tree, 1);
	flaggy.AttachSubcommand(fast_match, 1);
	flaggy.Parse()

	if len(os.Args) <= 1 {
		flaggy.ShowHelpAndExit("No inputs passed");
	}
}

func main() {
	//	Parallel distances
	cli()
	output_buffer, file_out := CreateOutputBuffer(OUTPUT_FILE)
	defer file_out.Close()
	
	if distance_matrix.Used {

		if len(os.Args) <= 2 {
			flaggy.ShowHelpAndExit("No commands selected.");
		}
		data_ := LoadProfile(INPUT_PROFILE)
		data := *data_
		RunData(&data, output_buffer);
		log.Println("All threads depleted.")
		output_buffer.Flush()
	}

	if convert_matrix.Used {
		if len(os.Args) <= 2 {
			flaggy.ShowHelpAndExit("No commands selected.");
		}
		PairwiseToMatrix(INPUT_PROFILE, OUTPUT_FILE)
	}


	if fast_match.Used {
		if len(os.Args) <= 2 {
			flaggy.ShowHelpAndExit("No commands selected.");
		}
		if distance_functions[DIST_FUNC].assignment < 2 && MATCH_THRESHOLD < 1 {
			flaggy.ShowHelpAndExit("Distance function selected requires a value >1 for selection.")
		}
		IdentifyMatches(REFERENCE_PROFILES, INPUT_PROFILE, MATCH_THRESHOLD, output_buffer)
		output_buffer.Flush()
	}

	if tree.Used {
		if len(os.Args) <= 2 {
			flaggy.ShowHelpAndExit("No commands selected.");
		}
		if INPUT_PROFILE == "" {
			flaggy.ShowHelpAndExit("No input file selected");
		}
		if LINKAGE_METHOD > LINKAGE_METHODS[len(LINKAGE_METHODS)-1].match_value || LINKAGE_METHOD < 0 {
			flaggy.ShowHelpAndExit("Invalid linkage method selected.");
		}
		Cluster(INPUT_PROFILE, LINKAGE_METHOD, output_buffer)
		output_buffer.Flush()
	}
	
	log.Println("Done")

}