package lexer

import (
	"dot-parser/iterator"
	"dot-parser/result"
	"errors"
)

func (lexer *Lexer) matchComment(firstChar rune, iter iterator.PeekableIterator[rune]) result.Result[iterator.PeekableIterator[rune]] {
	switch firstChar {
	case '/':
		next := result.FromOption(iter.GetNext(), errors.New("invalid comment"))

		return result.FlatMap(next, func(char rune) result.Result[iterator.PeekableIterator[rune]] {
			if char == '*' {
				return result.Ok(lexer.skipMultiLineComment(iter))
			} else if char == '/' {
				return result.Ok(lexer.skipLine(iter))
			} else {
				return result.Err[iterator.PeekableIterator[rune]](errors.New("invalid comment"))
			}
		})
	case '#':
		if lexer.startPosition.column == 1 {
			return result.Ok(lexer.skipLine(iter))
		}
		fallthrough
	default:
		return result.Err[iterator.PeekableIterator[rune]](errors.New("invalid comment"))
	}
}

func (lexer *Lexer) skipLine(iter iterator.PeekableIterator[rune]) iterator.PeekableIterator[rune] {
	return iterator.SkipWhile(iter, func(char rune) bool {
		return char != '\n' && char != '\x03'
	})
}

func (lexer *Lexer) skipMultiLineComment(iter iterator.PeekableIterator[rune]) iterator.PeekableIterator[rune] {
	var lastChar rune
	return iterator.SkipWhile(iter, func(char rune) bool {
		res := lastChar == '*' && char == '/'
		lastChar = char
		return !res && char != '\x03'
	})
}
