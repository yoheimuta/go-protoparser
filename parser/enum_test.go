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

func TestParser_ParseEnum(t *testing.T) {
	tests := []struct {
		name                       string
		input                      string
		inputBodyIncludingComments bool
		permissive                 bool
		wantEnum                   *parser.Enum
		wantErr                    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an invalid option",
			input: `enum EnumAllowingAlias {
  allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
`,
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 52,
								Line:   2,
								Column: 28,
							},
						},
					},
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 56,
								Line:   3,
								Column: 3,
							},
						},
					},
					&parser.EnumField{
						Ident:  "STARTED",
						Number: "1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 71,
								Line:   4,
								Column: 3,
							},
						},
					},
					&parser.EnumField{
						Ident:  "RUNNING",
						Number: "2",
						EnumValueOptions: []*parser.EnumValueOption{
							{
								OptionName: "(custom_option)",
								Constant:   `"hello world"`,
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 86,
								Line:   5,
								Column: 3,
							},
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
						Offset: 133,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing enumValueOptions",
			input: `enum EnumAllowingAlias {
  RUNNING = 0 [(custom_option) = "hello world", (custom_option2) = "hello world2"];
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.EnumField{
						Ident:  "RUNNING",
						Number: "0",
						EnumValueOptions: []*parser.EnumValueOption{
							{
								OptionName: "(custom_option)",
								Constant:   `"hello world"`,
							},
							{
								OptionName: "(custom_option2)",
								Constant:   `"hello world2"`,
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
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
						Offset: 109,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing comments",
			input: `enum EnumAllowingAlias {
  // option
  option allow_alias = true;
  // UNKNOWN
  UNKNOWN = 0;
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						Comments: []*parser.Comment{
							{
								Raw: `// option`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 27,
										Line:   2,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 35,
										Line:   2,
										Column: 11,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 39,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 64,
								Line:   3,
								Column: 28,
							},
						},
					},
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
						Comments: []*parser.Comment{
							{
								Raw: `// UNKNOWN`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 68,
										Line:   4,
										Column: 3,
									},
									LastPos: meta.Position{
										Offset: 77,
										Line:   4,
										Column: 12,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 81,
								Line:   5,
								Column: 3,
							},
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
						Offset: 94,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing inline comments",
			input: `enum EnumAllowingAlias { // TODO: implementation
  option allow_alias = true; // option
  UNKNOWN = 0; // UNKNOWN
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						InlineComment: &parser.Comment{
							Raw: `// option`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 78,
									Line:   2,
									Column: 30,
								},
								LastPos: meta.Position{
									Offset: 86,
									Line:   2,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 51,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 76,
								Line:   2,
								Column: 28,
							},
						},
					},
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
						InlineComment: &parser.Comment{
							Raw: `// UNKNOWN`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 103,
									Line:   3,
									Column: 16,
								},
								LastPos: meta.Position{
									Offset: 112,
									Line:   3,
									Column: 25,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 90,
								Line:   3,
								Column: 3,
							},
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
						LastPos: meta.Position{
							Offset: 47,
							Line:   1,
							Column: 48,
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
						Offset: 114,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "skipping a last comment",
			input: `enum EnumAllowingAlias {
  option allow_alias = true;
  // last line
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 52,
								Line:   2,
								Column: 28,
							},
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
						Offset: 69,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing last comments",
			input: `enum EnumAllowingAlias {
  option allow_alias = true;
  // last first comment
  /* last second comment */
}
`,
			inputBodyIncludingComments: true,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 52,
								Line:   2,
								Column: 28,
							},
						},
					},
					&parser.Comment{
						Raw: `// last first comment`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 56,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 76,
								Line:   3,
								Column: 23,
							},
						},
					},
					&parser.Comment{
						Raw: `/* last second comment */`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 80,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 104,
								Line:   4,
								Column: 27,
							},
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
						Offset: 106,
						Line:   5,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an excerpt containing reserved, from the official reference",
			input: `enum Foo {
  reserved 2, 15, 9 to 11, 40 to max;
  reserved "FOO", "BAR";
}
`,
			wantEnum: &parser.Enum{
				EnumName: "Foo",
				EnumBody: []parser.Visitee{
					&parser.Reserved{
						Ranges: []*parser.Range{
							{
								Begin: "2",
							},
							{
								Begin: "15",
							},
							{
								Begin: "9",
								End:   "11",
							},
							{
								Begin: "40",
								End:   "max",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 13,
								Line:   2,
								Column: 3,
							},
						},
					},
					&parser.Reserved{
						FieldNames: []string{
							`"FOO"`,
							`"BAR"`,
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 51,
								Line:   3,
								Column: 3,
							},
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
						Offset: 74,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `enum EnumAllowingAlias {
};
`,
			permissive: true,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 26,
						Line:   2,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `enum EnumAllowingAlias {
}
`,
			permissive: true,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 25,
						Line:   2,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing an enum with multiple string literals. See #35",
			input: `enum EnumAllowingAlias {
  UNKNOWN = 0 [(custom_option) = "this is a "
                                 "string on two lines"
              ];
}
`,
			permissive: true,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
						EnumValueOptions: []*parser.EnumValueOption{
							{
								OptionName: "(custom_option)",
								Constant:   `"this is a string on two lines"`,
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
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
						Offset: 143,
						Line:   5,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a negative enum field",
			input: `enum TestNegativeValue {
  NEGATIVE_CONSTANT = -1;
  ZERO_CONSTANT = 0;
  POSITIVE_CONSTANT = 1;
}
`,
			wantEnum: &parser.Enum{
				EnumName: "TestNegativeValue",
				EnumBody: []parser.Visitee{
					&parser.EnumField{
						Ident:  "NEGATIVE_CONSTANT",
						Number: "-1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
							},
						},
					},
					&parser.EnumField{
						Ident:  "ZERO_CONSTANT",
						Number: "0",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 53,
								Line:   3,
								Column: 3,
							},
						},
					},
					&parser.EnumField{
						Ident:  "POSITIVE_CONSTANT",
						Number: "1",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 74,
								Line:   4,
								Column: 3,
							},
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
						Offset: 97,
						Line:   5,
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
			got, err := p.ParseEnum()
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

			if !reflect.DeepEqual(got, test.wantEnum) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantEnum))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
