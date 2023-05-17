package lexer

import (
	"bufio"
	"dot-parser/iterator"
	"dot-parser/option"
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
		startPosition:   Position{line: 1, column: 1},
		currentPosition: Position{line: 1, column: 1},
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

func (lexer *Lexer) Next() result.Result[TokenData] {
	// keep looping until we return a token
	lexer.startPosition = lexer.currentPosition

	for {
		res_char := lexer.advance()
		if !res_char.IsOk() {
			err := res_char.UnwrapErr()
			if err == io.EOF {
				return lexer.makeTokenData(EOF, "")
			} else {
				return result.Err[TokenData](err)
			}
		}

		char := res_char.Unwrap()
		if unicode.IsSpace(char) || char == '\n' {
			lexer.startPosition = lexer.currentPosition
		} else {
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
			// match comments
			case '#':
				fallthrough
			case '/':
				commentMatched := lexer.matchComment(char)
				if commentMatched.IsOk() {
					lexer.startPosition = lexer.currentPosition
				} else {
					err := commentMatched.UnwrapErr().Error()
					return lexer.makeTokenError(err)
				}
			// identifiers
			case '-':
				return result.FlatMap(lexer.peek(), func(char rune) result.Result[TokenData] {
					switch char {
					case '-':
						lexer.advance()
						return lexer.makeTokenData(ARC, "")
					case '>':
						lexer.advance()
						return lexer.makeTokenData(DIRECTED_ARC, "")
					default:
						return lexer.matchIdentifier('-')
					}
				})
			case '"':
				fallthrough
			default:
				return lexer.matchIdentifier(char)
			}
		}
	}
}

func (lexer *Lexer) advance() result.Result[rune] {
	char, _, err := lexer.reader.ReadRune()
	res := result.Make(char, err)

	if res.IsOk() {
		if res.Unwrap() == '\n' {
			lexer.currentPosition.line += 1
			lexer.currentPosition.column = 1
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

// Lexer rune iterator
type lexerIterator struct {
	lexer *Lexer
}

func (iter *lexerIterator) HasNext() bool {
	return iter.lexer.peek().IsOk()
}

func (iter *lexerIterator) GetNext() option.Option[rune] {
	return result.ToOption(iter.lexer.advance())
}

func (lexer *Lexer) iter() iterator.Iterator[rune] {
	return &lexerIterator{lexer: lexer}
}

// comments matching
func (lexer *Lexer) matchComment(firstChar rune) result.Result[any] {
	switch firstChar {
	case '/':
		next := lexer.advance()
		return result.FlatMap(next, func(char rune) result.Result[any] {
			if char == '*' {
				lexer.skipMultiLineComment()
				return result.Ok[any](nil)
			} else if char == '/' {
				lexer.skipLine()
				return result.Ok[any](nil)
			} else {
				return result.Err[any](errors.New("invalid comment"))
			}
		})
	case '#':
		if lexer.startPosition.column == 1 {
			lexer.skipLine()
			return result.Ok[any](nil)
		}
		fallthrough
	default:
		return result.Err[any](errors.New("invalid comment"))
	}
}

func (lexer *Lexer) skipLine() {
	iter := iterator.TakeWhile(lexer.iter(), func(char rune) bool {
		return char != '\n'
	})

	iterator.Consume(iter)
}

func (lexer *Lexer) skipMultiLineComment() {
	var lastChar rune
	iter := iterator.TakeWhile(lexer.iter(), func(char rune) bool {
		res := lastChar == '*' && char == '/'
		lastChar = char
		return !res
	})

	iterator.Consume(iter)
}

// identifiers matching
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

func (lexer *Lexer) matchString() result.Result[TokenData] {
	iter := iterator.TakeWhile(lexer.iter(), func(char rune) bool {
		return char != '"'
	})

	lexeme := iterator.Fold("", iter, func(accum string, char rune) string {
		return accum + string(char)
	})

	return lexer.makeTokenData(ID, Lexeme(lexeme))
}

func (lexer *Lexer) matchAlphaNumeric(char rune) result.Result[TokenData] {
	iter := iterator.TakeWhile(lexer.iter(), func(char rune) bool {
		return char == '_' || unicode.IsDigit(char) || unicode.IsLetter(char)
	})

	lexeme := iterator.Fold(string(char), iter, func(accum string, char rune) string {
		return accum + string(char)
	})

	return lexer.matchKeyword(lexeme)
}

func (lexer *Lexer) matchKeyword(ide string) result.Result[TokenData] {
	token, exist := keywords[ide]
	if exist {
		return lexer.makeTokenData(token, "")
	} else {
		return lexer.makeTokenData(ID, Lexeme(ide))
	}
}

func (lexer *Lexer) matchNumeral(char rune) result.Result[TokenData] {
	var canBeDot = true
	iter := iterator.TakeWhile(lexer.iter(), func(char rune) bool {
		if char == '.' && canBeDot {
			canBeDot = false
			return true
		} else if unicode.IsDigit(char) {
			return true
		} else {
			return false
		}
	})

	lexeme := iterator.Fold(string(char), iter, func(accum string, char rune) string {
		return accum + string(char)
	})

	return lexer.makeTokenData(ID, Lexeme(lexeme))
}
