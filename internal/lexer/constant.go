package lexer

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// ReadConstant reads a constant. If permissive is true, accepts multiline string literals.
// constant = fullIdent | ( [ "-" | "+" ] intLit ) | ( [ "-" | "+" ] floatLit ) | strLit | boolLit
func (lex *Lexer) ReadConstant(permissive bool) (string, scanner.Position, error) {
	lex.NextLit()

	startPos := lex.Pos
	cons := lex.Text

	switch {
	case lex.Token == scanner.TSTRLIT:
		if !permissive {
			return cons, startPos, nil
		}
		var b strings.Builder
		b.WriteString("\"")
		for lex.Token == scanner.TSTRLIT {
			strippedString := strings.Trim(lex.Text, "\"")
			b.WriteString(strippedString)
			lex.NextLit()
		}
		lex.UnNext()
		b.WriteString("\"")
		return b.String(), startPos, nil
	case lex.Token == scanner.TBOOLLIT:
		return cons, startPos, nil
	case lex.Token == scanner.TIDENT:
		lex.UnNext()
		fullIdent, pos, err := lex.ReadFullIdent()
		if err != nil {
			return "", scanner.Position{}, err
		}
		return fullIdent, pos, nil
	case lex.Token == scanner.TINTLIT, lex.Token == scanner.TFLOATLIT:
		return cons, startPos, nil
	case lex.Text == "-" || lex.Text == "+":
		lex.NextLit()

		switch lex.Token {
		case scanner.TINTLIT, scanner.TFLOATLIT:
			cons += lex.Text
			return cons, startPos, nil
		default:
			return "", scanner.Position{}, lex.unexpected(lex.Text, "TINTLIT or TFLOATLIT")
		}
	default:
		return "", scanner.Position{}, lex.unexpected(lex.Text, "constant")
	}
}
