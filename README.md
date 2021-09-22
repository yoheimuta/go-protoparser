# go-protoparser [![GoDoc](https://godoc.org/github.com/yoheimuta/go-protoparser/v4?status.svg)](https://pkg.go.dev/github.com/yoheimuta/go-protoparser/v4)[![CircleCI](https://circleci.com/gh/yoheimuta/go-protoparser/tree/master.svg?style=svg)](https://circleci.com/gh/yoheimuta/go-protoparser/tree/master)[![Go Report Card](https://goreportcard.com/badge/github.com/yoheimuta/go-protoparser/v4)](https://goreportcard.com/report/github.com/yoheimuta/go-protoparser/v4)[![Release](http://img.shields.io/github/release/yoheimuta/go-protoparser.svg?style=flat)](https://github.com/yoheimuta/go-protoparser/releases/latest)[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/yoheimuta/go-protoparser/blob/master/LICENSE.md)

go-protoparser is a yet another Go package which parses a Protocol Buffer file (proto2+proto3).

- Conforms to the exactly [official spec](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec).
- Undergone rigorous testing. The parser can parses all examples of the official spec well.
- Easy to use the parser. You can just call the [Parse function](https://godoc.org/github.com/yoheimuta/go-protoparser/v4#Parse) and receive the [Proto struct](https://godoc.org/github.com/yoheimuta/go-protoparser/v4/parser#Proto).
  - If you don't care about the order of body elements, consider to use the [unordered.Proto struct](https://godoc.org/github.com/yoheimuta/go-protoparser/v4/interpret/unordered#Proto).
  - Or if you want to use the visitor pattern, use the [Visitor struct](https://godoc.org/github.com/yoheimuta/go-protoparser/v4/parser#Visitor).

### Installation

```
GO111MODULE=on go get github.com/yoheimuta/go-protoparser/v4
```

### Example

A Protocol Buffer file versioned 3 which is [an example of the official reference](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file).

```proto
syntax = "proto3";
// An example of the official reference
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file
package examplepb;
import public "other.proto";
option java_package = "com.example.foo";
enum EnumAllowingAlias {
    option allow_alias = true;
    UNKNOWN = 0;
    STARTED = 1;
    RUNNING = 2 [(custom_option) = "this is a "
                                   "string on two lines"
                ];
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
service HelloService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {};
}
```

The Parsed result is a Go typed struct. The below output is encoded to JSON for simplicity.

```json
{
  "Syntax": {
    "ProtobufVersion": "proto3",
    "ProtobufVersionQuote": "\"proto3\"",
    "Comments": null,
    "InlineComment": null,
    "Meta": {
      "Pos": {
        "Filename": "simple.proto",
        "Offset": 0,
        "Line": 1,
        "Column": 1
      },
      "LastPos": {
        "Filename": "",
        "Offset": 0,
        "Line": 0,
        "Column": 0
      }
    }
  },
  "ProtoBody": [
    {
      "Name": "examplepb",
      "Comments": [
        {
          "Raw": "// An example of the official reference",
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 19,
              "Line": 2,
              "Column": 1
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "Raw": "// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file",
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 59,
              "Line": 3,
              "Column": 1
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        }
      ],
      "InlineComment": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 151,
          "Line": 4,
          "Column": 1
        },
        "LastPos": {
          "Filename": "",
          "Offset": 0,
          "Line": 0,
          "Column": 0
        }
      }
    },
    {
      "Modifier": 1,
      "Location": "\"other.proto\"",
      "Comments": null,
      "InlineComment": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 170,
          "Line": 5,
          "Column": 1
        },
        "LastPos": {
          "Filename": "",
          "Offset": 0,
          "Line": 0,
          "Column": 0
        }
      }
    },
    {
      "OptionName": "java_package",
      "Constant": "\"com.example.foo\"",
      "Comments": null,
      "InlineComment": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 199,
          "Line": 6,
          "Column": 1
        },
        "LastPos": {
          "Filename": "",
          "Offset": 0,
          "Line": 0,
          "Column": 0
        }
      }
    },
    {
      "EnumName": "EnumAllowingAlias",
      "EnumBody": [
        {
          "OptionName": "allow_alias",
          "Constant": "true",
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 269,
              "Line": 8,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "Ident": "UNKNOWN",
          "Number": "0",
          "EnumValueOptions": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 300,
              "Line": 9,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "Ident": "STARTED",
          "Number": "1",
          "EnumValueOptions": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 317,
              "Line": 10,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "Ident": "RUNNING",
          "Number": "2",
          "EnumValueOptions": [
            {
              "OptionName": "(custom_option)",
              "Constant": "\"this is a string on two lines\""
            }
          ],
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 334,
              "Line": 11,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        }
      ],
      "Comments": null,
      "InlineComment": null,
      "InlineCommentBehindLeftCurly": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 240,
          "Line": 7,
          "Column": 1
        },
        "LastPos": {
          "Filename": "simple.proto",
          "Offset": 454,
          "Line": 14,
          "Column": 1
        }
      }
    },
    {
      "MessageName": "outer",
      "MessageBody": [
        {
          "OptionName": "(my_option).a",
          "Constant": "true",
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 476,
              "Line": 16,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "MessageName": "inner",
          "MessageBody": [
            {
              "IsRepeated": false,
              "IsRequired": false,
              "IsOptional": false,
              "Type": "int64",
              "FieldName": "ival",
              "FieldNumber": "1",
              "FieldOptions": null,
              "Comments": null,
              "InlineComment": null,
              "Meta": {
                "Pos": {
                  "Filename": "simple.proto",
                  "Offset": 544,
                  "Line": 18,
                  "Column": 7
                },
                "LastPos": {
                  "Filename": "",
                  "Offset": 0,
                  "Line": 0,
                  "Column": 0
                }
              }
            }
          ],
          "Comments": null,
          "InlineComment": null,
          "InlineCommentBehindLeftCurly": {
            "Raw": "// Level 2",
            "Meta": {
              "Pos": {
                "Filename": "simple.proto",
                "Offset": 527,
                "Line": 17,
                "Column": 23
              },
              "LastPos": {
                "Filename": "",
                "Offset": 0,
                "Line": 0,
                "Column": 0
              }
            }
          },
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 509,
              "Line": 17,
              "Column": 5
            },
            "LastPos": {
              "Filename": "simple.proto",
              "Offset": 564,
              "Line": 19,
              "Column": 5
            }
          }
        },
        {
          "IsRepeated": true,
          "IsRequired": false,
          "IsOptional": false,
          "Type": "inner",
          "FieldName": "inner_message",
          "FieldNumber": "2",
          "FieldOptions": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 570,
              "Line": 20,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "IsRepeated": false,
          "IsRequired": false,
          "IsOptional": false,
          "Type": "EnumAllowingAlias",
          "FieldName": "enum_field",
          "FieldNumber": "3",
          "FieldOptions": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 608,
              "Line": 21,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        },
        {
          "KeyType": "int32",
          "Type": "string",
          "MapName": "my_map",
          "FieldNumber": "4",
          "FieldOptions": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 645,
              "Line": 22,
              "Column": 5
            },
            "LastPos": {
              "Filename": "",
              "Offset": 0,
              "Line": 0,
              "Column": 0
            }
          }
        }
      ],
      "Comments": null,
      "InlineComment": null,
      "InlineCommentBehindLeftCurly": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 456,
          "Line": 15,
          "Column": 1
        },
        "LastPos": {
          "Filename": "simple.proto",
          "Offset": 676,
          "Line": 23,
          "Column": 1
        }
      }
    },
    {
      "ServiceName": "HelloService",
      "ServiceBody": [
        {
          "RPCName": "SayHello",
          "RPCRequest": {
            "IsStream": false,
            "MessageType": "HelloRequest",
            "Meta": {
              "Pos": {
                "Filename": "simple.proto",
                "Offset": 716,
                "Line": 25,
                "Column": 16
              },
              "LastPos": {
                "Filename": "",
                "Offset": 0,
                "Line": 0,
                "Column": 0
              }
            }
          },
          "RPCResponse": {
            "IsStream": false,
            "MessageType": "HelloResponse",
            "Meta": {
              "Pos": {
                "Filename": "simple.proto",
                "Offset": 739,
                "Line": 25,
                "Column": 39
              },
              "LastPos": {
                "Filename": "",
                "Offset": 0,
                "Line": 0,
                "Column": 0
              }
            }
          },
          "Options": null,
          "Comments": null,
          "InlineComment": null,
          "Meta": {
            "Pos": {
              "Filename": "simple.proto",
              "Offset": 703,
              "Line": 25,
              "Column": 3
            },
            "LastPos": {
              "Filename": "simple.proto",
              "Offset": 757,
              "Line": 25,
              "Column": 57
            }
          }
        }
      ],
      "Comments": null,
      "InlineComment": null,
      "InlineCommentBehindLeftCurly": null,
      "Meta": {
        "Pos": {
          "Filename": "simple.proto",
          "Offset": 678,
          "Line": 24,
          "Column": 1
        },
        "LastPos": {
          "Filename": "simple.proto",
          "Offset": 759,
          "Line": 26,
          "Column": 1
        }
      }
    }
  ],
  "Meta": {
    "Filename": "simple.proto"
  }
}
```

### Usage

See also `_example/dump`.

```go
func run() int {
	flag.Parse()

	reader, err := os.Open(*proto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s, err %v\n", *proto, err)
		return 1
	}
	defer reader.Close()

	got, err := protoparser.Parse(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse, err %v\n", err)
		return 1
	}

	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	fmt.Print(string(gotJSON))
	return 0
}

func main() {
	os.Exit(run())
}
```

### Users

- [protolint](https://github.com/yoheimuta/protolint)
- [Known importers](https://pkg.go.dev/github.com/yoheimuta/go-protoparser/v4?tab=importedby)

### Motivation

There exists the similar protobuf parser packages in Go.

But I could not find the parser which just return a parsing result well-typed value.
A parser which supports a visitor pattern is useful to implement like linter, but it may be difficult to use.
It can be sufficient for most parsing situations to just return a parsing result well-typed value.
This is easier to use.

### Contributing

- Fork it
- Create your feature branch: git checkout -b your-new-feature
- Commit changes: git commit -m 'Add your feature'
- Pass all tests
- Push to the branch: git push origin your-new-feature
- Submit a pull request

### License

The MIT License (MIT)

### Acknowledgement

Thank you to the proto package: https://github.com/emicklei/proto

I referred to the package for the good proven design, interface and some source code.
