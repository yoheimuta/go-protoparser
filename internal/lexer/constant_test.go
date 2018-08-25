package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestLexer2_ReadConstant(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantText string
		wantErr  bool
	}{
		{
			name:     "fullIdent",
			input:    "foo.bar",
			wantText: "foo.bar",
		},
		{
			name:     "intLit",
			input:    "1928",
			wantText: "1928",
		},
		{
			name:     "+intLit",
			input:    "+1928",
			wantText: "+1928",
		},
		{
			name:     "-intLit",
			input:    "-1928",
			wantText: "-1928",
		},
		{
			name:     "floatLit",
			input:    "1928.123",
			wantText: "1928.123",
		},
		{
			name:     "+floatLit",
			input:    "+1928e10",
			wantText: "+1928e10",
		},
		{
			name:     "-floatLit",
			input:    "-1928E-3",
			wantText: "-1928E-3",
		},
		{
			name:     "strLit",
			input:    `"あいうえお''"`,
			wantText: `"あいうえお''"`,
		},
		{
			name:     "boolLit",
			input:    "true",
			wantText: "true",
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
			lex := lexer.NewLexer2(strings.NewReader(test.input))
			got, err := lex.ReadConstant()

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
