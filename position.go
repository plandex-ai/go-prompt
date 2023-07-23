package prompt

import (
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/mattn/go-runewidth"
)

// Position stores the coordinates
// of a p.
//
// (0, 0) represents the top-left corner of the prompt,
// while (n, n) the bottom-right corner.
type Position struct {
	X istrings.Width
	Y int
}

// Join two positions and return a new position.
func (p Position) Join(other Position) Position {
	if other.Y == 0 {
		p.X += other.X
	} else {
		p.X = other.X
		p.Y += other.Y
	}
	return p
}

// Add two positions and return a new position.
func (p Position) Add(other Position) Position {
	return Position{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Subtract two positions and return a new position.
func (p Position) Subtract(other Position) Position {
	return Position{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

// positionAtEndOfString calculates the position of the
// p at the end of the given string.
func positionAtEndOfStringLine(str string, columns istrings.Width, line int) Position {
	return positionAtEndOfReaderLine(strings.NewReader(str), columns, line)
}

// positionAtEndOfString calculates the position
// at the end of the given string.
func positionAtEndOfString(str string, columns istrings.Width) Position {
	return positionAtEndOfReader(strings.NewReader(str), columns)
}

// Returns the index of the first character on the specified line (terminal row).
// If the line wraps because its contents are longer than the current columns in the terminal
// then the index of the first character of the first word of the specified line gets returned
// (the word may begin the line before or a few lines before and it's taken into consideration).
//
// The unique behaviour is intentional, this function has been designed for use in lexing.
// In order to improve performance the lexer only receives the visible part of the text.
// But the tokens could be incorrect if a token spanning multiple lines (because of wrapping)
// gets divided. This functions is meant to alleviate this effect.
func indexOfFirstTokenOnLine(input string, columns istrings.Width, line int) istrings.ByteNumber {
	if len(input) == 0 || line == 0 {
		return 0
	}

	str := input
	var indexOfWord istrings.ByteNumber
	var lastCharSize istrings.ByteNumber
	var down int
	var right istrings.Width
	var i istrings.ByteNumber

charLoop:
	for {
		char, size := utf8.DecodeRuneInString(str)
		i += lastCharSize
		if size == 0 {
			break charLoop
		}
		str = str[size:]
		lastCharSize = istrings.ByteNumber(size)

		switch char {
		case '\r':
			char, size := utf8.DecodeRuneInString(str)
			i += lastCharSize
			if size == 0 {
				break charLoop
			}
			str = str[size:]
			lastCharSize = istrings.ByteNumber(size)

			if char == '\n' {
				down++
				right = 0
				indexOfWord = i + 1
				if down >= line {
					break charLoop
				}
			}
		case '\n':
			down++
			right = 0
			indexOfWord = i + 1
			if down >= line {
				break charLoop
			}
		default:
			right += istrings.Width(runewidth.RuneWidth(char))
			if right > columns {
				right = istrings.Width(runewidth.RuneWidth(char))
				down++
				if down >= line {
					break charLoop
				}
			}
			if unicode.IsSpace(char) {
				indexOfWord = i + 1
			}
		}
	}

	if indexOfWord >= istrings.ByteNumber(len(input)) {
		return istrings.ByteNumber(len(input)) - 1
	}
	return indexOfWord
}

// positionAtEndOfReader calculates the position
// at the end of the given io.Reader.
func positionAtEndOfReader(reader io.RuneReader, columns istrings.Width) Position {
	var down int
	var right istrings.Width

charLoop:
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break charLoop
		}

		switch char {
		case '\r':
			char, _, err := reader.ReadRune()
			if err != nil {
				break charLoop
			}

			if char == '\n' {
				down++
				right = 0
			}
		case '\n':
			down++
			right = 0
		default:
			right += istrings.Width(runewidth.RuneWidth(char))
			if right == columns {
				right = 0
				down++
			}
		}
	}

	return Position{
		X: right,
		Y: down,
	}
}

// positionAtEndOfReaderLine calculates the position
// at the given line of the given io.Reader.
func positionAtEndOfReaderLine(reader io.RuneReader, columns istrings.Width, line int) Position {
	var down int
	var right istrings.Width

charLoop:
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break charLoop
		}

		switch char {
		case '\r':
			char, _, err := reader.ReadRune()
			if err != nil {
				break charLoop
			}

			if char == '\n' {
				if down == line {
					break charLoop
				}
				down++
				right = 0
			}
		case '\n':
			if down == line {
				break charLoop
			}
			down++
			right = 0
		default:
			right += istrings.Width(runewidth.RuneWidth(char))
			if right > columns {
				if down == line {
					right = columns - 1
					break charLoop
				}
				right = istrings.Width(runewidth.RuneWidth(char))
				down++
			}
		}
	}

	return Position{
		X: right,
		Y: down,
	}
}
