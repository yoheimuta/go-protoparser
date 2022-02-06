package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/internal/util_test"
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

func TestParser_ParseExtend(t *testing.T) {
	tests := []struct {
		name                       string
		input                      string
		inputBodyIncludingComments bool
		permissive                 bool
		wantExtend                 *parser.Extend
		wantErr                    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
extend Foo {
  int32 bar = 126;
}
`,
			wantExtend: &parser.Extend{
				MessageType: "Foo",
				ExtendBody: []parser.Visitee{
					&parser.Field{
						Type:        "int32",
						FieldName:   "bar",
						FieldNumber: "126",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 16,
								Line:   3,
								Column: 3,
							},
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 33,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an excerpt from the google/api/annotations.proto",
			input: `
extend google.protobuf.MethodOptions {
  // See HttpRule.
  HttpRule http = 72295728;
}`,
			wantExtend: &parser.Extend{
				MessageType: "google.protobuf.MethodOptions",
				ExtendBody: []parser.Visitee{
					&parser.Field{
						Type:        "HttpRule",
						FieldName:   "http",
						FieldNumber: "72295728",
						Comments: []*parser.Comment{
							{
								Raw: "// See HttpRule.",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 42,
										Line:   3,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 57,
										Line:   3,
										Column: 18,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 61,
								Line:   4,
								Column: 3,
							},
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 87,
						Line:   5,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `
extend Foo {
};
`,
			permissive: true,
			wantExtend: &parser.Extend{
				MessageType: "Foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 15,
						Line:   3,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `
extend Foo {
}
`,
			permissive: true,
			wantExtend: &parser.Extend{
				MessageType: "Foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 14,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(
				lexer.NewLexer(strings.NewReader(test.input)),
				parser.WithBodyIncludingComments(test.inputBodyIncludingComments),
				parser.WithPermissive(test.permissive),
			)
			got, err := p.ParseExtend()
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err, parsed=%v", got)
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantExtend) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantExtend))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
