package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

func TestParser_ParseOneof(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantOneof *parser.Oneof
		wantErr   bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an invalid; without oneof",
			input: `foo {
    string name = 4;
    SubMessage sub_message = 9;
}
`,
			wantErr: true,
		},
		{
			name: "parsing an invalid; without }",
			input: `oneof foo {
    string name = 4;
    SubMessage sub_message = 9;
`,
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `oneof foo {
    string name = 4;
    SubMessage sub_message = 9;
}
`,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 17,
								Line:   2,
								Column: 5,
							},
						},
					},
					{
						Type:        "SubMessage",
						FieldName:   "sub_message",
						FieldNumber: "9",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 38,
								Line:   3,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   1,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an emptyStatement",
			input: `oneof foo {
    string name = 4;
    ;
    SubMessage sub_message = 9;
}
`,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 17,
								Line:   2,
								Column: 5,
							},
						},
					},
					{
						Type:        "SubMessage",
						FieldName:   "sub_message",
						FieldNumber: "9",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 44,
								Line:   4,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   1,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing comments",
			input: `oneof foo {
    // name
    string name = 4;
    // sub_message
    SubMessage sub_message = 9;
}
`,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						Comments: []*parser.Comment{
							{
								Raw: `// name`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 17,
										Line:   2,
										Column: 5,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 29,
								Line:   3,
								Column: 5,
							},
						},
					},
					{
						Type:        "SubMessage",
						FieldName:   "sub_message",
						FieldNumber: "9",
						Comments: []*parser.Comment{
							{
								Raw: `// sub_message`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 50,
										Line:   4,
										Column: 5,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 69,
								Line:   5,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   1,
						Column: 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got, err := p.ParseOneof()
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantOneof) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantOneof))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
