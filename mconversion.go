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

TODO Create buffer for to contain sorted writes
TODO write to array in parallel, then sort to create sequential writes to a file
	- This is currently not in parallel, but a buffer can be written to sorting writes output writes


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
TODO The output of this adds a final tab charactar that needs to be removed.
Matthew Wells: 2024-02-07
*/

package main

import (
	"bufio"
	"container/heap"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	profile_1_pos  = 0
	profile_2_pos  = 1
	comparison_pos = 2
	separator      = '\t'
)

// / This value is set up so that values can be stored before writing out to disc
// / The index field is used exclusively by the min-heap structure as it is needed in some of its operations
type WriteValue struct {
	key   int64
	value []byte
	index int // needed to update the heap interface
}

func open_file(file_path string, open_type int) *os.File {
	file, err := os.OpenFile(file_path, int(open_type), 0o666)
	if err != nil {
		if os.IsNotExist(err) && open_type == os.O_WRONLY {
			_, err := os.Create(file_path)
			if err != nil {
				log.Fatalf("Failed to create output file. %s", err)
			}
			file = open_file(file_path, open_type)
		} else {
			log.Fatalf("Failed to open file. %s", err)
		}

	}
	return file
}

func get_keys(value *map[string]bool) (*[]string, int) {
	map_vals := make([]string, len(*value))
	vals := 0
	longest_key := 0
	for k := range *value {
		if len(k) > longest_key {
			longest_key = len(k)
		}
		map_vals[vals] = k
		vals++
	}
	map_vals = map_vals[0:vals]
	// Leaving the sort as an option, but as the output buffer of the distance calculation step is now
	// written to a buffer concurrently maintaining order, the inputs always enter this process in the order of
	// the lower triangle
	sort.Strings(map_vals) // TODO need to use a stable sort

	return &map_vals, longest_key
}

func unique_values(file_path string) (*[]string, int) {
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
	if longest_key > longest_val {
		return sorted_keys, longest_key
	}
	return sorted_keys, longest_val

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
	mask[modulus] = byte(COLUMN_DELIMITER[0]) // convert delimiter into byte string
	return mask
}

func WriteQueueToFile(queue *WriteQueue, output_file *os.File) {
	output_file.Seek(0, io.SeekStart)
	for queue.Len() > 0 {
		output_value := heap.Pop(queue).(*WriteValue)
		name_out, err := output_file.WriteAt(output_value.value, output_value.key)
		_ = name_out
		if err != nil {
			log.Fatal(err)
		}
	}
}

func write_matrix(input_path string, output_path string, positions *map[string]int, longest_val int) {
	/*

		TODO optimize for sequential writes, priority queue is implemented now, to finish off the implementation
			1. Fill array containing data pairs of output position, and text out
			2. Sort array on position out
			3. Subtract difference in location from each sequential write.
				- e.g. File pointer is going to be increasing each time, so the next write should be relevant to that offset
			4. after flushing buffer, return file pointer to 0
			6. After all data is iterated through, flush buffer and remaining entries

		? The inputs are now always in the lower triangle order, so there can be some alterations to this method
	*/

	// input data fields
	file := open_file(input_path, os.O_RDONLY)
	reader := bufio.NewReader(io.Reader(file))

	// output data fields
	output := open_file(output_path, os.O_WRONLY|os.O_CREATE) // making this a buffered output may be easier

	// columns size
	modulus := len(*positions) + 1 // increase length by one to include data name row
	modulus_64 := int64(modulus)

	mask := make_mask(longest_val)
	pad_len := int64(len(mask))

	var buffered_writes int = 1000
	write_heap := make(WriteQueue, 0, buffered_writes) // Set capacity to write buffer size
	heap.Init(&write_heap)

	/*
		For optimizing the outputs, an AVL tree can be used to a balance them as the
		positions used are calculated. Then the buffer can be purged afterwards.
	*/
	rows := modulus_64
	for {
		if write_heap.Len() == buffered_writes {
			WriteQueueToFile(&write_heap, output)
		}

		rl, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// Get data positions
		data := strings.Split(rl, " ")
		string_val_up := data[comparison_pos]
		string_val_up = string_val_up[:len(string_val_up)-1] // drop new line character
		string_val := pad_value(string_val_up, mask)

		// Get data positions, + 1 to offset their position in the matrix
		p1 := (*positions)[data[profile_1_pos]] + 1
		p2 := (*positions)[data[profile_2_pos]] + 1

		sp1 := calculate_buffer_position(p1, p2, modulus)
		sp2 := calculate_buffer_position(p2, p1, modulus)

		// TODO making more writes than nessecary
		profile_1_name := pad_value(data[profile_1_pos], mask)
		heap.Push(&write_heap, &WriteValue{key: int64(p1) * modulus_64 * pad_len, value: profile_1_name, index: 0})
		heap.Push(&write_heap, &WriteValue{key: int64(p1) * pad_len, value: profile_1_name, index: 0}) // Write the columns position

		profile_2_name := pad_value(data[profile_2_pos], mask)
		heap.Push(&write_heap, &WriteValue{key: int64(p2) * modulus_64 * pad_len, value: profile_2_name, index: 0})
		heap.Push(&write_heap, &WriteValue{key: int64(p2) * pad_len, value: profile_2_name, index: 0}) // Column Position to write

		// * name pad_len should only be applied to one value, this will differ for the top row
		heap.Push(&write_heap, &WriteValue{key: sp1 * pad_len, value: string_val, index: 0})
		heap.Push(&write_heap, &WriteValue{key: sp2 * pad_len, value: string_val, index: 0})

	}

	if write_heap.Len() > 0 {
		WriteQueueToFile(&write_heap, output)
	}

	// Add byte mask to start of file, to prevent binary inclusion
	output.Seek(0, io.SeekStart)
	_, err := output.WriteAt(mask, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Replace tabs with new line characters in output file
	output.Seek(0, io.SeekStart)
	log.Println("Adding new line characters to reformatted matrix.")
	newline := []byte("\n")
	for i := modulus_64; i < modulus_64*modulus_64; i = i + modulus_64 {
		b, err := output.WriteAt(newline, (i*pad_len)-1)
		_ = b
		if err != nil {
			log.Fatal(err)
		}
	}

	// Remove final tab character added
	output.Seek(-1, io.SeekEnd)
	output.Write([]byte{' '})

	log.Printf("Rows output: %d", rows)

	defer file.Close()
	defer output.Close()
}

func calculate_buffer_position(p1 int, p2 int, modulus int) int64 {
	/*
		rows and columns provided here, the value can go to two positions,
		and the modulus is used to calculate where based on the rows and columns.

		e.g. to get rows (p1) * modulus + p2 (columns) and flip the location for the other value
	*/
	return int64((p1 * modulus) + p2)
}

/*
Function used to create a pairwise distance matrix from a previously generated molten output.
*/
func PairwiseToMatrix(input_file string, output_file string) {
	/* Old main function for the mconversion routine.

	Needs to be optimized to use buffers for output, so the buffer can be sorted and sequential writes are read to disk
	*/
	sorted_keys, longest_val := unique_values(input_file)
	key_positions := map[string]int{}

	vals := 0
	for _, v := range *sorted_keys {
		key_positions[v] = vals
		vals++
	}

	var output string = output_file
	write_matrix(input_file, output, &key_positions, longest_val)
}
