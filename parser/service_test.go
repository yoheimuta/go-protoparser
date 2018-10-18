package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/parser"
)

func TestParser_ParseService(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantService *parser.Service
		wantErr     bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []interface{}{
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
			},
		},
		{
			name: "parsing a rpc option",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse) { option (my_option).a = true; }
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []interface{}{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
						},
						Options: []*parser.Option{
							{
								OptionName: "(my_option).a",
								Constant:   "true",
							},
						},
					},
				},
			},
		},
		{
			name: "parsing rpcs",
			input: `
// ItemService is a service to manage items.
service ItemService {
    // CreateUserItem is a method to create a user's item.
    rpc CreateUserItem(CreateUserItemRequest) returns (aggregatespb.UserItemAggregate) {}

    // UpdateUserItem is a method to update a user's item.
    rpc UpdateUserItem(UpdateUserItemRequest) returns (entitiespb.UserItem) {}
}
`,
			wantService: &parser.Service{
				ServiceName: "ItemService",
				ServiceBody: []interface{}{
					&parser.RPC{
						RPCName: "CreateUserItem",
						RPCRequest: &parser.RPCRequest{
							MessageType: "CreateUserItemRequest",
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "aggregatespb.UserItemAggregate",
						},
					},
					&parser.RPC{
						RPCName: "UpdateUserItem",
						RPCRequest: &parser.RPCRequest{
							MessageType: "UpdateUserItemRequest",
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "entitiespb.UserItem",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(lexer.NewLexer2(strings.NewReader(test.input)))
			got, err := p.ParseService()
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

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
