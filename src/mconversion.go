/*


Read in output molten distance matrix and convert it into a distance matrix
in a memory respectful way e.g. memory mapping and maximizing concurrent writes

Initial algo:
1. create a buffer in memory
	- buffer needs to be (# of ids1 * # of ids2) e.g. square
	- ids should be sorted
2. Calculate the place each value should go in said buffer
	- Symmetry can be exploited here e.g. the mirror of coordinates also equals the value in
3. write Buffer (later will be a file)

For creating file in memory:

1. get the larges int so we can pad the all the other entries when writing
	- the padding will be the number of characters in the number + 1, as the padding will include the tab character
2. Write will be calculated from the offset at the start.
	- each write will be the offset + pad size to prevent over writing entries
	- if padding characters pose an issue they can be stripped in another pass of writing
	- each line will have to be padded with a new line character

TODO write to array in parallel, then sort to create sequential writes to a file
	- This is currently no in parallel, but a buffer can be written to sorting writes output writes


apparently file systems to not like writing to files in parallel

=============================================================================
~~~~~~~~~~~~~~~~~~~ In the future however ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
A buffer can be written to in parallel then dumped to a file

- The buffer can contain structs of positions to writen too
- The buffer can be sorted to make writes sequential (offsets can be altered)
- buffer can be flushed to the file
- new lines can be added after again
=============================================================================

TODO this currently wont work as output format changed

Matthew Wells: 2024-02-07
*/

package main


import ( 
	"os"
	"io"
	"fmt"
	"strconv"
	"log"
	"bufio"
	"strings"
	"sort"
)


const (
	profile_1_pos = 0
	profile_2_pos = 1
	comparison_pos = 2
)

func open_file(file_path string, open_type int) *os.File {
	file, err := os.OpenFile(file_path, int(open_type), 0555)
	if err != nil {
		if os.IsNotExist(err) && open_type == os.O_WRONLY {
			_, err := os.Create(file_path)
			if err != nil {
				log.Fatal(err)
				file = open_file(file_path, open_type)
			}
		}else{
			log.Fatal(err)
		}
		
	}
	return file
}


func parse_int(value string) int {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(val)
}


func get_keys(value *map[string]bool) (*[]string, int ){
	map_vals := make([]string, len(*value))
	vals := 0
	longest_key := 0
	for k, _ := range *value {
		if len(k) > longest_key {
			longest_key = len(k)
		}
		map_vals[vals] = k
		vals++
	}
	map_vals = map_vals[0:vals]
	sort.Strings(map_vals)
	
	return &map_vals, longest_key
}



func unique_values(file_path string) (*[]string, int, int) {
	set := map[string]bool{}
	file := open_file(file_path, os.O_RDONLY)
	reader := bufio.NewReader(io.Reader(file))
	longest_val := 0

	for {
		rl, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		data := strings.Split(rl, " ")
		len_int := len(data[comparison_pos])
		if len_int > longest_val {
			longest_val = len_int
		}
		set[data[profile_1_pos]] = true
		set[data[profile_2_pos]] = true
	}
	defer file.Close()
	sorted_keys, longest_key := get_keys(&set)
	return sorted_keys, longest_val, longest_key
	
}

func pad_value(characters string, mask []byte) []byte {
	/* Pad the string to not offset the file writing locations
	characters: the characters to pad
	mask: the byte mask to re-use
	*/
	values := make([]byte, len(mask))
	copy(values, mask)
	for i, v := range characters {
		
		values[i] = byte(v)
	}
	return values
}

func make_mask(modulus int) []byte {
	// Create a buffer of spaces for each value to fill in
	mask := make([]byte, modulus+1)
	for i := range mask {
		mask[i] = ' '
	}
	mask[modulus] = '\t'
	return mask
}

func write_matrix(input_path string, output_path string, positions *map[string]int, longest_val int){

	// input data fields
	file := open_file(input_path, os.O_RDONLY)
	reader := bufio.NewReader(io.Reader(file))
	
	// output data fields
	output := open_file(output_path, os.O_WRONLY | os.O_CREATE) // making this a buffered output may be easier

	// columns size
	modulus := len(*positions)
	modulus_64 := int64(modulus)

	mask := make_mask(longest_val)
	pad_len := int64(len(mask))
	rows := modulus_64
	for{
		rl, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// Get data positions
		data := strings.Split(rl, " ")
		string_val_up := data[comparison_pos]
		string_val_up = string_val_up[:len(string_val_up)-1] // drop new line character
		string_val := pad_value(string_val_up, mask)
		
		// Id locations
		p1 := (*positions)[data[profile_1_pos]]
		p2 := (*positions)[data[profile_2_pos]]
		sp1 := calculate_buffer_position(p1, p2, modulus)
		sp2 := calculate_buffer_position(p2, p1, modulus)


		// Write at offsets
		output.Seek(0, io.SeekStart)
		b1, err := output.WriteAt(string_val, sp1 * pad_len)
		_ = b1
		if err != nil {
			log.Fatal(err)
		}

		output.Seek(0, io.SeekStart)
		b2, err := output.WriteAt(string_val, sp2 * pad_len)
		_ = b2
		
		if err != nil {
			log.Fatal(err)
		}
	}

	// Replace tabs with new line characters in output file
	newline := []byte("\n")
	output.Seek(0, io.SeekStart)
	fmt.Fprintf(os.Stderr, "Adding new line characters to reformatted matrix.\n")
	for i := modulus_64; i < modulus_64 * modulus_64; i = i + modulus_64 {
		b, err := output.WriteAt(newline, (i * pad_len) - 1)
		_ = b
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(os.Stderr, "Rows: %d\n", rows)
	defer file.Close()
	defer output.Close()

}

func get_matrix_values(file_path string, positions *map[string]int, buffer *[]int){
	file := open_file(file_path, os.O_RDONLY)
	reader := bufio.NewReader(io.Reader(file))
	modulus := len(*positions)
	for {
		rl, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data := strings.Split(rl, " ")
		string_val := data[comparison_pos]
		string_val = string_val[:len(string_val)-1] // drop new line character
		
		// Id locations
		p1 := (*positions)[data[profile_1_pos]]
		p2 := (*positions)[data[profile_2_pos]]

		int_val := parse_int(string_val)
		sp1 := calculate_buffer_position(p1, p2, modulus)
		sp2 := calculate_buffer_position(p2, p1, modulus)
		(*buffer)[sp1] = int_val
		(*buffer)[sp2] = int_val

	}
	defer file.Close()
}


func calculate_buffer_size(key_len int) int{
	size := key_len * key_len
	return size
}

func calculate_buffer_position(p1 int, p2 int, modulus int) int64 {
	/*
	rows and columns provided here, the value can go to two positions,
	and the modulus is used to calculate where based on the rows and columns.

	e.g. to get rows (p1) * modulus + p2 (columns) and flip the location for the other value
	*/
	//fmt.Fprintf(os.Stdout, "%d %d\n", p1, p2)
	return int64((p1 * modulus) + p2)
}

func print_buffer(buffer *[]int, modulus int, buff_size int){
	// ! This will go once memory mapping is implemented
	fmt.Fprintf(os.Stdout, "\n")
	for i := 1; i < buff_size; i++{
		fmt.Fprintf(os.Stdout, "%d\t", (*buffer)[i-1])
		if i % modulus == 0 {
			fmt.Fprintf(os.Stdout, "\n")
		}
	}
	fmt.Fprintf(os.Stdout, "\n")
}



func pariwise_to_matrix(input_file string, output_file string) {
	/* Old main function for the mconversion routine.

	Needs to be optimized to use buffers for output, so the buffer can be sorted and sequential writes are read to disk

	TODO need to include sample names when reading them in to add to annotate the matrix
	
	*/
	sorted_keys, longest_val, longest_key := unique_values(input_file)
	key_positions := map[string]int{}
	longest_in := fmt.Sprintf("Longest key: %d", longest_key)
	
	log.Println(longest_in)

	vals := 0
	for _, v := range *sorted_keys {
		key_positions[v] = vals
		vals++
	}

	var output string = output_file
	write_matrix(input_file, output, &key_positions, longest_val)

	log.Println("Done", longest_val)

}



