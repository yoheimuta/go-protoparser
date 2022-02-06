package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
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
				ProtobufVersion:      "proto3",
				ProtobufVersionQuote: `"proto3"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
		},
		{
			name:  "parsing a single-quote string",
			input: `syntax = 'proto3';`,
			wantSyntax: &parser.Syntax{
				ProtobufVersion:      "proto3",
				ProtobufVersionQuote: `'proto3'`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
		},
		{
			name:  "parsing an excerpt from the official reference(proto2)",
			input: `syntax = "proto2";`,
			wantSyntax: &parser.Syntax{
				ProtobufVersion:      "proto2",
				ProtobufVersionQuote: `"proto2"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 17,
						Line:   1,
						Column: 18,
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
				t.Errorf("got %v, but want %v", got, test.wantSyntax)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
