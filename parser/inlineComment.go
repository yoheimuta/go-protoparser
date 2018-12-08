package parser

// HasInlineCommentSetter requires to have a setter for an InlineComment field.
type HasInlineCommentSetter interface {
	SetInlineComment(comment *Comment)
}

// MaybeScanInlineComment tries to scan a comment on the current line. If present then set it with setter.
func (p *Parser) MaybeScanInlineComment(
	hasSetter HasInlineCommentSetter,
) {
	currentPos := p.lex.Pos

	comment, err := p.parseComment()
	if err != nil {
		return
	}

	firstCommentPos := comment.Meta.Pos
	if currentPos.Line != firstCommentPos.Line {
		p.lex.UnNext()
		return
	}

	hasSetter.SetInlineComment(comment)
	return
}
