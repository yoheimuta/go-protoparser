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

func TestParser_ParseOneof(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		permissive bool
		wantOneof  *parser.Oneof
		wantErr    bool
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
								Offset: 16,
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
								Offset: 37,
								Line:   3,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 65,
						Line:   4,
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
								Offset: 16,
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
								Offset: 43,
								Line:   4,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 71,
						Line:   5,
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
										Offset: 16,
										Line:   2,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 22,
										Line:   2,
										Column: 11,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 28,
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
										Offset: 49,
										Line:   4,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 62,
										Line:   4,
										Column: 18,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 68,
								Line:   5,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 96,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing inline comments",
			input: `oneof foo { // TODO: implementation
    string name = 4; // name
    SubMessage sub_message = 9; // sub_message
}
`,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						InlineComment: &parser.Comment{
							Raw: `// name`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 57,
									Line:   2,
									Column: 22,
								},
								LastPos: meta.Position{
									Offset: 63,
									Line:   2,
									Column: 28,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 40,
								Line:   2,
								Column: 5,
							},
						},
					},
					{
						Type:        "SubMessage",
						FieldName:   "sub_message",
						FieldNumber: "9",
						InlineComment: &parser.Comment{
							Raw: `// sub_message`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 97,
									Line:   3,
									Column: 33,
								},
								LastPos: meta.Position{
									Offset: 110,
									Line:   3,
									Column: 46,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 69,
								Line:   3,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 12,
							Line:   1,
							Column: 13,
						},
						LastPos: meta.Position{
							Offset: 34,
							Line:   1,
							Column: 35,
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 112,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `oneof foo {
    string name = 4;
};
`,
			permissive: true,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 16,
								Line:   2,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 34,
						Line:   3,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `oneof foo {
    string name = 4;
}
`,
			permissive: true,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "string",
						FieldName:   "name",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 16,
								Line:   2,
								Column: 5,
							},
						},
					},
				},
				OneofName: "foo",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 33,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
		{
			name: "accept options. See https://github.com/yoheimuta/go-protoparser/issues/39",
			input: `oneof something {
  option (validator.oneof) = {required: true};
  uint32 three_int = 5 [(validator.field) = {int_gt: 20}];
  uint32 four_int = 6 [(validator.field) = {int_gt: 100}];
  string five_regex = 7 [(validator.field) = {regex: "^[a-z]{2,5}$"}];
}
`,
			permissive: true,
			wantOneof: &parser.Oneof{
				OneofFields: []*parser.OneofField{
					{
						Type:        "uint32",
						FieldName:   "three_int",
						FieldNumber: "5",
						FieldOptions: []*parser.FieldOption{
							{
								OptionName: "(validator.field)",
								Constant:   "{int_gt:20}",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 67,
								Line:   3,
								Column: 3,
							},
						},
					},
					{
						Type:        "uint32",
						FieldName:   "four_int",
						FieldNumber: "6",
						FieldOptions: []*parser.FieldOption{
							{
								OptionName: "(validator.field)",
								Constant:   "{int_gt:100}",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 126,
								Line:   4,
								Column: 3,
							},
						},
					},
					{
						Type:        "string",
						FieldName:   "five_regex",
						FieldNumber: "7",
						FieldOptions: []*parser.FieldOption{
							{
								OptionName: "(validator.field)",
								Constant:   "{regex:\"^[a-z]{2,5}$\"}",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 185,
								Line:   5,
								Column: 3,
							},
						},
					},
				},
				Options: []*parser.Option{
					{
						OptionName: "(validator.oneof)",
						Constant:   "{required:true}",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 20,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 63,
								Line:   2,
								Column: 46,
							},
						},
					},
				},
				OneofName: "something",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 254,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "accept a simple option without a permissive option. See https://github.com/yoheimuta/go-protoparser/issues/57",
			input: `oneof something {
  option (my_option).a = true;
}
`,
			wantOneof: &parser.Oneof{
				Options: []*parser.Option{
					{
						OptionName: "(my_option).a",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 20,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 47,
								Line:   2,
								Column: 30,
							},
						},
					},
				},
				OneofName: "something",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 49,
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
				parser.WithPermissive(test.permissive),
			)
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
