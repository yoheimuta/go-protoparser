package lexer

// ReadEnumType reads a messageType.
// enumType = [ "." ] { ident "." } enumName
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#identifiers
func (lex *Lexer) ReadEnumType() (string, error) {
	return lex.ReadMessageType()
}
