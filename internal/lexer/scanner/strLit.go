package scanner

// strLit = ( "'" { charValue } "'" ) |  ( '"' { charValue } '"' )
func (s *Scanner) scanStrLit() (string, error) {
	quote := s.read()
	lit := string(quote)

	ch := s.peek()
	for ch != quote {
		cv, err := s.scanCharValue()
		if err != nil {
			return "", err
		}
		lit += cv
		ch = s.peek()
	}

	// consume quote
	lit += string(s.read())
	return lit, nil
}

// charValue = hexEscape | octEscape | charEscape | /[^\0\n\\]/
func (s *Scanner) scanCharValue() (string, error) {
	ch := s.peek()

	switch ch {
	case eof, '\n':
		return "", s.unexpected(ch, `/[^\0\n\\]`)
	case '\\':
		return s.tryScanEscape(), nil
	default:
		return string(s.read()), nil
	}
}

// hexEscape = '\' ( "x" | "X" ) hexDigit hexDigit
// octEscape = '\' octalDigit octalDigit octalDigit
// charEscape = '\' ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | '\' | "'" | '"' )
func (s *Scanner) tryScanEscape() string {
	lit := string(s.read())

	isCharEscape := func(r rune) bool {
		cs := []rune{'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '\'', '"'}
		for _, c := range cs {
			if r == c {
				return true
			}
		}
		return false
	}

	ch := s.peek()
	switch {
	case ch == 'x' || ch == 'X':
		lit += string(s.read())

		for i := 0; i < 2; i++ {
			if !isHexDigit(s.peek()) {
				return lit
			}
			lit += string(s.read())
		}
	case isOctalDigit(ch):
		for i := 0; i < 3; i++ {
			if !isOctalDigit(s.peek()) {
				return lit
			}
			lit += string(s.read())
		}
	case isCharEscape(ch):
		lit += string(s.read())
		return lit
	}
	return lit
}
