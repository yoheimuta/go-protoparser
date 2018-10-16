package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer"

// Parser is a parser.
type Parser struct {
	lex *lexer.Lexer2
}

// NewParser creates a new Parser.
func NewParser(lex *lexer.Lexer2) *Parser {
	return &Parser{
		lex: lex,
	}
}

// IsEOF checks whether the lex's read buffer is empty.
func (p *Parser) IsEOF() bool {
	p.lex.Next()
	defer p.lex.UnNext()
	return p.lex.IsEOF()
}
