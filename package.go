package protoparser

import (
	"fmt"
	"text/scanner"
)

// 'package' var';'
func parsePackage(lex *lexer) (string, error) {
	text := lex.text()
	if text != "package" {
		return "", fmt.Errorf("[BUG] not found package, text=%s", text)
	}

	// consume 'package' {
	lex.next()
	// }

	var packageName string
	for lex.text() != ";" && lex.token != scanner.EOF {
		packageName += lex.text()

		// consume {
		lex.next()
		// }
	}

	// consume ';' {
	lex.next()
	// }
	return packageName, nil
}
