package protoparser

import (
	"fmt"
	"text/scanner"
)

// 'package' var';'
func parsePackage(lex *Lexer) (string, error) {
	text := lex.Text()
	if text != "package" {
		return "", fmt.Errorf("[BUG] not found package, Text=%s", text)
	}

	// consume 'package' {
	lex.Next()
	// }

	var packageName string
	for lex.Text() != ";" && lex.token != scanner.EOF {
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
