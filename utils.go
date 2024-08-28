/*
	Utility functions for comparing distances.
*/

package main

import (
	"bufio"
	"bytes"
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

// Get The size of the first line to determine the maximum amount of charactars to hold.
func GetHeaderSize(file_path string, new_line_char string) int {
	const header_increase_size int = 8
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	header_len := 0

	new_line_bytes := []byte(new_line_char)
	data_buffer := make([]byte, len(new_line_bytes))
	s := bufio.NewReader(file)
	for {
		if _, err := s.Read(data_buffer); err != nil {
			if err == io.EOF {
				log.Fatal("Reached end of file without finding new line character.")
			} else {
				log.Fatal(err)
			}
		}
		if bytes.Equal(data_buffer, new_line_bytes) {
			break
		}
		header_len = header_len + len(new_line_bytes)
	}

	file.Close()
	return header_len * header_increase_size

}

// Create profiles for data profiles
func CreateProfiles(file_scanner *bufio.Scanner, lookup *[]*ProfileLookup, new_line_char string, line_delimiter string, missing_value string) *[]*Profile {

	var data []*Profile
	for file_scanner.Scan() {
		input_text := SplitLine(file_scanner.Text(), new_line_char, line_delimiter)
		data_in := make([]int, len(*input_text)-1) // Create an array to populate
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
	for f, x := range input_data[1:] {
		if no_allele != x {
			(*data_in)[f] = (*lookup)[f].InsertValue(&x)
		} else {
			(*data_in)[f] = MissingAlleleValue
		}
	}
	out_profile := NewProfile((*input_text)[0], data_in)
	return out_profile
}
