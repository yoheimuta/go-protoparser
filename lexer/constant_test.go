package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
)

func TestLexer2_ReadConstant(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantIsEOF bool
		wantErr   bool
	}{
		{
			name:      "fullIdent",
			input:     "foo.bar",
			wantText:  "foo.bar",
			wantIsEOF: true,
		},
		{
			name:      "intLit",
			input:     "1928",
			wantText:  "1928",
			wantIsEOF: true,
		},
		{
			name:      "+intLit",
			input:     "+1928",
			wantText:  "+1928",
			wantIsEOF: true,
		},
		{
			name:      "-intLit",
			input:     "-1928",
			wantText:  "-1928",
			wantIsEOF: true,
		},
		{
			name:      "floatLit",
			input:     "1928.123",
			wantText:  "1928.123",
			wantIsEOF: true,
		},
		{
			name:      "+floatLit",
			input:     "+1928e10",
			wantText:  "+1928e10",
			wantIsEOF: true,
		},
		{
			name:      "-floatLit",
			input:     "-1928E-3",
			wantText:  "-1928E-3",
			wantIsEOF: true,
		},
		{
			name:      "single line strLit",
			input:     `"あいうえお''"`,
			wantText:  `"あいうえお''"`,
			wantIsEOF: true,
		},
		{
			name:      "multiline strLit with double quotes",
			input:     "\"line1 \"\n\"line2 \" \n\"line3\" ",
			wantText:  `"line1 line2 line3"`,
			wantIsEOF: true,
		},
		{
			name:      "multiline strLit with single quotes",
			input:     "'line1 '\n'line2 ' \n'line3' ",
			wantText:  `'line1 line2 line3'`,
			wantIsEOF: true,
		},
		{
			name:      "boolLit",
			input:     "true",
			wantText:  "true",
			wantIsEOF: true,
		},
		{
			name:     "boolLit.",
			input:    "false.",
			wantText: "false",
		},
		{
			name:    "ident.",
			input:   "rpc.",
			wantErr: true,
		},
		{
			name:    `left quote`,
			input:   `"`,
			wantErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, pos, err := lex.ReadConstant(true)

			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if got != test.wantText {
				t.Errorf("got %s, but want %s", got, test.wantText)
			}

			if pos.Offset != 0 {
				t.Errorf("got %d, but want 0", pos.Offset)
			}
			if pos.Line != 1 {
				t.Errorf("got %d, but want 1", pos.Line)
			}
			if pos.Column != 1 {
				t.Errorf("got %d, but want 1", pos.Column)
			}

			lex.Next()
			if lex.IsEOF() != test.wantIsEOF {
				t.Errorf("got %v, but want %v", lex.IsEOF(), test.wantIsEOF)
			}
		})
	}
}
