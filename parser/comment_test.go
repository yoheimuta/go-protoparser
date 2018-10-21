package parser_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
)

func TestComment_IsCStyle(t *testing.T) {
	tests := []struct {
		name         string
		inputComment *parser.Comment
		wantIsCStyle bool
	}{
		{
			name: "parsing a C-style comment",
			inputComment: &parser.Comment{
				Raw: `/*
comment
*/
`,
			},
			wantIsCStyle: true,
		},
		{
			name: "parsing a C++-style comment",
			inputComment: &parser.Comment{
				Raw: "// comment",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := test.inputComment.IsCStyle()
			if got != test.wantIsCStyle {
				t.Errorf("got %v, but want %v", got, test.wantIsCStyle)
			}
		})
	}
}

func TestComment_Lines(t *testing.T) {
	tests := []struct {
		name         string
		inputComment *parser.Comment
		wantLines    []string
	}{
		{
			name: "parsing a C-style comment line",
			inputComment: &parser.Comment{
				Raw: `/*comment*/`,
			},
			wantLines: []string{
				"comment",
			},
		},
		{
			name: "parsing C-style comment lines",
			inputComment: &parser.Comment{
				Raw: `/* comment1
comment2
*/`,
			},
			wantLines: []string{
				" comment1",
				"comment2",
				"",
			},
		},
		{
			name: "parsing a C++-style comment",
			inputComment: &parser.Comment{
				Raw: "// comment",
			},
			wantLines: []string{
				" comment",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := test.inputComment.Lines()
			if !reflect.DeepEqual(got, test.wantLines) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantLines))
			}
		})
	}

}
