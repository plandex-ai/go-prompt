package prompt

import (
	"reflect"
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
)

func TestNewBuffer(t *testing.T) {
	b := NewBuffer()
	if b.workingIndex != 0 {
		t.Errorf("workingIndex should be %#v, got %#v", 0, b.workingIndex)
	}
	if !reflect.DeepEqual(b.workingLines, []string{""}) {
		t.Errorf("workingLines should be %#v, got %#v", []string{""}, b.workingLines)
	}
}

func TestBuffer_InsertText(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("some_text", DefColCount, DefRowCount, false)

	if b.Text() != "some_text" {
		t.Errorf("Text should be %#v, got %#v", "some_text", b.Text())
	}

	if b.cursorPosition != istrings.RuneCount("some_text") {
		t.Errorf("cursorPosition should be %#v, got %#v", istrings.RuneCount("some_text"), b.cursorPosition)
	}
}

func TestBuffer_InsertText_Overwrite(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("ABC", DefColCount, DefRowCount, false)

	if b.Text() != "ABC" {
		t.Errorf("Text should be %#v, got %#v", "ABC", b.Text())
	}

	if b.cursorPosition != istrings.RuneCount("ABC") {
		t.Errorf("cursorPosition should be %#v, got %#v", istrings.RuneCount("ABC"), b.cursorPosition)
	}

	b.CursorLeft(1, DefColCount, DefRowCount)
	// Replace C with DEF in ABC
	b.InsertTextMoveCursor("DEF", DefColCount, DefRowCount, true)

	if b.Text() != "ABDEF" {
		t.Errorf("Text should be %#v, got %#v", "ABDEF", b.Text())
	}

	b.CursorLeft(100, DefColCount, DefRowCount)
	// Replace ABD with GHI in ABDEF
	b.InsertTextMoveCursor("GHI", DefColCount, DefRowCount, true)

	if b.Text() != "GHIEF" {
		t.Errorf("Text should be %#v, got %#v", "GHIEF", b.Text())
	}

	b.CursorLeft(100, DefColCount, DefRowCount)
	// Replace GHI with J\nK in GHIEF
	b.InsertTextMoveCursor("J\nK", DefColCount, DefRowCount, true)

	if b.Text() != "J\nKEF" {
		t.Errorf("Text should be %#v, got %#v", "J\nKEF", b.Text())
	}

	b.CursorUp(100, DefColCount, DefRowCount)
	b.CursorLeft(100, DefColCount, DefRowCount)
	// Replace J with LMN in J\nKEF test end of line
	b.InsertTextMoveCursor("LMN", DefColCount, DefRowCount, true)

	if b.Text() != "LMN\nKEF" {
		t.Errorf("Text should be %#v, got %#v", "LMN\nKEF", b.Text())
	}
}

func TestBuffer_CursorMovement(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("some_text", DefColCount, DefRowCount, false)

	b.CursorLeft(1, DefColCount, DefRowCount)
	b.CursorLeft(2, DefColCount, DefRowCount)
	b.CursorRight(1, DefColCount, DefRowCount)
	b.InsertTextMoveCursor("A", DefColCount, DefRowCount, false)
	if b.Text() != "some_teAxt" {
		t.Errorf("Text should be %#v, got %#v", "some_teAxt", b.Text())
	}
	if b.cursorPosition != istrings.RuneCount("some_teA") {
		t.Errorf("Text should be %#v, got %#v", istrings.RuneCount("some_teA"), b.cursorPosition)
	}

	// Moving over left character counts.
	b.CursorLeft(100, DefColCount, DefRowCount)
	b.InsertTextMoveCursor("A", DefColCount, DefRowCount, false)
	if b.Text() != "Asome_teAxt" {
		t.Errorf("Text should be %#v, got %#v", "some_teAxt", b.Text())
	}
	if b.cursorPosition != istrings.RuneCount("A") {
		t.Errorf("Text should be %#v, got %#v", istrings.RuneCount("some_teA"), b.cursorPosition)
	}

	// TODO: Going right already at right end.
}

func TestBuffer_CursorMovement_WithMultiByte(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("あいうえお", DefColCount, DefRowCount, false)
	b.CursorLeft(1, DefColCount, DefRowCount)
	if l := b.Document().TextAfterCursor(); l != "お" {
		t.Errorf("Should be 'お', but got %s", l)
	}
	b.InsertTextMoveCursor("żółć", DefColCount, DefRowCount, true)
	if b.Text() != "あいうえżółć" {
		t.Errorf("Text should be %#v, got %#v", "あいうえżółć", b.Text())
	}
}

func TestBuffer_CursorUp(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("long line1\nline2", DefColCount, DefRowCount, false)
	b.CursorUp(1, DefColCount, DefRowCount)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}

	// Going up when already at the top.
	b.CursorUp(1, DefColCount, DefRowCount)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}

	// Going up to a line that's shorter.
	b.setDocument(&Document{}, DefColCount, DefRowCount)
	b.InsertTextMoveCursor("line1\nlong line2", DefColCount, DefRowCount, false)
	b.CursorUp(1, DefColCount, DefRowCount)
	if b.Document().cursorPosition != 5 {
		t.Errorf("Should be %#v, got %#v", 5, b.Document().cursorPosition)
	}
}

func TestBuffer_CursorDown(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("line1\nline2", DefColCount, DefRowCount, false)
	b.cursorPosition = 3

	// Normally going down
	b.CursorDown(1, DefColCount, DefRowCount)
	if b.Document().cursorPosition != istrings.RuneCount("line1\nlin") {
		t.Errorf("Should be %#v, got %#v", istrings.RuneCount("line1\nlin"), b.Document().cursorPosition)
	}

	// Going down to a line that's storter.
	b = NewBuffer()
	b.InsertTextMoveCursor("long line1\na\nb", DefColCount, DefRowCount, false)
	b.cursorPosition = 3
	b.CursorDown(1, DefColCount, DefRowCount)
	if b.Document().cursorPosition != istrings.RuneCount("long line1\na") {
		t.Errorf("Should be %#v, got %#v", istrings.RuneCount("long line1\na"), b.Document().cursorPosition)
	}
}

func TestBuffer_DeleteBeforeCursor(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("some_text", DefColCount, DefRowCount, false)
	b.CursorLeft(2, DefColCount, DefRowCount)
	deleted := b.DeleteBeforeCursor(1, DefColCount, DefRowCount)

	if b.Text() != "some_txt" {
		t.Errorf("Should be %#v, got %#v", "some_txt", b.Text())
	}
	if deleted != "e" {
		t.Errorf("Should be %#v, got %#v", deleted, "e")
	}
	if b.cursorPosition != istrings.RuneCount("some_t") {
		t.Errorf("Should be %#v, got %#v", istrings.RuneCount("some_t"), b.cursorPosition)
	}

	// Delete over the characters length before cursor.
	deleted = b.DeleteBeforeCursor(100, DefColCount, DefRowCount)
	if deleted != "some_t" {
		t.Errorf("Should be %#v, got %#v", "some_t", deleted)
	}
	if b.Text() != "xt" {
		t.Errorf("Should be %#v, got %#v", "xt", b.Text())
	}

	// If cursor position is a beginning of line, it has no effect.
	deleted = b.DeleteBeforeCursor(1, DefColCount, DefRowCount)
	if deleted != "" {
		t.Errorf("Should be empty, got %#v", deleted)
	}
}

func TestBuffer_NewLine(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("  hello", DefColCount, DefRowCount, false)
	b.NewLine(DefColCount, DefRowCount, false)
	ac := b.Text()
	ex := "  hello\n"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	b = NewBuffer()
	b.InsertTextMoveCursor("  hello", DefColCount, DefRowCount, false)
	b.NewLine(DefColCount, DefRowCount, true)
	ac = b.Text()
	ex = "  hello\n  "
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestBuffer_JoinNextLine(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("line1\nline2\nline3", DefColCount, DefRowCount, false)
	b.CursorUp(1, DefColCount, DefRowCount)
	b.JoinNextLine(" ", DefColCount, DefRowCount)

	ac := b.Text()
	ex := "line1\nline2 line3"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}

	// Test when there is no '\n' in the text
	b = NewBuffer()
	b.InsertTextMoveCursor("line1", DefColCount, DefRowCount, false)
	b.cursorPosition = 0
	b.JoinNextLine(" ", DefColCount, DefRowCount)
	ac = b.Text()
	ex = "line1"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}

func TestBuffer_SwapCharactersBeforeCursor(t *testing.T) {
	b := NewBuffer()
	b.InsertTextMoveCursor("hello world", DefColCount, DefRowCount, false)
	b.CursorLeft(2, DefColCount, DefRowCount)
	b.SwapCharactersBeforeCursor(DefColCount, DefRowCount)
	ac := b.Text()
	ex := "hello wrold"
	if ac != ex {
		t.Errorf("Should be %#v, got %#v", ex, ac)
	}
}
