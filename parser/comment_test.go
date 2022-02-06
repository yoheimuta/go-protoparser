package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/internal/util_test"
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
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
			name: "parsing a C-style comment",
			inputComment: &parser.Comment{
				Raw: `/*comment*/`,
			},
			wantLines: []string{
				"comment",
			},
		},
		{
			name: "parsing C-style comments",
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

func TestParser_ParseComments(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantComments []*parser.Comment
	}{
		{
			name: "parsing an empty",
		},
		{
			name: "parsing a C++-style comment",
			input: `// comment
`,
			wantComments: []*parser.Comment{
				{
					Raw: `// comment`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 0,
							Line:   1,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 9,
							Line:   1,
							Column: 10,
						},
					},
				},
			},
		},
		{
			name: "parsing C++-style comments",
			input: `// comment
// comment2
`,
			wantComments: []*parser.Comment{
				{
					Raw: `// comment`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 0,
							Line:   1,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 9,
							Line:   1,
							Column: 10,
						},
					},
				},
				{
					Raw: `// comment2`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 11,
							Line:   2,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 21,
							Line:   2,
							Column: 11,
						},
					},
				},
			},
		},
		{
			name: "parsing a C-style comment",
			input: `/*
comment
*/`,
			wantComments: []*parser.Comment{
				{
					Raw: `/*
comment
*/`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 0,
							Line:   1,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 12,
							Line:   3,
							Column: 2,
						},
					},
				},
			},
		},
		{
			name: "parsing C-style comments",
			input: `/*
comment
*/
/*
comment2
*/`,
			wantComments: []*parser.Comment{
				{
					Raw: `/*
comment
*/`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 0,
							Line:   1,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 12,
							Line:   3,
							Column: 2,
						},
					},
				},
				{
					Raw: `/*
comment2
*/`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 14,
							Line:   4,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 27,
							Line:   6,
							Column: 2,
						},
					},
				},
			},
		},
		{
			name: "parsing a C-style comment and a C++-style comment",
			input: `/*
comment
*/

// comment2
`,
			wantComments: []*parser.Comment{
				{
					Raw: `/*
comment
*/`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 0,
							Line:   1,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 12,
							Line:   3,
							Column: 2,
						},
					},
				},
				{
					Raw: `// comment2`,
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 15,
							Line:   5,
							Column: 1,
						},
						LastPos: meta.Position{
							Offset: 25,
							Line:   5,
							Column: 11,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got := p.ParseComments()

			if !reflect.DeepEqual(got, test.wantComments) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantComments))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
