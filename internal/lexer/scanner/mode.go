package scanner

// Mode is an enum type to control recognition of tokens.
type Mode uint

// Predefined mode bits to control recognition of tokens.
const (
	ScanIdent Mode = 1 << iota
	ScanNumberLit
	ScanStrLit
	ScanBoolLit
	ScanKeyword
	ScanComment
	ScanLit = ScanNumberLit | ScanStrLit | ScanBoolLit
)
