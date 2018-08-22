package protoparser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantMessage       *Message
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parse an empty",
			wantErr: true,
		},
		{
			name: "parse a normal message",
			input: `
message Outer {
  // Inner is an inner message.
  message Inner {
    // ival is an ival.
    int64 ival = 1;
  }
  // inner_message is an inner.
  repeated inner inner_message = 1;
}
            `,
			wantMessage: &Message{
				Name: "Outer",
				Fields: []*Field{
					{
						Name: "inner_message",
						Comments: []string{
							"// inner_message is an inner.",
						},
						Type: &Type{
							Name: "inner",
						},
						HasRepeated: true,
					},
				},
				Messages: []*Message{
					{
						Name: "Inner",
						Fields: []*Field{
							{
								Name: "ival",
								Comments: []string{
									"// ival is an ival.",
								},
								Type: &Type{
									Name: "int64",
								},
							},
						},
						Comments: []string{
							"// Inner is an inner message.",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.NewLexer(strings.NewReader(test.input))
			got, err := parseMessage(lex)
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

			if !reflect.DeepEqual(got, test.wantMessage) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantMessage))
			}
			if lex.Text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.Text(), test.wantRecentScanned)
			}
		})
	}
}
