package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

type parseEnumBodyStatementErr struct {
	parseEnumFieldErr      error
	parseEmptyStatementErr error
}

func (e *parseEnumBodyStatementErr) Error() string {
	return fmt.Sprintf(
		"%v:%v",
		e.parseEnumFieldErr,
		e.parseEmptyStatementErr,
	)
}

// EnumValueOption is an option of a enumField.
type EnumValueOption struct {
	OptionName string
	Constant   string
}

// EnumField is a field of enum.
type EnumField struct {
	Ident            string
	Number           string
	EnumValueOptions []*EnumValueOption

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (f *EnumField) SetInlineComment(comment *Comment) {
	f.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (f *EnumField) Accept(v Visitor) {
	if !v.VisitEnumField(f) {
		return
	}

	for _, comment := range f.Comments {
		comment.Accept(v)
	}
	if f.InlineComment != nil {
		f.InlineComment.Accept(v)
	}
}

// Enum consists of a name and an enum body.
type Enum struct {
	EnumName string
	// EnumBody can have options and enum fields.
	// The element of this is the union of an option, enumField, reserved, and emptyStatement.
	EnumBody []Visitee

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
func (e *Enum) SetInlineComment(comment *Comment) {
	e.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (e *Enum) Accept(v Visitor) {
	if !v.VisitEnum(e) {
		return
	}

	for _, body := range e.EnumBody {
		body.Accept(v)
	}
	for _, comment := range e.Comments {
		comment.Accept(v)
	}
	if e.InlineComment != nil {
		e.InlineComment.Accept(v)
	}
	if e.InlineCommentBehindLeftCurly != nil {
		e.InlineCommentBehindLeftCurly.Accept(v)
	}
}

// ParseEnum parses the enum.
//  enum = "enum" enumName enumBody
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) ParseEnum() (*Enum, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TENUM {
		return nil, p.unexpected("enum")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("enumName")
	}
	enumName := p.lex.Text

	enumBody, inlineLeftCurly, lastPos, err := p.parseEnumBody()
	if err != nil {
		return nil, err
	}

	return &Enum{
		EnumName:                     enumName,
		EnumBody:                     enumBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

// enumBody = "{" { option | enumField | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) parseEnumBody() (
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
		case scanner.TRESERVED:
			// See https://developers.google.com/protocol-buffers/docs/proto3#enum_reserved
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			reserved.Comments = comments
			stmt = reserved
		default:
			enumField, enumFieldErr := p.parseEnumField()
			if enumFieldErr == nil {
				enumField.Comments = comments
				stmt = enumField
				break
			}
			p.lex.UnNext()

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &EmptyStatement{}
				break
			}

			return nil, nil, scanner.Position{}, &parseEnumBodyStatementErr{
				parseEnumFieldErr:      enumFieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}

// enumField = [ "-" ] ident "=" intLit [ "[" enumValueOption { ","  enumValueOption } "]" ]";"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) parseEnumField() (*EnumField, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("ident")
	}
	startPos := p.lex.Pos
	ident := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	var intLit string
	p.lex.ConsumeToken(scanner.TMINUS)
	if p.lex.Token == scanner.TMINUS {
		intLit = "-"
	}

	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		return nil, p.unexpected("intLit")
	}
	intLit += p.lex.Text

	enumValueOptions, err := p.parseEnumValueOptions()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &EnumField{
		Ident:            ident,
		Number:           intLit,
		EnumValueOptions: enumValueOptions,
		Meta:             meta.Meta{Pos: startPos.Position},
	}, nil
}

// enumValueOptions = "[" enumValueOption { ","  enumValueOption } "]"
func (p *Parser) parseEnumValueOptions() ([]*EnumValueOption, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTSQUARE {
		p.lex.UnNext()
		return nil, nil
	}

	opt, err := p.parseEnumValueOption()
	if err != nil {
		return nil, p.unexpected("enumValueOption")
	}

	var opts []*EnumValueOption
	opts = append(opts, opt)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		opt, err = p.parseEnumValueOption()
		if err != nil {
			return nil, p.unexpected("enumValueOption")
		}
		opts = append(opts, opt)
	}

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTSQUARE {
		return nil, p.unexpected("]")
	}
	return opts, nil
}

// enumValueOption = optionName "=" constant
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) parseEnumValueOption() (*EnumValueOption, error) {
	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	constant, err := p.parseOptionConstant()
	if err != nil {
		return nil, err
	}

	return &EnumValueOption{
		OptionName: optionName,
		Constant:   constant,
	}, nil
}
