package main

import (
	"testing"
)

func TestMaskLinks(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// just a link
		{
			"Here's my spammy page: http://hehefouls.netHAHAHA see you.",
			"Here's my spammy page: http://******************* see you.",
		},

		// without any link
		{
			"No links here!",
			"No links here!",
		},

		// two links
		{
			"Check this out: http://example.com and http://test.com",
			"Check this out: http://*********** and http://********",
		},

		// russian symbols
		{
			"Check this out: http://пара",
			"Check this out: http://********",
		},
	}

	for _, test := range tests {
		result := maskLinks(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', but got '%s'", test.input, test.expected, result)
		}
	}
}
