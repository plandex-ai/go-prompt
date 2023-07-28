package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
)

func TestEmacsKeyBindings(t *testing.T) {
	p := New(NoopExecutor)
	buf := p.buffer
	buf.InsertTextMoveCursor("abcde", DefColCount, DefRowCount, false)
	if buf.cursorPosition != istrings.RuneNumber(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}

	// Go to the beginning of the line
	applyEmacsKeyBind(p, ControlA)
	if buf.cursorPosition != 0 {
		t.Errorf("Want %d, but got %d", 0, buf.cursorPosition)
	}

	// Go to the end of the line
	applyEmacsKeyBind(p, ControlE)
	if buf.cursorPosition != istrings.RuneNumber(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}
}

func applyEmacsKeyBind(p *Prompt, key Key) {
	for i := range emacsKeyBindings {
		kb := emacsKeyBindings[i]
		if kb.Key == key {
			kb.Fn(p)
		}
	}
}
