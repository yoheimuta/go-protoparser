package protoparser

import (
	"io"
)

// Enum は Enum 型を表す。
type Enum struct {
	Comments   []string
	Name       string
	EnumFields []*EnumField
}

// Message は独自に定義した型情報を表す。
type Message struct {
	Comments []string
	Name     string
	Fields   []*Field
	Nests    []*Message
	Enums    []*Enum
	Oneofs   []*Oneof
}

// ProtocolBuffer は Protocol Buffers ファイルをパースした結果を表す。
type ProtocolBuffer struct {
	Package  string
	Service  *Service
	Messages []*Message
}

// Parse は Protocol Bufffers ファイルをパースする。
func Parse(input io.Reader) (*ProtocolBuffer, error) {
	lex := newlexer(input)
	return parse(lex)
}
