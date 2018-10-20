package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
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
}

// EmptyStatement represents ";".
type EmptyStatement struct{}

// Enum consists of a name and an enum body.
type Enum struct {
	EnumName string
	// EnumBody can have options and enum fields.
	// The element of this is the union of an option, enumField and emptyStatement.
	EnumBody []interface{}
}

// ParseEnum parses the enum.
// enum = "enum" enumName enumBody
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) ParseEnum() (*Enum, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TENUM {
		return nil, p.unexpected("enum")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("enumName")
	}
	enumName := p.lex.Text

	enumBody, err := p.parseEnumBody()
	if err != nil {
		return nil, err
	}

	return &Enum{
		EnumName: enumName,
		EnumBody: enumBody,
	}, nil
}

// enumBody = "{" { option | enumField | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) parseEnumBody() ([]interface{}, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	var stmts []interface{}

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
			stmts = append(stmts, option)
		default:
			enumField, enumFieldErr := p.parseEnumField()
			if enumFieldErr == nil {
				stmts = append(stmts, enumField)
				break
			}
			p.lex.UnNext()

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmts = append(stmts, EmptyStatement{})
				break
			}

			return nil, &parseEnumBodyStatementErr{
				parseEnumFieldErr:      enumFieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return stmts, nil
		}
		p.lex.UnNext()
	}
}

// enumField = ident "=" intLit [ "[" enumValueOption { ","  enumValueOption } "]" ]";"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#enum_definition
func (p *Parser) parseEnumField() (*EnumField, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("ident")
	}
	ident := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		return nil, p.unexpected("intLit")
	}
	number := p.lex.Text

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
		Number:           number,
		EnumValueOptions: enumValueOptions,
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

	constant, err := p.lex.ReadConstant()
	if err != nil {
		return nil, err
	}

	return &EnumValueOption{
		OptionName: optionName,
		Constant:   constant,
	}, nil
}
