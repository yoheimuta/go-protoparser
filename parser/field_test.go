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
		{
			name:       "parsing fieldOption constant with { and a trailing comma by permissive mode",
			input:      "int64 display_order = 1 [(validator.field) = {int_gt: 0,}];",
			permissive: true,
			wantField: &parser.Field{
				Type:        "int64",
				FieldName:   "display_order",
				FieldNumber: "1",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "(validator.field)",
						Constant:   "{int_gt:0,}",
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
			name:       "parsing fieldOption constant with { and , by permissive mode. Required by go-proto-validators",
			input:      `string email = 2 [(validator.field) = {length_gt: 0, length_lt: 1025},(validator.field) = {regex: "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"}];`,
			permissive: true,
			wantField: &parser.Field{
				Type:        "string",
				FieldName:   "email",
				FieldNumber: "2",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "(validator.field)",
						Constant:   "{length_gt:0,length_lt:1025}",
					},
					{
						OptionName: "(validator.field)",
						Constant:   `{regex:"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"}`,
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
			name:  "parsing an excerpt from the official reference(proto2)",
			input: `optional foo.bar nested_message = 2;`,
			wantField: &parser.Field{
				IsOptional:  true,
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
			name:  "parsing a required label(proto2)",
			input: `required int32 samples = 4 [packed=true];`,
			wantField: &parser.Field{
				IsRequired:  true,
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
			name: "parsing fieldOption constant meaning a swagger annotation. Fix #52",
			input: `string email_id = 1[(grpc.gateway.protoc_gen_swagger.options.openapiv2_field) = {
	pattern: "^-!#$%&'*+\/0-9=?A-Z^_a-z{|}~@a-zA-Z0-9\.a-zA-Z+$"
	max_length: 254
	min_length: 1
	description: "Enter user email"
}];
`,
			permissive: true,
			wantField: &parser.Field{
				Type:        "string",
				FieldName:   "email_id",
				FieldNumber: "1",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "(grpc.gateway.protoc_gen_swagger.options.openapiv2_field)",
						Constant: `{pattern:"^-!#$%&'*+\/0-9=?A-Z^_a-z{|}~@a-zA-Z0-9\.a-zA-Z+$"
max_length:254
min_length:1
description:"Enter user email"}`,
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
			name:       "parsing fieldOption constant contained in a_bit_of_everything.proto provided by grpc-gateway. Fix #52",
			input:      `float float_value = 3 [(grpc.gateway.protoc_gen_swagger.options.openapiv2_field) = {description: "Float value field", default: "0.2", required: ['float_value']}];`,
			permissive: true,
			wantField: &parser.Field{
				Type:        "float",
				FieldName:   "float_value",
				FieldNumber: "3",
				FieldOptions: []*parser.FieldOption{
					{
						OptionName: "(grpc.gateway.protoc_gen_swagger.options.openapiv2_field)",
						Constant:   `{description:"Float value field",default:"0.2",required:['float_value']}`,
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
