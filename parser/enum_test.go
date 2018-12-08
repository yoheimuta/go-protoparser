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

func TestParser_ParseEnum(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantEnum *parser.Enum
		wantErr  bool
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
				EnumBody: []interface{}{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 3,
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
				EnumBody: []interface{}{
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
				EnumBody: []interface{}{
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
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 39,
								Line:   3,
								Column: 3,
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
				},
			},
		},
		{
			name: "parsing inline comments",
			input: `enum EnumAllowingAlias {
  option allow_alias = true; // option
  UNKNOWN = 0; // UNKNOWN
}
`,
			wantEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []interface{}{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
						InlineComment: &parser.Comment{
							Raw: `// option`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 54,
									Line:   2,
									Column: 30,
								},
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
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
						InlineComment: &parser.Comment{
							Raw: `// UNKNOWN`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 79,
									Line:   3,
									Column: 16,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 66,
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
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
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
