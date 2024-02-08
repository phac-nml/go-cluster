/* Ingest an allele profiles and convert all values into integers

Matthew Wells: 2024-02-06
*/


package main


import (
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
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
	//defer file.Close()
	scanner := bufio.NewScanner(file) // ! Caps out at 64k line
	return scanner, file
}


func load_profile(file_path string) *[]*Profile {
	/*Split a tab delimited profile and convert it into allelic profile
	*/
	new_line_char := "\n"
	line_delimiter := "\t"
	numeric_base := 10 // could be 16 for hex
	file_scanner, file := _create_scanner(file_path)
	defer file.Close()
	var data []*Profile


	for file_scanner.Scan() {
		input_text := strings.Split(strings.TrimSuffix(file_scanner.Text(), new_line_char), line_delimiter)


		data_in := make([]int, len(input_text) - 1) // Create an array to populate

		for f, x := range input_text[1:len(input_text)] { // may require offset by -1
			// ? Can probasbly just call everything a base 64 to handle strings
			i, err := strconv.ParseInt(x, numeric_base, 64) // Converts integer in, to base 10 64bit number TODO may need to handle hex values for hashes
			if err != nil {
				// TODO trigger clean up here
				log.Println("Improperly formatted allele: " + string(x), err)
				i = 0 // overwrite bad code with allele profile of 0
			}
			data_in[f] = int(i)
		}
		//data = append(data, &data_in)// pop data back array
		new_profile := newProfile(input_text[0], &data_in)
		data = append(data, new_profile)// pop data back array
	}

	if err := file_scanner.Err(); err != nil {
		log.Fatal(err)
		os.Exit(5)
	}

	return &data

}