// Package typedef has type definitions used throughout the hclust package.
package typedef

// SubCluster stores the node, distance and names of leafs for a subcluster.
type SubCluster struct {
	Leafa   int
	Leafb   int
	Lengtha float64
	Lengthb float64
	Node    int
}
