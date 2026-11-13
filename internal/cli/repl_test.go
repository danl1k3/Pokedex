package cli

import (
	"testing"

)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "CHarLEmange   is  not  pokemon",
			expected: []string{"charlemange", "is", "not", "pokemon"},
		},
		{
			input: "h e l l o world",
			expected: []string{"h", "e", "l", "l", "o", "world"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "  ",
			expected: []string{},
		},
		{
			input: "PIKACHU",
			expected: []string{"pikachu"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		
		if len(actual) != len(c.expected) {
			t.Errorf("\nLength doesn`t match! \nexp: %v, \ngot: %v", c.expected, actual)
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("\nexpected: %v, got: %v", c.expected, actual)
			}

		}
	}
}
