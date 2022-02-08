package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestParser_ParseImport(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantImport *parser.Import
		wantErr    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing the invalid statement without import",
			input:   `"other.proto";`,
			wantErr: true,
		},
		{
			name:    "parsing the invalid statement without strLit",
			input:   `import 'other.proto";`,
			wantErr: true,
		},
		{
			name:  "parsing the statement without a modifier",
			input: `import "google/protobuf/timestamp.proto";`,
			wantImport: &parser.Import{
				Modifier: parser.ImportModifierNone,
				Location: `"google/protobuf/timestamp.proto"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 40,
						Line:   1,
						Column: 41,
					},
				},
			},
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: `import public "other.proto";`,
			wantImport: &parser.Import{
				Modifier: parser.ImportModifierPublic,
				Location: `"other.proto"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 27,
						Line:   1,
						Column: 28,
					},
				},
			},
		},
		{
			name:  "parsing the statement with weak",
			input: `import weak "other.proto";`,
			wantImport: &parser.Import{
				Modifier: parser.ImportModifierWeak,
				Location: `"other.proto"`,
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 25,
						Line:   1,
						Column: 26,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got, err := p.ParseImport()
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

			if !reflect.DeepEqual(got, test.wantImport) {
				t.Errorf("got %v, but want %v", got, test.wantImport)
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
