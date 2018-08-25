package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestLexer2_ReadFullIdent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantText string
		wantErr  bool
	}{
		{
			name:     "ident",
			input:    "foo",
			wantText: "foo",
		},
		{
			name:     "ident.ident",
			input:    "foo.true",
			wantText: "foo.true",
		},
		{
			name:     "ident.ident.ident.ident",
			input:    "foo.bar.rpc.fuga",
			wantText: "foo.bar.rpc.fuga",
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
			lex := lexer.NewLexer2(strings.NewReader(test.input))
			got, err := lex.ReadFullIdent()

			switch {
			case test.wantErr && err == nil:
				t.Errorf("got nil but want err")
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if got != test.wantText {
				t.Errorf("got %s, but want %s", got, test.wantText)
			}
		})
	}
}
