package parser

import "github.com/yoheimuta/go-protoparser/v4/lexer/scanner"

// ProtoMeta represents a meta information about the Proto.
type ProtoMeta struct {
	// Filename is a name of file, if any.
	Filename string
}

// Proto represents a protocol buffer definition.
type Proto struct {
	Syntax *Syntax
	// ProtoBody is a slice of sum type consisted of *Import, *Package, *Option, *Message, *Enum, *Service, *Extend and *EmptyStatement.
	ProtoBody []Visitee
	Meta      *ProtoMeta
}

// Accept dispatches the call to the visitor.
func (p *Proto) Accept(v Visitor) {
	if p.Syntax != nil {
		p.Syntax.Accept(v)
	}

	for _, body := range p.ProtoBody {
		body.Accept(v)
	}
}

// ParseProto parses the proto.
//  proto = syntax { import | package | option | topLevelDef | emptyStatement }
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file
func (p *Parser) ParseProto() (*Proto, error) {
	syntaxComments := p.ParseComments()
	syntax, err := p.ParseSyntax()
	if err != nil {
		return nil, err
	}
	syntax.Comments = syntaxComments
	p.MaybeScanInlineComment(syntax)

	protoBody, err := p.parseProtoBody()
	if err != nil {
		return nil, err
	}

	return &Proto{
		Syntax:    syntax,
		ProtoBody: protoBody,
		Meta: &ProtoMeta{
			Filename: p.lex.Pos.Filename,
		},
	}, nil
}

// protoBody = { import | package | option | topLevelDef | emptyStatement }
// topLevelDef = message | enum | service | extend
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file
func (p *Parser) parseProtoBody() ([]Visitee, error) {
	var protoBody []Visitee

	for {
		comments := p.ParseComments()

		if p.IsEOF() {
			if p.bodyIncludingComments {
				for _, comment := range comments {
					protoBody = append(protoBody, Visitee(comment))
				}
			}
			return protoBody, nil
		}

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			HasInlineCommentSetter
			Visitee
		}

		switch token {
		case scanner.TIMPORT:
			importValue, err := p.ParseImport()
			if err != nil {
				return nil, err
			}
			importValue.Comments = comments
			stmt = importValue
		case scanner.TPACKAGE:
			packageValue, err := p.ParsePackage()
			if err != nil {
				return nil, err
			}
			packageValue.Comments = comments
			stmt = packageValue
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, err
			}
			message.Comments = comments
			stmt = message
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, err
			}
			enum.Comments = comments
			stmt = enum
		case scanner.TSERVICE:
			service, err := p.ParseService()
			if err != nil {
				return nil, err
			}
			service.Comments = comments
			stmt = service
		case scanner.TEXTEND:
			extend, err := p.ParseExtend()
			if err != nil {
				return nil, err
			}
			extend.Comments = comments
			stmt = extend
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, &EmptyStatement{})
		}

		p.MaybeScanInlineComment(stmt)
		protoBody = append(protoBody, stmt)
	}
}
