# go-protoparser [![GoDoc](https://godoc.org/github.com/yoheimuta/go-protoparser?status.svg)](https://godoc.org/github.com/yoheimuta/go-protoparser)[![CircleCI](https://circleci.com/gh/yoheimuta/go-protoparser/tree/master.svg?style=svg)](https://circleci.com/gh/yoheimuta/go-protoparser/tree/master)

go-protoparser is a yet another Go package which parses a Protocol Buffer file (proto3).

- Conforms to the exactly [official spec](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec). 
- Undergone rigorous testing. The parser can parses all examples of the official spec well.
- Easy to use the parser. You can just call the [Parse function](https://godoc.org/github.com/yoheimuta/go-protoparser#Parse) and receive the [Proto struct](https://godoc.org/github.com/yoheimuta/go-protoparser/parser#Proto).

### Installation

```
go get github.com/yoheimuta/go-protoparser
```

### Example

A Protocol Buffer file versioned 3 which is [an example of the official reference](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file).

```proto
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
```

The Parsed result is a Go typed struct. The below output is encoded to JSON for simplicity.

```json
{
  "Syntax": {
    "ProtobufVersion": "proto3"
  },
  "ProtoBody": [
    {
      "Modifier": 1,
      "Location": "\"other.proto\""
    },
    {
      "OptionName": "java_package",
      "Constant": "\"com.example.foo\""
    },
    {
      "EnumName": "EnumAllowingAlias",
      "EnumBody": [
        {
          "OptionName": "allow_alias",
          "Constant": "true"
        },
        {
          "Ident": "UNKNOWN",
          "Number": "0",
          "EnumValueOptions": null
        },
        {
          "Ident": "STARTED",
          "Number": "1",
          "EnumValueOptions": null
        },
        {
          "Ident": "RUNNING",
          "Number": "2",
          "EnumValueOptions": [
            {
              "OptionName": "(custom_option)",
              "Constant": "\"hello world\""
            }
          ]
        }
      ]
    },
    {
      "MessageName": "outer",
      "MessageBody": [
        {
          "OptionName": "(my_option).a",
          "Constant": "true"
        },
        {
          "MessageName": "inner",
          "MessageBody": [
            {
              "IsRepeated": false,
              "Type": "int64",
              "FieldName": "ival",
              "FieldNumber": "1",
              "FieldOptions": null
            }
          ]
        },
        {
          "IsRepeated": true,
          "Type": "inner",
          "FieldName": "inner_message",
          "FieldNumber": "2",
          "FieldOptions": null
        },
        {
          "IsRepeated": false,
          "Type": "EnumAllowingAlias",
          "FieldName": "enum_field",
          "FieldNumber": "3",
          "FieldOptions": null
        },
        {
          "KeyType": "int32",
          "Type": "string",
          "MapName": "my_map",
          "FieldNumber": "4",
          "FieldOptions": null
        }
      ]
    }
  ]
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

### TODO

- [ ] Parse comments. Currently, the parser ignores comments.

### Motivation

There exists the similar protobuf parser packages in Go.

For example, https://github.com/emicklei/proto is good and flexible, but the visitor pattern is not easy to use.
It can be sufficient for most parsing situations to just return a parsing result well-typed value. This is easier to use.