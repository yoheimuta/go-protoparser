package scanner

import (
	"bufio"
	"io"
	"unicode"
)

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader

	// The Mode field controls which tokens are recognized.
	Mode Mode
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) peek() rune {
	ch := s.read()
	s.unread()
	return ch
}

// Scan returns the next token and text value.
func (s *Scanner) Scan() (Token, string, error) {
	ch := s.peek()

	switch {
	case unicode.IsSpace(ch):
		s.read()
		return s.Scan()
	case s.isEOF():
		return TEOF, "", nil
	case isLetter(ch):
		ident := s.scanIdent()
		if s.Mode&ScanBoolLit != 0 && isBoolLit(ident) {
			return TBOOLLIT, ident, nil
		}
		if s.Mode&ScanNumberLit != 0 && isFloatLitKeyword(ident) {
			return TFLOATLIT, ident, nil
		}
		if s.Mode&ScanKeyword != 0 && asKeywordToken(ident) != TILLEGAL {
			return asKeywordToken(ident), ident, nil
		}
		return TIDENT, ident, nil
	case ch == '/':
		lit, err := s.scanComment()
		if err != nil {
			return TILLEGAL, "", err
		}
		if s.Mode&ScanComment != 0 {
			return TCOMMENT, lit, nil
		}
		return s.Scan()
	case isQuote(ch) && s.Mode&ScanStrLit != 0:
		lit, err := s.scanStrLit()
		if err != nil {
			return TILLEGAL, "", err
		}
		return TSTRLIT, lit, nil
	case (isDecimalDigit(ch) || ch == '.') && s.Mode&ScanNumberLit != 0:
		tok, lit, err := s.scanNumberLit()
		if err != nil {
			return TILLEGAL, "", err
		}
		return tok, lit, nil
	default:
		return asMiscToken(ch), string(s.read()), nil
	}
}
