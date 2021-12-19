package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestParser_ParseMapField(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantMapField *parser.MapField
		wantErr      bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without map",
			input:   "<string, Project> projects = 3;",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; not keyType constant",
			input:   "map<customType, Project> projects = 3;",
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: "map<string, Project> projects = 3;",
			wantMapField: &parser.MapField{
				KeyType:     "string",
				Type:        "Project",
				MapName:     "projects",
				FieldNumber: "3",
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
			got, err := p.ParseMapField()
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

			if !reflect.DeepEqual(got, test.wantMapField) {
				t.Errorf("got %v, but want %v", got, test.wantMapField)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
