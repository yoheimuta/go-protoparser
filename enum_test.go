package protoparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseEnum(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantEnum          *Enum
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parse an empty",
			wantErr: true,
		},
		{
			name: "parse a normal",
			input: `
enum EnumAllowingAlias {
  // UNKNOWN is an unknown state.
  UNKNOWN = 0;
  // STARTED is a started state.
  STARTED = 1;
  // RUNNING is a running state.
  RUNNING = 2 [(custom_option) = "hello world"];
}
            `,
			wantEnum: &Enum{
				Name: "EnumAllowingAlias",
				EnumFields: []*EnumField{
					{
						Name: "UNKNOWN",
						Comments: []string{
							"// UNKNOWN is an unknown state.",
						},
					},
					{
						Name: "STARTED",
						Comments: []string{
							"// STARTED is a started state.",
						},
					},
					{
						Name: "RUNNING",
						Comments: []string{
							"// RUNNING is a running state.",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := NewLexer(strings.NewReader(test.input))
			got, err := parseEnum(lex)
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

			if !reflect.DeepEqual(got, test.wantEnum) {
				t.Errorf("got %v, but want %v", got, test.wantEnum)
			}
			if lex.Text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.Text(), test.wantRecentScanned)
			}
		})
	}
}
