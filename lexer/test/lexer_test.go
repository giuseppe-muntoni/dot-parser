package test

import (
	"dot-parser/iterator"
	"dot-parser/lexer"
	"dot-parser/result"
	"strings"
	"testing"
)

func getLexer(input string) iterator.Iterator[result.Result[lexer.TokenData]] {
	return lexer.MakeLexer(strings.NewReader(input))
}

func printToken(t *testing.T, error string, position lexer.Position, token lexer.Token, lexeme lexer.Lexeme) {
	t.Errorf("%s; Got line %d, column %d, token %s, lexeme %s", error, position.Line(), position.Column(), token.String(), lexeme)
}

func testSingleToken(t *testing.T, input string, expectedToken lexer.Token, errorMsg string) {
	var lex = getLexer(input)

	res := lex.Next().Unwrap().Unwrap()
	if res.Token() != expectedToken {
		printToken(t, errorMsg, res.Position(), res.Token(), res.Lexeme())
	}
}

func testIdToken(t *testing.T, input string, expectedLexeme string, errorMsg string) {
	var lex = getLexer(input)

	res := lex.Next().Unwrap().Unwrap()
	if res.Token() != lexer.ID || res.Lexeme() != lexer.Lexeme(expectedLexeme) {
		printToken(t, errorMsg, res.Position(), res.Token(), res.Lexeme())
	}
}

func TestEOF(t *testing.T) {
	testSingleToken(t, "", lexer.EOF, "Expected EOF")
}

// Test single character tokens
func TestOpenBrace(t *testing.T) {
	testSingleToken(t, "{", lexer.OPEN_BRACE, "Expected Open Brace")
}

func TestClosedBrace(t *testing.T) {
	testSingleToken(t, "}", lexer.CLOSE_BRACE, "Expected Closed Brace")
}

func TestSemicolon(t *testing.T) {
	testSingleToken(t, ";", lexer.SEMICOLON, "Expected Semicolon")
}

func TestColon(t *testing.T) {
	testSingleToken(t, ":", lexer.COLON, "Expected Colon")
}

func TestOpenSquareBracket(t *testing.T) {
	testSingleToken(t, "[", lexer.OPEN_SQUARE_BRACKET, "Expected Open Square Bracket")
}

func TestClosedSquareBracket(t *testing.T) {
	testSingleToken(t, "]", lexer.CLOSE_SQUARE_BRACKET, "Expected Closed Square Bracket")
}

func TestEqual(t *testing.T) {
	testSingleToken(t, "=", lexer.EQUAL, "Expected Equal")
}

func TestComma(t *testing.T) {
	testSingleToken(t, ",", lexer.COMMA, "Expected Comma")
}

// Test two-character tokens
func TestArc(t *testing.T) {
	testSingleToken(t, "--", lexer.ARC, "Expected Arc (--)")
}

func TestDirectedArc(t *testing.T) {
	testSingleToken(t, "->", lexer.DIRECTED_ARC, "Expected Directed Arc (->)")
}

// Test keywords
func TestKeywordGraph(t *testing.T) {
	testSingleToken(t, "graph", lexer.GRAPH, "Expected Keyword 'graph'")
}

func TestKeywordDigraph(t *testing.T) {
	testSingleToken(t, "digraph", lexer.DIGRAPH, "Expected Keyword 'digraph'")
}

func TestKeywordStrict(t *testing.T) {
	testSingleToken(t, "strict", lexer.STRICT, "Expected Keyword 'strict'")
}

func TestKeywordNode(t *testing.T) {
	testSingleToken(t, "node", lexer.NODE, "Expected Keyword 'node'")
}

func TestKeywordEdge(t *testing.T) {
	testSingleToken(t, "edge", lexer.EDGE, "Expected Keyword 'edge'")
}

func TestKeywordSubgraph(t *testing.T) {
	testSingleToken(t, "subgraph", lexer.SUBGRAPH, "Expected Keyword 'subgraph'")
}

// Test whitespaces and comments
func TestIgnoreWhitespaces(t *testing.T) {
	testSingleToken(t, " \n\t\r  ", lexer.EOF, "Expected EOF after ignoring whitespaces")
}

func TestIgnoreSingleLineComment(t *testing.T) {
	testSingleToken(t, "//comment", lexer.EOF, "Expected EOF after ignoring single line comment")
}

func TestIgnoreSingleLineComment2(t *testing.T) {
	testSingleToken(t, "#comment", lexer.EOF, "Expected EOF after ignoring single line # comment")
}

func TestIgnoreMultiLineComment(t *testing.T) {
	testSingleToken(t, "/*comment \n comment \n comment */", lexer.EOF, "Expected EOF after ignoring multi line comment")
}

// Test identifiers
func TestNumeralIdentifier(t *testing.T) {
	testIdToken(t, "1234", "1234", "Expected Number")
}

func TestNumeralIdentifierFloat(t *testing.T) {
	testIdToken(t, "1234.44", "1234.44", "Expected Number")
}

func TestNumeralIdentifierNegative(t *testing.T) {
	testIdToken(t, "-1234", "-1234", "Expected Number")
}

func TestNumeralIdentifierDot(t *testing.T) {
	testIdToken(t, ".1234", ".1234", "Expected Number")
}

func TestNumeralIdentifierNegativeDot(t *testing.T) {
	testIdToken(t, "-.1234", "-.1234", "Expected Number")
}

func TestStringIdentifier(t *testing.T) {
	testIdToken(t, "\"my string\"", "my string", "Expected String")
}

func TestAlphaNumericIdentifier(t *testing.T) {
	testIdToken(t, "identifier123", "identifier123", "Expected Identifier")
}

func TestAlphaNumericIdentifier2(t *testing.T) {
	testIdToken(t, "_identi_fier", "_identi_fier", "Expected Identifier")
}

// Test examples
func TestTokenizeExample(t *testing.T) {
	var lex = getLexer("graph graphname {\n    a -- b -> c;\n    b -- d;\n}")
	var expectedPositions = []lexer.Position{
		*lexer.MakePosition(1, 1),
		*lexer.MakePosition(1, 7),
		*lexer.MakePosition(1, 17),
		*lexer.MakePosition(2, 5),
		*lexer.MakePosition(2, 7),
		*lexer.MakePosition(2, 10),
		*lexer.MakePosition(2, 12),
		*lexer.MakePosition(2, 15),
		*lexer.MakePosition(2, 16),
		*lexer.MakePosition(3, 5),
		*lexer.MakePosition(3, 7),
		*lexer.MakePosition(3, 10),
		*lexer.MakePosition(3, 11),
		*lexer.MakePosition(4, 1),
		*lexer.MakePosition(4, 2),
	}

	var expectedTokens = []lexer.Token{
		lexer.GRAPH,
		lexer.ID,
		lexer.OPEN_BRACE,
		lexer.ID,
		lexer.ARC,
		lexer.ID,
		lexer.DIRECTED_ARC,
		lexer.ID,
		lexer.SEMICOLON,
		lexer.ID,
		lexer.ARC,
		lexer.ID,
		lexer.SEMICOLON,
		lexer.CLOSE_BRACE,
		lexer.EOF,
	}

	var expectedLexemes = []string{
		"",
		"graphname",
		"",
		"a",
		"",
		"b",
		"",
		"c",
		"",
		"b",
		"",
		"d",
		"",
		"",
		"",
	}

	var i = 0
	for {
		res := lex.Next().Unwrap().Unwrap()
		if res.Position() != expectedPositions[i] || res.Token() != expectedTokens[i] || res.Lexeme() != lexer.Lexeme(expectedLexemes[i]) {
			printToken(t, "Expected", expectedPositions[i], expectedTokens[i], lexer.Lexeme(expectedLexemes[i]))
			printToken(t, "Got", res.Position(), res.Token(), res.Lexeme())
		}
		if res.Token() == lexer.EOF {
			return
		}

		i += 1
	}
}
