package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// Proto represents a protocol buffer definition.
type Proto struct {
	Syntax *Syntax
	// ProtoBody is a slice of sum type consisted of *Import, *Package, *Option, *Message, *Enum, *Service and *EmptyStatement.
	ProtoBody []interface{}
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

	protoBody, err := p.parseProtoBody()
	if err != nil {
		return nil, err
	}

	return &Proto{
		Syntax:    syntax,
		ProtoBody: protoBody,
	}, nil
}

// protoBody = { import | package | option | topLevelDef | emptyStatement }
// topLevelDef = message | enum | service
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#proto_file
func (p *Parser) parseProtoBody() ([]interface{}, error) {
	var protoBody []interface{}

	for {
		comments := p.ParseComments()

		if p.IsEOF() {
			return protoBody, nil
		}

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TIMPORT:
			importValue, err := p.ParseImport()
			if err != nil {
				return nil, err
			}
			importValue.Comments = comments
			protoBody = append(protoBody, importValue)
		case scanner.TPACKAGE:
			packageValue, err := p.ParsePackage()
			if err != nil {
				return nil, err
			}
			packageValue.Comments = comments
			protoBody = append(protoBody, packageValue)
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			option.Comments = comments
			protoBody = append(protoBody, option)
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, message)
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, enum)
		case scanner.TSERVICE:
			service, err := p.ParseService()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, service)
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, &EmptyStatement{})
		}
	}
}
