package scanner

import (
	"bufio"
	"io"
	"unicode"
)

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r              *bufio.Reader
	lastReadBuffer []rune
	lastScanRaw    []rune

	// pos is a current source position.
	pos *Position

	// The Mode field controls which tokens are recognized.
	Mode Mode
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(r),
		pos: NewPosition(),
	}
}

func (s *Scanner) read() (r rune) {
	defer func() {
		if r == eof {
			return
		}
		s.lastScanRaw = append(s.lastScanRaw, r)

		s.pos.Advance(r)
	}()

	if 0 < len(s.lastReadBuffer) {
		var ch rune
		ch, s.lastReadBuffer = s.lastReadBuffer[len(s.lastReadBuffer)-1], s.lastReadBuffer[:len(s.lastReadBuffer)-1]
		return ch
	}
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread(ch rune) {
	s.lastReadBuffer = append(s.lastReadBuffer, ch)

	s.pos.Revert(ch)
}

func (s *Scanner) peek() rune {
	ch := s.read()
	if ch != eof {
		s.lastScanRaw = s.lastScanRaw[0 : len(s.lastScanRaw)-1]
		s.unread(ch)
	}
	return ch
}

// UnScan put the specified text back to the read buffer.
func (s *Scanner) UnScan() {
	var reversedRunes []rune
	for _, ch := range s.lastScanRaw {
		reversedRunes = append([]rune{ch}, reversedRunes...)
	}
	for _, ch := range reversedRunes {
		s.unread(ch)
	}
}

// Scan returns the next token and text value.
func (s *Scanner) Scan() (Token, string, Position, error) {
	s.lastScanRaw = s.lastScanRaw[:0]
	return s.scan()
}

func (s *Scanner) scan() (Token, string, Position, error) {
	ch := s.peek()

	startPos := *s.pos
	startPos.Advance(ch)

	switch {
	case unicode.IsSpace(ch):
		s.read()
		return s.scan()
	case s.isEOF():
		return TEOF, "", startPos, nil
	case isLetter(ch):
		ident := s.scanIdent()
		if s.Mode&ScanBoolLit != 0 && isBoolLit(ident) {
			return TBOOLLIT, ident, startPos, nil
		}
		if s.Mode&ScanNumberLit != 0 && isFloatLitKeyword(ident) {
			return TFLOATLIT, ident, startPos, nil
		}
		if s.Mode&ScanKeyword != 0 && asKeywordToken(ident) != TILLEGAL {
			return asKeywordToken(ident), ident, startPos, nil
		}
		return TIDENT, ident, startPos, nil
	case ch == '/':
		lit, err := s.scanComment()
		if err != nil {
			return TILLEGAL, "", startPos, err
		}
		if s.Mode&ScanComment != 0 {
			return TCOMMENT, lit, startPos, nil
		}
		return s.scan()
	case isQuote(ch) && s.Mode&ScanStrLit != 0:
		lit, err := s.scanStrLit()
		if err != nil {
			return TILLEGAL, "", startPos, err
		}
		return TSTRLIT, lit, startPos, nil
	case (isDecimalDigit(ch) || ch == '.') && s.Mode&ScanNumberLit != 0:
		tok, lit, err := s.scanNumberLit()
		if err != nil {
			return TILLEGAL, "", startPos, err
		}
		return tok, lit, startPos, nil
	default:
		return asMiscToken(ch), string(s.read()), startPos, nil
	}
}
