package lexer

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ReadFullIdent reads a fullIdent.
// fullIdent = ident { "." ident }
func (lex *Lexer) ReadFullIdent() (string, error) {
	ident, err := lex.readIdent()
	if err != nil {
		return "", err
	}
	fullIdent := ident
	lex.Next()

	for !lex.IsEOF() {
		if lex.Token != scanner.TDOT {
			lex.UnNext()
			break
		}

		ident, err = lex.readIdent()
		if err != nil {
			return "", err
		}
		fullIdent += "." + ident
		lex.Next()
	}
	return fullIdent, nil
}

func (lex *Lexer) readIdent() (string, error) {
	lex.Next()

	switch lex.Token {
	case scanner.TIDENT:
		return lex.Text, nil
	case scanner.TLEFTCURLY:
		// go-proto-validators requires this exceptions.
		if lex.permissive {
			ident := lex.Text
			for {
				lex.Next()
				ident += lex.Text
				if lex.Token == scanner.TRIGHTCURLY {
					return ident, nil
				}
			}
		}
	}
	return "", lex.unexpected(lex.Text, "TIDENT")
}
