package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
)

func TestLexer2_ReadFullIdent(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantIsEOF bool
		wantErr   bool
	}{
		{
			name:      "ident",
			input:     "foo",
			wantText:  "foo",
			wantIsEOF: true,
		},
		{
			name:     "ident;",
			input:    "foo;",
			wantText: "foo",
		},
		{
			name:      "ident.ident",
			input:     "foo.true",
			wantText:  "foo.true",
			wantIsEOF: true,
		},
		{
			name:      "ident.ident.ident.ident",
			input:     "foo.bar.rpc.fuga",
			wantText:  "foo.bar.rpc.fuga",
			wantIsEOF: true,
		},
		{
			name:     "read invalid {.",
			input:    "{int_gt: 0}",
			wantText: "{int_gt:0}",
			wantErr:  true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "ident.",
			input:   "foo.",
			wantErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, pos, err := lex.ReadFullIdent()

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
