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

type mockHasInlineCommentSetter struct {
	inlineComment *parser.Comment
}

func (m *mockHasInlineCommentSetter) SetInlineComment(comment *parser.Comment) {
	m.inlineComment = comment
}

func TestParser_MaybeScanInlineComment(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantInlineComment *parser.Comment
	}{
		{
			name: "parsing an empty",
		},
		{
			name: "parsing a C++-style comment on the current line",
			input: `int32 page_number = 2;  // Which page number do we want?
`,
			wantInlineComment: &parser.Comment{
				Raw: "// Which page number do we want?",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 24,
						Line:   1,
						Column: 25,
					},
					LastPos: meta.Position{
						Offset: 55,
						Line:   1,
						Column: 56,
					},
				},
			},
		},
		{
			name: "parsing a C-style comment on the current line",
			input: `int32 page_number = 2;  /* Which page number do we want?
*/
`,
			wantInlineComment: &parser.Comment{
				Raw: `/* Which page number do we want?
*/`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 24,
						Line:   1,
						Column: 25,
					},
					LastPos: meta.Position{
						Offset: 58,
						Line:   2,
						Column: 2,
					},
				},
			},
		},
		{
			name: "parsing a C++-style comment on the next line",
			input: `int32 page_number = 2;
// Which page number do we want?
`,
		},
		{
			name: "parsing a C-style comment on the next line",
			input: `int32 page_number = 2;
/* Which page number do we want?
*/
`,
		},
		{
			name: "parsing C++-style comments on the current and next line",
			input: `int32 page_number = 2;  // Which page number do we want?
// Number of results to return per page.
`,
			wantInlineComment: &parser.Comment{
				Raw: "// Which page number do we want?",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 24,
						Line:   1,
						Column: 25,
					},
					LastPos: meta.Position{
						Offset: 55,
						Line:   1,
						Column: 56,
					},
				},
			},
		},
		{
			name: "parsing C-style comments on the current and next line",
			input: `int32 page_number = 2;  /* Which page number do we want?
*/
/* Number of results to return per page.
*/
`,
			wantInlineComment: &parser.Comment{
				Raw: `/* Which page number do we want?
*/`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 24,
						Line:   1,
						Column: 25,
					},
					LastPos: meta.Position{
						Offset: 58,
						Line:   2,
						Column: 2,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			_, _ = p.ParseField()

			hasSetter := &mockHasInlineCommentSetter{}
			p.MaybeScanInlineComment(hasSetter)
			got := hasSetter.inlineComment

			if !reflect.DeepEqual(got, test.wantInlineComment) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantInlineComment))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
