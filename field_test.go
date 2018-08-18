package protoparser

import (
	"reflect"
	"testing"
)

func TestParseField(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantType          *Field
		wantRecentScanned string
	}{
		{
			name:     "parsing an empty",
			wantType: &Field{},
		},
		{
			name:  "parsing a normal field",
			input: "foo.bar nested_message = 2;",
			wantType: &Field{
				Type: &Type{
					Name: "foo.bar",
				},
				Name: "nested_message",
			},
		},
		{
			name:  "parsing a normal field with repreated and a field option",
			input: "repeated int32 samples = 4 [packed=true];",
			wantType: &Field{
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
			lex := lex(test.input)
			got := parseField(lex)

			if !reflect.DeepEqual(got, test.wantType) {
				t.Errorf("got %v, but want %v", got, test.wantType)
			}
			if lex.text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.text(), test.wantRecentScanned)
			}
		})
	}
}
