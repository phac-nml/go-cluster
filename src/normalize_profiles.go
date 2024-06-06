/*
Normalize the inputs of the profiles from strings to ints
*/

package main


//import (
//	"os"
//	"log"
//	"strings"
//)


type ProfileNormalizer interface {
	InsertValue() int
}

type Node struct {
	Next *Node
	Char rune
	Value int64
}

type Trie struct {
	Root Node
	Counter int // Increment everytime a value is inserted
}

/*
	Insert a new value into the trie if it does not exist
*/
//func (e *Trie) InsertValue(new_string *string){
//
//}


type ProfileLookup struct {
	Counter int;
	Lookup map[string]int;
}

/*
	Insert a value into the map

*/
func (p *ProfileLookup) InsertValue(new_value *string) int {
	value, ok := p.Lookup[*new_value]
	if ok {
		return value;
	} else {
		p.Counter++;
		p.Lookup[*new_value] =  p.Counter;
	}
	return p.Counter;
} 

/*
	Create a new type of ProfileLookup
	* In the future It may be worth calculating an initial value size for the map
*/
func NewProfile() *ProfileLookup {
	return &ProfileLookup{Counter: 0, Lookup: make(map[string]int)};
}
