package lexer

import (
	"testing"
)

func TestNextToken_BasicTokens(t *testing.T) {
	input := `. | : := ^ ( ) [ ] # #(`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenPeriod, "."},
		{TokenPipe, "|"},
		{TokenColon, ":"},
		{TokenAssign, ":="},
		{TokenCaret, "^"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenLBracket, "["},
		{TokenRBracket, "]"},
		{TokenHash, "#"},
		{TokenHashLParen, "#("},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Operators(t *testing.T) {
	input := `+ - * / % < > <= >= = ~=`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenPlus, "+"},
		{TokenMinus, "-"},
		{TokenStar, "*"},
		{TokenSlash, "/"},
		{TokenPercent, "%"},
		{TokenLess, "<"},
		{TokenGreater, ">"},
		{TokenLessEq, "<="},
		{TokenGreaterEq, ">="},
		{TokenEqual, "="},
		{TokenNotEqual, "~="},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Numbers(t *testing.T) {
	input := `42 3.14 -17 -2.5 100`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenInteger, "42"},
		{TokenFloat, "3.14"},
		{TokenInteger, "-17"},
		{TokenFloat, "-2.5"},
		{TokenInteger, "100"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Strings(t *testing.T) {
	input := `'Hello, World!' 'test' ''`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenString, "Hello, World!"},
		{TokenString, "test"},
		{TokenString, ""},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `true false nil`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenTrue, "true"},
		{TokenFalse, "false"},
		{TokenNil, "nil"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	input := `x count Point println ifTrue`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdentifier, "x"},
		{TokenIdentifier, "count"},
		{TokenIdentifier, "Point"},
		{TokenIdentifier, "println"},
		{TokenIdentifier, "ifTrue"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Comments(t *testing.T) {
	input := `x " this is a comment " y`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdentifier, "x"},
		{TokenIdentifier, "y"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_HelloWorld(t *testing.T) {
	input := `'Hello, World!' println.`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenString, "Hello, World!"},
		{TokenIdentifier, "println"},
		{TokenPeriod, "."},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_VariableDeclaration(t *testing.T) {
	input := `| x y |
x := 10.
y := 20.`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenPipe, "|"},
		{TokenIdentifier, "x"},
		{TokenIdentifier, "y"},
		{TokenPipe, "|"},
		{TokenIdentifier, "x"},
		{TokenAssign, ":="},
		{TokenInteger, "10"},
		{TokenPeriod, "."},
		{TokenIdentifier, "y"},
		{TokenAssign, ":="},
		{TokenInteger, "20"},
		{TokenPeriod, "."},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Arithmetic(t *testing.T) {
	input := `3 + 4 * 5`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenInteger, "3"},
		{TokenPlus, "+"},
		{TokenInteger, "4"},
		{TokenStar, "*"},
		{TokenInteger, "5"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTokenize_ValidInput(t *testing.T) {
	input := `'Hello' println.`

	l := New(input)
	tokens, err := l.Tokenize()

	if err != nil {
		t.Fatalf("Tokenize returned error: %v", err)
	}

	if len(tokens) != 4 { // STRING, IDENTIFIER, PERIOD, EOF
		t.Fatalf("Expected 4 tokens, got %d", len(tokens))
	}

	expectedTypes := []TokenType{
		TokenString,
		TokenIdentifier,
		TokenPeriod,
		TokenEOF,
	}

	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Fatalf("Token %d: expected type %q, got %q",
				i, expectedType, tokens[i].Type)
		}
	}
}

func TestTokenize_IllegalToken(t *testing.T) {
	input := `x ~ y` // ~ without = is illegal

	l := New(input)
	tokens, err := l.Tokenize()

	if err == nil {
		t.Fatal("Expected error for illegal token, got nil")
	}

	// Should still return tokens up to the illegal one
	if len(tokens) < 2 {
		t.Fatalf("Expected at least 2 tokens, got %d", len(tokens))
	}
}

func TestLineAndColumn_Tracking(t *testing.T) {
	input := `x
y
z`

	l := New(input)

	tok1 := l.NextToken()
	if tok1.Line != 1 {
		t.Errorf("Expected token on line 1, got line %d", tok1.Line)
	}

	tok2 := l.NextToken()
	if tok2.Line != 2 {
		t.Errorf("Expected token on line 2, got line %d", tok2.Line)
	}

	tok3 := l.NextToken()
	if tok3.Line != 3 {
		t.Errorf("Expected token on line 3, got line %d", tok3.Line)
	}
}

func TestNextToken_MultilineComment(t *testing.T) {
	input := `x " this is
a multi-line
comment " y`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdentifier, "x"},
		{TokenIdentifier, "y"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_NumberBeforePeriod(t *testing.T) {
	input := `42.`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenInteger, "42"},
		{TokenPeriod, "."},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
