package lexer

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ReadConstant reads a constant.
// constant = fullIdent | ( [ "-" | "+" ] intLit ) | ( [ "-" | "+" ] floatLit ) | strLit | boolLit
func (lex *Lexer2) ReadConstant() (string, error) {
	lex.NextLit()

	cons := lex.Text

	switch {
	case lex.Token == scanner.TSTRLIT:
		return cons, nil
	case lex.Token == scanner.TBOOLLIT:
		return cons, nil
	case lex.Token == scanner.TIDENT,
		// go-proto-validators requires this exceptions.
		lex.Token == scanner.TLEFTCURLY && lex.permissive:
		lex.UnNext()
		fullIdent, err := lex.ReadFullIdent()
		if err != nil {
			return "", err
		}
		return fullIdent, nil
	case lex.Token == scanner.TINTLIT, lex.Token == scanner.TFLOATLIT:
		return cons, nil
	case lex.Text == "-" || lex.Text == "+":
		lex.NextLit()

		switch lex.Token {
		case scanner.TINTLIT, scanner.TFLOATLIT:
			cons += lex.Text
			return cons, nil
		default:
			return "", lex.unexpected(lex.Text, "TINTLIT or TFLOATLIT")
		}
	default:
		return "", lex.unexpected(lex.Text, "constant")
	}
}
