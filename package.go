package protoparser

import (
	"fmt"
	"text/scanner"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// 'package' var';'
func parsePackage(lex *lexer.Lexer) (string, error) {
	text := lex.Text()
	if text != "package" {
		return "", fmt.Errorf("[BUG] not found package, Text=%s", text)
	}

	// consume 'package' {
	lex.Next()
	// }

	var packageName string
	for lex.Text() != ";" && lex.Token != scanner.EOF {
		packageName += lex.Text()

		// consume {
		lex.Next()
		// }
	}

	// consume ';' {
	lex.Next()
	// }
	return packageName, nil
}
