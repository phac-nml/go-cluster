// Test for heap writes


package main

import (
	"testing"
	"container/heap"
	"reflect"
)

type WriteValueT struct {
	value []WriteValue
	expected []int64
}

var s func(s string) []byte = func(s string) []byte { return []byte(s) } // Prop to make string

var QueueWriteTests = []WriteValueT{
	WriteValueT{
		[]WriteValue{
			WriteValue{0, s("1"), 1},
			WriteValue{2, s("2"), 2},
			WriteValue{3, s("3"), 3},
			WriteValue{4, s("4"), 4},
		},
		[]int64{0, 2, 3, 4},
	}, 
	WriteValueT{
		[]WriteValue{
			WriteValue{0, s("1"), 1},
			WriteValue{2, s("2"), 2},
			WriteValue{2, s("3"), 2},
			WriteValue{3, s("3"), 3},
			WriteValue{4, s("4"), 4},
		},
		[]int64{0, 2, 2, 3, 4}, // Need to add behaviour for the offset in the writes
	},
}

func TestWriteQueue(t *testing.T){

	for _, test := range QueueWriteTests {
		
		wheap := make(WriteQueue, len(test.value))
		for idx, _ := range test.value {
			wheap[idx] = &test.value[idx]
		}

		heap.Init(&wheap)
		outputs := make([]int64, 0, len(test.value))
		for wheap.Len() > 0 {
			item := heap.Pop(&wheap).(*WriteValue)
			outputs = append(outputs, item.key)
		}
		
		if !reflect.DeepEqual(outputs, test.expected) {
			t.Errorf("Output: does not match expected %v %v", outputs, test.expected)
		}
	}
}

// TODO need to benchmark heap vs append to array and sort