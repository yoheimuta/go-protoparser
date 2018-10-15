package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/parser"
)

func TestParser_ParseReserved(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantReserved *parser.Reserved
		wantErr      bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without to",
			input:   "reserved 2, 15, 9 11;",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; including both ranges and fieldNames",
			input:   `reserved 2, "foo", 9 to 11;`,
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: "reserved 2, 15, 9 to 11;",
			wantReserved: &parser.Reserved{
				Ranges: []*parser.Range{
					{
						Begin: "2",
					},
					{
						Begin: "15",
					},
					{
						Begin: "9",
						End:   "11",
					},
				},
			},
		},
		{
			name:  "parsing another excerpt from the official reference",
			input: `reserved "foo", "bar";`,
			wantReserved: &parser.Reserved{
				FieldNames: []string{
					`"foo"`,
					`"bar"`,
				},
			},
		},
		{
			name:  "parsing an input with max",
			input: "reserved 9 to max;",
			wantReserved: &parser.Reserved{
				Ranges: []*parser.Range{
					{
						Begin: "9",
						End:   "max",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := parser.ParseReserved()
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

			if !reflect.DeepEqual(got, test.wantReserved) {
				t.Errorf("got %v, but want %v", got, test.wantReserved)
			}

			if !parser.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
