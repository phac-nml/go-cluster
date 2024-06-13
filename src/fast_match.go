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
	"sync/atomic"
	"fmt"
	"bufio"
)

type FastMatch struct {
	reference *string;
	query *string;
	distance float64;
}

/*
	Waitgroup count struct and functions come from: https://stackoverflow.com/a/68995552
*/
type WaitGroupCount struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroupCount) Add(delta int) {
    atomic.AddInt64(&wg.count, int64(delta))
    wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupCount) Done() {
    atomic.AddInt64(&wg.count, -1)
    wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) GetCount() int64 {
    return int64(atomic.LoadInt64(&wg.count))
}

func identify_matches(reference_profiles string, query_profiles string, match_threshold float64, output *bufio.Writer) {
	/*
		Fast match isolates
		reference_profiles string: Input profiles for query against with the reference profiles
		query_profiles string: Profiles to query against the references with
		match_threshold uint: integer threshold to use in a match
	*/

	reference, query := load_profiles(reference_profiles, query_profiles)
	wg := WaitGroupCount{count: 0}
	dist_function := distance_functions[DIST_FUNC].function
	outputs := make([]*[]*FastMatch, len(*query))
	default_capacity := int(0.05 * float64(len(*reference)) + 1) // Create capacity at 5% of reference values

	// TODO add threading limit
	for idx, profile := range *query {
		output_arr := make([]*FastMatch, 0, default_capacity) 
		outputs[idx] = &output_arr
		profile_comp := profile
		wg.Add(1)
		log.Printf("Querying distances for %s", profile_comp.name)
		go func() {
			get_distances(profile_comp, reference, dist_function, match_threshold, &output_arr)
			wg.Done()
		}()

		if wg.GetCount() == FM_THREAD_LIMIT {
			log.Printf("Waiting for active threads to finish. %d", wg.GetCount())
			wg.Wait()
		}
	}
	wg.Wait()

	format_string := get_format_string()
	// TODO need to add log message for this as there could be no matches.
	for idx, matches := range outputs {
		prof_matches := *matches
		if len(prof_matches) == 0 {
			log.Printf("No matches identified for profile: %s", (*query)[idx].name)
			continue
		}
		for _, match := range prof_matches {
			fmt.Fprintf(output, format_string, *match.reference, *match.query, match.distance)
		}
	}
}


func get_distances(query_profile  *Profile, reference_profiles *[]*Profile, dist_fn func(*[]int, *[]int) float64, match_threshold float64, output_values *[]*FastMatch) {
	/*
		Tabulate all distances for a profile
	*/
	
	
	for _, r_profile := range *reference_profiles {
		output := dist_fn(query_profile.profile, r_profile.profile)
		if output <= match_threshold {
			*output_values = append(*output_values, &FastMatch{&r_profile.name, &query_profile.name, output})
		}
	}
	sort.Slice(*output_values, func(i, j int) bool {
		return (*output_values)[i].distance < (*output_values)[j].distance
	})
}