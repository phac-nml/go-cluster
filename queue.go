/*
	Store outputs to a heap so that the min element can be popped from the heap each time
	to order output writes.
*/

package main

// / Creating a queue using a min-heap in order to make writes to out to the file system sequential instead of relying on random access
type WriteQueue []*WriteValue

func (h WriteQueue) Len() int           { return len(h) }
func (h WriteQueue) Less(i, j int) bool { return h[i].key < h[j].key }
func (h WriteQueue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *WriteQueue) Push(x any) {
	n := len(*h)
	item := x.(*WriteValue)
	item.index = n
	*h = append(*h, item)
}

func (h *WriteQueue) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Prevents a memory leak
	item.index = -1
	*h = old[0 : n-1]
	return item
}
