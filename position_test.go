//go:build !windows
// +build !windows

package prompt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	istrings "github.com/plandex-ai/go-prompt/strings"
)

func TestPositionAtEndOfString(t *testing.T) {
	tests := map[string]struct {
		input   string
		columns istrings.Width
		want    Position
	}{
		"empty": {
			input:   "",
			columns: 20,
			want: Position{
				X: 0,
				Y: 0,
			},
		},
		"one letter": {
			input:   "f",
			columns: 20,
			want: Position{
				X: 1,
				Y: 0,
			},
		},
		"one word": {
			input:   "foo",
			columns: 20,
			want: Position{
				X: 3,
				Y: 0,
			},
		},
		"complex emoji": {
			input:   "🙆🏿‍♂️",
			columns: 20,
			want: Position{
				X: 2,
				Y: 0,
			},
		},
		"one-line fits in columns": {
			input:   "foo bar",
			columns: 20,
			want: Position{
				X: 7,
				Y: 0,
			},
		},
		"multiline": {
			input:   "foo\nbar\n",
			columns: 20,
			want: Position{
				X: 0,
				Y: 2,
			},
		},
		"one-line wrapping": {
			input:   "foobar",
			columns: 3,
			want: Position{
				X: 0,
				Y: 2,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := positionAtEndOfString(tc.input, tc.columns)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionAtEndOfStringLine(t *testing.T) {
	tests := map[string]struct {
		input string
		cols  istrings.Width
		line  int
		want  Position
	}{
		"last line overflows": {
			input: `hi
foobar`,
			cols: 3,
			line: 1,
			want: Position{
				X: 2,
				Y: 1,
			},
		},
		"last line is in the middle": {
			input: `hi
foo hey
bar boo ba
baz`,
			cols: 20,
			line: 2,
			want: Position{
				X: 10,
				Y: 2,
			},
		},
		"last line is out of bounds": {
			input: `hi
foo hey
bar boo ba
baz`,
			cols: 20,
			line: 20,
			want: Position{
				X: 3,
				Y: 3,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := positionAtEndOfStringLine(tc.input, tc.cols, tc.line)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionAdd(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: 3},
		},
		"only Y": {
			left:  Position{Y: 1},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 1},
			right: Position{Y: 2},
			want:  Position{X: 1, Y: 2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: 11, Y: 7},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Add(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionSubtract(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: -1},
		},
		"only Y": {
			left:  Position{Y: 5},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 1},
			right: Position{Y: 2},
			want:  Position{X: 1, Y: -2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: -9, Y: 3},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Subtract(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionJoin(t *testing.T) {
	tests := map[string]struct {
		left  Position
		right Position
		want  Position
	}{
		"empty": {
			left:  Position{},
			right: Position{},
			want:  Position{},
		},
		"only X": {
			left:  Position{X: 1},
			right: Position{X: 2},
			want:  Position{X: 3},
		},
		"only Y": {
			left:  Position{Y: 1},
			right: Position{Y: 2},
			want:  Position{Y: 3},
		},
		"different coordinates": {
			left:  Position{X: 5},
			right: Position{Y: 2},
			want:  Position{X: 0, Y: 2},
		},
		"both X and Y": {
			left:  Position{X: 1, Y: 5},
			right: Position{X: 10, Y: 2},
			want:  Position{X: 10, Y: 7},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
