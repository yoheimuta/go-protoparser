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

func TestParser_ParseProto(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantProto *parser.Proto
		wantErr   bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
syntax = "proto3";
import public "other.proto";
option java_package = "com.example.foo";
enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
message outer {
  option (my_option).a = true;
  message inner {
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
`,
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 2,
							Line:   2,
							Column: 1,
						},
					},
				},
				ProtoBody: []interface{}{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 21,
								Line:   3,
								Column: 1,
							},
						},
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 50,
								Line:   4,
								Column: 1,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []interface{}{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 118,
										Line:   6,
										Column: 3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "UNKNOWN",
								Number: "0",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 147,
										Line:   7,
										Column: 3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "STARTED",
								Number: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 162,
										Line:   8,
										Column: 3,
									},
								},
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
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 177,
										Line:   9,
										Column: 3,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 91,
								Line:   5,
								Column: 1,
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						MessageBody: []interface{}{
							&parser.Option{
								OptionName: "(my_option).a",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 244,
										Line:   12,
										Column: 3,
									},
								},
							},
							&parser.Message{
								MessageName: "inner",
								MessageBody: []interface{}{
									&parser.Field{
										Type:        "int64",
										FieldName:   "ival",
										FieldNumber: "1",
										Meta: meta.Meta{
											Pos: meta.Position{
												Offset: 295,
												Line:   14,
												Column: 5,
											},
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 275,
										Line:   13,
										Column: 3,
									},
								},
							},
							&parser.Field{
								IsRepeated:  true,
								Type:        "inner",
								FieldName:   "inner_message",
								FieldNumber: "2",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 317,
										Line:   16,
										Column: 3,
									},
								},
							},
							&parser.Field{
								Type:        "EnumAllowingAlias",
								FieldName:   "enum_field",
								FieldNumber: "3",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 353,
										Line:   17,
										Column: 3,
									},
								},
							},
							&parser.MapField{
								KeyType:     "int32",
								Type:        "string",
								MapName:     "my_map",
								FieldNumber: "4",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 388,
										Line:   18,
										Column: 3,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 226,
								Line:   11,
								Column: 1,
							},
						},
					},
				},
			},
		},
		{
			name: "parsing a service",
			input: `
syntax = "proto3";
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
`,
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 2,
							Line:   2,
							Column: 1,
						},
					},
				},
				ProtoBody: []interface{}{
					&parser.Service{
						ServiceName: "SearchService",
						ServiceBody: []interface{}{
							&parser.RPC{
								RPCName: "Search",
								RPCRequest: &parser.RPCRequest{
									MessageType: "SearchRequest",
									Meta: meta.Meta{
										Pos: meta.Position{
											Offset: 58,
											Line:   4,
											Column: 14,
										},
									},
								},
								RPCResponse: &parser.RPCResponse{
									MessageType: "SearchResponse",
									Meta: meta.Meta{
										Pos: meta.Position{
											Offset: 82,
											Line:   4,
											Column: 38,
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 47,
										Line:   4,
										Column: 3,
									},
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
			},
		},
		{
			name: "parsing comments",
			input: `
// syntax
/*
syntax2
*/
syntax = "proto3";
// import
import public "other.proto";
/* package */
package foo.bar;
// option
option java_package = "com.example.foo";
// message
message outer {
}
// enum
enum EnumAllowingAlias {
  option allow_alias = true;
}
// service
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
`,
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Comments: []*parser.Comment{
						{
							Raw: `// syntax`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 2,
									Line:   2,
									Column: 1,
								},
							},
						},
						{
							Raw: `/*
syntax2
*/`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Offset: 12,
									Line:   3,
									Column: 1,
								},
							},
						},
					},
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 26,
							Line:   6,
							Column: 1,
						},
					},
				},
				ProtoBody: []interface{}{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
						Comments: []*parser.Comment{
							{
								Raw: `// import`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 45,
										Line:   7,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 55,
								Line:   8,
								Column: 1,
							},
						},
					},
					&parser.Package{
						Name: `foo.bar`,
						Comments: []*parser.Comment{
							{
								Raw: `/* package */`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 84,
										Line:   9,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 98,
								Line:   10,
								Column: 1,
							},
						},
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
						Comments: []*parser.Comment{
							{
								Raw: `// option`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 115,
										Line:   11,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 125,
								Line:   12,
								Column: 1,
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						Comments: []*parser.Comment{
							{
								Raw: `// message`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 166,
										Line:   13,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 177,
								Line:   14,
								Column: 1,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []interface{}{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 230,
										Line:   18,
										Column: 3,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// enum`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 195,
										Line:   16,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 203,
								Line:   17,
								Column: 1,
							},
						},
					},
					&parser.Service{
						ServiceName: "SearchService",
						ServiceBody: []interface{}{
							&parser.RPC{
								RPCName: "Search",
								RPCRequest: &parser.RPCRequest{
									MessageType: "SearchRequest",
									Meta: meta.Meta{
										Pos: meta.Position{
											Offset: 307,
											Line:   22,
											Column: 14,
										},
									},
								},
								RPCResponse: &parser.RPCResponse{
									MessageType: "SearchResponse",
									Meta: meta.Meta{
										Pos: meta.Position{
											Offset: 331,
											Line:   22,
											Column: 38,
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 296,
										Line:   22,
										Column: 3,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// service`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Offset: 259,
										Line:   20,
										Column: 1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Offset: 270,
								Line:   21,
								Column: 1,
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
			p := parser.NewParser(lexer.NewLexer(strings.NewReader(test.input)))
			got, err := p.ParseProto()
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
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantProto))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
