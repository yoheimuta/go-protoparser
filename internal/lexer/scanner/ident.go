package scanner

// ident = letter { letter | decimalDigit | "_" }
func (s *Scanner) scanIdent() string {
	ident := string(s.read())

	for {
		next := s.peek()
		switch {
		case isLetter(next), isDecimalDigit(next), next == '_':
			ident += string(s.read())
		default:
			return ident
		}
	}
}
