package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

func TestParser_ParseField(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		permissive bool
		wantField  *parser.Field
		wantErr    bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; without fieldNumber",
			input:   "foo.bar nested_message = ;",
			wantErr: true,
		},
		{
			name:    "parsing an invalid; string fieldNumber",
			input:   "foo.bar nested_message = a;",
			wantErr: true,
		},
		{
			name:  "parsing an excerpt from the official reference",
			input: "foo.bar nested_message = 2;",
			wantField: &parser.Field{
				Type:        "foo.bar",
				FieldName:   "nested_message",
				FieldNumber: "2",
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
			input: "repeated int32 samples = 4 [packed=true];",
			wantField: &parser.Field{
				IsRepeated:  true,
				Type:        "int32",
				FieldName:   "samples",
				FieldNumber: "4",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "packed",
						Constant:   "true",
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
			name:  "parsing fieldOptions",
			input: "repeated int32 samples = 4 [packed=true, required=false];",
			wantField: &parser.Field{
				IsRepeated:  true,
				Type:        "int32",
				FieldName:   "samples",
				FieldNumber: "4",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "packed",
						Constant:   "true",
					},
					{
						OptionName: "required",
						Constant:   "false",
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
			name:    "parsing an invalid fieldOption constant",
			input:   "int64 display_order = 1 [(validator.field) = {int_gt: 0}];",
			wantErr: true,
		},
		{
			name:       "parsing fieldOption constant with { by permissive mode. Required by go-proto-validators",
			input:      "int64 display_order = 1 [(validator.field) = {int_gt: 0}];",
			permissive: true,
			wantField: &parser.Field{
				Type:        "int64",
				FieldName:   "display_order",
				FieldNumber: "1",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "(validator.field)",
						Constant:   "{int_gt:0}",
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
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)), parser.WithPermissive(test.permissive))
			got, err := p.ParseField()
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

			if !reflect.DeepEqual(got, test.wantField) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantField))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
