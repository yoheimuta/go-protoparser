package protoparser

import (
	"reflect"
	"strings"
	"testing"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func TestParseService(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantService       *Service
		wantRecentScanned string
		wantErr           bool
	}{
		{
			name:    "parse an empty",
			wantErr: true,
		},
		{
			name: "parse normal service",
			input: `
service SearchService {
  // Search searches items.
  // See reference page in detail.
  rpc Search (SearchRequest) returns (SearchResponse);
  // SearchWithPrefix searches items with prefix.
  rpc SearchWithPrefix (SearchWithPrefixRequest) returns (SearchWithPrefixResponse);
}
            `,
			wantService: &Service{
				Name: "SearchService",
				RPCs: []*RPC{
					{
						Name: "Search",
						Argument: &Type{
							Name: "SearchRequest",
						},
						Return: &Type{
							Name: "SearchResponse",
						},
						Comments: []string{
							"// Search searches items.",
							"// See reference page in detail.",
						},
					},
					{
						Name: "SearchWithPrefix",
						Argument: &Type{
							Name: "SearchWithPrefixRequest",
						},
						Return: &Type{
							Name: "SearchWithPrefixResponse",
						},
						Comments: []string{
							"// SearchWithPrefix searches items with prefix.",
						},
					},
				},
			},
		},
		{
			name: "parse normal service with the emptyStatement Option",
			input: `
service SearchService {
  // Search searches items.
  // See reference page in detail.
  rpc Search (SearchRequest) returns (SearchResponse) {}
  // SearchWithPrefix searches items with prefix.
  rpc SearchWithPrefix (SearchWithPrefixRequest) returns (SearchWithPrefixResponse);
}
            `,
			wantService: &Service{
				Name: "SearchService",
				RPCs: []*RPC{
					{
						Name: "Search",
						Argument: &Type{
							Name: "SearchRequest",
						},
						Return: &Type{
							Name: "SearchResponse",
						},
						Comments: []string{
							"// Search searches items.",
							"// See reference page in detail.",
						},
					},
					{
						Name: "SearchWithPrefix",
						Argument: &Type{
							Name: "SearchWithPrefixRequest",
						},
						Return: &Type{
							Name: "SearchWithPrefixResponse",
						},
						Comments: []string{
							"// SearchWithPrefix searches items with prefix.",
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
			got, err := parseService(lex)
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

			if !reflect.DeepEqual(got, test.wantService) {
				t.Errorf("got %v, but want %v", got, test.wantService)
			}
			if lex.Text() != test.wantRecentScanned {
				t.Errorf("got %v, but want %v", lex.Text(), test.wantRecentScanned)
			}
		})
	}
}
