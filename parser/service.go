package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// RPCRequest is a request of RPC.
type RPCRequest struct {
	IsStream    bool
	MessageType string
}

// RPCResponse is a response of RPC.
type RPCResponse struct {
	IsStream    bool
	MessageType string
}

// RPC is a Remote Procedure Call.
type RPC struct {
	RPCName     string
	RPCRequest  *RPCRequest
	RPCResponse *RPCResponse
	Options     []*Option

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
}

// Service consists of RPCs.
type Service struct {
	ServiceName string
	// ServiceBody can have options and rpcs.
	ServiceBody []interface{}

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
}

// ParseService parses the service.
//  service = "service" serviceName "{" { option | rpc | emptyStatement } "}"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) ParseService() (*Service, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TSERVICE {
		return nil, p.unexpected("service")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("serviceName")
	}
	serviceName := p.lex.Text

	serviceBody, err := p.parseServiceBody()
	if err != nil {
		return nil, err
	}

	return &Service{
		ServiceName: serviceName,
		ServiceBody: serviceBody,
	}, nil
}

// serviceBody = "{" { option | rpc | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseServiceBody() ([]interface{}, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	var stmts []interface{}
	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			option.Comments = comments
			stmts = append(stmts, option)
		case scanner.TRPC:
			rpc, err := p.parseRPC()
			if err != nil {
				return nil, err
			}
			rpc.Comments = comments
			stmts = append(stmts, rpc)
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, err
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return stmts, nil
		}
		p.lex.UnNext()
	}
}

// rpc = "rpc" rpcName "(" [ "stream" ] messageType ")" "returns" "(" [ "stream" ]
// messageType ")" (( "{" {option | emptyStatement } "}" ) | ";")
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPC() (*RPC, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRPC {
		return nil, p.unexpected("rpc")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("serviceName")
	}
	rpcName := p.lex.Text

	rpcRequest, err := p.parseRPCRequest()
	if err != nil {
		return nil, err
	}

	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRETURNS {
		return nil, p.unexpected("returns")
	}

	rpcResponse, err := p.parseRPCResponse()
	if err != nil {
		return nil, err
	}

	var opts []*Option
	p.lex.Next()
	switch p.lex.Token {
	case scanner.TLEFTCURLY:
		p.lex.UnNext()
		opts, err = p.parseRPCOptions()
		if err != nil {
			return nil, err
		}
	case scanner.TSEMICOLON:
		break
	default:
		return nil, p.unexpected("{ or ;")
	}

	return &RPC{
		RPCName:     rpcName,
		RPCRequest:  rpcRequest,
		RPCResponse: rpcResponse,
		Options:     opts,
	}, nil
}

// rpcRequest = "(" [ "stream" ] messageType ")"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCRequest() (*RPCRequest, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTPAREN {
		return nil, p.unexpected(")")
	}

	return &RPCRequest{
		IsStream:    isStream,
		MessageType: messageType,
	}, nil
}

// rpcResponse = "(" [ "stream" ] messageType ")"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCResponse() (*RPCResponse, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTPAREN {
		return nil, p.unexpected(")")
	}

	return &RPCResponse{
		IsStream:    isStream,
		MessageType: messageType,
	}, nil
}

// rpcOptions = ( "{" {option | emptyStatement } "}" )
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCOptions() ([]*Option, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	var options []*Option
	for {
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			options = append(options, option)
		case scanner.TRIGHTCURLY:
			// This spec is not documented, but allowed in general.
			break
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, err
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return options, nil
		}
	}
}
