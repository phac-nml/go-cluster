/* Distance functions used to calculate allele profiles in go

To supply different return types a type and an interface may need to be added, as unions are not an option in go
TODO offer int or float return

Matthew Wells: 2024-04-06
*/

package main

type DistFunc struct {
	function func(*[]int, *[]int) float64
	assignment int
	help string
}


var ham = DistFunc{function: hamming_distance, assignment: 0, help: "Hamming Distance"}
var ham_missing = DistFunc{function: hamming_distance_missing, assignment: 1, help: "Hamming distance skipping missing values"}
var scaled = DistFunc{function: scaled_distance, assignment: 2, help: "Scaled Distance"}
var scaled_missing = DistFunc{function: scaled_distance_missing, assignment: 3, help: "Scaled distance skipping missing values"}


var DIST_FUNC = 0 // Distance function default

func select_dist_func() func(*[]int, *[]int) float64 {
	/*
	Select  distance function based on a cmd-line paramter

	This logic can definately be cleaned up to make things easier to maintain
	*/

	switch DIST_FUNC {
		case ham.assignment:
			return ham.function
		case ham_missing.assignment:
			return ham_missing.function
		case scaled.assignment:
			return scaled.function
		case scaled_missing.assignment:
			return scaled_missing.function
	}
	return ham_missing.function
}

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
		if  (p1[idx] ^ p2[idx]) != 0 {
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
	for idx := 0; idx < profile_len; idx++  {
		if  (p1[idx] ^ p2[idx]) != 0 {
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
	for idx := 0; idx < profile_len; idx++  {
		if p1[idx] == 0 || p2[idx] == 0 {
			continue
		}
		count_compared_sites++
		if  (p1[idx] ^ p2[idx]) != 0 {
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
	for idx := 0; idx < profile_len; idx++  {
		if  (p1[idx] ^ p2[idx]) != 0 { // skip if the same
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

