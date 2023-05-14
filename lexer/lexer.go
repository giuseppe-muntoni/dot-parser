package lexer

import (
	"bufio"
	"errors"
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

func (lexer *Lexer) Lex() (Position, Token, Lexeme) {
	// keep looping until we return a token
	lexer.startPosition = lexer.currentPosition
	for {
		char, err := lexer.advance()

		if err != nil {
			if err == io.EOF {
				return lexer.currentPosition, EOF, ""
			}
			panic(err)
		}

		if char == '\n' {
			lexer.newLine()
			lexer.startPosition = lexer.currentPosition
			continue
		} else if unicode.IsSpace(char) {
			lexer.startPosition = lexer.currentPosition
			continue
		}

		switch char {
		// match single-char tokens
		case '{':
			return lexer.addToken(OPEN_BRACE)
		case '}':
			return lexer.addToken(CLOSE_BRACE)
		case ';':
			return lexer.addToken(SEMICOLON)
		case ':':
			return lexer.addToken(COLON)
		case ',':
			return lexer.addToken(COMMA)
		case '[':
			return lexer.addToken(OPEN_SQUARE_BRACKET)
		case ']':
			return lexer.addToken(CLOSE_SQUARE_BRACKET)
		case '=':
			return lexer.addToken(EQUAL)
		// match comments
		case '#':
			fallthrough
		case '/':
			lexer.matchComment(char)
			lexer.startPosition = lexer.currentPosition
		// identifiers
		case '-':
			nextChar, nextErr := lexer.peek()
			if nextErr != nil {
				panic(nextErr)
			} else if nextChar == '-' {
				lexer.advance()
				return lexer.addToken(ARC)
			} else if nextChar == '>' {
				lexer.advance()
				return lexer.addToken(DIRECTED_ARC)
			}
			fallthrough
		case '"':
			fallthrough
		default:
			return lexer.matchIdentifier(char)
		}
	}
}

func (lexer *Lexer) matchComment(firstChar rune) error {
	if firstChar == '#' {
		if lexer.currentPosition.column != 1 {
			return errors.New("")
		} else {
			lexer.skipLine()
			return nil
		}
	} else if firstChar == '/' {
		char, err := lexer.advance()
		if err != nil {
			return errors.New("")
		}
		if char == '/' {
			lexer.skipLine()
			return nil
		} else if char == '*' {
			lexer.skipMultiLineComment()
			return nil
		} else {
			return errors.New("")
		}
	} else {
		return errors.New("")
	}
}

func (lexer *Lexer) matchKeyword(ide string) (Position, Token, Lexeme) {
	token, exist := keywords[ide]
	if exist {
		return lexer.startPosition, token, ""
	} else {
		return lexer.startPosition, ID, Lexeme(ide)
	}
}

func (lexer *Lexer) addToken(token Token) (Position, Token, Lexeme) {
	return lexer.startPosition, token, ""
}

func (lexer *Lexer) advance() (char rune, err error) {
	char, _, err = lexer.reader.ReadRune()
	lexer.currentPosition.column += 1
	return
}

func (lexer *Lexer) peek() (char rune, err error) {
	char, _, err = lexer.reader.ReadRune()
	lexer.reader.UnreadRune()
	return
}

func (lexer *Lexer) skipLine() {
	for {
		char, err := lexer.advance()
		if err != nil {
			lexer.reader.UnreadRune()
			return
		}
		if char == '\n' {
			lexer.newLine()
			return
		}
	}
}

func (lexer *Lexer) skipMultiLineComment() {
	for {
		char, err := lexer.advance()
		if err != nil {
			lexer.reader.UnreadRune()
			return
		}
		if char == '\n' {
			lexer.newLine()
		} else if char == '*' {
			lexer.advance()
			char, err := lexer.advance()
			if err != nil {
				lexer.reader.UnreadRune()
				return
			}
			if char == '/' {
				return
			}
		}
	}
}

func (lexer *Lexer) newLine() {
	lexer.currentPosition.line += 1
	lexer.currentPosition.column = 0
}

func (lexer *Lexer) matchString() (Position, Token, Lexeme) {
	var lexeme string

	for {
		char, err := lexer.advance()
		if err != nil {
			lexer.reader.UnreadRune()
			break
		} else if char == '"' {
			break
		}
		lexeme += string(char)
	}

	return lexer.startPosition, ID, Lexeme(lexeme)
}

func (lexer *Lexer) matchAlphaNumeric(char rune) (Position, Token, Lexeme) {
	var lexeme = string(char)

	for {
		char, err := lexer.advance()
		if err != nil {
			lexer.reader.UnreadRune()
			break
		} else if char == '_' || unicode.IsDigit(char) || unicode.IsLetter(char) {
			lexeme += string(char)
		} else {
			lexer.reader.UnreadRune()
			break
		}
	}

	return lexer.matchKeyword(lexeme)
}

func (lexer *Lexer) matchNumeral(char rune) (Position, Token, Lexeme) {
	var lexeme = string(char)
	var canBeDot = true

	for {
		char, err := lexer.advance()
		if err != nil {
			lexer.reader.UnreadRune()
			break
		} else if char == '.' && canBeDot {
			canBeDot = false
			lexeme += string(char)
		} else if unicode.IsDigit(char) {
			lexeme += string(char)
		} else {
			lexer.reader.UnreadRune()
			break
		}
	}

	return lexer.startPosition, ID, Lexeme(lexeme)
}

func (lexer *Lexer) matchIdentifier(char rune) (Position, Token, Lexeme) {
	if char == '"' {
		return lexer.matchString()
	} else if unicode.IsDigit(char) || char == '-' || char == '.' {
		return lexer.matchNumeral(char)
	} else if unicode.IsLetter(char) || char == '_' {
		return lexer.matchAlphaNumeric(char)
	} else {
		panic(nil)
	}
}
