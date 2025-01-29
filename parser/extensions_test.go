package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/internal/util_test"
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
					LastPos: meta.Position{
						Offset: 21,
						Line:   1,
						Column: 22,
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
					LastPos: meta.Position{
						Offset: 23,
						Line:   1,
						Column: 24,
					},
				},
			},
		},
		{
			name: "parsing an excerpt with extension declarations from the official reference",
			input: `extensions 4 to 1000 [
    declaration = {
      number: 4,
      full_name: ".my.package.event_annotations",
      type: ".logs.proto.ValidationAnnotations",
      repeated: true },
    declaration = {
      number: 999,
      full_name: ".foo.package.bar",
      type: "int32"}];`,
			wantExtensions: &parser.Extensions{
				Ranges: []*parser.Range{
					{
						Begin: "4",
						End:   "1000",
					},
				},
				Declarations: []*parser.Declaration{
					{
						Number:   "4",
						FullName: `".my.package.event_annotations"`,
						Type:     `".logs.proto.ValidationAnnotations"`,
						Repeated: true,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   2,
								Column: 5,
							},
							LastPos: meta.Position{
								Offset: 180,
								Line:   6,
								Column: 22,
							},
						},
					},
					{
						Number:   "999",
						FullName: `".foo.package.bar"`,
						Type:     `"int32"`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 187,
								Line:   7,
								Column: 5,
							},
							LastPos: meta.Position{
								Offset: 278,
								Line:   10,
								Column: 20,
							},
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 280,
						Line:   10,
						Column: 22,
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
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantExtensions))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
