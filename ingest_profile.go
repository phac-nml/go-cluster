/* Ingest an allele profiles and convert all values into integers

Matthew Wells: 2024-02-06
*/

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// A rows profile information
type Profile struct {
	name    string
	profile *[]int
}

/*
Create a new profile struct for the passed input value.
*/
func NewProfile(name string, profile *[]int) *Profile {
	p := Profile{name: name, profile: profile}
	return &p
}

// A utility function for creating a new bufio.Scanner the same way each time.
func _create_scanner(file_path string, new_line_char string) (*bufio.Scanner, *os.File) {

	header_chars := GetHeaderSize(file_path, new_line_char)
	log.Println("Determined header file size.")
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file) // ! Caps out at 64k line
	if header_chars > bufio.MaxScanTokenSize {
		buffer := make([]byte, header_chars)
		scanner.Buffer(buffer, header_chars)
	}
	log.Println("Created file scanner.")

	return scanner, file
}

/*
Initialize the look up for each allele type per a given column. This takes the header of the input allele
profile to create a Map for each column. This allows for a unique identifier to be assigned to each unique
value passed into the program regardless of its format as the Map uses strings for its key values.
*/
func InitializeLookup(scanner *bufio.Scanner, new_line_char string, line_delimiter string) (*[]*ProfileLookup, *[]string) {

	first_line := scanner.Scan() // get header line
	scanner_err := scanner.Err()
	if scanner_err != nil {
		log.Printf("%+v, increasing buffer size for scanner.", scanner.Err())
	}
	if !first_line {
		log.Fatal("Input File appears to be empty.")
	}

	separated_line := SplitLine(scanner.Text(), new_line_char, line_delimiter)
	new_array := make([]*ProfileLookup, len(*separated_line))
	for idx := range new_array {
		new_array[idx] = NewProfileLookup()
	}
	return &new_array, separated_line
}

/*
Split a tab delimited profile and convert it into allelic profile.

The input is a string (to an existing file) of horizontally listed allele profiles, The contents
of the allele profiles can be string, hashes, integers etc. As we normalize the inputs by assigning
a unique ID per a column input.
*/
func LoadProfile(file_path string) *[]*Profile {
	new_line_char := NEWLINE_CHARACTER
	line_delimiter := COLUMN_DELIMITER
	log.Printf("Column Delimiter used: %s", line_delimiter)
	file_scanner, file := _create_scanner(file_path, new_line_char)
	defer file.Close()
	//var data []*Profile;
	log.Println("Ingesting profile and normalizing allele inputs.")
	var missing_value string = MISSING_ALLELE_STRING
	normalization_lookup, _ := InitializeLookup(file_scanner, new_line_char, line_delimiter)
	data := CreateProfiles(file_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value)

	normalization_lookup = nil // Flag objects for GC
	log.Println("Finished ingesting profile.")

	return data
}

/*
Reference and query profile are loaded for fast matching against each other. This function is slightly
different than the one used for loading a single profile for distance matrix generation as both the
of the profiles need to be "normalized" using the same data structure to make sure both profiles
receive the same allele code between two files.

reference_profiles string: Input profiles for query against with the reference profiles
query_profiles: Profiles to query against the references with
*/
func LoadProfiles(reference_profiles string, query_profiles string) (*[]*Profile, *[]*Profile) {

	var missing_value string = MISSING_ALLELE_STRING
	new_line_char := NEWLINE_CHARACTER
	line_delimiter := COLUMN_DELIMITER
	log.Printf("Column Delimiter used: %s", line_delimiter)
	reference_scanner, ref_file := _create_scanner(reference_profiles, new_line_char)
	query_scanner, query_file := _create_scanner(query_profiles, new_line_char)
	defer ref_file.Close()
	defer query_file.Close()

	log.Println("Ingesting and normalizing reference profiles.")
	normalization_lookup, reference_headers := InitializeLookup(reference_scanner, new_line_char, line_delimiter)
	ref_data := CreateProfiles(reference_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value)

	log.Println("Ingesting and normalizing query profiles.")
	// Get first line of scanner to verify inputs are the same
	first_line_query := query_scanner.Scan()
	if !first_line_query {
		log.Fatal("Query File appears to be empty.")
	}
	// Get first line to skip header and get profiles
	query_headers := SplitLine(query_scanner.Text(), new_line_char, line_delimiter)
	CompareProfileHeaders(query_headers, reference_headers)

	query_data := CreateProfiles(query_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value)

	// Append query profiles to reference profiles
	*ref_data = append(*ref_data, *query_data...)

	normalization_lookup = nil
	log.Println("Finished ingesting and normalizing profiles.")
	return ref_data, query_data
}

/*
Compare the columns from your queries vs the references. This is a rudimentary
check to verify that you are passing allelic profiles generated with the same inputs.
*/
func CompareProfileHeaders(query_headers *[]string, reference_headers *[]string) {
	len_query := len(*query_headers)
	len_ref := len(*reference_headers)
	if len_query != len_ref {
		log.Fatalf("Different number of columns present in query (%d) vs reference (%d).", len_ref, len_query)
	}
	for idx := range *query_headers {
		if q_h, r_h := (*query_headers)[idx], (*reference_headers)[idx]; q_h != r_h {
			log.Fatalf("Mismatch in column names between query (%s) and reference (%s).", q_h, r_h)
		}
	}

}

/*
Ingest a previously generated symmetric data matrix for use in clustering. All data
within the matrix will be cast into a float.
*/
func IngestMatrix(input string) ([][]float64, []string) {
	scanner, _ := _create_scanner(input, NEWLINE_CHARACTER)
	scanner.Scan() // Throw away first line
	var matrix [][]float64
	var ids []string
	line := 0
	for scanner.Scan() {
		separated_line := *SplitLine(scanner.Text(), NEWLINE_CHARACTER, COLUMN_DELIMITER)
		row := make([]float64, len(separated_line)-1) // Separated Line length includes column ID
		ids = append(ids, separated_line[0])
		for idx, value := range separated_line[1:] {
			parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err != nil {
				log.Fatalf("Could not convert value %s into a float64 value. Did you use this program to generate your distance matrix? [Line: %d Column %d]", value, line, idx)
			}
			row[idx] = parsed
		}
		matrix = append(matrix, row)
		line++
	}
	return matrix, ids
}
