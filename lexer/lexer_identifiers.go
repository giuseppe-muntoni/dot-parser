package lexer

import (
	"dot-parser/iterator"
	"dot-parser/result"
	"unicode"
)

func (lexer *Lexer) matchIdentifier(char rune, iter iterator.PeekableIterator[rune]) (result.Result[TokenData], iterator.PeekableIterator[rune]) {
	if char == '"' {
		return lexer.matchString(iter)
	} else if unicode.IsDigit(char) || char == '-' || char == '.' {
		return lexer.matchNumeral(char, iter)
	} else if unicode.IsLetter(char) || char == '_' {
		return lexer.matchAlphaNumeric(char, iter)
	} else {
		return lexer.makeTokenError("invalid identifier"), iter
	}
}

func (lexer *Lexer) matchString(iter iterator.PeekableIterator[rune]) (result.Result[TokenData], iterator.PeekableIterator[rune]) {
	lexeme, iter := iterator.FoldWhile("", iter, func(accum string, char rune) (bool, string) {
		if char != '"' {
			return true, accum + string(char)
		} else {
			return false, accum
		}
	})

	return lexer.makeTokenData(ID, Lexeme(lexeme)), iter
}

func (lexer *Lexer) matchAlphaNumeric(char rune, iter iterator.PeekableIterator[rune]) (result.Result[TokenData], iterator.PeekableIterator[rune]) {
	lexeme, iter := iterator.FoldWhile(string(char), iter, func(accum string, char rune) (bool, string) {
		if char == '_' || unicode.IsDigit(char) || unicode.IsLetter(char) {
			return true, accum + string(char)
		} else {
			return false, accum
		}
	})

	return lexer.matchKeyword(lexeme), iter
}

func (lexer *Lexer) matchKeyword(ide string) result.Result[TokenData] {
	token, exist := keywords[ide]
	if exist {
		return lexer.makeTokenData(token, "")
	} else {
		return lexer.makeTokenData(ID, Lexeme(ide))
	}
}

func (lexer *Lexer) matchNumeral(char rune, iter iterator.PeekableIterator[rune]) (result.Result[TokenData], iterator.PeekableIterator[rune]) {
	var canBeDot = true
	lexeme, iter := iterator.FoldWhile(string(char), iter, func(accum string, char rune) (bool, string) {
		if char == '.' && canBeDot {
			canBeDot = false
			return true, accum + string(char)
		} else if unicode.IsDigit(char) {
			return true, accum + string(char)
		} else {
			return false, accum
		}
	})

	return lexer.makeTokenData(ID, Lexeme(lexeme)), iter
}
