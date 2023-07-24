package prompt

import (
	istrings "github.com/elk-language/go-prompt/strings"
)

// GoLineEnd Go to the End of the line
func GoLineEnd(p *Prompt) bool {
	x := []rune(p.Buffer.Document().TextAfterCursor())
	return p.CursorRight(istrings.RuneNumber(len(x)))
}

// GoLineBeginning Go to the beginning of the line
func GoLineBeginning(p *Prompt) bool {
	x := []rune(p.Buffer.Document().TextBeforeCursor())
	return p.CursorLeft(istrings.RuneNumber(len(x)))
}

// DeleteChar Delete character under the cursor
func DeleteChar(p *Prompt) bool {
	p.Buffer.Delete(1, p.renderer.col, p.renderer.row)
	return true
}

// DeleteBeforeChar Go to Backspace
func DeleteBeforeChar(p *Prompt) bool {
	p.Buffer.DeleteBeforeCursor(1, p.renderer.col, p.renderer.row)
	return true
}

// GoRightChar Forward one character
func GoRightChar(p *Prompt) bool {
	return p.CursorRight(1)
}

// GoLeftChar Backward one character
func GoLeftChar(p *Prompt) bool {
	return p.CursorLeft(1)
}

func DeleteWordBeforeCursor(p *Prompt) bool {
	p.Buffer.DeleteBeforeCursor(
		istrings.RuneCount(p.Buffer.Document().GetWordBeforeCursorWithSpace()),
		p.renderer.col,
		p.renderer.row,
	)
	return true
}
