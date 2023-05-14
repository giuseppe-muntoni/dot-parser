package lexer

import (
	"strings"
	"testing"
)

func getLexer(input string) *Lexer {
	return New(strings.NewReader(input))
}

func tokenToString(token Token) string {
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
		return "default!?"
	}
}

func printToken(t *testing.T, error string, position Position, token Token, lexeme Lexeme) {
	t.Errorf("%s; Got line %d, column %d, token %s, lexeme %s", error, position.line, position.column, tokenToString(token), lexeme)
}

func testSingleToken(t *testing.T, input string, expectedToken Token, errorMsg string) {
	var lexer = getLexer(input)

	position, token, lexeme := lexer.Lex()
	if token != expectedToken {
		printToken(t, errorMsg, position, token, lexeme)
	}
}

func testIdToken(t *testing.T, input string, expectedLexeme string, errorMsg string) {
	var lexer = getLexer(input)

	position, token, lexeme := lexer.Lex()
	if token != ID || lexeme != Lexeme(expectedLexeme) {
		printToken(t, errorMsg, position, token, lexeme)
	}
}

func TestEOF(t *testing.T) {
	testSingleToken(t, "", EOF, "Expected EOF")
}

// Test single character tokens
func TestOpenBrace(t *testing.T) {
	testSingleToken(t, "{", OPEN_BRACE, "Expected Open Brace")
}

func TestClosedBrace(t *testing.T) {
	testSingleToken(t, "}", CLOSE_BRACE, "Expected Closed Brace")
}

func TestSemicolon(t *testing.T) {
	testSingleToken(t, ";", SEMICOLON, "Expected Semicolon")
}

func TestColon(t *testing.T) {
	testSingleToken(t, ":", COLON, "Expected Colon")
}

func TestOpenSquareBracket(t *testing.T) {
	testSingleToken(t, "[", OPEN_SQUARE_BRACKET, "Expected Open Square Bracket")
}

func TestClosedSquareBracket(t *testing.T) {
	testSingleToken(t, "]", CLOSE_SQUARE_BRACKET, "Expected Closed Square Bracket")
}

func TestEqual(t *testing.T) {
	testSingleToken(t, "=", EQUAL, "Expected Equal")
}

func TestComma(t *testing.T) {
	testSingleToken(t, ",", COMMA, "Expected Comma")
}

// Test two-character tokens
func TestArc(t *testing.T) {
	testSingleToken(t, "--", ARC, "Expected Arc (--)")
}

func TestDirectedArc(t *testing.T) {
	testSingleToken(t, "->", DIRECTED_ARC, "Expected Directed Arc (->)")
}

// Test keywords
func TestKeywordGraph(t *testing.T) {
	testSingleToken(t, "graph", GRAPH, "Expected Keyword 'graph'")
}

func TestKeywordDigraph(t *testing.T) {
	testSingleToken(t, "digraph", DIGRAPH, "Expected Keyword 'digraph'")
}

func TestKeywordStrict(t *testing.T) {
	testSingleToken(t, "strict", STRICT, "Expected Keyword 'strict'")
}

func TestKeywordNode(t *testing.T) {
	testSingleToken(t, "node", NODE, "Expected Keyword 'node'")
}

func TestKeywordEdge(t *testing.T) {
	testSingleToken(t, "edge", EDGE, "Expected Keyword 'edge'")
}

func TestKeywordSubgraph(t *testing.T) {
	testSingleToken(t, "subgraph", SUBGRAPH, "Expected Keyword 'subgraph'")
}

// Test whitespaces and comments
func TestIgnoreWhitespaces(t *testing.T) {
	testSingleToken(t, " \n\t\r  ", EOF, "Expected EOF after ignoring whitespaces")
}

func TestIgnoreSingleLineComment(t *testing.T) {
	testSingleToken(t, "//comment", EOF, "Expected EOF after ignoring single line comment")
}

func TestIgnoreSingleLineComment2(t *testing.T) {
	testSingleToken(t, "#comment", EOF, "Expected EOF after ignoring single line # comment")
}

func TestIgnoreMultiLineComment(t *testing.T) {
	testSingleToken(t, "/*comment \n comment \n comment */", EOF, "Expected EOF after ignoring multi line comment")
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
func makePosition(line int, column int) *Position {
	return &Position{line: line, column: column}
}

func TestTokenizeExample(t *testing.T) {
	var lexer = getLexer("graph graphname {\n    a -- b -> c;\n    b -- d;\n}")
	var expectedPositions = []Position{
		*makePosition(1, 1),
		*makePosition(1, 7),
		*makePosition(1, 17),
		*makePosition(2, 5),
		*makePosition(2, 7),
		*makePosition(2, 10),
		*makePosition(2, 12),
		*makePosition(2, 15),
		*makePosition(2, 16),
		*makePosition(3, 5),
		*makePosition(3, 7),
		*makePosition(3, 10),
		*makePosition(3, 11),
		*makePosition(4, 1),
		*makePosition(4, 2),
	}

	var expectedTokens = []Token{
		GRAPH,
		ID,
		OPEN_BRACE,
		ID,
		ARC,
		ID,
		DIRECTED_ARC,
		ID,
		SEMICOLON,
		ID,
		ARC,
		ID,
		SEMICOLON,
		CLOSE_BRACE,
		EOF,
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
		position, token, lexeme := lexer.Lex()
		if position != expectedPositions[i] || token != expectedTokens[i] || lexeme != Lexeme(expectedLexemes[i]) {
			printToken(t, "Error", position, token, lexeme)
		}
		if token == EOF {
			return
		}

		i += 1
	}
}
