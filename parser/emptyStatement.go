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
