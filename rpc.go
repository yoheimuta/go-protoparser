package protoparser

import (
	"fmt"
	"text/scanner"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// RPC is the basic element of a protocol buffer service.
type RPC struct {
	Comments []string
	Name     string
	Argument *Type
	Return   *Type
}

// Name'('Argument')' 'returns' '('Return')' ('{''}'|';')
func parseRPC(lex *lexer.Lexer) (*RPC, error) {
	text := lex.Text()
	if text != "rpc" {
		return nil, fmt.Errorf("not found rpc, Text=%s", text)
	}
	// consume 'rpc' {
	lex.Next()
	// }

	rpc := &RPC{}

	for lex.Text() != "}" && lex.Text() != ";" && lex.Token != scanner.EOF {
		token := lex.Text()
		if rpc.Name == "" {
			rpc.Name = token
			lex.Next()
			continue
		}
		if rpc.Argument == nil {
			// consume '(' {
			lex.Next()
			// }

			rpc.Argument = parseType(lex)

			// consume ')' {
			lex.Next()
			// }
			continue
		}
		if rpc.Return == nil {
			// consume 'returns' {
			lex.Next()
			// }
			// consume '(' {
			lex.Next()
			// }

			rpc.Return = parseType(lex)

			// consume ')' {
			lex.Next()
			// }
			continue
		}

		if token == "{" {
			lex.Next()
		}
	}

	// consume '}' or ';' {
	lex.Next()
	// }

	return rpc, nil
}
