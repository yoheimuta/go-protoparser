package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// SimpleVisitor is a simple implementation of the parser.Visitor interface.
type SimpleVisitor struct{}

func (v *SimpleVisitor) VisitComment(*parser.Comment) {}
func (v *SimpleVisitor) VisitDeclaration(*parser.Declaration) bool {
	return true
}
func (v *SimpleVisitor) VisitEdition(*parser.Edition) bool {
	return true
}
func (v *SimpleVisitor) VisitEmptyStatement(*parser.EmptyStatement) bool {
	return true
}
func (v *SimpleVisitor) VisitEnum(*parser.Enum) bool {
	return true
}
func (v *SimpleVisitor) VisitEnumField(*parser.EnumField) bool {
	return true
}
func (v *SimpleVisitor) VisitExtend(*parser.Extend) bool {
	return true
}
func (v *SimpleVisitor) VisitExtensions(*parser.Extensions) bool {
	return true
}
func (v *SimpleVisitor) VisitField(*parser.Field) bool {
	return true
}
func (v *SimpleVisitor) VisitGroupField(*parser.GroupField) bool {
	return true
}
func (v *SimpleVisitor) VisitImport(*parser.Import) bool {
	return true
}
func (v *SimpleVisitor) VisitMapField(*parser.MapField) bool {
	return true
}
func (v *SimpleVisitor) VisitMessage(*parser.Message) bool {
	return true
}
func (v *SimpleVisitor) VisitOneof(*parser.Oneof) bool {
	return true
}
func (v *SimpleVisitor) VisitOneofField(*parser.OneofField) bool {
	return true
}
func (v *SimpleVisitor) VisitOption(*parser.Option) bool {
	return true
}
func (v *SimpleVisitor) VisitPackage(*parser.Package) bool {
	return true
}
func (v *SimpleVisitor) VisitReserved(*parser.Reserved) bool {
	return true
}
func (v *SimpleVisitor) VisitRPC(*parser.RPC) bool {
	return true
}
func (v *SimpleVisitor) VisitService(*parser.Service) bool {
	return true
}
func (v *SimpleVisitor) VisitSyntax(*parser.Syntax) bool {
	return true
}

func main() {
	// Proto file with double semicolon in service option
	input := `
syntax = "proto3";
import "google/protobuf/descriptor.proto";

extend google.protobuf.ServiceOptions  {
  string service_description = 51000;
}

service MyService{
  option(service_description) = "description";;
}
`

	// Parse the proto file
	p := parser.NewParser(lexer.NewLexer(strings.NewReader(input)))
	proto, err := p.ParseProto()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse proto: %v\n", err)
		os.Exit(1)
	}

	// Use the visitor pattern to process the parsed proto
	visitor := &SimpleVisitor{}
	proto.Accept(visitor)

	fmt.Println("Successfully processed proto with double semicolon")
}
