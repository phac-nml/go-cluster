/* Distance functions used to calculate allele profiles in go

To supply different return types a type and an interface may need to be added, as unions are not an option in go
TODO offer int or float return

TODO treating missing alleles (zeroes) as missing should be optional

Matthew Wells: 2024-04-06
*/

package main

type DistFunc struct {
	function   func(*[]int, *[]int) float64
	assignment int
	help       string
	truncate   bool // to truncate the output value to an integer or to remain as a float
}

var ham = DistFunc{function: hamming_distance, assignment: 0, help: "Hamming Distance", truncate: true}
var ham_missing = DistFunc{function: hamming_distance_missing, assignment: 1, help: "Hamming distance skipping missing values", truncate: true}
var scaled = DistFunc{function: scaled_distance, assignment: 2, help: "Scaled Distance", truncate: false}
var scaled_missing = DistFunc{function: scaled_distance_missing, assignment: 3, help: "Scaled distance skipping missing values", truncate: false}

// update distance functions, with their position in the array pertaining to their calling
var distance_functions = []DistFunc{ham, ham_missing, scaled, scaled_missing}

var DIST_FUNC = 0 // Distance function default

func hamming_distance(profile_1 *[]int, profile_2 *[]int) float64 {
	/* Hamming distance not including missing data

	 */
	p1 := *profile_1
	p2 := *profile_2
	count := 0
	profile_len := len(p1)
	for idx := 0; idx < profile_len; idx++ {
		if p1[idx] == 0 || p2[idx] == 0 {
			continue
		}
		if (p1[idx] ^ p2[idx]) != 0 {
			count++
		}
	}
	return float64(count)
}

func hamming_distance_missing(profile_1 *[]int, profile_2 *[]int) float64 {
	/* Returns hamming distance, with missing values counted as differences
	 */
	p1 := *profile_1
	p2 := *profile_2
	count := 0
	profile_len := len(p1)
	for idx := 0; idx < profile_len; idx++ {
		if (p1[idx] ^ p2[idx]) != 0 {
			count++
		}
	}
	return float64(count)
}

func scaled_distance(profile_1 *[]int, profile_2 *[]int) float64 {
	/* Scaled distance with missing data skipped, increment counter

	 */
	p1 := *profile_1
	p2 := *profile_2
	count_compared_sites := 0
	count_match := 0
	profile_len := len(p1)
	default_return := 100.0
	for idx := 0; idx < profile_len; idx++ {
		if p1[idx] == 0 || p2[idx] == 0 {
			continue
		}
		count_compared_sites++
		if (p1[idx] ^ p2[idx]) != 0 {
			// If not equal skip
			continue
		}
		count_match++
	}

	if count_compared_sites != 0 {
		cc_sites_f64 := float64(count_compared_sites)
		count_match_f64 := float64(count_match)
		scaled_value := default_return * ((cc_sites_f64 - count_match_f64) / cc_sites_f64)
		return scaled_value
	}
	return default_return
}

func scaled_distance_missing(profile_1 *[]int, profile_2 *[]int) float64 {
	/* Scaled distance with missing data counted as differences

	 */
	p1 := *profile_1
	p2 := *profile_2
	count_compared_sites := 0
	count_match := 0
	default_return := 100.0
	profile_len := len(p1)
	for idx := 0; idx < profile_len; idx++ {
		if (p1[idx] ^ p2[idx]) != 0 { // skip if the same
			continue
		}
		count_match++
	}

	if count_compared_sites != 0 {
		cc_sites_f64 := float64(profile_len)
		count_match_f64 := float64(count_match)
		scaled_value := default_return * ((cc_sites_f64 - count_match_f64) / cc_sites_f64)
		return scaled_value
	}
	return default_return
}
