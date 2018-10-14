package scanner_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

func TestScanner_Scan(t *testing.T) {
	type want struct {
		token scanner.Token
		text  string
		isErr bool
	}

	tests := []struct {
		name  string
		mode  scanner.Mode
		input string
		wants []want
	}{
		{
			name: "scan an empty string",
		},
		{
			name:  "skip whitespaces",
			input: "  ",
		},
		{
			name:  "scan idents",
			input: "service s1928 s_a 1ac-",
			wants: []want{
				{
					token: scanner.TIDENT,
					text:  "service",
				},
				{
					token: scanner.TIDENT,
					text:  "s1928",
				},
				{
					token: scanner.TIDENT,
					text:  "s_a",
				},
				{
					token: scanner.TILLEGAL,
					text:  "1",
				},
				{
					token: scanner.TIDENT,
					text:  "ac",
				},
				{
					token: scanner.TILLEGAL,
					text:  "-",
				},
			},
		},
		{
			name:  "scan boolLits",
			input: "true.false,talse",
			mode:  scanner.ScanBoolLit,
			wants: []want{
				{
					token: scanner.TBOOLLIT,
					text:  "true",
				},
				{
					token: scanner.TDOT,
					text:  ".",
				},
				{
					token: scanner.TBOOLLIT,
					text:  "false",
				},
				{
					token: scanner.TCOMMA,
					text:  ",",
				},
				{
					token: scanner.TIDENT,
					text:  "talse",
				},
			},
		},
		{
			name:  "scan keywords",
			input: "true service rpc",
			mode:  scanner.ScanKeyword,
			wants: []want{
				{
					token: scanner.TIDENT,
					text:  "true",
				},
				{
					token: scanner.TSERVICE,
					text:  "service",
				},
				{
					token: scanner.TRPC,
					text:  "rpc",
				},
			},
		},
		{
			name: "scan comments",
			input: `
// hogehoge
hogehoge
//
/*
fugafuga
*/
/**/
`,
			mode: scanner.ScanComment,
			wants: []want{
				{
					token: scanner.TCOMMENT,
					text:  "// hogehoge",
				},
				{
					token: scanner.TIDENT,
					text:  "hogehoge",
				},
				{
					token: scanner.TCOMMENT,
					text:  "//",
				},
				{
					token: scanner.TCOMMENT,
					text: `/*
fugafuga
*/`,
				},
				{
					token: scanner.TCOMMENT,
					text:  "/**/",
				},
			},
		},
		{
			name:  "scan strLits",
			input: `"" '' "abc" 'あいう' "\x1fzz" '\123\n\\'`,
			mode:  scanner.ScanStrLit,
			wants: []want{
				{
					token: scanner.TSTRLIT,
					text:  `""`,
				},
				{
					token: scanner.TSTRLIT,
					text:  `''`,
				},
				{
					token: scanner.TSTRLIT,
					text:  `"abc"`,
				},
				{
					token: scanner.TSTRLIT,
					text:  `'あいう'`,
				},
				{
					token: scanner.TSTRLIT,
					text:  `"\x1fzz"`,
				},
				{
					token: scanner.TSTRLIT,
					text:  `'\123\n\\'`,
				},
			},
		},
		{
			name:  "scan intLits",
			input: "1 10 9999 07 0123 0xf 0X123",
			mode:  scanner.ScanNumberLit,
			wants: []want{
				{
					token: scanner.TINTLIT,
					text:  "1",
				},
				{
					token: scanner.TINTLIT,
					text:  "10",
				},
				{
					token: scanner.TINTLIT,
					text:  "9999",
				},
				{
					token: scanner.TINTLIT,
					text:  "07",
				},
				{
					token: scanner.TINTLIT,
					text:  "0123",
				},
				{
					token: scanner.TINTLIT,
					text:  "0xf",
				},
				{
					token: scanner.TINTLIT,
					text:  "0X123",
				},
			},
		},
		{
			name:  "scan floatLits",
			input: "1.0 99.9 99.999 0.11 .101 1.234e5 1928e10 100.234E+15 1.234e-5 inf nan",
			mode:  scanner.ScanNumberLit,
			wants: []want{
				{
					token: scanner.TFLOATLIT,
					text:  "1.0",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "99.9",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "99.999",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "0.11",
				},
				{
					token: scanner.TFLOATLIT,
					text:  ".101",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1.234e5",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1928e10",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "100.234E+15",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1.234e-5",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "inf",
				},
				{
					token: scanner.TFLOATLIT,
					text:  "nan",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			s := scanner.NewScanner(strings.NewReader(test.input))
			s.Mode = test.mode

			for _, want := range test.wants {
				gtok, gtxt, gerr := s.Scan()
				if gtok != want.token {
					t.Errorf("got %v, but want %v", gtok, want.token)
				}
				if gtxt != want.text {
					t.Errorf("got %v, but want %v", gtxt, want.text)
				}
				switch {
				case want.isErr && gerr == nil:
					t.Errorf("got nil but want err")
					return
				case !want.isErr && gerr != nil:
					t.Errorf("got err %v, but want nil", gerr)
					return
				}
			}

			gtok, _, _ := s.Scan()
			if gtok != scanner.TEOF {
				t.Errorf("got %v, but want TEOF", gtok)
			}
		})
	}
}
