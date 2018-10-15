package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/parser"
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
					},
					{
						Type:        "SubMessage",
						FieldName:   "sub_message",
						FieldNumber: "9",
					},
				},
				OneofName: "foo",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := parser.ParseOneof()
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
				t.Errorf("got %v, but want %v", got, test.wantOneof)
			}

			if !parser.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
