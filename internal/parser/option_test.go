package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/parser"
)

func TestParser_ParseOption(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOption *parser.Option
		wantErr    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without option",
			input:   `java_package = "com.example.foo";`,
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without =",
			input:   `option java_package "com.example.foo";`,
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without ;",
			input:   `option java_package = "com.example.foo"`,
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `option java_package = "com.example.foo";`,
			wantOption: &parser.Option{
				OptionName: "java_package",
				Constant:   `"com.example.foo"`,
			},
		},
		{
			name:  "parsing fullIdent",
			input: `option java_package.baz.bar = "com.example.foo";`,
			wantOption: &parser.Option{
				OptionName: "java_package.baz.bar",
				Constant:   `"com.example.foo"`,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := parser.ParseOption()
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

			if !reflect.DeepEqual(got, test.wantOption) {
				t.Errorf("got %v, but want %v", got, test.wantOption)
			}

			if !parser.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
