package protoparser

import (
	"io"
	"text/scanner"
)

// EnumField は Enum の値を表す。
type EnumField struct {
	Comments []string
	Name     string
}

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

// RPC は関数を表す。
type RPC struct {
	Comments []string
	Name     string
	Argument *Type
	Return   *Type
}

// Service は複数の RPC を定義するサービスを表す。
type Service struct {
	Comments []string
	Name     string
	RPCs     []*RPC
}

// ProtocolBuffer は Protocol Buffers ファイルをパースした結果を表す。
type ProtocolBuffer struct {
	Package  string
	Service  *Service
	Messages []*Message
}

// Parse は Protocol Bufffers ファイルをパースする。
func Parse(input io.Reader) (*ProtocolBuffer, error) {
	lex := new(lexer)
	lex.scan.Init(input)
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanComments
	lex.next()
	return parse(lex)
}
