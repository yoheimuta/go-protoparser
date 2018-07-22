package protoparser

import (
	"reflect"
	"testing"
)

func TestParseType(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantType          *Type
		wantRecentScanned string
	}{
		{
			name:     "parsing an empty creates an empty type",
			wantType: &Type{},
		},
		{
			name:  "parsing a normal type creates a type",
			input: `bytes binary = 2 [(validator.field) = {length_gt: 0}];`,
			wantType: &Type{
				Name: "bytes",
			},
			wantRecentScanned: "binary",
		},
		{
			name:  "parsing a normal type from other package creates a type",
			input: `entitiespb.UserItem item = 1 [(validator.field) = {msg_exists : true}];`,
			wantType: &Type{
				Name: "entitiespb.UserItem",
			},
			wantRecentScanned: "item",
		},
		{
			name:  "parsing a normal type from an inner of other package creates a type",
			input: `entitiespb.inner.UserItem item = 1 [(validator.field) = {msg_exists : true}];`,
			wantType: &Type{
				Name: "entitiespb.inner.UserItem",
			},
			wantRecentScanned: "item",
		},
	}

	for _, test := range tests {
		lex := lex(test.input)
		got := parseType(lex)
		if !reflect.DeepEqual(got, test.wantType) {
			t.Errorf("[%s] got %v, but want %v", test.name, got, test.wantType)
		}
		if lex.text() != test.wantRecentScanned {
			t.Errorf("[%s] got %v, but want %v", test.name, lex.text(), test.wantRecentScanned)
		}
	}
}
