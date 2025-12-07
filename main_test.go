package main

import (	
	
	"testing"
)	




func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "GoLang  is  Awesome",
			expected: []string{"golang", "is", "awesome"},
		},
		{
			input:    "   Multiple    Spaces   ",
			expected: []string{"multiple", "spaces"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("For input '%s', expected length %d but got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("For input '%s', at index %d, expected '%s' but got '%s'", c.input, i, expectedWord, word)
			}
		}
	}

}
