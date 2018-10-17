package parser_test

import (
	"strings"
	"testing"

	"reflect"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/parser"
)

func TestParser_ParsePackage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPackage *parser.Package
		wantErr     bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `package foo.bar;`,
			wantPackage: &parser.Package{
				Name: "foo.bar",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := p.ParsePackage()
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

			if !reflect.DeepEqual(got, test.wantPackage) {
				t.Errorf("got %v, but want %v", got, test.wantPackage)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
