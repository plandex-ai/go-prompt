package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// GoLineEnd Go to the End of the line
func GoLineEnd(buf *Buffer, cols istrings.Width, rows int) {
	x := []rune(buf.Document().TextAfterCursor())
	buf.CursorRight(istrings.RuneNumber(len(x)), cols, rows)
}

// GoLineBeginning Go to the beginning of the line
func GoLineBeginning(buf *Buffer, cols istrings.Width, rows int) {
	x := []rune(buf.Document().TextBeforeCursor())
	buf.CursorLeft(istrings.RuneNumber(len(x)), cols, rows)
}

// DeleteChar Delete character under the cursor
func DeleteChar(buf *Buffer, cols istrings.Width, rows int) {
	buf.Delete(1, cols, rows)
}

// DeleteBeforeChar Go to Backspace
func DeleteBeforeChar(buf *Buffer, cols istrings.Width, rows int) {
	buf.DeleteBeforeCursor(1, cols, rows)
}

// GoRightChar Forward one character
func GoRightChar(buf *Buffer, cols istrings.Width, rows int) {
	buf.CursorRight(1, cols, rows)
}

// GoLeftChar Backward one character
func GoLeftChar(buf *Buffer, cols istrings.Width, rows int) {
	buf.CursorLeft(1, cols, rows)
}
