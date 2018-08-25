package scanner

// intLit     = decimalLit | octalLit | hexLit
// decimalLit = ( "1" … "9" ) { decimalDigit }
// octalLit   = "0" { octalDigit }
// hexLit     = "0" ( "x" | "X" ) hexDigit { hexDigit }
//
// floatLit = ( decimals "." [ decimals ] [ exponent ] | decimals exponent | "."decimals [ exponent ] ) | "inf" | "nan"
func (s *Scanner) scanNumberLit() (Token, string, error) {
	lit := string(s.read())
	ch := s.peek()

	switch {
	case lit == "0" && (ch == 'x' || ch == 'X'):
		// hexLit
		lit += string(s.read())
		if !isHexDigit(s.peek()) {
			return TILLEGAL, "", s.unexpected(s.peek(), "hexDigit")
		}
		lit += string(s.read())

		for !s.isEOF() {
			if !isHexDigit(s.peek()) {
				break
			}
			lit += string(s.read())
		}
		return TINTLIT, lit, nil
	case lit == ".":
		// floatLit
		fractional, err := s.scanFractionPartNoOmit()
		if err != nil {
			return TILLEGAL, "", err
		}
		return TFLOATLIT, lit + fractional, nil
	case ch == '.':
		// floatLit
		lit += string(s.read())
		fractional, err := s.scanFractionPart()
		if err != nil {
			return TILLEGAL, "", err
		}
		return TFLOATLIT, lit + fractional, nil
	case ch == 'e' || ch == 'E':
		// floatLit
		exp, err := s.scanExponent()
		if err != nil {
			return TILLEGAL, "", err
		}
		return TFLOATLIT, lit + exp, nil
	case lit == "0":
		// octalLit
		for !s.isEOF() {
			if !isOctalDigit(s.peek()) {
				break
			}
			lit += string(s.read())
		}
		return TINTLIT, lit, nil
	default:
		// decimalLit or floatLit
		for !s.isEOF() {
			if !isDecimalDigit(s.peek()) {
				break
			}
			lit += string(s.read())
		}

		switch s.peek() {
		case '.':
			// floatLit
			lit += string(s.read())
			fractional, err := s.scanFractionPart()
			if err != nil {
				return TILLEGAL, "", err
			}
			return TFLOATLIT, lit + fractional, nil
		case 'e', 'E':
			// floatLit
			exp, err := s.scanExponent()
			if err != nil {
				return TILLEGAL, "", err
			}
			return TFLOATLIT, lit + exp, nil
		default:
			// decimalLit
			return TINTLIT, lit, nil
		}
	}
}

// [ decimals ] [ exponent ]
func (s *Scanner) scanFractionPart() (string, error) {
	lit := ""

	ch := s.peek()
	switch {
	case isDecimalDigit(ch):
		decimals, err := s.scanDecimals()
		if err != nil {
			return "", err
		}
		lit += decimals
	}

	switch s.peek() {
	case 'e', 'E':
		exp, err := s.scanExponent()
		if err != nil {
			return "", err
		}
		lit += exp
	}
	return lit, nil
}

// decimals [ exponent ]
func (s *Scanner) scanFractionPartNoOmit() (string, error) {
	decimals, err := s.scanDecimals()
	if err != nil {
		return "", err
	}

	switch s.peek() {
	case 'e', 'E':
		exp, err := s.scanExponent()
		if err != nil {
			return "", err
		}
		return decimals + exp, nil
	default:
		return decimals, nil
	}
}

// exponent  = ( "e" | "E" ) [ "+" | "-" ] decimals
func (s *Scanner) scanExponent() (string, error) {
	ch := s.peek()
	switch ch {
	case 'e', 'E':
		lit := string(s.read())

		switch s.peek() {
		case '+', '-':
			lit += string(s.read())
		}
		decimals, err := s.scanDecimals()
		if err != nil {
			return "", err
		}
		return lit + decimals, nil
	default:
		return "", s.unexpected(ch, "e or E")
	}
}

// decimals  = decimalDigit { decimalDigit }
func (s *Scanner) scanDecimals() (string, error) {
	ch := s.peek()
	if !isDecimalDigit(ch) {
		return "", s.unexpected(ch, "decimalDigit")
	}
	lit := string(s.read())

	for !s.isEOF() {
		if !isDecimalDigit(s.peek()) {
			break
		}
		lit += string(s.read())
	}
	return lit, nil
}
