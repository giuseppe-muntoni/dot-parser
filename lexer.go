package dotparser

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
	position Position
	reader   *bufio.Reader
}

func New(reader io.Reader) *Lexer {
	return &Lexer{
		position: Position{line: 1, column: 1},
		reader:   bufio.NewReader(reader),
	}
}

func (lexer *Lexer) Lex() (Position, Token, Lexeme) {
	// keep looping until we return a token
	for {
		char, err := lexer.advance()

		if err != nil {
			if err == io.EOF {
				return lexer.position, EOF, ""
			}
			panic(err)
		}

		if unicode.IsSpace(char) {
			continue
		}

		switch char {
		// match single-char tokens
		case '{':
			lexer.addToken(OPEN_BRACE)
		case '}':
			lexer.addToken(CLOSE_BRACE)
		case ';':
			lexer.addToken(SEMICOLON)
		case ':':
			lexer.addToken(COLON)
		case ',':
			lexer.addToken(COMMA)
		case '[':
			lexer.addToken(OPEN_SQUARE_BRACKET)
		case ']':
			lexer.addToken(CLOSE_SQUARE_BRACKET)
		case '=':
			lexer.addToken(EQUAL)
		// match comments
		case '#':
			fallthrough
		case '/':
			lexer.matchComment(char)
		}
	}
}

func (lexer *Lexer) matchComment(firstChar rune) error {
	if firstChar == '#' {
		if lexer.position.column != 1 {
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

func matchKeyword(ide string) (Token, Lexeme) {
	token, exist := keywords[ide]
	if exist {
		return token, ""
	} else {
		return ID, Lexeme(ide)
	}
}

func (lexer *Lexer) addToken(token Token) (Position, Token, Lexeme) {
	return lexer.position, token, ""
}

func (lexer *Lexer) advance() (char rune, err error) {
	char, _, err = lexer.reader.ReadRune()
	lexer.position.column += 1
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
	lexer.position.line++
	lexer.position.column = 1
}
