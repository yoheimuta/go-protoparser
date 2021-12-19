package lexer

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
)

// ReadConstant reads a constant. If permissive is true, accepts multiline string literals.
// constant = fullIdent | ( [ "-" | "+" ] intLit ) | ( [ "-" | "+" ] floatLit ) | strLit | boolLit
func (lex *Lexer) ReadConstant(permissive bool) (string, scanner.Position, error) {
	lex.NextLit()

	startPos := lex.Pos
	cons := lex.Text

	switch {
	case lex.Token == scanner.TSTRLIT:
		if permissive {
			return lex.mergeMultilineStrLit(), startPos, nil
		}
		return cons, startPos, nil
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

// Merges a multiline string literal into a single string.
func (lex *Lexer) mergeMultilineStrLit() string {
	q := "'"
	if strings.HasPrefix(lex.Text, "\"") {
		q = "\""
	}
	var b strings.Builder
	b.WriteString(q)
	for lex.Token == scanner.TSTRLIT {
		strippedString := strings.Trim(lex.Text, q)
		b.WriteString(strippedString)
		lex.NextLit()
	}
	lex.UnNext()
	b.WriteString(q)
	return b.String()
}
