package unordered_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

func TestInterpretService(t *testing.T) {
	tests := []struct {
		name         string
		inputService *parser.Service
		wantService  *unordered.Service
		wantErr      bool
	}{
		{
			name: "interpreting a nil",
		},
		{
			name: "interpreting an excerpt from the official reference with a option and comments",
			inputService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []interface{}{
					&parser.Option{
						OptionName: "case-sensitive",
						Constant:   "true",
					},
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
						},
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// service",
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
			wantService: &unordered.Service{
				ServiceName: "SearchService",
				ServiceBody: &unordered.ServiceBody{
					Options: []*parser.Option{
						{
							OptionName: "case-sensitive",
							Constant:   "true",
						},
					},
					RPCs: []*parser.RPC{
						{
							RPCName: "Search",
							RPCRequest: &parser.RPCRequest{
								MessageType: "SearchRequest",
							},
							RPCResponse: &parser.RPCResponse{
								MessageType: "SearchResponse",
							},
						},
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// service",
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
			got, err := unordered.InterpretService(test.inputService)
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

			if !reflect.DeepEqual(got, test.wantService) {
				t.Errorf("got %v, but want %v", got, test.wantService)
			}
		})
	}

}
