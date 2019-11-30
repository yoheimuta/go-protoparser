package unordered_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestInterpretMessage(t *testing.T) {
	tests := []struct {
		name         string
		inputMessage *parser.Message
		wantMessage  *unordered.Message
		wantErr      bool
	}{
		{
			name: "interpreting a nil",
		},
		{
			name: "interpreting an excerpt from the official reference with comments",
			inputMessage: &parser.Message{
				MessageName: "Outer",
				MessageBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(my_option).a",
						Constant:   "true",
					},
					&parser.Message{
						MessageName: "Inner",
						MessageBody: []parser.Visitee{
							&parser.Field{
								Type:        "int64",
								FieldName:   "ival",
								FieldNumber: "1",
							},
						},
					},
					&parser.MapField{
						KeyType:     "int32",
						Type:        "string",
						MapName:     "my_map",
						FieldNumber: "2",
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// message",
					},
				},
				InlineComment: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   2,
							Column: 26,
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation2",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 21,
						Line:   3,
						Column: 1,
					},
				},
			},
			wantMessage: &unordered.Message{
				MessageName: "Outer",
				MessageBody: &unordered.MessageBody{
					Options: []*parser.Option{
						{
							OptionName: "(my_option).a",
							Constant:   "true",
						},
					},
					Messages: []*unordered.Message{
						{
							MessageName: "Inner",
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
					Maps: []*parser.MapField{
						{
							KeyType:     "int32",
							Type:        "string",
							MapName:     "my_map",
							FieldNumber: "2",
						},
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// message",
					},
				},
				InlineComment: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   2,
							Column: 26,
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation2",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 21,
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
			got, err := unordered.InterpretMessage(test.inputMessage)
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

			if !reflect.DeepEqual(got, test.wantMessage) {
				t.Errorf("got %v, but want %v", got, test.wantMessage)
			}
		})
	}

}
