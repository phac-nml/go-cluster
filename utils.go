/*
	Utility functions for comparing distances.
*/

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)
const MissingAlleleValue int = 0

// Create an output buffer for writing too, if no input file is passed stdout will be used.
func CreateOutputBuffer(file_in string) (*bufio.Writer, os.File) {
	var f *bufio.Writer
	var file os.File
	if file_in != "" {
		file := open_file(file_in, os.O_WRONLY)
		f = bufio.NewWriterSize(io.Writer(file), BUFFER_SIZE)
	} else {
		f = bufio.NewWriterSize(os.Stdout, BUFFER_SIZE)
	}
	return f, file
}

// Split a profile lines on columns
func SplitLine(string_in string, new_line_char string, line_delimiter string) *[]string {
	output := strings.Split(strings.TrimRight(string_in, new_line_char), line_delimiter)
	return &output
}


// Users will expect a given distance matrix to have either decimals or just scalar values.
// Depending on the distance function used, To truncate the values accordingly the format string
// for Sprintf is returned by this function.
func GetFormatString() string {
	var format_expression string = "%s %s %.2f\n"
	if distance_functions[DIST_FUNC].truncate {
		format_expression = "%s %s %.0f\n"
	}
	return format_expression
}

// Create profiles for data profiles
func CreateProfiles(file_scanner *bufio.Scanner, lookup *[]*ProfileLookup, new_line_char string, line_delimiter string, missing_value string) *[]*Profile {
	
	var data []*Profile
	for file_scanner.Scan() {
		input_text := SplitLine(file_scanner.Text(), new_line_char, line_delimiter)
		data_in := make([]int, len(*input_text)-1)         // Create an array to populate
		new_profile := LineToProfile(input_text, &data_in, lookup, &missing_value)
		data = append(data, new_profile) // pop data back array
	}

	if err := file_scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Data contains: %d profiles.", len(data))
	return &data
}


// Convert a text line into a profile use in distance calculations
func LineToProfile(input_text *[]string, data_in *[]int, lookup *[]*ProfileLookup, missing_value *string) *Profile {

	input_data := *input_text
	no_allele := *missing_value
	for f, x := range input_data[1:len(input_data)] {
		if no_allele != x {
			(*data_in)[f] = (*lookup)[f].InsertValue(&x)
		}else {
			(*data_in)[f] = MissingAlleleValue
		}
	}
	out_profile := NewProfile((*input_text)[0], data_in)
	return out_profile
}