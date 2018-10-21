package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
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
  message inner {   // Level 2
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
				},
				ProtoBody: []interface{}{
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
						EnumBody: []interface{}{
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
						MessageBody: []interface{}{
							&parser.Option{
								OptionName: "(my_option).a",
								Constant:   "true",
							},
							&parser.Message{
								MessageName: "inner",
								MessageBody: []interface{}{
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
				},
				ProtoBody: []interface{}{
					&parser.Service{
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
`,
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Comments: []*parser.Comment{
						{
							Raw: `// syntax`,
						},
						{
							Raw: `/*
syntax2
*/`,
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
