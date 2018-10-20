package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	protoparser "github.com/yoheimuta/go-protoparser"
)

var (
	proto      = flag.String("proto", "_testdata/parser.proto", "path to the Protocol Buffer file")
	debug      = flag.Bool("debug", false, "debug flag to output more parsing process detail")
	permissive = flag.Bool("permissive", true, "permissive flag to allow the permissive parsing rather than the just documented spec")
)

func run() int {
	flag.Parse()

	reader, err := os.Open(*proto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s, err %v\n", *proto, err)
		return 1
	}
	defer reader.Close()

	got, err := protoparser.Parse(reader, protoparser.WithDebug(*debug), protoparser.WithPermissive(*permissive))
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
