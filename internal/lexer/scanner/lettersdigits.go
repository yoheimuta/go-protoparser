package scanner

// See
//  https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#letters_and_digits
//  https://ascii.cl/

// letter = "A" … "Z" | "a" … "z"
func isLetter(r rune) bool {
	if r < 'A' {
		return false
	}

	if r > 'z' {
		return false
	}

	if r > 'Z' && r < 'a' {
		return false
	}

	return true
}

// decimalDigit = "0" … "9"
func isDecimalDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// octalDigit   = "0" … "7"
func isOctalDigit(r rune) bool {
	return '0' <= r && r <= '7'
}

// hexDigit     = "0" … "9" | "A" … "F" | "a" … "f"
func isHexDigit(r rune) bool {
	if '0' <= r && r <= '9' {
		return true
	}
	if 'A' <= r && r <= 'F' {
		return true
	}
	if 'a' <= r && r <= 'f' {
		return true
	}
	return false
}
