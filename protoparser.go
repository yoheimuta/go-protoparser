package protoparser

import (
	"io"
	"text/scanner"
)

// ProtocolBuffer is the parsed result from a Protocol Buffer file.
type ProtocolBuffer struct {
	Package  string
	Service  *Service
	Messages []*Message
}

// Parse parses a Protocol Buffer file.
func Parse(input io.Reader) (*ProtocolBuffer, error) {
	lex := newlexer(input)
	return parse(lex)
}

// comment\npackage...
// comment\nservice...
// comment\nmessage...
func parse(lex *lexer) (*ProtocolBuffer, error) {
	var pkg string
	service := &Service{}
	var messages []*Message
	for lex.token != scanner.EOF {
		comments := parseComments(lex)

		switch lex.text() {
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
		default:
			lex.next()
			continue
		}
	}
	return &ProtocolBuffer{
		Package:  pkg,
		Service:  service,
		Messages: messages,
	}, nil
}
