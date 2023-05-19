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
