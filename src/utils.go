/*
	Utility functions for comparing distances.
*/

package main

import (
	"strings"
	"log"
	"os"
	"bufio"
)


/*
	Split a profile lines on columns
*/
func split_line(string_in string, new_line_char string, line_delimiter string) *[]string {
	output := strings.Split(strings.TrimSuffix(string_in, new_line_char), line_delimiter)
	return &output
}

/*
	get format output string
*/
func get_format_string() string {
	var format_expression string = "%s %s %.2f\n";
	if distance_functions[DIST_FUNC].truncate {
		format_expression = "%s %s %.0f\n";
	}
	return format_expression
}


/*
	Create profiles for data profiles
*/
func create_profiles(file_scanner *bufio.Scanner, lookup *[]*ProfileLookup, new_line_char string, line_delimiter string, missing_value string) *[]*Profile {
	const missing_allele_value int = 0;
	var data []*Profile;
	for file_scanner.Scan() {
		input_text := *split_line(file_scanner.Text(), new_line_char, line_delimiter);
		data_in := make([]int, len(input_text) - 1) // Create an array to populate
		for f, x := range input_text[1:len(input_text)] { // starting at position 1 as first value is the sample ID
			if missing_value != x {
				data_in[f] = (*lookup)[f].InsertValue(&x);
			}else{
				data_in[f] = missing_allele_value;
			}
			
		}
		new_profile := newProfile(input_text[0], &data_in);
		data = append(data, new_profile); // pop data back array
	}

	if err := file_scanner.Err(); err != nil {
		log.Fatal(err);
		os.Exit(5);
	}
	log.Printf("Data contains: %d profiles.", len(data))
	return &data;
}