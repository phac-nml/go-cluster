/*
	Perform some basic clustering on generated distance matrices.
*/

package main

import (
	"bufio"
	"fmt"
	hclust "github.com/knightjdr/hclust"
	"log"
)

type LinkageMethod struct {
	identifier  string
	match_value int
}

var average LinkageMethod = LinkageMethod{"average", 0}
var centroid LinkageMethod = LinkageMethod{"centroid", 1}
var complete LinkageMethod = LinkageMethod{"complete", 2}
var mcquitty LinkageMethod = LinkageMethod{"mcquitty", 3}
var median LinkageMethod = LinkageMethod{"median", 4}
var single LinkageMethod = LinkageMethod{"single", 5}
var ward LinkageMethod = LinkageMethod{"ward", 6}

var LINKAGE_METHODS []LinkageMethod = []LinkageMethod{
	average,
	centroid,
	complete,
	mcquitty,
	median,
	single,
	ward}

var linkage_methods_help string = func() string {
	start_message := "Please enter an integer corresponding to one of the linkage method of your choice: "
	for _, value := range LINKAGE_METHODS {
		start_message += fmt.Sprintf("%s: %d ", value.identifier, value.match_value)
	}
	return start_message
}()

func GetLinkageMethod(value int) string {
	var linkage_method LinkageMethod
	switch value {
	case average.match_value:
		linkage_method = average
	case centroid.match_value:
		linkage_method = centroid
	case complete.match_value:
		linkage_method = complete
	case mcquitty.match_value:
		linkage_method = mcquitty
	case median.match_value:
		linkage_method = median
	case single.match_value:
		linkage_method = single
	case ward.match_value:
		linkage_method = ward
	default:
		log.Fatalf("Invalid linkage method specified: %d", value)
	}
	return linkage_method.identifier
}

/*
Cluster the profiles and create a dendrogram output
*/
func Cluster(input_file string, linkage_value int, f *bufio.Writer) {

	matrix, ids := IngestMatrix(input_file)
	linkage_method := GetLinkageMethod(linkage_value)
	log.Printf("Using %s as the linkage method for clustering", linkage_method)
	clustered_data, err := hclust.Cluster(matrix, linkage_method)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Adding labels to dendrogram.")
	newick, err := hclust.Tree(clustered_data, ids)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(f, "%s\n", newick.Newick)
}
