package scanner

// Predefined mode bits to control recognition of tokens.
const (
	ScanIdent = 1 << iota
	ScanNumberLit
	ScanStrLit
	ScanBoolLit
	ScanKeyword
	ScanComment
	ScanLit = ScanNumberLit | ScanStrLit | ScanBoolLit
)
