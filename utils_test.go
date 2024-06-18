package main

import (
	"testing"
	"reflect"
)

type SplitLineTest struct {
	string_in, newline_char, line_delimiter string;
	expected []string
}


var SplitLineTests = []SplitLineTest{
	SplitLineTest{ "1\t2\n", "\n", "\t", []string{"1", "2"}},
	SplitLineTest{ "1,2\n", "\n", ",", []string{"1", "2"}},
	SplitLineTest{ "1,2\n\n", "\n", ",", []string{"1", "2"}},
	SplitLineTest{ "1,,3\n\n", "\n", ",", []string{"1", "", "3"}},
	SplitLineTest{ "1\t\t3\n", "\n", "\t", []string{"1", "", "3"}},
	SplitLineTest{ "1\t\t\t3\n", "\n", "\t", []string{"1", "", "", "3"}},
	SplitLineTest{ "1      3\n", "\n", "  ", []string{"1", "", "", "3"}}, // double spaces as delimiter
}


func TestSplitLine(t *testing.T){
	for _, test := range SplitLineTests {
		if output := SplitLine(test.string_in, test.newline_char, test.line_delimiter); !reflect.DeepEqual(*output, test.expected) {
			t.Errorf("Output not equal to expected")
			t.Errorf("Output %+v", output)
		}
	}
}
