package parser

// HasInlineCommentSetter requires to have a setter for an InlineComment field.
type HasInlineCommentSetter interface {
	SetInlineComment(comment *Comment)
}

// MaybeScanInlineComment tries to scan a comment on the current line. If present then set it with setter.
func (p *Parser) MaybeScanInlineComment(
	hasSetter HasInlineCommentSetter,
) {
	inlineComment := p.parseInlineComment()
	if inlineComment == nil {
		return
	}
	hasSetter.SetInlineComment(inlineComment)
}

func (p *Parser) parseInlineComment() *Comment {
	currentPos := p.lex.Pos

	comment, err := p.parseComment()
	if err != nil {
		return nil
	}

	if currentPos.Line != comment.Meta.Pos.Line {
		p.lex.UnNext()
		return nil
	}

	return comment
}
