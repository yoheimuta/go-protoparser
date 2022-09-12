package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"path/filepath"

	protoparser "github.com/yoheimuta/go-protoparser/v4"
)

var (
	proto      = flag.String("proto", "_testdata/simple.proto", "path to the Protocol Buffer file")
	debug      = flag.Bool("debug", false, "debug flag to output more parsing process detail")
	permissive = flag.Bool("permissive", true, "permissive flag to allow the permissive parsing rather than the just documented spec")
	unordered  = flag.Bool("unordered", false, "unordered flag to output another one without interface{}")
)

func run() int {
	flag.Parse()

	reader, err := os.Open(*proto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s, err %v\n", *proto, err)
		return 1
	}
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	got, err := protoparser.Parse(
		reader,
		protoparser.WithDebug(*debug),
		protoparser.WithPermissive(*permissive),
		protoparser.WithFilename(filepath.Base(*proto)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse, err %v\n", err)
		return 1
	}

	var v interface{}
	v = got
	if *unordered {
		v, err = protoparser.UnorderedInterpret(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to interpret, err %v\n", err)
			return 1
		}
	}

	gotJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	fmt.Print(string(gotJSON))
	return 0
}

func main() {
	os.Exit(run())
}
