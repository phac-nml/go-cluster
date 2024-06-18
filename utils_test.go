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


type LineToProfileTest struct {
	input_text []string;
	missing_value string;
	expected Profile;
}


var LineToProfileTests = []LineToProfileTest{
	LineToProfileTest{
		[]string{"profile1", "name2", "name1", "0"}, // outputs for each column should be one as they are in the same columns
		"0",
		Profile{ "profile1", &[]int{1, 1, MissingAlleleValue}},
	}, 
}

func TestLineToProfile(t *testing.T) {
	
	for _, test := range LineToProfileTests {
		profiles_len := len(test.input_text) - 1
		data_in := make([]int, profiles_len)
		lookups := make([]*ProfileLookup, profiles_len)
		val := 0
		for val < profiles_len {
			pLookup := NewProfileLookup()
			lookups[val] = pLookup
			val += 1
		}
		if output := LineToProfile(&test.input_text,  &data_in, &lookups, &test.missing_value); (*output).name != test.expected.name  || !reflect.DeepEqual(*(*output).profile, test.expected.profile) {
			t.Errorf("Output not equal to expected")
			t.Errorf("Output %+v expected %+v", (*output).profile, test.expected.profile)
		}
	}
}