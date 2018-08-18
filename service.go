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
func parseService(lex *lexer) (*Service, error) {
	text := lex.text()
	if text != "service" {
		return nil, fmt.Errorf("[BUG] not found service, text=%s", text)
	}
	// consume 'service'
	lex.next()

	// get the service name {
	name := lex.text()
	lex.next()
	// }

	// get the rpcs {
	/// consume '{' {
	lex.next()
	/// }
	rpcs, err := parseServiceContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.next()
	// }

	return &Service{
		Name: name,
		RPCs: rpcs,
	}, nil
}

// rpc
func parseServiceContent(lex *lexer) ([]*RPC, error) {
	var rpcs []*RPC
	for lex.text() != "}" && lex.token != scanner.EOF {
		if lex.token != scanner.Comment {
			return nil, fmt.Errorf("not found comment, text=%s", lex.text())
		}
		comments := parseComments(lex)

		switch lex.text() {
		case "rpc":
			var rpc *RPC
			rpc, err := parseRPC(lex)
			if err != nil {
				return nil, err
			}
			rpc.Comments = append(rpc.Comments, comments...)
			rpcs = append(rpcs, rpc)
		default:
			return nil, fmt.Errorf("not found rpc, text=%s", lex.text())
		}
	}
	return rpcs, nil
}
