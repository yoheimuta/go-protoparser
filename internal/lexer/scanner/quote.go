package scanner

// isQuote checks ch is the quote.
// quote = "'" | '"'
func isQuote(ch rune) bool {
	switch ch {
	case '\'', '"':
		return true
	default:
		return false
	}
}
