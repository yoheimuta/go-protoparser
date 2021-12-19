package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
)

func TestLexer2_ReadEnumType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantIsEOF bool
		wantErr   bool
	}{
		{
			name:      "ident",
			input:     "EnumAllowingAlias",
			wantText:  "EnumAllowingAlias",
			wantIsEOF: true,
		},
		{
			name:      ".ident",
			input:     ".EnumAllowingAlias",
			wantText:  ".EnumAllowingAlias",
			wantIsEOF: true,
		},
		{
			name:      ".ident.ident",
			input:     ".search.EnumAllowingAlias",
			wantText:  ".search.EnumAllowingAlias",
			wantIsEOF: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, pos, err := lex.ReadEnumType()

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
