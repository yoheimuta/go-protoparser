package unordered_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestInterpretEnum(t *testing.T) {
	tests := []struct {
		name      string
		inputEnum *parser.Enum
		wantEnum  *unordered.Enum
		wantErr   bool
	}{
		{
			name: "interpreting a nil",
		},
		{
			name: "interpreting an excerpt from the official reference with comments and reserved",
			inputEnum: &parser.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: []parser.Visitee{
					&parser.Option{
						OptionName: "allow_alias",
						Constant:   "true",
					},
					&parser.EnumField{
						Ident:  "UNKNOWN",
						Number: "0",
					},
					&parser.EnumField{
						Ident:  "STARTED",
						Number: "1",
					},
					&parser.EnumField{
						Ident:  "RUNNING",
						Number: "2",
						EnumValueOptions: []*parser.EnumValueOption{
							{
								OptionName: "(custom_option)",
								Constant:   `"hello world"`,
							},
						},
					},
					&parser.Reserved{
						FieldNames: []string{
							`"FOO"`,
							`"BAR"`,
						},
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// enum",
					},
				},
				InlineComment: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   2,
							Column: 26,
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation2",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 21,
						Line:   3,
						Column: 1,
					},
				},
			},
			wantEnum: &unordered.Enum{
				EnumName: "EnumAllowingAlias",
				EnumBody: &unordered.EnumBody{
					Options: []*parser.Option{
						{
							OptionName: "allow_alias",
							Constant:   "true",
						},
					},
					EnumFields: []*parser.EnumField{
						{
							Ident:  "UNKNOWN",
							Number: "0",
						},
						{
							Ident:  "STARTED",
							Number: "1",
						},
						{
							Ident:  "RUNNING",
							Number: "2",
							EnumValueOptions: []*parser.EnumValueOption{
								{
									OptionName: "(custom_option)",
									Constant:   `"hello world"`,
								},
							},
						},
					},
					Reserveds: []*parser.Reserved{
						{
							FieldNames: []string{
								`"FOO"`,
								`"BAR"`,
							},
						},
					},
				},
				Comments: []*parser.Comment{
					{
						Raw: "// enum",
					},
				},
				InlineComment: &parser.Comment{
					Raw: "// TODO: implementation",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   2,
							Column: 26,
						},
					},
				},
				InlineCommentBehindLeftCurly: &parser.Comment{
					Raw: "// TODO: implementation2",
					Meta: meta.Meta{
						Pos: meta.Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
					},
				},
				Meta: meta.Meta{
					Pos: meta.Position{
						Offset: 21,
						Line:   3,
						Column: 1,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := unordered.InterpretEnum(test.inputEnum)
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err, parsed=%v", got)
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantEnum) {
				t.Errorf("got %v, but want %v", got, test.wantEnum)
			}
		})
	}

}
