package lexer

import (
	"bufio"
	"dot-parser/result"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type Token int
type Lexeme string

const (
	// Single-char tokens
	OPEN_BRACE Token = iota
	CLOSE_BRACE
	SEMICOLON
	COLON
	COMMA
	OPEN_SQUARE_BRACKET
	CLOSE_SQUARE_BRACKET
	EQUAL

	// Two-char tokens
	ARC
	DIRECTED_ARC

	// Literals
	ID

	// Keywords
	GRAPH
	DIGRAPH
	STRICT
	NODE
	EDGE
	SUBGRAPH

	EOF
)

var keywords = map[string]Token{
	"graph":    GRAPH,
	"digraph":  DIGRAPH,
	"strict":   STRICT,
	"node":     NODE,
	"edge":     EDGE,
	"subgraph": SUBGRAPH,
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	startPosition   Position
	currentPosition Position
	reader          *bufio.Reader
}

func New(reader io.Reader) *Lexer {
	return &Lexer{
		startPosition:   Position{line: 1, column: 0},
		currentPosition: Position{line: 1, column: 0},
		reader:          bufio.NewReader(reader),
	}
}

// Success
type TokenData struct {
	position Position
	token    Token
	lexeme   Lexeme
}

func (lexer *Lexer) makeTokenData(token Token, lexeme Lexeme) result.Result[TokenData] {
	return result.Ok(
		TokenData{
			position: lexer.startPosition,
			token:    token,
			lexeme:   lexeme,
		},
	)
}

// Error
type TokenError struct {
	position Position
	message  string
}

func (err *TokenError) Error() string {
	return fmt.Sprintf(
		"Lexing error at line %d column %d: %s",
		err.position.line,
		err.position.column,
		err.message)
}

func (lexer *Lexer) makeTokenError(message string) result.Result[TokenData] {
	return result.Err[TokenData](
		&TokenError{
			position: lexer.startPosition,
			message:  message,
		},
	)
}

// Iterator
type Iterator[T any] interface {
	hasNext() bool
	getNext() result.Result[T]
}

type TakeWhileIterator[T any] struct {
	iterator  Iterator[T]
	predicate func(T) bool
	buffer    result.Result[T]
}

func (iter *TakeWhileIterator[T]) hasNext() bool {
	if !iter.buffer.IsOk() {
		if !iter.iterator.hasNext() {
			return false
		} else {
			iter.buffer = iter.iterator.getNext()
		}
	}

	return iter.predicate(iter.buffer.Unwrap())
}

func (iter *TakeWhileIterator[T]) getNext() result.Result[T] {
	if !iter.hasNext() {
		return result.Err[T](nil)
	} else {
		res := iter.buffer
		iter.buffer = result.Err[T](nil)
		return res
	}
}

func takeWhile[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return &TakeWhileIterator[T]{
		iterator:  iter,
		predicate: predicate,
	}
}

func last[T any](iter Iterator[T]) result.Result[T] {
	for {
		res := iter.getNext()
		if !res.IsOk() {
			return res
		}
	}
}

type LexerIterator struct {
	lexer *Lexer
}

func (iter *LexerIterator) hasNext() bool {
	return iter.lexer.peek().IsOk()
}

func (iter *LexerIterator) getNext() result.Result[rune] {
	return iter.lexer.advance()
}

func (lexer *Lexer) Next() result.Result[TokenData] {
	// keep looping until we return a token
	lexer.startPosition = lexer.currentPosition

	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	res := result.Err[TokenData](errors.New(""))
	iter = takeWhile[rune](iter, func(char rune) bool {
		if unicode.IsSpace(char) || char == '\n' {
			return true
		} else {
			switch char {
			// match single-char tokens
			case '{':
				res = lexer.makeTokenData(OPEN_BRACE, "")
			case '}':
				res = lexer.makeTokenData(CLOSE_BRACE, "")
			case ';':
				res = lexer.makeTokenData(SEMICOLON, "")
			case ':':
				res = lexer.makeTokenData(COLON, "")
			case ',':
				res = lexer.makeTokenData(COMMA, "")
			case '[':
				res = lexer.makeTokenData(OPEN_SQUARE_BRACKET, "")
			case ']':
				res = lexer.makeTokenData(CLOSE_SQUARE_BRACKET, "")
			case '=':
				res = lexer.makeTokenData(EQUAL, "")
			// match comments
			case '#':
				fallthrough
			case '/':
				commentMatched := lexer.matchComment(char)
				if commentMatched.IsOk() {
					lexer.startPosition = lexer.currentPosition
				} else {
					_, err := commentMatched.Get()
					res = lexer.makeTokenError(err.Error())
				}
			// identifiers
			case '-':
				res = result.FlatMap(lexer.peek(), func(char rune) result.Result[TokenData] {
					switch char {
					case '-':
						lexer.advance()
						return lexer.makeTokenData(ARC, "")
					case '>':
						lexer.advance()
						return lexer.makeTokenData(DIRECTED_ARC, "")
					default:
						return lexer.matchIdentifier(char)
					}
				})
			case '"':
				fallthrough
			default:
				res = lexer.matchIdentifier(char)
			}
		}
		return res.IsOk()
	})

	last(iter)

	return res
}

func (lexer *Lexer) matchComment(firstChar rune) result.Result[rune] {
	switch firstChar {
	case '*':
		return lexer.skipMultiLineComment()
	case '/':
		return lexer.skipLine()
	case '#':
		if lexer.currentPosition.column == 1 {
			return lexer.skipLine()
		}
		fallthrough
	default:
		return result.Err[rune](errors.New("invalid comment"))
	}
}

func (lexer *Lexer) matchKeyword(ide string) result.Result[TokenData] {
	token, exist := keywords[ide]
	if exist {
		return lexer.makeTokenData(token, "")
	} else {
		return lexer.makeTokenData(ID, Lexeme(ide))
	}
}

func (lexer *Lexer) advance() result.Result[rune] {
	char, _, err := lexer.reader.ReadRune()
	res := result.Make(char, err)

	if res.IsOk() {
		if res.Unwrap() == '\n' {
			lexer.currentPosition.line += 1
			lexer.currentPosition.column = 0
			lexer.startPosition = lexer.currentPosition
		} else if unicode.IsSpace(res.Unwrap()) {
			lexer.startPosition = lexer.currentPosition
			lexer.currentPosition.column += 1
		} else {
			lexer.currentPosition.column += 1
		}
	}

	return res
}

func (lexer *Lexer) peek() result.Result[rune] {
	char, _, err := lexer.reader.ReadRune()
	res := result.Make(char, err)
	lexer.reader.UnreadRune()

	return res
}

func (lexer *Lexer) skipLine() result.Result[rune] {
	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	iter = takeWhile[rune](iter, func(char rune) bool {
		return char != '\n'
	})

	return last(iter)
}

func (lexer *Lexer) skipMultiLineComment() result.Result[rune] {
	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	var lastChar rune
	iter = takeWhile[rune](iter, func(char rune) bool {
		res := lastChar == '*' && char == '/'
		lastChar = char
		return !res
	})

	return last(iter)
}

func (lexer *Lexer) matchString() result.Result[TokenData] {
	var lexeme string

	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	iter = takeWhile[rune](iter, func(char rune) bool {
		res := char != '"'
		if res {
			lexeme += string(char)
		}

		return res
	})

	last(iter)

	return lexer.makeTokenData(ID, Lexeme(lexeme))
}

func (lexer *Lexer) matchAlphaNumeric(char rune) result.Result[TokenData] {
	var lexeme = string(char)

	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	iter = takeWhile[rune](iter, func(char rune) bool {
		res := char == '_' || unicode.IsDigit(char) || unicode.IsLetter(char)
		if res {
			lexeme += string(char)
		}

		return res
	})

	last(iter)

	return lexer.matchKeyword(lexeme)
}

func (lexer *Lexer) matchNumeral(char rune) result.Result[TokenData] {
	var lexeme = string(char)
	var canBeDot = true

	var iter Iterator[rune] = &LexerIterator{
		lexer: lexer,
	}

	iter = takeWhile[rune](iter, func(char rune) bool {
		if char == '.' && canBeDot {
			canBeDot = false
			lexeme += string(char)
			return true
		} else if unicode.IsDigit(char) {
			lexeme += string(char)
			return true
		} else {
			return false
		}
	})

	last(iter)

	return lexer.makeTokenData(ID, Lexeme(lexeme))
}

func (lexer *Lexer) matchIdentifier(char rune) result.Result[TokenData] {
	if char == '"' {
		return lexer.matchString()
	} else if unicode.IsDigit(char) || char == '-' || char == '.' {
		return lexer.matchNumeral(char)
	} else if unicode.IsLetter(char) || char == '_' {
		return lexer.matchAlphaNumeric(char)
	} else {
		return lexer.makeTokenError("invalid identifier")
	}
}
