package protoparser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestParseComments(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantComments      []string
		wantRecentScanned string
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
            bytes binary = 2 [(validator.field) = {length_gt: 0}];
            `,
			wantComments: []string{
				"// binary is an image binary. Required.",
				"// hogehoge",
			},
			wantRecentScanned: "bytes",
		},
	}

	for _, test := range tests {
		lex := lexer.NewLexer(strings.NewReader(test.input))
		got := parseComments(lex)
		if !reflect.DeepEqual(got, test.wantComments) {
			t.Errorf("[%s] got %v, but want %v", test.name, got, test.wantComments)
		}
		if lex.Text() != test.wantRecentScanned {
			t.Errorf("[%s] got %v, but want %v", test.name, lex.Text(), test.wantRecentScanned)
		}
	}
}
