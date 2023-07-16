package prompt

import (
	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

/*

========
PROGRESS
========

Moving the cursor
-----------------

* [x] Ctrl + a   Go to the beginning of the line (Home)
* [x] Ctrl + e   Go to the End of the line (End)
* [x] Ctrl + p   Previous command (Up arrow)
* [x] Ctrl + n   Next command (Down arrow)
* [x] Ctrl + f   Forward one character
* [x] Ctrl + b   Backward one character
* [x] Ctrl + xx  Toggle between the start of line and current cursor position

Editing
-------

* [x] Ctrl + L   Clear the Screen, similar to the clear command
* [x] Ctrl + d   Delete character under the cursor
* [x] Ctrl + h   Delete character before the cursor (Backspace)

* [x] Ctrl + w   Cut the Word before the cursor to the clipboard.
* [x] Ctrl + k   Cut the Line after the cursor to the clipboard.
* [x] Ctrl + u   Cut/delete the Line before the cursor to the clipboard.

* [ ] Ctrl + t   Swap the last two characters before the cursor (typo).
* [ ] Esc  + t   Swap the last two words before the cursor.

* [ ] ctrl + y   Paste the last thing to be cut (yank)
* [ ] ctrl + _   Undo

*/

var emacsKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: ControlE,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorRight(istrings.RuneCount(buf.Document().CurrentLineAfterCursor()), cols, rows)
		},
	},
	// Go to the beginning of the line
	{
		Key: ControlA,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorLeft(buf.Document().FindStartOfFirstWordOfLine(), cols, rows)
		},
	},
	// Cut the Line after the cursor
	{
		Key: ControlK,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.Delete(istrings.RuneCount(buf.Document().CurrentLineAfterCursor()), cols, rows)
		},
	},
	// Cut/delete the Line before the cursor
	{
		Key: ControlU,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.DeleteBeforeCursor(istrings.RuneCount(buf.Document().CurrentLineBeforeCursor()), cols, rows)
		},
	},
	// Delete character under the cursor
	{
		Key: ControlD,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			if buf.Text() != "" {
				buf.Delete(1, cols, rows)
			}
		},
	},
	// Backspace
	{
		Key: ControlH,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.DeleteBeforeCursor(1, cols, rows)
		},
	},
	// Right allow: Forward one character
	{
		Key: ControlF,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorRight(1, cols, rows)
		},
	},
	// Alt Right allow: Forward one word
	{
		Key: AltRight,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorRight(
				buf.Document().FindRuneNumberUntilEndOfCurrentWord(),
				cols,
				rows,
			)
		},
	},
	// Left allow: Backward one character
	{
		Key: ControlB,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorLeft(1, cols, rows)
		},
	},
	// Alt Left allow: Backward one word
	{
		Key: AltLeft,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.CursorLeft(
				buf.Document().FindRuneNumberUntilStartOfPreviousWord(),
				cols,
				rows,
			)
		},
	},
	// Cut the Word before the cursor.
	{
		Key: ControlW,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.DeleteBeforeCursor(
				istrings.RuneCount(buf.Document().GetWordBeforeCursorWithSpace()),
				cols,
				rows,
			)
		},
	},
	{
		Key: AltBackspace,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			buf.DeleteBeforeCursor(
				istrings.RuneCount(buf.Document().GetWordBeforeCursorWithSpace()),
				cols,
				rows,
			)
		},
	},
	// Clear the Screen, similar to the clear command
	{
		Key: ControlL,
		Fn: func(buf *Buffer, cols istrings.Width, rows int) {
			consoleWriter.EraseScreen()
			consoleWriter.CursorGoTo(0, 0)
			debug.AssertNoError(consoleWriter.Flush())
		},
	},
}
