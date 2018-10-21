package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
)

func TestParser_ParseSyntax(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantSyntax *parser.Syntax
		wantErr    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `syntax = "proto3";`,
			wantSyntax: &parser.Syntax{
				ProtobufVersion: "proto3",
			},
		},
		{
			name: "parsing a comment",
			input: `
// use proto3
syntax = "proto3";
`,
			wantSyntax: &parser.Syntax{
				ProtobufVersion: "proto3",
				Comments: []*parser.Comment{
					{
						Raw: `// use proto3`,
					},
				},
			},
		},
		{
			name: "parsing a comment",
			input: `
/*
use proto3
*/
syntax = "proto3";
`,
			wantSyntax: &parser.Syntax{
				ProtobufVersion: "proto3",
				Comments: []*parser.Comment{
					{
						Raw: `/*
use proto3
*/`,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got, err := p.ParseSyntax()
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

			if !reflect.DeepEqual(got, test.wantSyntax) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantSyntax))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
