/*
	Fast matching of profiles.

	! It is important to remember that we are making a lot of assumptions here
	! e.g. that profiles have equivalent values.
*/

package main


import (
	"log"
	"sort"
	"sync"
)

type FastMatch struct {
	reference *string;
	query *string;
	distance float64;
}

func identify_matches(reference_profiles string, query_profiles string, match_threshold float64) {
	/*
		Fast match isolates
		reference_profiles string: Input profiles for query against with the reference profiles
		query_profiles string: Profiles to query against the references with
		match_threshold uint: integer threshold to use in a match
	*/
	reference, query := load_profiles(reference_profiles, query_profiles)
	var wg sync.WaitGroup
	dist_function := distance_functions[DIST_FUNC].function
	outputs := make(map[string][]*FastMatch, len(*query))
	// TODO add threading limit
	for _, profile := range *query {
		outputs[profile.name] = make([]*FastMatch, 0, int(0.05 * float64(len(*reference)))) // Create capacity at 5% of reference values
		log.Printf("Querying distances for %s", profile.name)
		go get_distances(profile, reference, dist_function, match_threshold, outputs[profile.name], &wg)
		wg.Add(1)
	}
	wg.Wait()

}


func get_distances(query_profile  *Profile, reference_profiles *[]*Profile, dist_fn func(*[]int, *[]int) float64, match_threshold float64, output_values []*FastMatch, wg *sync.WaitGroup) {
	/*
		Tabulate all distances for a profile
	*/
	
	
	for _, r_profile := range *reference_profiles {
		output := dist_fn(query_profile.profile, r_profile.profile)
		if output <= match_threshold {
			output_values = append(output_values, &FastMatch{&r_profile.name, &query_profile.name, output})
		}
	}

	sort.Slice(output_values, func(i, j int) bool {
		return output_values[i].distance < output_values[j].distance
	})
	defer wg.Done()
	

}