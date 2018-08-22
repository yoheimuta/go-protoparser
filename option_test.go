package protoparser

import (
	"reflect"
	"strings"
	"testing"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestParseOption(t *testing.T) {
	// TODO: Fix
	t.Skip()
	tests := []struct {
		name              string
		input             string
		wantOption        *Option
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing a normal Option",
			input: `
Option java_package = "com.example.foo";
            `,
			wantOption: &Option{
				Name:     "java_package",
				Constant: "com.example.foo",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, err := parseOption(lex)
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

			if !reflect.DeepEqual(got, test.wantOption) {
				t.Errorf("got %v, but want %v", got, test.wantOption)
			}
			if lex.Text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.Text(), test.wantRecentScanned)
			}
		})
	}
}
