package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestParser_ParseExtensions(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantExtensions *parser.Extensions
		wantErr        bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without to",
			input:   "extensions 2, 15, 9 11;",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; including both ranges and fieldNames",
			input:   `extensions 2, "foo", 9 to 11;`,
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `extensions 100 to 199;`,
			wantExtensions: &parser.Extensions{
				Ranges: []*parser.Range{
					{
						Begin: "100",
						End:   "199",
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
			name:  "parsing another excerpt from the official reference",
			input: `extensions 4, 20 to max;`,
			wantExtensions: &parser.Extensions{
				Ranges: []*parser.Range{
					{
						Begin: "4",
					},
					{
						Begin: "20",
						End:   "max",
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
			got, err := p.ParseExtensions()
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

			if !reflect.DeepEqual(got, test.wantExtensions) {
				t.Errorf("got %v, but want %v", got, test.wantExtensions)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
