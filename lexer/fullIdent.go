package lexer

import "github.com/yoheimuta/go-protoparser/v4/lexer/scanner"

// ReadFullIdent reads a fullIdent.
// fullIdent = ident { "." ident }
func (lex *Lexer) ReadFullIdent() (string, scanner.Position, error) {
	lex.Next()
	if lex.Token != scanner.TIDENT {
		return "", scanner.Position{}, lex.unexpected(lex.Text, "TIDENT")
	}
	startPos := lex.Pos

	fullIdent := lex.Text
	lex.Next()

	for !lex.IsEOF() {
		if lex.Token != scanner.TDOT {
			lex.UnNext()
			break
		}

		lex.Next()
		if lex.Token != scanner.TIDENT {
			return "", scanner.Position{}, lex.unexpected(lex.Text, "TIDENT")
		}
		fullIdent += "." + lex.Text
		lex.Next()
	}
	return fullIdent, startPos, nil
}
