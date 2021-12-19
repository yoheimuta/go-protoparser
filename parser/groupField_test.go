package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/internal/util_test"
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

func TestParser_ParseGroupField(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		permissive     bool
		wantGroupField *parser.GroupField
		wantErr        bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an invalid: groupName is not capitalized.",
			input: `
repeated group result = 1 {
    required string url = 2;
    optional string title = 3;
    repeated string snippets = 4;
}
`,
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
repeated group Result = 1 {
    required string url = 2;
    optional string title = 3;
    repeated string snippets = 4;
}
`,
			wantGroupField: &parser.GroupField{
				IsRepeated:  true,
				GroupName:   "Result",
				FieldNumber: "1",
				MessageBody: []parser.Visitee{
					&parser.Field{
						IsRequired:  true,
						Type:        "string",
						FieldName:   "url",
						FieldNumber: "2",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 33,
								Line:   3,
								Column: 5,
							},
						},
					},
					&parser.Field{
						IsOptional:  true,
						Type:        "string",
						FieldName:   "title",
						FieldNumber: "3",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 62,
								Line:   4,
								Column: 5,
							},
						},
					},
					&parser.Field{
						IsRepeated:  true,
						Type:        "string",
						FieldName:   "snippets",
						FieldNumber: "4",
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 93,
								Line:   5,
								Column: 5,
							},
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 123,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `
group Result = 1 {
};
`,
			permissive: true,
			wantGroupField: &parser.GroupField{
				GroupName:   "Result",
				FieldNumber: "1",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 21,
						Line:   3,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `
group Result = 1 {
}
`,
			permissive: true,
			wantGroupField: &parser.GroupField{
				GroupName:   "Result",
				FieldNumber: "1",
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 20,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)), parser.WithPermissive(test.permissive))
			got, err := p.ParseGroupField()
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

			if !reflect.DeepEqual(got, test.wantGroupField) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantGroupField))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
