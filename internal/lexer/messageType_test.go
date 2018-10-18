package lexer_test

import (
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestLexer2_ReadMessageType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantIsEOF bool
		wantErr   bool
	}{
		{
			name:      "ident",
			input:     "SearchRequest",
			wantText:  "SearchRequest",
			wantIsEOF: true,
		},
		{
			name:      ".ident",
			input:     ".SearchRequest",
			wantText:  ".SearchRequest",
			wantIsEOF: true,
		},
		{
			name:      ".ident.ident",
			input:     ".search.SearchRequest",
			wantText:  ".search.SearchRequest",
			wantIsEOF: true,
		},
		{
			name:      "ident.ident",
			input:     "aggregatespb.UserItemAggregate",
			wantText:  "aggregatespb.UserItemAggregate",
			wantIsEOF: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, err := lex.ReadMessageType()

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
