package protoparser

import (
	"fmt"
	"text/scanner"
)

// Service is one of top level definition in a protocol buffer.
type Service struct {
	Comments []string
	Name     string
	RPCs     []*RPC
}

// "service var '{' serviceContent '}'
func parseService(lex *Lexer) (*Service, error) {
	text := lex.Text()
	if text != "service" {
		return nil, fmt.Errorf("[BUG] not found service, Text=%s", text)
	}
	// consume 'service'
	lex.Next()

	// get the service name {
	name := lex.Text()
	lex.Next()
	// }

	// get the rpcs {
	/// consume '{' {
	lex.Next()
	/// }
	rpcs, err := parseServiceContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.Next()
	// }

	return &Service{
		Name: name,
		RPCs: rpcs,
	}, nil
}

// rpc
func parseServiceContent(lex *Lexer) ([]*RPC, error) {
	var rpcs []*RPC
	for lex.Text() != "}" && lex.token != scanner.EOF {
		if lex.token != scanner.Comment {
			return nil, fmt.Errorf("not found comment, Text=%s", lex.Text())
		}
		comments := parseComments(lex)

		switch lex.Text() {
		case "rpc":
			var rpc *RPC
			rpc, err := parseRPC(lex)
			if err != nil {
				return nil, err
			}
			rpc.Comments = append(rpc.Comments, comments...)
			rpcs = append(rpcs, rpc)
		default:
			return nil, fmt.Errorf("not found rpc, Text=%s", lex.Text())
		}
	}
	return rpcs, nil
}
