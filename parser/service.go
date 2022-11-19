package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
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
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (r *RPC) SetInlineComment(comment *Comment) {
	r.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (r *RPC) Accept(v Visitor) {
	if !v.VisitRPC(r) {
		return
	}

	for _, comment := range r.Comments {
		comment.Accept(v)
	}
	if r.InlineComment != nil {
		r.InlineComment.Accept(v)
	}
}

// Service consists of RPCs.
type Service struct {
	ServiceName string
	// ServiceBody can have options and rpcs.
	ServiceBody []Visitee

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

// Accept dispatches the call to the visitor.
func (s *Service) Accept(v Visitor) {
	if !v.VisitService(s) {
		return
	}

	for _, body := range s.ServiceBody {
		body.Accept(v)
	}
	for _, comment := range s.Comments {
		comment.Accept(v)
	}
	if s.InlineComment != nil {
		s.InlineComment.Accept(v)
	}
	if s.InlineCommentBehindLeftCurly != nil {
		s.InlineCommentBehindLeftCurly.Accept(v)
	}
}

// ParseService parses the service.
//
//	service = "service" serviceName "{" { option | rpc | emptyStatement } "}"
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

	serviceBody, inlineLeftCurly, lastPos, err := p.parseServiceBody()
	if err != nil {
		return nil, err
	}

	return &Service{
		ServiceName:                  serviceName,
		ServiceBody:                  serviceBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

// serviceBody = "{" { option | rpc | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseServiceBody() (
	[]Visitee,
	*Comment,
	scanner.Position,
	error,
) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var stmts []Visitee
	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			HasInlineCommentSetter
			Visitee
		}

		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, Visitee(comment))
				}
			}
			p.lex.Next()

			lastPos := p.lex.Pos
			if p.permissive {
				// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TRPC:
			rpc, err := p.parseRPC()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			rpc.Comments = comments
			stmt = rpc
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
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
	var inlineLeftCurly *Comment
	p.lex.Next()
	lastPos := p.lex.Pos
	switch p.lex.Token {
	case scanner.TLEFTCURLY:
		p.lex.UnNext()
		opts, inlineLeftCurly, err = p.parseRPCOptions()
		if err != nil {
			return nil, err
		}
		lastPos = p.lex.Pos
		if p.permissive {
			// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
			p.lex.ConsumeToken(scanner.TSEMICOLON)
			if p.lex.Token == scanner.TSEMICOLON {
				lastPos = p.lex.Pos
			}
		}
	case scanner.TSEMICOLON:
		break
	default:
		return nil, p.unexpected("{ or ;")
	}

	return &RPC{
		RPCName:                      rpcName,
		RPCRequest:                   rpcRequest,
		RPCResponse:                  rpcResponse,
		Options:                      opts,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
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
		Meta:        meta.Meta{Pos: startPos.Position},
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
		Meta:        meta.Meta{Pos: startPos.Position},
	}, nil
}

// rpcOptions = ( "{" {option | emptyStatement } "}" )
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
func (p *Parser) parseRPCOptions() ([]*Option, *Comment, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var options []*Option
	for {
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, err
			}
			options = append(options, option)
		case scanner.TRIGHTCURLY:
			// This spec is not documented, but allowed in general.
			break
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, nil, err
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return options, inlineLeftCurly, nil
		}
		p.lex.UnNext()
	}
}
