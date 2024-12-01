package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestParser_ParseEdition(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantEdition *parser.Edition
		wantErr     bool
	}{
		{
			name: "parsing an empty",
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `edition = "2023";`,
			wantEdition: &parser.Edition{
				Edition:      "2023",
				EditionQuote: `"2023"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 16,
						Line:   1,
						Column: 17,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got, err := p.ParseEdition()
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

			if !reflect.DeepEqual(got, test.wantEdition) {
				t.Errorf("got %v, but want %v", got, test.wantEdition)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
