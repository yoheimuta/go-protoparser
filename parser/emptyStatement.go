package parser

// EmptyStatement represents ";".
type EmptyStatement struct {
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (e *EmptyStatement) SetInlineComment(comment *Comment) {
	e.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (e *EmptyStatement) Accept(v Visitor) {
	if !v.VisitEmptyStatement(e) {
		return
	}

	if e.InlineComment != nil {
		e.InlineComment.Accept(v)
	}
}
