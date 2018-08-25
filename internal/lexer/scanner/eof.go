package scanner

func (s *Scanner) isEOF() bool {
	ch := s.peek()
	return ch == eof
}
