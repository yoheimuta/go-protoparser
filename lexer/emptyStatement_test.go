package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
)

func TestLexer2_ReadEmptyStatement(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "read ;",
			input: ";",
		},
		{
			name:    "not found ;",
			input:   ":",
			wantErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			err := lex.ReadEmptyStatement()

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

			lex.Next()
			if !lex.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
