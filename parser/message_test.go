package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/parser"
)

func TestParser_ParseMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMessage *parser.Message
		wantErr     bool
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
  message Inner {   // Level 2
    int64 ival = 1;
  }
  map<int32, string> my_map = 2;
}
`,
			wantMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []interface{}{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
					},
					&parser.Message{
						MessageName: "Inner",
						MessageBody: []interface{}{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "2",
					},
				},
			},
		},
		{
			name: "parsing another excerpt from the official reference",
			input: `
message outer {
  option (my_option).a = true;
  message inner {   // Level 2
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
`,
			wantMessage: &parser.Message{
				MessageName: "outer",
				MessageBody: []interface{}{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
					},
					&parser.Message{
						MessageName: "Inner",
						MessageBody: []interface{}{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
							},
						},
					},
					&parser.Field{
						IsRepeated:   true,
						Type:         "",
						FieldName:    "",
						FieldNumber:  "",
						FieldOptions: nil,
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "4",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
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
				t.Errorf("got %v, but want %v", got, test.wantMessage)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
