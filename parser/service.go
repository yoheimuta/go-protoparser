package parser

import (
	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// RPCRequest is a request of RPC.
type RPCRequest struct {
	IsStream    bool
	MessageType string

	// Meta is the meta information.
	Meta meta.Meta
}

// RPCResponse is a response of RPC.
type RPCResponse struct {
	IsStream    bool
	MessageType string

	// Meta is the meta information.
	Meta meta.Meta
}

// RPC is a Remote Procedure Call.
type RPC struct {
	RPCName     string
	RPCRequest  *RPCRequest
	RPCResponse *RPCResponse
	Options     []*Option

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (r *RPC) SetInlineComment(comment *Comment) {
	r.InlineComment = comment
}

// Service consists of RPCs.
type Service struct {
	ServiceName string
	// ServiceBody can have options and rpcs.
	ServiceBody []interface{}

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (s *Service) SetInlineComment(comment *Comment) {
	s.InlineComment = comment
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
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("serviceName")
	}
	serviceName := p.lex.Text

	serviceBody, inlineLeftCurly, err := p.parseServiceBody()
	if err != nil {
		return nil, err
	}

	return &Service{
		ServiceName:                  serviceName,
		ServiceBody:                  serviceBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta:                         meta.NewMeta(startPos),
	}, nil
}

// serviceBody = "{" { option | rpc | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseServiceBody() ([]interface{}, *Comment, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var stmts []interface{}
	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			HasInlineCommentSetter
		}

		switch token {
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TRPC:
			rpc, err := p.parseRPC()
			if err != nil {
				return nil, nil, err
			}
			rpc.Comments = comments
			stmt = rpc
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, nil, err
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return stmts, inlineLeftCurly, nil
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
	startPos := p.lex.Pos

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
		Meta:        meta.NewMeta(startPos),
	}, nil
}

// rpcRequest = "(" [ "stream" ] messageType ")"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCRequest() (*RPCRequest, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}
	startPos := p.lex.Pos

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, _, err := p.lex.ReadMessageType()
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
		Meta:        meta.NewMeta(startPos),
	}, nil
}

// rpcResponse = "(" [ "stream" ] messageType ")"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCResponse() (*RPCResponse, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}
	startPos := p.lex.Pos

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, _, err := p.lex.ReadMessageType()
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
		Meta:        meta.NewMeta(startPos),
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
		p.lex.UnNext()
	}
}
