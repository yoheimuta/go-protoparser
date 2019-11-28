package unordered_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

func TestInterpretProto(t *testing.T) {
	tests := []struct {
		name       string
		inputProto *parser.Proto
		wantProto  *unordered.Proto
		wantErr    bool
	}{
		{
			name: "interpreting a nil",
		},
		{
			name: "interpreting an excerpt from the official reference",
			inputProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
				},
				ProtoBody: []parser.Visitee{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
							},
							&parser.EnumField{
								Ident:  "UNKNOWN",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "STARTED",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "RUNNING",
								Number: "2",
								EnumValueOptions: []*parser.EnumValueOption{
									{
										OptionName: "(custom_option)",
										Constant:   `"hello world"`,
									},
								},
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						MessageBody: []parser.Visitee{
							&parser.Option{
								OptionName: "(my_option).a",
								Constant:   "true",
							},
							&parser.Message{
								MessageName: "inner",
								MessageBody: []parser.Visitee{
									&parser.Field{
										Type:        "int64",
										FieldName:   "ival",
										FieldNumber: "1",
									},
								},
							},
							&parser.Field{
								IsRepeated:  true,
								Type:        "inner",
								FieldName:   "inner_message",
								FieldNumber: "2",
							},
							&parser.Field{
								Type:        "EnumAllowingAlias",
								FieldName:   "enum_field",
								FieldNumber: "3",
							},
							&parser.MapField{
								KeyType:     "int32",
								Type:        "string",
								MapName:     "my_map",
								FieldNumber: "4",
							},
						},
					},
				},
			},
			wantProto: &unordered.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
				},
				ProtoBody: &unordered.ProtoBody{
					Imports: []*parser.Import{
						{
							Modifier: parser.ImportModifierPublic,
							Location: `"other.proto"`,
						},
					},
					Options: []*parser.Option{
						{
							OptionName: "java_package",
							Constant:   `"com.example.foo"`,
						},
					},
					Enums: []*unordered.Enum{
						{

							EnumName: "EnumAllowingAlias",
							EnumBody: &unordered.EnumBody{
								Options: []*parser.Option{
									{
										OptionName: "allow_alias",
										Constant:   "true",
									},
								},
								EnumFields: []*parser.EnumField{
									{
										Ident:  "UNKNOWN",
										Number: "0",
									},
									{
										Ident:  "STARTED",
										Number: "1",
									},
									{
										Ident:  "RUNNING",
										Number: "2",
										EnumValueOptions: []*parser.EnumValueOption{
											{
												OptionName: "(custom_option)",
												Constant:   `"hello world"`,
											},
										},
									},
								},
							},
						},
					},
					Messages: []*unordered.Message{
						{
							MessageName: "outer",
							MessageBody: &unordered.MessageBody{
								Options: []*parser.Option{
									{
										OptionName: "(my_option).a",
										Constant:   "true",
									},
								},
								Messages: []*unordered.Message{
									{
										MessageName: "inner",
										MessageBody: &unordered.MessageBody{
											Fields: []*parser.Field{
												{
													Type:        "int64",
													FieldName:   "ival",
													FieldNumber: "1",
												},
											},
										},
									},
								},
								Fields: []*parser.Field{
									{
										IsRepeated:  true,
										Type:        "inner",
										FieldName:   "inner_message",
										FieldNumber: "2",
									},
									{
										Type:        "EnumAllowingAlias",
										FieldName:   "enum_field",
										FieldNumber: "3",
									},
								},
								Maps: []*parser.MapField{
									{
										KeyType:     "int32",
										Type:        "string",
										MapName:     "my_map",
										FieldNumber: "4",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := unordered.InterpretProto(test.inputProto)
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err, parsed=%v", got)
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantProto) {
				t.Errorf("got %v, but want %v", got, test.wantProto)
			}
		})
	}

}
