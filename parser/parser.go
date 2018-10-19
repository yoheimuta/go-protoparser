package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer"

// Parser is a parser.
type Parser struct {
	lex *lexer.Lexer

	permissive bool
}

// ConfigOption is an option for Parser.
type ConfigOption func(*Parser)

// WithPermissive is an option to allow the permissive parsing rather than the just documented spec.
func WithPermissive(permissive bool) ConfigOption {
	return func(p *Parser) {
		p.permissive = permissive
	}
}

// NewParser creates a new Parser.
func NewParser(lex *lexer.Lexer, opts ...ConfigOption) *Parser {
	p := &Parser{
		lex: lex,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// IsEOF checks whether the lex's read buffer is empty.
func (p *Parser) IsEOF() bool {
	p.lex.Next()
	defer p.lex.UnNext()
	return p.lex.IsEOF()
}
