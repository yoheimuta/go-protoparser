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

func TestParser_ParseMessage(t *testing.T) {
	tests := []struct {
		name                       string
		input                      string
		inputBodyIncludingComments bool
		permissive                 bool
		wantMessage                *parser.Message
		wantErr                    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
message Outer {
  option (my_option).a = true;
  message Inner {
    int64 ival = 1;
  }
  map<int32, string> my_map = 2;
}
`,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 46,
								Line:   3,
								Column: 30,
							},
						},
					},
					&parser.Message{
						MessageName: "Inner",
						MessageBody: []parser.Visitee{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 70,
										Line:   5,
										Column: 5,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 50,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 88,
								Line:   6,
								Column: 3,
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "2",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 92,
								Line:   7,
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
						Offset: 123,
						Line:   8,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing another excerpt from the official reference",
			input: `
message outer {
  option (my_option).a = true;
  message inner {
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
`,
			wantMessage: &parser.Message{
				MessageName: "outer",
				MessageBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 46,
								Line:   3,
								Column: 30,
							},
						},
					},
					&parser.Message{
						MessageName: "inner",
						MessageBody: []parser.Visitee{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 70,
										Line:   5,
										Column: 5,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 50,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 88,
								Line:   6,
								Column: 3,
							},
						},
					},
					&parser.Field{
						IsRepeated:  true,
						Type:        "inner",
						FieldName:   "inner_message",
						FieldNumber: "2",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 92,
								Line:   7,
								Column: 3,
							},
						},
					},
					&parser.Field{
						Type:        "EnumAllowingAlias",
						FieldName:   "enum_field",
						FieldNumber: "3",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 128,
								Line:   8,
								Column: 3,
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 163,
								Line:   9,
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
						Offset: 194,
						Line:   10,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an empty MessageBody",
			input: `
message Outer {
}
`,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 17,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing comments",
			input: `
message outer {
  // option
  option (my_option).a = true;
  // message
  message inner {   // Level 2
    int64 ival = 1;
  }
  // field
  repeated inner inner_message = 2;
  // enum
  enum EnumAllowingAlias {
    option allow_alias = true;
  }
  EnumAllowingAlias enum_field =3;
  // map
  map<int32, string> my_map = 4;
  // oneof
  oneof foo {
    string name = 5;
    SubMessage sub_message = 6;
  }
  // reserved
  reserved "bar";
}
`,
			wantMessage: &parser.Message{
				MessageName: "outer",
				MessageBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
						Comments: []*parser.Comment{
							{
								Raw: `// option`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 19,
										Line:   3,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 27,
										Line:   3,
										Column: 11,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 31,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 58,
								Line:   4,
								Column: 30,
							},
						},
					},
					&parser.Message{
						MessageName: "inner",
						MessageBody: []parser.Visitee{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 108,
										Line:   7,
										Column: 5,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// message`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 62,
										Line:   5,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 71,
										Line:   5,
										Column: 12,
									},
								},
							},
						},
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// Level 2",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 93,
									Line:   6,
									Column: 21,
								},
								LastPos: meta.Position{
									Offset: 102,
									Line:   6,
									Column: 30,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 75,
								Line:   6,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 126,
								Line:   8,
								Column: 3,
							},
						},
					},
					&parser.Field{
						IsRepeated:  true,
						Type:        "inner",
						FieldName:   "inner_message",
						FieldNumber: "2",
						Comments: []*parser.Comment{
							{
								Raw: `// field`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 130,
										Line:   9,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 137,
										Line:   9,
										Column: 10,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 141,
								Line:   10,
								Column: 3,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 216,
										Line:   13,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 241,
										Line:   13,
										Column: 30,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// enum`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 177,
										Line:   11,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 183,
										Line:   11,
										Column: 9,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 187,
								Line:   12,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 245,
								Line:   14,
								Column: 3,
							},
						},
					},
					&parser.Field{
						Type:        "EnumAllowingAlias",
						FieldName:   "enum_field",
						FieldNumber: "3",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 249,
								Line:   15,
								Column: 3,
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "4",
						Comments: []*parser.Comment{
							{
								Raw: `// map`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 284,
										Line:   16,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 289,
										Line:   16,
										Column: 8,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 293,
								Line:   17,
								Column: 3,
							},
						},
					},
					&parser.Oneof{
						OneofFields: []*parser.OneofField{
							{
								Type:        "string",
								FieldName:   "name",
								FieldNumber: "5",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 353,
										Line:   20,
										Column: 5,
									},
								},
							},
							{
								Type:        "SubMessage",
								FieldName:   "sub_message",
								FieldNumber: "6",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 374,
										Line:   21,
										Column: 5,
									},
								},
							},
						},
						OneofName: "foo",
						Comments: []*parser.Comment{
							{
								Raw: `// oneof`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 326,
										Line:   18,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 333,
										Line:   18,
										Column: 10,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 337,
								Line:   19,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 404,
								Line:   22,
								Column: 3,
							},
						},
					},
					&parser.Reserved{
						FieldNames: []string{
							`"bar"`,
						},
						Comments: []*parser.Comment{
							{
								Raw: `// reserved`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 408,
										Line:   23,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 418,
										Line:   23,
										Column: 13,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 422,
								Line:   24,
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
						Offset: 438,
						Line:   25,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing inline comments",
			input: `
message SearchRequest {
  string query = 1;
  int32 page_number = 2;  // Which page number do we want?
  int32 result_per_page = 3;  // Number of results to return per page.
  enum EnumAllowingAlias {
    option allow_alias = true;
  } // Alias
}
`,
			wantMessage: &parser.Message{
				MessageName: "SearchRequest",
				MessageBody: []parser.Visitee{
					&parser.Field{
						Type:        "string",
						FieldName:   "query",
						FieldNumber: "1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
						},
					},
					&parser.Field{
						Type:        "int32",
						FieldName:   "page_number",
						FieldNumber: "2",
						InlineComment: &parser.Comment{
							Raw: `// Which page number do we want?`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 71,
									Line:   4,
									Column: 27,
								},
								LastPos: meta.Position{
									Offset: 102,
									Line:   4,
									Column: 58,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 47,
								Line:   4,
								Column: 3,
							},
						},
					},
					&parser.Field{
						Type:        "int32",
						FieldName:   "result_per_page",
						FieldNumber: "3",
						InlineComment: &parser.Comment{
							Raw: `// Number of results to return per page.`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 134,
									Line:   5,
									Column: 31,
								},
								LastPos: meta.Position{
									Offset: 173,
									Line:   5,
									Column: 70,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 106,
								Line:   5,
								Column: 3,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 206,
										Line:   7,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 231,
										Line:   7,
										Column: 30,
									},
								},
							},
						},
						InlineComment: &parser.Comment{
							Raw: `// Alias`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 237,
									Line:   8,
									Column: 5,
								},
								LastPos: meta.Position{
									Offset: 244,
									Line:   8,
									Column: 12,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 177,
								Line:   6,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 235,
								Line:   8,
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
						Offset: 246,
						Line:   9,
						Column: 1,
					},
				},
			},
		},
		{
			name: "skipping a last comment",
			input: `
message SearchRequest {
  string query = 1;
  // last comment
}
`,
			wantMessage: &parser.Message{
				MessageName: "SearchRequest",
				MessageBody: []parser.Visitee{
					&parser.Field{
						Type:        "string",
						FieldName:   "query",
						FieldNumber: "1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
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
						Offset: 63,
						Line:   5,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing last comments",
			input: `
message SearchRequest {
  string query = 1;
  // last first comment
  /* last second comment */
}
`,
			inputBodyIncludingComments: true,
			wantMessage: &parser.Message{
				MessageName: "SearchRequest",
				MessageBody: []parser.Visitee{
					&parser.Field{
						Type:        "string",
						FieldName:   "query",
						FieldNumber: "1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
						},
					},
					&parser.Comment{
						Raw: `// last first comment`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 47,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 67,
								Line:   4,
								Column: 23,
							},
						},
					},
					&parser.Comment{
						Raw: `/* last second comment */`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 71,
								Line:   5,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 95,
								Line:   5,
								Column: 27,
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
						Offset: 97,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an extend",
			input: `
message Outer {
  extend Foo {
    int32 bar = 126;
  }
}
`,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Extend{
						MessageType: "Foo",
						ExtendBody: []parser.Visitee{
							&parser.Field{
								Type:        "int32",
								FieldName:   "bar",
								FieldNumber: "126",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 36,
										Line:   4,
										Column: 5,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 55,
								Line:   5,
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
						Offset: 57,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an excerpt from the official reference(proto2)",
			input: `
message Outer {
  option (my_option).a = true;
  message Inner {   // Level 2
    required int64 ival = 1;
  }
  map<int32, string> my_map = 2;
  extensions 20 to 30;
}`,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 46,
								Line:   3,
								Column: 30,
							},
						},
					},
					&parser.Message{
						MessageName: "Inner",
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRequired:  true,
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 83,
										Line:   5,
										Column: 5,
									},
								},
							},
						},
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// Level 2",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 68,
									Line:   4,
									Column: 21,
								},
								LastPos: meta.Position{
									Offset: 77,
									Line:   4,
									Column: 30,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 50,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 110,
								Line:   6,
								Column: 3,
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "2",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 114,
								Line:   7,
								Column: 3,
							},
						},
					},
					&parser.Extensions{
						Ranges: []*parser.Range{
							{
								Begin: "20",
								End:   "30",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 147,
								Line:   8,
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
						Offset: 168,
						Line:   9,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `
message Outer {
  message Inner {};
};
`,
			permissive: true,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Message{
						MessageName: "Inner",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 35,
								Line:   3,
								Column: 19,
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
						Offset: 38,
						Line:   4,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `
message Outer {
  message Inner {}
}
`,
			permissive: true,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Message{
						MessageName: "Inner",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 19,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 34,
								Line:   3,
								Column: 18,
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
						Offset: 36,
						Line:   4,
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
			got, err := p.ParseMessage()
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

			if !reflect.DeepEqual(got, test.wantMessage) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantMessage))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
