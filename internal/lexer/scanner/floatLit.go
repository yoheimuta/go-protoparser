package scanner

// floatLit = ( decimals "." [ decimals ] [ exponent ] | decimals exponent | "."decimals [ exponent ] ) | "inf" | "nan"
func isFloatLitKeyword(ident string) bool {
	switch ident {
	case "inf", "nan":
		return true
	default:
		return false
	}
}
