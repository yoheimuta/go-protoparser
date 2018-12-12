package parser

// Visitee is implemented by all Protocol Buffer elements.
type Visitee interface {
	Accept(v Visitor)
}
