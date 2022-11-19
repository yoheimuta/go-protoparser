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

func TestParser_ParseService(t *testing.T) {
	tests := []struct {
		name                       string
		input                      string
		inputBodyIncludingComments bool
		permissive                 bool
		wantService                *parser.Service
		wantErr                    bool
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
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 78,
								Line:   3,
								Column: 54,
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
						Offset: 80,
						Line:   4,
						Column: 1,
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
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Options: []*parser.Option{
							{
								OptionName: "(my_option).a",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 81,
										Line:   3,
										Column: 57,
									},
									LastPos: meta.Position{
										Offset: 108,
										Line:   3,
										Column: 84,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 110,
								Line:   3,
								Column: 86,
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
						Offset: 112,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing multiple rpc options",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse) { 
	option (my_option).a = true; 
	option (my_option).b = false;
  }
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Options: []*parser.Option{
							{
								OptionName: "(my_option).a",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 83,
										Line:   4,
										Column: 2,
									},
									LastPos: meta.Position{
										Offset: 110,
										Line:   4,
										Column: 29,
									},
								},
							},
							{
								OptionName: "(my_option).b",
								Constant:   "false",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 114,
										Line:   5,
										Column: 2,
									},
									LastPos: meta.Position{
										Offset: 142,
										Line:   5,
										Column: 30,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 146,
								Line:   6,
								Column: 3,
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
						Offset: 148,
						Line:   7,
						Column: 1,
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
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "CreateUserItem",
						RPCRequest: &parser.RPCRequest{
							MessageType: "CreateUserItemRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 149,
									Line:   5,
									Column: 23,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "aggregatespb.UserItemAggregate",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 181,
									Line:   5,
									Column: 55,
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: "// CreateUserItem is a method to create a user's item.",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 72,
										Line:   4,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 125,
										Line:   4,
										Column: 58,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 131,
								Line:   5,
								Column: 5,
							},
							LastPos: meta.Position{
								Offset: 215,
								Line:   5,
								Column: 89,
							},
						},
					},
					&parser.RPC{
						RPCName: "UpdateUserItem",
						RPCRequest: &parser.RPCRequest{
							MessageType: "UpdateUserItemRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 299,
									Line:   8,
									Column: 23,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "entitiespb.UserItem",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 331,
									Line:   8,
									Column: 55,
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: "// UpdateUserItem is a method to update a user's item.",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 222,
										Line:   7,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 275,
										Line:   7,
										Column: 58,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 281,
								Line:   8,
								Column: 5,
							},
							LastPos: meta.Position{
								Offset: 354,
								Line:   8,
								Column: 78,
							},
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 46,
						Line:   3,
						Column: 1,
					},
					LastPos: meta.Position{
						Offset: 356,
						Line:   9,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a inline comment",
			input: `
service SearchService { // TODO: Search is not implemented yet.
  rpc Search (SearchRequest) returns (SearchResponse); // TODO: implementation
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 78,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 102,
									Line:   3,
									Column: 38,
								},
							},
						},
						InlineComment: &parser.Comment{
							Raw: `// TODO: implementation`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 120,
									Line:   3,
									Column: 56,
								},
								LastPos: meta.Position{
									Offset: 142,
									Line:   3,
									Column: 78,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 67,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 118,
								Line:   3,
								Column: 54,
							},
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: Search is not implemented yet.",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   2,
							Column: 25,
						},
						LastPos: meta.Position{
							Offset: 63,
							Line:   2,
							Column: 63,
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
						Offset: 144,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "skipping a last comment",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
  // last comment
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 78,
								Line:   3,
								Column: 54,
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
						Offset: 98,
						Line:   5,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing last comments",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
  // last first comment
  /* last second comment */
}
`,
			inputBodyIncludingComments: true,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 78,
								Line:   3,
								Column: 54,
							},
						},
					},
					&parser.Comment{
						Raw: `// last first comment`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 82,
								Line:   4,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 102,
								Line:   4,
								Column: 23,
							},
						},
					},
					&parser.Comment{
						Raw: `/* last second comment */`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 106,
								Line:   5,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 130,
								Line:   5,
								Column: 27,
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
						Offset: 132,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing a block followed by semicolon",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse) {};
};
`,
			permissive: true,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 81,
								Line:   3,
								Column: 57,
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
						Offset: 84,
						Line:   4,
						Column: 2,
					},
				},
			},
		},
		{
			name: "set LastPos to the correct position when a semicolon doesn't follow the last block",
			input: `
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse) {}
}
`,
			permissive: true,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "Search",
						RPCRequest: &parser.RPCRequest{
							MessageType: "SearchRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 38,
									Line:   3,
									Column: 14,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "SearchResponse",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 62,
									Line:   3,
									Column: 38,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 80,
								Line:   3,
								Column: 56,
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
						Offset: 82,
						Line:   4,
						Column: 1,
					},
				},
			},
		},
		{
			name: "parsing the rpc with a trailing comment followd by the left curly",
			input: `
service SearchService {
  rpc GetAll(GetRequest) returns(GetReply) { // get the global address table
    option(requestreply.Nats).Subject = "get.addrs";
  }
}
`,
			wantService: &parser.Service{
				ServiceName: "SearchService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "GetAll",
						RPCRequest: &parser.RPCRequest{
							MessageType: "GetRequest",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 37,
									Line:   3,
									Column: 13,
								},
							},
						},
						RPCResponse: &parser.RPCResponse{
							MessageType: "GetReply",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 57,
									Line:   3,
									Column: 33,
								},
							},
						},
						Options: []*parser.Option{
							{
								OptionName: "(requestreply.Nats).Subject",
								Constant:   `"get.addrs"`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 106,
										Line:   4,
										Column: 5,
									},
									LastPos: meta.Position{
										Offset: 153,
										Line:   4,
										Column: 52,
									},
								},
							},
						},
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// get the global address table",
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 70,
									Line:   3,
									Column: 46,
								},
								LastPos: meta.Position{
									Offset: 100,
									Line:   3,
									Column: 76,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 27,
								Line:   3,
								Column: 3,
							},
							LastPos: meta.Position{
								Offset: 157,
								Line:   5,
								Column: 3,
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
						Offset: 159,
						Line:   6,
						Column: 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(
				lexer.NewLexer(strings.NewReader(test.input)),
				parser.WithBodyIncludingComments(test.inputBodyIncludingComments),
				parser.WithPermissive(test.permissive),
			)
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
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantService))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}

}
