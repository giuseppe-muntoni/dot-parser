package dotparser

import (
	"bufio"
	"io"
)

type Token int

const (
	GRAPH Token = iota
	DIGRAPH
	STRICT
	OPEN_BRACE
	CLOSE_BRACE
	SEMICOLON
	COLON
	COMMA
	NODE
	EDGE
	OPEN_SQUARE_BRACKET
	CLOSE_SQUARE_BRACKET
	EQUAL
	SUBGRAPH
	EOF
)

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
		position: Position{line: 1, column: 0},
		reader:   bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}
	}
}

func (lexer *Lexer) advance() {
	lexer.position.column += 1
}

func (l *Lexer) newLine() {
	l.pos.line++
	l.pos.column = 1
}
