package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestLexer2_ReadFullIdent(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		permissive bool
		wantText   string
		wantIsEOF  bool
		wantErr    bool
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
			name:       "read { by permissive mode. Required by go-proto-validators",
			input:      "{int_gt: 0}",
			permissive: true,
			wantText:   "{int_gt:0}",
			wantIsEOF:  true,
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
			lex := lexer.NewLexer2(strings.NewReader(test.input), lexer.WithPermissive(test.permissive))
			got, err := lex.ReadFullIdent()

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

			lex.Next()
			if lex.IsEOF() != test.wantIsEOF {
				t.Errorf("got %v, but want %v", lex.IsEOF(), test.wantIsEOF)
			}
		})
	}
}
