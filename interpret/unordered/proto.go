package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// ProtoBody is unordered in nature, but each slice field preserves the original order.
type ProtoBody struct {
	Imports         []*parser.Import
	Packages        []*parser.Package
	Options         []*parser.Option
	Messages        []*Message
	Extends         []*Extend
	Enums           []*Enum
	Services        []*Service
	EmptyStatements []*parser.EmptyStatement
}

// Proto represents a protocol buffer definition.
type Proto struct {
	Syntax    *parser.Syntax
	ProtoBody *ProtoBody
}

// InterpretProto interprets *parser.Proto to *Proto.
func InterpretProto(src *parser.Proto) (*Proto, error) {
	if src == nil {
		return nil, nil
	}

	enumBody, err := interpretProtoBody(src.ProtoBody)
	if err != nil {
		return nil, err
	}
	return &Proto{
		Syntax:    src.Syntax,
		ProtoBody: enumBody,
	}, nil
}

func interpretProtoBody(src []parser.Visitee) (
	*ProtoBody,
	error,
) {
	var imports []*parser.Import
	var packages []*parser.Package
	var options []*parser.Option
	var messages []*Message
	var extends []*Extend
	var enums []*Enum
	var services []*Service
	var emptyStatements []*parser.EmptyStatement
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Import:
			imports = append(imports, t)
		case *parser.Package:
			packages = append(packages, t)
		case *parser.Option:
			options = append(options, t)
		case *parser.Message:
			message, err := InterpretMessage(t)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		case *parser.Extend:
			extend, err := InterpretExtend(t)
			if err != nil {
				return nil, err
			}
			extends = append(extends, extend)
		case *parser.Enum:
			enum, err := InterpretEnum(t)
			if err != nil {
				return nil, err
			}
			enums = append(enums, enum)
		case *parser.Service:
			service, err := InterpretService(t)
			if err != nil {
				return nil, err
			}
			services = append(services, service)
		case *parser.EmptyStatement:
			emptyStatements = append(emptyStatements, t)
		default:
			return nil, fmt.Errorf("invalid ProtoBody type %T of %v", t, t)
		}
	}
	return &ProtoBody{
		Imports:         imports,
		Packages:        packages,
		Options:         options,
		Messages:        messages,
		Extends:         extends,
		Enums:           enums,
		Services:        services,
		EmptyStatements: emptyStatements,
	}, nil
}
