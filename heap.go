/*
	Store outputs to a heap so that the min element can be popped from the heap each time
	to order output writes.
*/


package main 

import (
	"container/heap"
)


type WriteQueue []*WriteValue

func (h WriteHeap) Len() int {return len(h)}
// need to implement these properly
func (h WriteHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h WriteHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WriteHeap) Push(x *WriteValue) {
	*h = append(*h, x.key)
}

func (h *WriteHeap) Pop() *WriteValue {
	old := *h 
	n := len(old)
	x := old[n-1]
	*h = old[0: n-1]
	return x
}
