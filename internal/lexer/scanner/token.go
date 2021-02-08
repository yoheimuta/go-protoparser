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

	// Keywords
	TSYNTAX
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
	TENUM
	TSTREAM
	TGROUP
)

func asMiscToken(ch rune) Token {
	m := map[rune]Token{
		';':  TSEMICOLON,
		':':  TCOLON,
		'=':  TEQUALS,
		'"':  TQUOTE,
		'\'': TQUOTE,
		'(':  TLEFTPAREN,
		')':  TRIGHTPAREN,
		'{':  TLEFTCURLY,
		'}':  TRIGHTCURLY,
		'[':  TLEFTSQUARE,
		']':  TRIGHTSQUARE,
		'<':  TLESS,
		'>':  TGREATER,
		',':  TCOMMA,
		'.':  TDOT,
		'-':  TMINUS,
	}
	if t, ok := m[ch]; ok {
		return t
	}
	return TILLEGAL
}

func asKeywordToken(st string) Token {
	m := map[string]Token{
		"syntax":     TSYNTAX,
		"service":    TSERVICE,
		"rpc":        TRPC,
		"returns":    TRETURNS,
		"message":    TMESSAGE,
		"extend":     TEXTEND,
		"import":     TIMPORT,
		"package":    TPACKAGE,
		"option":     TOPTION,
		"repeated":   TREPEATED,
		"required":   TREQUIRED,
		"optional":   TOPTIONAL,
		"weak":       TWEAK,
		"public":     TPUBLIC,
		"oneof":      TONEOF,
		"map":        TMAP,
		"reserved":   TRESERVED,
		"extensions": TEXTENSIONS,
		"enum":       TENUM,
		"stream":     TSTREAM,
		"group":      TGROUP,
	}

	if t, ok := m[st]; ok {
		return t
	}
	return TILLEGAL
}
