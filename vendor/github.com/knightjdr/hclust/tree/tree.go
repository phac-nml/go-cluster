package tree

import (
	"errors"

	"github.com/knightjdr/hclust/typedef"
)

// Tree references a tree in newick format and the leaf order.
type Tree struct {
	Newick string
	Order  []string
}

// Create generates a newick tree in string format and returns the order
// of the clustering.
func Create(dendrogram []typedef.SubCluster, names []string) (tree Tree, err error) {
	// Return if names length does not match matix length.
	if len(names) != len(dendrogram)+1 {
		err = errors.New("The names vector must have the same dimension as the leaf number")
		return
	}

	// Dendrogram clusters/leaf number.
	n := len(dendrogram)

	// Create map of nodes to dendrogram indicies.
	nodeMap := make(map[int]int, n)
	for i, cluster := range dendrogram {
		nodeMap[cluster.Node] = i
	}

	// Begin with top node, iterate through left and right branches and add to
	// ordering.
	level := Descend(n, 2*n, nodeMap, dendrogram, names)
	tree.Newick = level.Newick
	tree.Order = level.Order
	return
}
