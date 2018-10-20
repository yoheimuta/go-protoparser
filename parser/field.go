package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// FieldOption is an option for the field.
type FieldOption struct {
	OptionName string
	Constant   string
}

// Field is a normal field that is the basic element of a protocol buffer message.
type Field struct {
	IsRepeated   bool
	Type         string
	FieldName    string
	FieldNumber  string
	FieldOptions []*FieldOption
}

// ParseField parses the field.
//  field = [ "repeated" ] type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#field
func (p *Parser) ParseField() (*Field, error) {
	var isRepeated bool
	p.lex.NextKeyword()
	if p.lex.Token == scanner.TREPEATED {
		isRepeated = true
	} else {
		p.lex.UnNext()
	}

	typeValue, err := p.parseType()
	if err != nil {
		return nil, p.unexpected("type")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("fieldName")
	}
	fieldName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}

	fieldOptions, err := p.parseFieldOptionsOption()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Field{
		IsRepeated:   isRepeated,
		Type:         typeValue,
		FieldName:    fieldName,
		FieldNumber:  fieldNumber,
		FieldOptions: fieldOptions,
	}, nil
}

// [ "[" fieldOptions "]" ]
func (p *Parser) parseFieldOptionsOption() ([]*FieldOption, error) {
	p.lex.Next()
	if p.lex.Token == scanner.TLEFTSQUARE {
		fieldOptions, err := p.parseFieldOptions()
		if err != nil {
			return nil, err
		}

		p.lex.Next()
		if p.lex.Token != scanner.TRIGHTSQUARE {
			return nil, p.unexpected("]")
		}
		return fieldOptions, nil
	}
	p.lex.UnNext()
	return nil, nil
}

// fieldOptions = fieldOption { ","  fieldOption }
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#field
func (p *Parser) parseFieldOptions() ([]*FieldOption, error) {
	opt, err := p.parseFieldOption()
	if err != nil {
		return nil, err
	}

	var opts []*FieldOption
	opts = append(opts, opt)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		opt, err = p.parseFieldOption()
		if err != nil {
			return nil, p.unexpected("fieldOption")
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// fieldOption = optionName "=" constant
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#field
func (p *Parser) parseFieldOption() (*FieldOption, error) {
	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	var constant string
	p.lex.Next()
	token := p.lex.Token
	p.lex.UnNext()
	switch token {
	// go-proto-validators requires this exceptions.
	case scanner.TLEFTCURLY:
		if !p.permissive {
			return nil, p.unexpected("constant or permissive mode")
		}

		constant, err = p.parseGoProtoValidatorFieldOptionConstant()
		if err != nil {
			return nil, err
		}
	default:
		constant, err = p.lex.ReadConstant()
		if err != nil {
			return nil, err
		}
	}

	return &FieldOption{
		OptionName: optionName,
		Constant:   constant,
	}, nil
}

// goProtoValidatorFieldOptionConstant = "{" ident ":" constant "}"
func (p *Parser) parseGoProtoValidatorFieldOptionConstant() (string, error) {
	var ret string

	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return "", p.unexpected("{")
	}
	ret += p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return "", p.unexpected("ident")
	}
	ret += p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TCOLON {
		return "", p.unexpected(":")
	}
	ret += p.lex.Text

	constant, err := p.lex.ReadConstant()
	if err != nil {
		return "", err
	}
	ret += constant

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTCURLY {
		return "", p.unexpected("}")
	}
	ret += p.lex.Text
	return ret, nil
}

var typeConstants = map[string]struct{}{
	"double":   {},
	"float":    {},
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
	"bytes":    {},
}

// type = "double" | "float" | "int32" | "int64" | "uint32" | "uint64"
//      | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64"
//      | "bool" | "string" | "bytes" | messageType | enumType
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#fields
func (p *Parser) parseType() (string, error) {
	p.lex.Next()
	if _, ok := typeConstants[p.lex.Text]; ok {
		return p.lex.Text, nil
	}
	p.lex.UnNext()

	messageOrEnumType, err := p.lex.ReadMessageType()
	if err != nil {
		return "", err
	}
	return messageOrEnumType, nil
}

// fieldNumber = intLit;
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#fields
func (p *Parser) parseFieldNumber() (string, error) {
	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		return "", p.unexpected("intLit")
	}
	return p.lex.Text, nil
}
