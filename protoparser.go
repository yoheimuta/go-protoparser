package protoparser

import (
	"io"
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
