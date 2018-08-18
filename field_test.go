package protoparser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseField(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantField         *Field
		wantRecentScanned string
	}{
		{
			name:      "parsing an empty",
			wantField: &Field{},
		},
		{
			name:  "parsing a normal field",
			input: "foo.bar nested_message = 2;",
			wantField: &Field{
				Type: &Type{
					Name: "foo.bar",
				},
				Name: "nested_message",
			},
		},
		{
			name:  "parsing a normal field with repreated and a field option",
			input: "repeated int32 samples = 4 [packed=true];",
			wantField: &Field{
				Type: &Type{
					Name:       "int32",
					IsRepeated: true,
				},
				Name: "samples",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			lex := newlexer(strings.NewReader(test.input))
			got := parseField(lex)

			if !reflect.DeepEqual(got, test.wantField) {
				t.Errorf("got %v, but want %v", got, test.wantField)
			}
			if lex.text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.text(), test.wantRecentScanned)
			}
		})
	}
}
