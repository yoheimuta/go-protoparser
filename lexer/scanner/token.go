package scanner

// Token represents a lexical token.
type Token int

// The result of Scan is one of these tokens.
const (
	// Special tokens
	TILLEGAL Token = iota
	TEOF

	// Identifiers
	TIDENT

	// Literals
	TINTLIT
	TFLOATLIT
	TBOOLLIT
	TSTRLIT

	// Comment
	TCOMMENT

	// Misc characters
	TSEMICOLON   // ;
	TCOLON       // :
	TEQUALS      // =
	TQUOTE       // " or '
	TLEFTPAREN   // (
	TRIGHTPAREN  // )
	TLEFTCURLY   // {
	TRIGHTCURLY  // }
	TLEFTSQUARE  // [
	TRIGHTSQUARE // ]
	TLESS        // <
	TGREATER     // >
	TCOMMA       // ,
	TDOT         // .
	TMINUS       // -
	TBOM         // Byte Order Mark

	// Keywords
	TSYNTAX
	TEDITION
	TSERVICE
	TRPC
	TRETURNS
	TMESSAGE
	TEXTEND
	TIMPORT
	TPACKAGE
	TOPTION
	TREPEATED
	TREQUIRED
	TOPTIONAL
	TWEAK
	TPUBLIC
	TONEOF
	TMAP
	TRESERVED
	TEXTENSIONS
	TDECLARATION
	TNUMBER
	TFULLNAME
	TTYPE
	TENUM
	TSTREAM
	TGROUP
)

func asMiscToken(ch rune) Token {
	m := map[rune]Token{
		';':      TSEMICOLON,
		':':      TCOLON,
		'=':      TEQUALS,
		'"':      TQUOTE,
		'\'':     TQUOTE,
		'(':      TLEFTPAREN,
		')':      TRIGHTPAREN,
		'{':      TLEFTCURLY,
		'}':      TRIGHTCURLY,
		'[':      TLEFTSQUARE,
		']':      TRIGHTSQUARE,
		'<':      TLESS,
		'>':      TGREATER,
		',':      TCOMMA,
		'.':      TDOT,
		'-':      TMINUS,
		'\uFEFF': TBOM,
	}
	if t, ok := m[ch]; ok {
		return t
	}
	return TILLEGAL
}

func asKeywordToken(st string) Token {
	m := map[string]Token{
		"syntax":      TSYNTAX,
		"edition":     TEDITION,
		"service":     TSERVICE,
		"rpc":         TRPC,
		"returns":     TRETURNS,
		"message":     TMESSAGE,
		"extend":      TEXTEND,
		"import":      TIMPORT,
		"package":     TPACKAGE,
		"option":      TOPTION,
		"repeated":    TREPEATED,
		"required":    TREQUIRED,
		"optional":    TOPTIONAL,
		"weak":        TWEAK,
		"public":      TPUBLIC,
		"oneof":       TONEOF,
		"map":         TMAP,
		"reserved":    TRESERVED,
		"extensions":  TEXTENSIONS,
		"number":      TNUMBER,
		"full_name":   TFULLNAME,
		"type":        TTYPE,
		"declaration": TDECLARATION,
		"enum":        TENUM,
		"stream":      TSTREAM,
		"group":       TGROUP,
	}

	if t, ok := m[st]; ok {
		return t
	}
	return TILLEGAL
}
