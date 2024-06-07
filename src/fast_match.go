/*
	Fast matching of profiles.

	! It is important to remember that we are making a lot of assumptions here
	! e.g. that profiles have equivalent values.
*/

package main


//import (
//	"os"
//	"log"
//	_ "bufio"
//	"strings"
//)

//func load_profiles(profiles_path string, query_path string){
//	const new_line_char = "\n"
//	line_delimiter := COLUMN_DELIMITER;
//	file_scanner, file := _create_scanner(file_path)
//	defer file.Close()
//	var data []*Profile
//	log.Println("Ingesting profile and normalizing allele inputs.");
//	const missing_allele_value int = 0;
//	var missing_value string = MISSING_ALLELE_STRING;
//	// TODO verify that Scan moves file pointer up
//	normalization_lookup := initialize_lookup(file_scanner, new_line_char, line_delimiter);
//	for file_scanner.Scan() {
//		input_text := strings.Split(strings.TrimSuffix(file_scanner.Text(), new_line_char), line_delimiter)
//
//
//		data_in := make([]int, len(input_text) - 1) // Create an array to populate
//
//		for f, x := range input_text[1:len(input_text)] { // starting at position 1 as first value is the sample ID
//
//			if missing_value != x {
//				data_in[f] = (*normalization_lookup)[f].InsertValue(&x);
//			}else{
//				data_in[f] = missing_allele_value;
//			}
//			
//		}
//	
//		new_profile := newProfile(input_text[0], &data_in)
//		data = append(data, new_profile)// pop data back array
//	}
//
//	normalization_lookup = nil // Flag objects for GC
//	if err := file_scanner.Err(); err != nil {
//		log.Fatal(err);
//		os.Exit(5);
//	}
//	log.Println("Finished ingesting profile.");
//
//	return &data
//}