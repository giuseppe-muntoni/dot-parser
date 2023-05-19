package lexer

import (
	"bufio"
	"dot-parser/option"
	"io"
)

type lexerIterator struct {
	currentPosition *Position
	reader          *bufio.Reader
}

func (iter *lexerIterator) Next() option.Option[rune] {
	char, _, err := iter.reader.ReadRune()

	res := option.None[rune]()
	if err != nil {
		if err == io.EOF {
			char = '\x03'
			res = option.Some(char)
		}
	} else {
		res = option.Some(char)
	}

	if res.IsSome() {
		if res.Unwrap() == '\n' {
			iter.currentPosition.line += 1
			iter.currentPosition.column = 1
		} else {
			iter.currentPosition.column += 1
		}
	}

	return res
}

func (iter *lexerIterator) Peek() option.Option[rune] {
	char, _, err := iter.reader.ReadRune()
	iter.reader.UnreadRune()

	res := option.None[rune]()
	if err != nil {
		if err == io.EOF {
			char = '\x03'
			res = option.Some(char)
		}
	} else {
		res = option.Some(char)
	}

	return res
}
