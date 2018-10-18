package lexer

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ReadMessageType reads a messageType.
// messageType = [ "." ] { ident "." } messageName
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#identifiers
func (lex *Lexer) ReadMessageType() (string, error) {
	lex.Next()

	var messageType string
	if lex.Token == scanner.TDOT {
		messageType = lex.Text
	} else {
		lex.UnNext()
	}

	lex.Next()
	for !lex.IsEOF() {
		if lex.Token != scanner.TIDENT {
			return "", lex.unexpected(lex.Text, "ident")
		}
		messageType += lex.Text

		lex.Next()
		if lex.Token != scanner.TDOT {
			lex.UnNext()
			break
		}
		messageType += lex.Text

		lex.Next()
	}

	return messageType, nil
}
