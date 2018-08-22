package protoparser

import (
	"io"
	"text/scanner"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// ProtocolBuffer is the parsed result from a Protocol Buffer file.
type ProtocolBuffer struct {
	Package  string
	Service  *Service
	Messages []*Message
	Enums    []*Enum
}

// Parse parses a Protocol Buffer file.
func Parse(input io.Reader) (*ProtocolBuffer, error) {
	lex := lexer.NewLexer(input)
	return parse(lex)
}

// comment\npackage...
// comment\nservice...
// comment\nmessage...
// comment\nenum...
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file
func parse(lex *lexer.Lexer) (*ProtocolBuffer, error) {
	var pkg string
	service := &Service{}
	var messages []*Message
	var enums []*Enum
	for lex.Token != scanner.EOF {
		comments := parseComments(lex)

		switch lex.Text() {
		case "package":
			p, err := parsePackage(lex)
			if err != nil {
				return nil, err
			}
			pkg = p
		case "service":
			s, err := parseService(lex)
			if err != nil {
				return nil, err
			}
			s.Comments = append(s.Comments, comments...)
			service = s
		case "message":
			message, err := parseMessage(lex)
			if err != nil {
				return nil, err
			}
			message.Comments = append(message.Comments, comments...)
			messages = append(messages, message)
		case "enum":
			enum, err := parseEnum(lex)
			if err != nil {
				return nil, err
			}
			enum.Comments = append(enum.Comments, comments...)
			enums = append(enums, enum)
		default:
			lex.Next()
			continue
		}
	}
	return &ProtocolBuffer{
		Package:  pkg,
		Service:  service,
		Messages: messages,
		Enums:    enums,
	}, nil
}
