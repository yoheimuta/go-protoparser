package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ImportModifier is a modifier enum type for import behavior.
type ImportModifier uint

// Optional import modifier value to change a default behavior.
const (
	ImportModifierNone ImportModifier = iota
	ImportModifierPublic
	ImportModifierWeak
)

// Import is used to import another .proto's definitions.
type Import struct {
	Modifier ImportModifier
	Location string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
}

// ParseImport parses the import.
//  import = "import" [ "weak" | "public" ] strLit ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#import_statement
func (p *Parser) ParseImport() (*Import, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TIMPORT {
		return nil, p.unexpected(`"import"`)
	}

	var modifier ImportModifier
	p.lex.NextKeywordOrStrLit()
	switch p.lex.Token {
	case scanner.TPUBLIC:
		modifier = ImportModifierPublic
	case scanner.TWEAK:
		modifier = ImportModifierWeak
	case scanner.TSTRLIT:
		modifier = ImportModifierNone
		p.lex.UnNext()
	}

	p.lex.NextStrLit()
	if p.lex.Token != scanner.TSTRLIT {
		return nil, p.unexpected("strLit")
	}
	location := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Import{
		Modifier: modifier,
		Location: location,
	}, nil
}
