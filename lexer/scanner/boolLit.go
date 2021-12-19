package scanner

// boolLit = "true" | "false"
func isBoolLit(ident string) bool {
	switch ident {
	case "true", "false":
		return true
	default:
		return false
	}
}
