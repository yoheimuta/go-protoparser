package protoparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseEnumField(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantEnumField     *EnumField
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parse an empty",
			wantErr: true,
		},
		{
			name: "parse a normal without an option",
			input: `
// UNKNOWN is the unknown state.
UNKNOWN = 0;
`,
			wantEnumField: &EnumField{
				Name: "UNKNOWN",
				Comments: []string{
					"// UNKNOWN is the unknown state.",
				},
			},
		},
		{
			name: "parse a normal with an enumValueOption",
			input: `
// RUNNING is the running state.
RUNNING = 2 [(custom_option) = "hello world"];
`,
			wantEnumField: &EnumField{
				Name: "RUNNING",
				Comments: []string{
					"// RUNNING is the running state.",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := newlexer(strings.NewReader(test.input))
			got, err := parseEnumField(lex)
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

			if !reflect.DeepEqual(got, test.wantEnumField) {
				t.Errorf("got %v, but want %v", got, test.wantEnumField)
			}
			if lex.text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.text(), test.wantRecentScanned)
			}
		})
	}
}
