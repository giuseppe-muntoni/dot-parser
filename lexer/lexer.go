package lexer

import (
	"bufio"
	"dot-parser/iterator"
	"dot-parser/result"
	"errors"
	"io"
	"unicode"
)

type Position struct {
	line   int
	column int
}

type Lexer struct {
	iter            iterator.PeekableIterator[rune]
	startPosition   Position
	currentPosition Position
}

func New(reader io.Reader) *Lexer {
	lexer := &Lexer{
		iter:            nil,
		startPosition:   Position{line: 1, column: 1},
		currentPosition: Position{line: 1, column: 1},
	}

	iter := lexerIterator{
		reader:          bufio.NewReader(reader),
		currentPosition: &lexer.currentPosition,
	}

	lexer.iter = &iter

	return lexer
}

func (lexer *Lexer) Next() result.Result[TokenData] {
	lexer.startPosition = lexer.currentPosition

	for {
		res := lexer.iter.GetNext()

		if !res.IsSome() {
			return result.Err[TokenData](errors.New("IO error"))
		}

		char := res.Unwrap()

		if unicode.IsSpace(char) {
			lexer.startPosition = lexer.currentPosition
			continue
		}

		switch char {
		// match single-char tokens
		case '{':
			return lexer.makeTokenData(OPEN_BRACE, "")
		case '}':
			return lexer.makeTokenData(CLOSE_BRACE, "")
		case ';':
			return lexer.makeTokenData(SEMICOLON, "")
		case ':':
			return lexer.makeTokenData(COLON, "")
		case ',':
			return lexer.makeTokenData(COMMA, "")
		case '[':
			return lexer.makeTokenData(OPEN_SQUARE_BRACKET, "")
		case ']':
			return lexer.makeTokenData(CLOSE_SQUARE_BRACKET, "")
		case '=':
			return lexer.makeTokenData(EQUAL, "")
		case '\x03':
			return lexer.makeTokenData(EOF, "")
		// match comments
		case '#':
			fallthrough
		case '/':
			commentMatched := lexer.matchComment(char, lexer.iter)
			if commentMatched.IsOk() {
				lexer.iter = commentMatched.Unwrap()
				lexer.iter.GetNext()
				lexer.startPosition = lexer.currentPosition
			} else {
				err := commentMatched.UnwrapErr().Error()
				return lexer.makeTokenError(err)
			}
		// identifiers
		case '-':
			return result.FlatMap(result.FromOption(lexer.iter.Peek(), errors.New("IO error")), func(char rune) (res result.Result[TokenData]) {
				switch char {
				case '-':
					lexer.iter.GetNext()
					res = lexer.makeTokenData(ARC, "")
				case '>':
					lexer.iter.GetNext()
					res = lexer.makeTokenData(DIRECTED_ARC, "")
				default:
					res, lexer.iter = lexer.matchIdentifier('-', lexer.iter)
				}

				return res
			})
		case '"':
			fallthrough
		default:
			var res result.Result[TokenData] = nil
			res, lexer.iter = lexer.matchIdentifier(char, lexer.iter)
			return res
		}
	}
}
