package scanner_test

import (
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
)

func TestPosition_Advance(t *testing.T) {
	tests := []struct {
		name       string
		inputRunes []rune
		wantOffset int
		wantLine   int
		wantColumn int
	}{
		{
			name: "advance an ascii character",
			inputRunes: []rune{
				'a',
			},
			wantOffset: 1,
			wantLine:   1,
			wantColumn: 2,
		},
		{
			name: "advance ascii characters",
			inputRunes: []rune{
				'a',
				'b',
			},
			wantOffset: 2,
			wantLine:   1,
			wantColumn: 3,
		},
		{
			name: "advance an ascii character and a new line",
			inputRunes: []rune{
				'a',
				'\n',
			},
			wantOffset: 2,
			wantLine:   2,
			wantColumn: 1,
		},
		{
			name: "advance utf8 characters and a new line",
			inputRunes: []rune{
				'あ',
				'\n',
				'い',
			},
			wantOffset: 7,
			wantLine:   2,
			wantColumn: 2,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			pos := scanner.NewPosition()
			for _, r := range test.inputRunes {
				pos.Advance(r)
			}

			if pos.Offset != test.wantOffset {
				t.Errorf("got %d, but want %d", pos.Offset, test.wantOffset)
			}
			if pos.Line != test.wantLine {
				t.Errorf("got %d, but want %d", pos.Line, test.wantLine)
			}
			if pos.Column != test.wantColumn {
				t.Errorf("got %d, but want %d", pos.Column, test.wantColumn)
			}
		})
	}
}

func TestPosition_Revert(t *testing.T) {
	tests := []struct {
		name                string
		inputAdvancingRunes []rune
		inputRevertingRunes []rune
		wantOffset          int
		wantLine            int
		wantColumn          int
	}{
		{
			name: "advance and revert an ascii character",
			inputAdvancingRunes: []rune{
				'a',
			},
			inputRevertingRunes: []rune{
				'a',
			},
			wantOffset: 0,
			wantLine:   1,
			wantColumn: 1,
		},
		{
			name: "advance and revert ascii characters",
			inputAdvancingRunes: []rune{
				'a',
				'b',
			},
			inputRevertingRunes: []rune{
				'b',
				'a',
			},
			wantOffset: 0,
			wantLine:   1,
			wantColumn: 1,
		},
		{
			name: "advance and revert an ascii character and a new line",
			inputAdvancingRunes: []rune{
				'a',
				'\n',
			},
			inputRevertingRunes: []rune{
				'\n',
				'a',
			},
			wantOffset: 0,
			wantLine:   1,
			wantColumn: 1,
		},
		{
			name: "advance and revert utf8 characters and a new line",
			inputAdvancingRunes: []rune{
				'あ',
				'\n',
				'い',
			},
			inputRevertingRunes: []rune{
				'い',
				'\n',
				'あ',
			},
			wantOffset: 0,
			wantLine:   1,
			wantColumn: 1,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			pos := scanner.NewPosition()
			for _, r := range test.inputAdvancingRunes {
				pos.Advance(r)
			}
			for _, r := range test.inputRevertingRunes {
				pos.Revert(r)
			}

			if pos.Offset != test.wantOffset {
				t.Errorf("got %d, but want %d", pos.Offset, test.wantOffset)
			}
			if pos.Line != test.wantLine {
				t.Errorf("got %d, but want %d", pos.Line, test.wantLine)
			}
			if pos.Column != test.wantColumn {
				t.Errorf("got %d, but want %d", pos.Column, test.wantColumn)
			}
		})
	}
}
