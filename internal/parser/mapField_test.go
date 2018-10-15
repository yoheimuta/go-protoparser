package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/parser"
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
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := parser.ParseMapField()
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

			if !parser.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
