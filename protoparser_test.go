package protoparser_test

import (
	"os"
	"testing"

	"reflect"

	protoparser "github.com/yoheimuta/go-protoparser"
	"github.com/yoheimuta/go-protoparser/internal/config_test"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		wantProtocolBuffer *protoparser.ProtocolBuffer
	}{
		{
			name:          "parse a whole proto file",
			inputFilename: config_test.TestDataPath("parser.proto"),
			wantProtocolBuffer: &protoparser.ProtocolBuffer{
				Package: "parserpb",
				Service: &protoparser.Service{
					Comments: []string{
						"// ItemService is a service to manage items.",
					},
					Name: "ItemService",
					RPCs: []*protoparser.RPC{
						{
							Comments: []string{
								"// CreateUserItem is a method to create a user's item.",
							},
							Name: "CreateUserItem",
							Argument: &protoparser.Type{
								Name:       "CreateUserItemRequest",
								IsRepeated: false,
							},
							Return: &protoparser.Type{
								Name:       "aggregatespb.UserItemAggregate",
								IsRepeated: false,
							},
						},
						{
							Comments: []string{
								"// UpdateUserItem is a method to update a user's item.",
							},
							Name: "UpdateUserItem",
							Argument: &protoparser.Type{
								Name:       "UpdateUserItemRequest",
								IsRepeated: false,
							},
							Return: &protoparser.Type{
								Name:       "entitiespb.UserItem",
								IsRepeated: false,
							},
						},
					},
				},
				Messages: []*protoparser.Message{
					{
						Comments: []string{
							"// CreateUserItemRequest is a request message for CreateUserItem.",
						},
						Name: "CreateUserItemRequest",
						Fields: []*protoparser.Field{
							{
								Comments: []string{
									"// item is an item entity. Required.",
								},
								Type: &protoparser.Type{
									Name:       "entitiespb.UserItem",
									IsRepeated: false,
								},
								Name: "item",
							},
							{
								Comments: []string{
									"// images are item's images. Max count is 10. Optional.",
								},
								Type: &protoparser.Type{
									Name:       "Image",
									IsRepeated: true,
								},
								Name: "images",
							},
							{
								Comments: []string{
									"// mapping is a item's mapping information. Required.",
								},
								Type: &protoparser.Type{
									Name:       "Mapping",
									IsRepeated: false,
								},
								Name: "mapping",
							},
						},
						Messages: []*protoparser.Message{
							{
								Comments: []string{
									"// Image is an item's image information for create",
								},
								Name: "Image",
								Fields: []*protoparser.Field{
									{
										Comments: []string{
											"// display_order is an order of position. Starts 1 at left and increment by one. Required.",
										},
										Type: &protoparser.Type{
											Name:       "int64",
											IsRepeated: false,
										},
										Name: "display_order",
									},
									{
										Comments: []string{
											"// binary is an image binary. Required.",
										},
										Type: &protoparser.Type{
											Name:       "bytes",
											IsRepeated: false,
										},
										Name: "binary",
									},
								},
								Messages: nil,
								Enums:    nil,
							},
							{
								Comments: []string{
									"// Mapping is",
									"// an information of an item mapping.",
								},
								Name: "Mapping",
								Fields: []*protoparser.Field{
									{
										Comments: []string{
											"// product is an item master information.",
										},
										Type: &protoparser.Type{
											Name:       "entitiespb.UserItemMappingProduct",
											IsRepeated: false,
										},
										Name: "product",
									},
								},
								Messages: nil,
								Enums:    nil,
							},
						},
						Enums: []*protoparser.Enum{},
						Oneofs: []*protoparser.Oneof{
							{
								Comments: []string{
									"// condition_oneof is an item's condition. Required.",
								},
								Name: "condition_oneof",
								Fields: []*protoparser.Field{
									{
										Comments: []string{
											"// content_type_id is a condition ID of an item with content.",
										},
										Type: &protoparser.Type{
											Name:       "itemContentConditionpb.Type",
											IsRepeated: false,
										},
										Name: "content_type_id",
									},
									{
										Comments: []string{
											"// no_content_type_id is a condition ID of an item without content.",
										},
										Type: &protoparser.Type{
											Name:       "itemNoContentConditionpb.Type",
											IsRepeated: false,
										},
										Name: "no_content_type_id",
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
			proto, err := os.Open(test.inputFilename)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			defer func() {
				err = proto.Close()
				if err != nil {
					t.Errorf("failed to close a proto, err=%v", err)
					return
				}
			}()
			got, err := protoparser.Parse(proto)
			if err != nil {
				t.Errorf("failed to parse, err=%v", err)
				return
			}

			want := test.wantProtocolBuffer
			if got.Package != want.Package {
				t.Errorf("got %s, but want %s", got.Package, want.Package)
			}

			if len(got.Service.Comments) != len(want.Service.Comments) {
				t.Errorf("got %d, but want %d", len(got.Service.Comments), len(want.Service.Comments))
			}
			if !reflect.DeepEqual(got.Service.Comments, want.Service.Comments) {
				t.Errorf("got %v, but want %v", got.Service.Comments, want.Service.Comments)
			}

			if got.Service.Name != want.Service.Name {
				t.Errorf("got %s, but want %s", got.Service.Name, want.Service.Name)
			}

			if len(got.Service.RPCs) != len(want.Service.RPCs) {
				t.Errorf("got %d, but want %d", len(got.Service.RPCs), len(want.Service.RPCs))
			}
			for j, rpc := range want.Service.RPCs {
				if !reflect.DeepEqual(got.Service.RPCs[j], rpc) {
					t.Errorf("got %v, but want %v", got.Service.RPCs[j], rpc)
				}
			}

			if len(got.Messages) != len(want.Messages) {
				t.Errorf("got %d, but want %d", len(got.Messages), len(want.Messages))
			}
			for j, message := range want.Messages {
				gotMessage := got.Messages[j]
				if len(gotMessage.Comments) != len(message.Comments) {
					t.Errorf("got %d, but want %d", len(gotMessage.Comments), len(message.Comments))
				}
				if !reflect.DeepEqual(gotMessage.Comments, message.Comments) {
					t.Errorf("got %v, but want %v", gotMessage.Comments, message.Comments)
				}

				if gotMessage.Name != message.Name {
					t.Errorf("got %s, but want %s", gotMessage.Name, message.Name)
				}

				if len(gotMessage.Fields) != len(message.Fields) {
					t.Errorf("got %d, but want %d", len(gotMessage.Fields), len(message.Fields))
				}
				for k, field := range message.Fields {
					gotField := gotMessage.Fields[k]
					if !reflect.DeepEqual(gotField, field) {
						t.Errorf("got %v, but want %v", gotField, field)
					}
				}

				if len(gotMessage.Messages) != len(message.Messages) {
					t.Errorf("got %d, but want %d", len(gotMessage.Messages), len(message.Messages))
				}
				for k, nest := range message.Messages {
					gotNest := gotMessage.Messages[k]
					if !reflect.DeepEqual(gotNest, nest) {
						t.Errorf("got %v, but want %v", gotNest, nest)
					}
				}

				if len(gotMessage.Enums) != len(message.Enums) {
					t.Errorf("got %d, but want %d", len(gotMessage.Enums), len(message.Enums))
				}
				for k, enum := range message.Enums {
					gotEnum := gotMessage.Enums[k]
					if !reflect.DeepEqual(gotEnum, enum) {
						t.Errorf("got %v, but want %v", gotEnum, enum)
					}
				}

				if len(gotMessage.Oneofs) != len(message.Oneofs) {
					t.Errorf("got %d, but want %d", len(gotMessage.Oneofs), len(message.Oneofs))
				}
				for k, oneof := range message.Oneofs {
					gotOneof := gotMessage.Oneofs[k]
					if !reflect.DeepEqual(gotOneof, oneof) {
						t.Errorf("got %v, but want %v", gotOneof, oneof)
					}
				}
			}
		})
	}
}
