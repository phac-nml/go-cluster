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
	const new_line_char = "\n";
	line_delimiter := COLUMN_DELIMITER;;
	file_scanner, file := _create_scanner(file_path);
	defer file.Close();
	//var data []*Profile;
	log.Println("Ingesting profile and normalizing allele inputs.");
	
	var missing_value string = MISSING_ALLELE_STRING;
	// TODO verify that Scan moves file pointer up
	normalization_lookup := initialize_lookup(file_scanner, new_line_char, line_delimiter);
	data := create_profiles(file_scanner, normalization_lookup, new_line_char, line_delimiter, missing_value);
	//for file_scanner.Scan() {
	//	input_text := *split_line(file_scanner.Text(), new_line_char, line_delimiter);
	//	data_in := make([]int, len(input_text) - 1) // Create an array to populate
	//	for f, x := range input_text[1:len(input_text)] { // starting at position 1 as first value is the sample ID
	//		if missing_value != x {
	//			data_in[f] = (*normalization_lookup)[f].InsertValue(&x);
	//		}else{
	//			data_in[f] = missing_allele_value;
	//		}
	//		
	//	}
	//	new_profile := newProfile(input_text[0], &data_in);
	//	data = append(data, new_profile); // pop data back array
	//}
	//
	//if err := file_scanner.Err(); err != nil {
	//	log.Fatal(err);
	//	os.Exit(5);
	//}

	normalization_lookup = nil // Flag objects for GC
	log.Println("Finished ingesting profile.");

	return data

}