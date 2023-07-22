package prompt

import (
	"strings"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

// Buffer emulates the console buffer.
type Buffer struct {
	workingLines   []string // The working lines. Similar to history
	workingIndex   int      // index of the current line
	startLine      int      // Line number of the first visible line in the terminal (0-indexed)
	cursorPosition istrings.RuneNumber
	cacheDocument  *Document
	lastKeyStroke  Key
}

// Text returns string of the current line.
func (b *Buffer) Text() string {
	return b.workingLines[b.workingIndex]
}

// Document method to return document instance from the current text and cursor position.
func (b *Buffer) Document() (d *Document) {
	if b.cacheDocument == nil ||
		b.cacheDocument.Text != b.Text() ||
		b.cacheDocument.cursorPosition != b.cursorPosition {
		b.cacheDocument = &Document{
			Text:           b.Text(),
			cursorPosition: b.cursorPosition,
		}
	}
	b.cacheDocument.lastKey = b.lastKeyStroke
	return b.cacheDocument
}

// DisplayCursorPosition returns the cursor position on rendered text on terminal emulators.
// So if Document is "日本(cursor)語", DisplayedCursorPosition returns 4 because '日' and '本' are double width characters.
func (b *Buffer) DisplayCursorPosition(columns istrings.Width) Position {
	return b.Document().DisplayCursorPosition(columns)
}

// Insert string into the buffer and move the cursor.
func (b *Buffer) InsertTextMoveCursor(text string, columns istrings.Width, rows int, overwrite bool) {
	b.insertText(text, columns, rows, overwrite, true)
}

// Insert string into the buffer without moving the cursor.
func (b *Buffer) InsertText(text string, overwrite bool) {
	b.insertText(text, 0, 0, overwrite, false)
}

// insertText insert string from current line.
func (b *Buffer) insertText(text string, columns istrings.Width, rows int, overwrite bool, moveCursor bool) {
	currentTextRunes := []rune(b.Text())
	cursor := b.cursorPosition

	if overwrite {
		overwritten := string(currentTextRunes[cursor:])
		if len(overwritten) >= int(cursor)+len(text) {
			overwritten = string(currentTextRunes[cursor : cursor+istrings.RuneCount(text)])
		}
		if i := strings.IndexAny(overwritten, "\n"); i != -1 {
			overwritten = overwritten[:i]
		}
		b.setText(
			string(currentTextRunes[:cursor])+text+string(currentTextRunes[cursor+istrings.RuneCount(overwritten):]),
			columns,
			rows,
		)
	} else {
		b.setText(
			string(currentTextRunes[:cursor])+text+string(currentTextRunes[cursor:]),
			columns,
			rows,
		)
	}

	if moveCursor {
		b.cursorPosition += istrings.RuneCount(text)
		b.RecalculateStartLine(columns, rows)
	}
}

func (b *Buffer) ResetStartLine() {
	b.startLine = 0
}

// Calculates the startLine once again and returns true when it's been changed.
func (b *Buffer) RecalculateStartLine(columns istrings.Width, rows int) bool {
	origStartLine := b.startLine
	pos := b.DisplayCursorPosition(columns)
	if pos.Y > b.startLine+rows-1 {
		b.startLine = pos.Y - rows + 1
	} else if pos.Y < b.startLine {
		b.startLine = pos.Y
	}

	if b.startLine < 0 {
		b.startLine = 0
	}
	return origStartLine != b.startLine
}

// SetText method to set text and update cursorPosition.
// (When doing this, make sure that the cursor_position is valid for this text.
// text/cursor_position should be consistent at any time, otherwise set a Document instead.)
func (b *Buffer) setText(text string, col istrings.Width, row int) {
	debug.Assert(b.cursorPosition <= istrings.RuneCount(text), "length of input should be shorter than cursor position")
	b.workingLines[b.workingIndex] = text
	b.RecalculateStartLine(col, row)
}

// Set cursor position. Return whether it changed.
func (b *Buffer) setCursorPosition(p istrings.RuneNumber) {
	if p > 0 {
		b.cursorPosition = p
	} else {
		b.cursorPosition = 0
	}
}

func (b *Buffer) setDocument(d *Document, columns istrings.Width, rows int) {
	b.cacheDocument = d
	b.setCursorPosition(d.cursorPosition) // Call before setText because setText check the relation between cursorPosition and line length.
	b.setText(d.Text, columns, rows)
	b.RecalculateStartLine(columns, rows)
}

// DeleteBeforeCursor delete specified number of characters before cursor and return the deleted text.
func (b *Buffer) DeleteBeforeCursor(count istrings.RuneNumber, columns istrings.Width, rows int) (deleted string) {
	debug.Assert(count >= 0, "count should be positive")
	r := []rune(b.Text())

	if b.cursorPosition > 0 {
		start := b.cursorPosition - count
		if start < 0 {
			start = 0
		}
		deleted = string(r[start:b.cursorPosition])
		b.setDocument(&Document{
			Text:           string(r[:start]) + string(r[b.cursorPosition:]),
			cursorPosition: b.cursorPosition - istrings.RuneNumber(len([]rune(deleted))),
		}, columns, rows)
	}
	return
}

// NewLine means CR.
func (b *Buffer) NewLine(columns istrings.Width, rows int, copyMargin bool) {
	if copyMargin {
		b.InsertTextMoveCursor("\n"+b.Document().leadingWhitespaceInCurrentLine(), columns, rows, false)
	} else {
		b.InsertTextMoveCursor("\n", columns, rows, false)
	}
}

// Delete specified number of characters and Return the deleted text.
func (b *Buffer) Delete(count istrings.RuneNumber, col istrings.Width, row int) string {
	r := []rune(b.Text())
	if b.cursorPosition < istrings.RuneNumber(len(r)) {
		textAfterCursor := b.Document().TextAfterCursor()
		textAfterCursorRunes := []rune(textAfterCursor)
		deletedRunes := textAfterCursorRunes[:count]
		b.setText(
			string(r[:b.cursorPosition])+string(r[b.cursorPosition+istrings.RuneNumber(len(deletedRunes)):]),
			col,
			row,
		)

		deleted := string(deletedRunes)
		return deleted
	}

	return ""
}

// JoinNextLine joins the next line to the current one by deleting the line ending after the current line.
func (b *Buffer) JoinNextLine(separator string, col istrings.Width, row int) {
	if !b.Document().OnLastLine() {
		b.cursorPosition += b.Document().GetEndOfLinePosition()
		b.Delete(1, col, row)
		// Remove spaces
		b.setText(
			b.Document().TextBeforeCursor()+separator+strings.TrimLeft(b.Document().TextAfterCursor(), " "),
			col,
			row,
		)
	}
}

// SwapCharactersBeforeCursor swaps the last two characters before the cursor.
func (b *Buffer) SwapCharactersBeforeCursor(col istrings.Width, row int) {
	if b.cursorPosition >= 2 {
		x := b.Text()[b.cursorPosition-2 : b.cursorPosition-1]
		y := b.Text()[b.cursorPosition-1 : b.cursorPosition]
		b.setText(
			b.Text()[:b.cursorPosition-2]+y+x+b.Text()[b.cursorPosition:],
			col,
			row,
		)
	}
}

// NewBuffer is constructor of Buffer struct.
func NewBuffer() (b *Buffer) {
	b = &Buffer{
		workingLines: []string{""},
		workingIndex: 0,
		startLine:    0,
	}
	return
}
