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

// Option is an option for scanner.NewScanner.
type Option func(*Scanner)

// WithFilename is an option to set filename to the pos.
func WithFilename(filename string) Option {
	return func(l *Scanner) {
		l.pos.Filename = filename
	}
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader, opts ...Option) *Scanner {
	s := &Scanner{
		r:   bufio.NewReader(r),
		pos: NewPosition(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
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

// UnScan put the last scanned text back to the read buffer.
func (s *Scanner) UnScan() Position {
	var reversedRunes []rune
	for _, ch := range s.lastScanRaw {
		reversedRunes = append([]rune{ch}, reversedRunes...)
	}
	for _, ch := range reversedRunes {
		s.unread(ch)
	}
	return *s.pos
}

// Scan returns the next token and text value.
func (s *Scanner) Scan() (Token, string, Position, error) {
	s.lastScanRaw = s.lastScanRaw[:0]
	return s.scan()
}

// LastScanRaw returns the deep-copied lastScanRaw.
func (s *Scanner) LastScanRaw() []rune {
	r := make([]rune, len(s.lastScanRaw))
	copy(r, s.lastScanRaw)
	return r
}

// SetLastScanRaw sets lastScanRaw to the given raw.
func (s *Scanner) SetLastScanRaw(raw []rune) {
	s.lastScanRaw = raw
}

func (s *Scanner) scan() (Token, string, Position, error) {
	ch := s.peek()

	startPos := *s.pos

	switch {
	case unicode.IsSpace(ch):
		s.read()
		return s.scan()
	case s.isEOF():
		return TEOF, "", startPos, nil
	case isLetter(ch), ch == '_':
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
