package strings_test

import (
	"fmt"
	"testing"

	"github.com/plandex-ai/go-prompt/strings"
)

func TestGetWidth(t *testing.T) {
	tests := []struct {
		in   string
		want strings.Width
	}{
		{
			in:   "foo",
			want: 3,
		},
		{
			in:   "🇵🇱",
			want: 2,
		},
		{
			in:   "🙆🏿‍♂️",
			want: 2,
		},
		{
			in:   "日本語",
			want: 6,
		},
	}

	for _, tc := range tests {
		if got := strings.GetWidth(tc.in); got != tc.want {
			t.Errorf("Should be %#v, but got %#v, for %#v", tc.want, got, tc.in)
		}
	}
}

func TestRuneIndexNthColumn(t *testing.T) {
	tests := []struct {
		text string
		n    strings.Width
		want strings.RuneNumber
	}{
		{
			text: "foo",
			n:    2,
			want: 2,
		},
		{
			text: "foo",
			n:    10,
			want: 3,
		},
		{
			text: "foo",
			n:    0,
			want: 0,
		},
		{
			text: "foo日本bar",
			n:    7,
			want: 5,
		},
		{
			text: "foo🇵🇱🙆🏿‍♂️bar",
			n:    7,
			want: 10,
		},
	}

	for _, tc := range tests {
		if got := strings.RuneIndexNthColumn(tc.text, tc.n); got != tc.want {
			t.Errorf("Should be %#v, but got %#v, for %#v", tc.want, got, tc.text)
		}
	}
}

func TestRuneIndexNthGrapheme(t *testing.T) {
	tests := []struct {
		text string
		n    strings.GraphemeNumber
		want strings.RuneNumber
	}{
		{
			text: "foo",
			n:    2,
			want: 2,
		},
		{
			text: "foo",
			n:    10,
			want: 3,
		},
		{
			text: "foo",
			n:    0,
			want: 0,
		},
		{
			text: "foo日本bar",
			n:    7,
			want: 7,
		},
		{
			text: "foo🇵🇱🙆🏿‍♂️bar",
			n:    7,
			want: 12,
		},
	}

	for _, tc := range tests {
		if got := strings.RuneIndexNthGrapheme(tc.text, tc.n); got != tc.want {
			t.Errorf("Should be %#v, but got %#v, for %#v", tc.want, got, tc.text)
		}
	}
}

func ExampleIndexNotByte() {
	fmt.Println(strings.IndexNotByte("golang", 'g'))
	fmt.Println(strings.IndexNotByte("golang", 'x'))
	fmt.Println(strings.IndexNotByte("gggggg", 'g'))
	// Output:
	// 1
	// 0
	// -1
}

func ExampleLastIndexNotByte() {
	fmt.Println(strings.LastIndexNotByte("golang", 'g'))
	fmt.Println(strings.LastIndexNotByte("golang", 'x'))
	fmt.Println(strings.LastIndexNotByte("gggggg", 'g'))
	// Output:
	// 4
	// 5
	// -1
}

func ExampleIndexNotAny() {
	fmt.Println(strings.IndexNotAny("golang", "glo"))
	fmt.Println(strings.IndexNotAny("golang", "gl"))
	fmt.Println(strings.IndexNotAny("golang", "golang"))
	// Output:
	// 3
	// 1
	// -1
}

func ExampleLastIndexNotAny() {
	fmt.Println(strings.LastIndexNotAny("golang", "agn"))
	fmt.Println(strings.LastIndexNotAny("golang", "an"))
	fmt.Println(strings.LastIndexNotAny("golang", "golang"))
	// Output:
	// 2
	// 5
	// -1
}
