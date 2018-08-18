package protoparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePackage(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantPackage       string
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name:        "parsing a normal package",
			input:       "package foo.bar;",
			wantPackage: "foo.bar",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := newlexer(strings.NewReader(test.input))
			got, err := parsePackage(lex)
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

			if !reflect.DeepEqual(got, test.wantPackage) {
				t.Errorf("got %v, but want %v", got, test.wantPackage)
			}
			if lex.text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.text(), test.wantRecentScanned)
			}
		})
	}
}
