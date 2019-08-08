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
		pos   scanner.Position
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
			input: "service s1928 s_a 1ac- _s_a",
			wants: []want{
				{
					token: scanner.TIDENT,
					text:  "service",
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "s1928",
					pos: scanner.Position{
						Offset: 8,
						Line:   1,
						Column: 9,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "s_a",
					pos: scanner.Position{
						Offset: 14,
						Line:   1,
						Column: 15,
					},
				},
				{
					token: scanner.TILLEGAL,
					text:  "1",
					pos: scanner.Position{
						Offset: 18,
						Line:   1,
						Column: 19,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "ac",
					pos: scanner.Position{
						Offset: 19,
						Line:   1,
						Column: 20,
					},
				},
				{
					token: scanner.TILLEGAL,
					text:  "-",
					pos: scanner.Position{
						Offset: 21,
						Line:   1,
						Column: 22,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "_s_a",
					pos: scanner.Position{
						Offset: 23,
						Line:   1,
						Column: 24,
					},
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
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TDOT,
					text:  ".",
					pos: scanner.Position{
						Offset: 4,
						Line:   1,
						Column: 5,
					},
				},
				{
					token: scanner.TBOOLLIT,
					text:  "false",
					pos: scanner.Position{
						Offset: 5,
						Line:   1,
						Column: 6,
					},
				},
				{
					token: scanner.TCOMMA,
					text:  ",",
					pos: scanner.Position{
						Offset: 10,
						Line:   1,
						Column: 11,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "talse",
					pos: scanner.Position{
						Offset: 11,
						Line:   1,
						Column: 12,
					},
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
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TSERVICE,
					text:  "service",
					pos: scanner.Position{
						Offset: 5,
						Line:   1,
						Column: 6,
					},
				},
				{
					token: scanner.TRPC,
					text:  "rpc",
					pos: scanner.Position{
						Offset: 13,
						Line:   1,
						Column: 14,
					},
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
					pos: scanner.Position{
						Offset: 1,
						Line:   2,
						Column: 1,
					},
				},
				{
					token: scanner.TIDENT,
					text:  "hogehoge",
					pos: scanner.Position{
						Offset: 13,
						Line:   3,
						Column: 1,
					},
				},
				{
					token: scanner.TCOMMENT,
					text:  "//",
					pos: scanner.Position{
						Offset: 22,
						Line:   4,
						Column: 1,
					},
				},
				{
					token: scanner.TCOMMENT,
					text: `/*
fugafuga
*/`,
					pos: scanner.Position{
						Offset: 25,
						Line:   5,
						Column: 1,
					},
				},
				{
					token: scanner.TCOMMENT,
					text:  "/**/",
					pos: scanner.Position{
						Offset: 40,
						Line:   8,
						Column: 1,
					},
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
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TSTRLIT,
					text:  `''`,
					pos: scanner.Position{
						Offset: 3,
						Line:   1,
						Column: 4,
					},
				},
				{
					token: scanner.TSTRLIT,
					text:  `"abc"`,
					pos: scanner.Position{
						Offset: 6,
						Line:   1,
						Column: 7,
					},
				},
				{
					token: scanner.TSTRLIT,
					text:  `'あいう'`,
					pos: scanner.Position{
						Offset: 12,
						Line:   1,
						Column: 13,
					},
				},
				{
					token: scanner.TSTRLIT,
					text:  `"\x1fzz"`,
					pos: scanner.Position{
						Offset: 24,
						Line:   1,
						Column: 19,
					},
				},
				{
					token: scanner.TSTRLIT,
					text:  `'\123\n\\'`,
					pos: scanner.Position{
						Offset: 33,
						Line:   1,
						Column: 28,
					},
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
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "10",
					pos: scanner.Position{
						Offset: 2,
						Line:   1,
						Column: 3,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "9999",
					pos: scanner.Position{
						Offset: 5,
						Line:   1,
						Column: 6,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "07",
					pos: scanner.Position{
						Offset: 10,
						Line:   1,
						Column: 11,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "0123",
					pos: scanner.Position{
						Offset: 13,
						Line:   1,
						Column: 14,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "0xf",
					pos: scanner.Position{
						Offset: 18,
						Line:   1,
						Column: 19,
					},
				},
				{
					token: scanner.TINTLIT,
					text:  "0X123",
					pos: scanner.Position{
						Offset: 22,
						Line:   1,
						Column: 23,
					},
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
					pos: scanner.Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "99.9",
					pos: scanner.Position{
						Offset: 4,
						Line:   1,
						Column: 5,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "99.999",
					pos: scanner.Position{
						Offset: 9,
						Line:   1,
						Column: 10,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "0.11",
					pos: scanner.Position{
						Offset: 16,
						Line:   1,
						Column: 17,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  ".101",
					pos: scanner.Position{
						Offset: 21,
						Line:   1,
						Column: 22,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1.234e5",
					pos: scanner.Position{
						Offset: 26,
						Line:   1,
						Column: 27,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1928e10",
					pos: scanner.Position{
						Offset: 34,
						Line:   1,
						Column: 35,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "100.234E+15",
					pos: scanner.Position{
						Offset: 42,
						Line:   1,
						Column: 43,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "1.234e-5",
					pos: scanner.Position{
						Offset: 54,
						Line:   1,
						Column: 55,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "inf",
					pos: scanner.Position{
						Offset: 63,
						Line:   1,
						Column: 64,
					},
				},
				{
					token: scanner.TFLOATLIT,
					text:  "nan",
					pos: scanner.Position{
						Offset: 67,
						Line:   1,
						Column: 68,
					},
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
				gtok, gtxt, gpos, gerr := s.Scan()
				if gtok != want.token {
					t.Errorf("got %v, but want %v", gtok, want.token)
				}
				if gtxt != want.text {
					t.Errorf("got %v, but want %v", gtxt, want.text)
				}
				if gpos.Offset != want.pos.Offset {
					t.Errorf("got %d, but want %d", gpos.Offset, want.pos.Offset)
				}
				if gpos.Line != want.pos.Line {
					t.Errorf("got %d, but want %d", gpos.Line, want.pos.Line)
				}
				if gpos.Column != want.pos.Column {
					t.Errorf("got %d, but want %d", gpos.Column, want.pos.Column)
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

			gtok, _, _, _ := s.Scan()
			if gtok != scanner.TEOF {
				t.Errorf("got %v, but want TEOF", gtok)
			}
		})
	}
}

func TestScanner_UnScan(t *testing.T) {
	tests := []struct {
		name  string
		mode  scanner.Mode
		input string
	}{
		{
			name:  "unscan ident",
			input: "service",
		},
		{
			name:  "unscan boolLit",
			mode:  scanner.ScanBoolLit,
			input: "true",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			s := scanner.NewScanner(strings.NewReader(test.input))
			s.Mode = test.mode
			token, text, pos, err := s.Scan()
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			s.UnScan()
			token2, text2, pos2, err := s.Scan()
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if token != token2 {
				t.Errorf("got %v, but want %v", token, token2)
			}
			if text != text2 {
				t.Errorf("got %v, but want %v", text, text2)
			}
			if pos.Offset != pos2.Offset {
				t.Errorf("got %d, but want %d", pos.Offset, pos2.Offset)
			}
			if pos.Line != pos2.Line {
				t.Errorf("got %d, but want %d", pos.Line, pos2.Line)
			}
			if pos.Column != pos2.Column {
				t.Errorf("got %d, but want %d", pos.Column, pos2.Column)
			}

			eof, _, _, _ := s.Scan()
			if eof != scanner.TEOF {
				t.Errorf("got %v, but want TEOF", eof)
			}
		})
	}
}
