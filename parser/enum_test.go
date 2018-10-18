package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/parser"
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
					},
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
					},
					&parser.EnumField{
						Ident:  "STARTED",
						Number: "1",
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
				t.Errorf("got %v, but want %v", got, test.wantEnum)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
