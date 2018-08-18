package protoparser

import (
	"fmt"
	"text/scanner"
)

// RPC is the basic element of a protocol buffer service.
type RPC struct {
	Comments []string
	Name     string
	Argument *Type
	Return   *Type
}

// Name'('Argument')' 'returns' '('Return')' ('{''}'|';')
func parseRPC(lex *lexer) (*RPC, error) {
	text := lex.text()
	if text != "rpc" {
		return nil, fmt.Errorf("not found rpc, text=%s", text)
	}
	// consume 'rpc' {
	lex.next()
	// }

	rpc := &RPC{}

	for lex.text() != "}" && lex.text() != ";" && lex.token != scanner.EOF {
		token := lex.text()
		if rpc.Name == "" {
			rpc.Name = token
			lex.next()
			continue
		}
		if rpc.Argument == nil {
			// consume '(' {
			lex.next()
			// }

			rpc.Argument = parseType(lex)

			// consume ')' {
			lex.next()
			// }
			continue
		}
		if rpc.Return == nil {
			// consume 'returns' {
			lex.next()
			// }
			// consume '(' {
			lex.next()
			// }

			rpc.Return = parseType(lex)

			// consume ')' {
			lex.next()
			// }
			continue
		}

		if token == "{" {
			lex.next()
		}
	}

	// consume '}' or ';' {
	lex.next()
	// }

	return rpc, nil
}
