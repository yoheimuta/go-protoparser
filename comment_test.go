package protoparser

import (
	"reflect"
	"testing"
)

func TestParseComments(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantComments []string
	}{
		{
			name: "parsing empty creates no comments",
		},
		{
			name:  "parsing one line creates a comment",
			input: `// binary is an image binary. Required.`,
			wantComments: []string{
				"// binary is an image binary. Required.",
			},
		},
		{
			name: "parsing lines creates comments",
			input: `// binary is an image binary. Required.
            // hogehoge
            `,
			wantComments: []string{
				"// binary is an image binary. Required.",
				"// hogehoge",
			},
		},
	}

	for _, test := range tests {
		lex := lex(test.input)
		got := parseComments(lex)
		if !reflect.DeepEqual(got, test.wantComments) {
			t.Errorf("[%s] got %v, but want %v", test.name, got, test.wantComments)
		}
	}
}
