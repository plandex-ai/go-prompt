package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEagerLexerNext(t *testing.T) {
	tests := map[string]struct {
		lexer *EagerLexer
		want  Token
		ok    bool
	}{
		"return the first token when at the beginning": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lastByteIndex: 0},
					&SimpleToken{lastByteIndex: 1},
				},
				currentIndex: 0,
			},
			want: &SimpleToken{lastByteIndex: 0},
			ok:   true,
		},
		"return the second token": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lastByteIndex: 3},
					&SimpleToken{lastByteIndex: 5},
					&SimpleToken{lastByteIndex: 6},
				},
				currentIndex: 1,
			},
			want: &SimpleToken{lastByteIndex: 5},
			ok:   true,
		},
		"return false when at the end": {
			lexer: &EagerLexer{
				tokens: []Token{
					&SimpleToken{lastByteIndex: 0},
					&SimpleToken{lastByteIndex: 4},
					&SimpleToken{lastByteIndex: 5},
				},
				currentIndex: 3,
			},
			want: nil,
			ok:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.lexer.Next()
			opts := []cmp.Option{
				cmp.AllowUnexported(SimpleToken{}, EagerLexer{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func charLex(s string) []Token {
	var result []Token
	for i := range s {
		result = append(result, NewSimpleToken(0, istrings.ByteNumber(i), istrings.ByteNumber(i)))
	}

	return result
}

func TestEagerLexerInit(t *testing.T) {
	tests := map[string]struct {
		lexer *EagerLexer
		input string
		want  *EagerLexer
	}{
		"reset the lexer's state": {
			lexer: &EagerLexer{
				lexFunc: charLex,
				tokens: []Token{
					&SimpleToken{firstByteIndex: 2, lastByteIndex: 2},
					&SimpleToken{firstByteIndex: 10, lastByteIndex: 10},
				},
				currentIndex: 11,
			},
			input: "foo",
			want: &EagerLexer{
				lexFunc: charLex,
				tokens: []Token{
					&SimpleToken{firstByteIndex: 0, lastByteIndex: 0},
					&SimpleToken{firstByteIndex: 1, lastByteIndex: 1},
					&SimpleToken{firstByteIndex: 2, lastByteIndex: 2},
				},
				currentIndex: 0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.lexer.Init(tc.input)
			opts := []cmp.Option{
				cmp.AllowUnexported(SimpleToken{}, EagerLexer{}),
				cmpopts.IgnoreFields(EagerLexer{}, "lexFunc"),
			}
			if diff := cmp.Diff(tc.want, tc.lexer, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
