package scanner

// comment = ( "//" { [^\n] } "\n" ) |  ( "/*" { any } "*/" )
func (s *Scanner) scanComment() (string, error) {
	lit := string(s.read())

	ch := s.read()
	switch ch {
	case '/':
		for ch != '\n' {
			lit += string(ch)

			if s.isEOF() {
				return lit, nil
			}
			ch = s.read()
		}
	case '*':
		for {
			if s.isEOF() {
				return lit, s.unexpected(eof, "\n")
			}
			lit += string(ch)

			ch = s.read()
			chn := s.peek()
			if ch == '*' && chn == '/' {
				lit += string(ch)
				lit += string(s.read())
				break
			}
		}
	default:
		return "", s.unexpected(ch, "/ or *")
	}

	return lit, nil
}
