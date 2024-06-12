/* Ingest an allele profiles and convert all values into integers

Matthew Wells: 2024-02-06
*/


package main


import (
	"os"
	"log"
	"bufio"
)

type Profile struct {
	name string
	profile *[]int
}

func newProfile(name string, profile *[]int) *Profile {
	p := Profile{name: name, profile: profile}
	return &p
}

func _create_scanner(file_path string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file) // ! Caps out at 64k line
	return scanner, file
}


func initialize_lookup(scanner *bufio.Scanner, new_line_char string, line_delimiter string) *[]*ProfileLookup {
	first_line := scanner.Scan();
	if !first_line {
		log.Fatal("Input File appears to be empty.");
		os.Exit(1);
	}

	split_line := split_line(scanner.Text(), new_line_char, line_delimiter);
	new_array := make([]*ProfileLookup, len(*split_line))
	for idx, _ := range new_array {
		new_array[idx] = NewProfile();
	}
	return &new_array;
}

func load_profile(file_path string) *[]*Profile {
	/*
		Split a tab delimited profile and convert it into allelic profile
	*/
	new_line_char := NEWLINE_CHARACTER;
	line_delimiter := COLUMN_DELIMITER;
	file_scanner, file := _create_scanner(file_path);
	defer file.Close();
	//var data []*Profile;
	log.Println("Ingesting profile and normalizing allele inputs.");
	var missing_value string = MISSING_ALLELE_STRING;
	// TODO verify that Scan moves file pointer up
	normalization_lookup := initialize_lookup(file_scanner, new_line_char, line_delimiter);
	data := create_profiles(file_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value);


	normalization_lookup = nil // Flag objects for GC
	log.Println("Finished ingesting profile.");

	return data
}

func load_profiles(reference_profiles string, query_profiles string) (*[]*Profile, *[]*Profile) {
	/*
		reference_profiles string: Input profiles for query against with the reference profiles
		query_profiles: Profiles to query against the references with
	*/
	var missing_value string = MISSING_ALLELE_STRING;
	new_line_char := NEWLINE_CHARACTER;
	line_delimiter := COLUMN_DELIMITER;
	reference_scanner, ref_file := _create_scanner(reference_profiles);
	query_scanner, query_file := _create_scanner(query_profiles);
	defer ref_file.Close();
	defer query_file.Close();

	log.Println("Ingesting and normalizing reference profiles.")
	normalization_lookup := initialize_lookup(reference_scanner, new_line_char, line_delimiter)
	ref_data := create_profiles(reference_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value)

	log.Println("Ingesting and normalizing query profiles.")
	query_data := create_profiles(query_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value)

	normalization_lookup = nil
	log.Println("Finished ingesting and noramlizing profiles.")
	return ref_data, query_data
}