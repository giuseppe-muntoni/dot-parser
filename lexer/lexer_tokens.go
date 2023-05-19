package lexer

import (
	"dot-parser/result"
	"fmt"
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

func (token TokenData) Position() Position {
	return token.position
}
func (token TokenData) Token() Token {
	return token.token
}
func (token TokenData) Lexeme() Lexeme {
	return token.lexeme
}

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

func (token Token) String() string {
	switch token {
	case OPEN_BRACE:
		return "{"
	case CLOSE_BRACE:
		return "}"
	case SEMICOLON:
		return ";"
	case COLON:
		return ":"
	case COMMA:
		return ","
	case OPEN_SQUARE_BRACKET:
		return "["
	case CLOSE_SQUARE_BRACKET:
		return "]"
	case EQUAL:
		return "="
	case ARC:
		return "--"
	case DIRECTED_ARC:
		return "->"
	case ID:
		return "ID"
	case GRAPH:
		return "'graph'"
	case DIGRAPH:
		return "'digraph'"
	case STRICT:
		return "'strict'"
	case NODE:
		return "'node'"
	case EDGE:
		return "'edge'"
	case SUBGRAPH:
		return "'subgraph'"
	case EOF:
		return "EOF"
	default:
		panic(nil)
	}
}
