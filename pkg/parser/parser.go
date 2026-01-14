// Package parser implements the smog language parser.
//
// The parser is responsible for converting a stream of tokens (from the lexer)
// into an Abstract Syntax Tree (AST). It performs syntactic analysis to ensure
// the code follows the grammar rules of the smog language.
//
// Parser Architecture:
//
// The parser uses a recursive descent parsing strategy, which means:
//   1. Each grammar rule corresponds to a parsing function
//   2. The parser looks ahead one token (via peekTok) to decide what to parse
//   3. Functions call each other recursively to handle nested structures
//
// Token Management:
//
// The parser maintains two tokens at all times:
//   - curTok: The current token being examined
//   - peekTok: The next token (one token lookahead)
//
// This two-token window allows the parser to make decisions without consuming
// tokens prematurely. For example, when seeing an identifier, we can peek ahead
// to see if it's followed by `:=` (assignment) or `:` (keyword message).
//
// Example Parse Flow:
//
//   Source: x := 5.
//
//   Token stream: [IDENT("x"), ASSIGN(":="), INTEGER(5), PERIOD("."), EOF]
//
//   Parse steps:
//     1. parseStatement() sees IDENT
//     2. parseExpression() sees IDENT + ASSIGN (peeking ahead)
//     3. parseAssignment() consumes IDENT, ASSIGN, parses 5
//     4. Returns Assignment{Name: "x", Value: IntegerLiteral{5}}
//
// Grammar Overview (Simplified):
//
//   Program      := Statement*
//   Statement    := VariableDecl | ExpressionStmt
//   VariableDecl := "|" Identifier* "|"
//   ExpressionStmt := Expression "."?
//   Expression   := Assignment | MessageSend
//   Assignment   := Identifier ":=" Expression
//   MessageSend  := Primary (UnaryMsg | BinaryMsg | KeywordMsg)?
//   Primary      := Literal | Identifier
//
// Error Handling:
//
// The parser accumulates errors in the `errors` slice rather than stopping
// at the first error. This allows reporting multiple syntax errors in one pass.
//
// Operator Precedence:
//
// Smog follows Smalltalk's message precedence rules:
//   1. Unary messages (highest precedence): object method
//   2. Binary messages: object + other
//   3. Keyword messages (lowest precedence): obj key: arg
//
// Within each category, messages are left-associative.
package parser

import (
	"fmt"
	"strconv"

	"github.com/kristofer/smog/pkg/ast"
	"github.com/kristofer/smog/pkg/lexer"
)

// Parser represents the smog parser.
//
// The parser maintains state during the parsing process:
//   - l: The lexer that provides tokens
//   - curTok: The current token being processed
//   - peekTok: The next token (lookahead)
//   - peekTok2: Second lookahead token (for distinguishing unary vs keyword messages)
//   - errors: Accumulated syntax errors
//
// The parser is stateful and single-use: create a new parser for each
// source file or code snippet.
//
// Note on lookahead: The parser uses two tokens of lookahead to distinguish
// between unary messages (identifier) and keyword messages (identifier followed by colon).
type Parser struct {
	l        *lexer.Lexer    // Token source
	curTok   lexer.Token     // Current token
	peekTok  lexer.Token     // Next token (1st lookahead)
	peekTok2 lexer.Token     // Token after next (2nd lookahead)
	errors   []string        // Accumulated error messages
}

// New creates a new parser for the given source code.
//
// The parser is initialized with the first three tokens from the lexer,
// setting up the two-token lookahead window that the parser uses for
// decision making (needed to distinguish unary vs keyword messages).
//
// Parameters:
//   - input: The source code string to parse
//
// Returns:
//   - A new Parser ready to parse the input
//
// Example:
//   p := parser.New("x := 5. x + 3.")
//   program, err := p.Parse()
func New(input string) *Parser {
	p := &Parser{
		l:      lexer.New(input),
		errors: []string{},
	}

	// Read three tokens to populate curTok, peekTok, and peekTok2.
	p.nextToken()
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token.
//
// This moves the lookahead window forward by one token:
//   - curTok becomes the old peekTok
//   - peekTok becomes the old peekTok2
//   - peekTok2 becomes the next token from the lexer
//
// This is called after successfully processing curTok to move to the next token.
func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.peekTok2
	p.peekTok2 = p.l.NextToken()
}

// peekIsKeywordStart checks if peekTok starts a keyword message.
//
// A keyword message starts with an identifier followed by a colon.
// With two-token lookahead, we can check this directly.
//
// Returns true if peekTok is an identifier and peekTok2 is a colon.
func (p *Parser) peekIsKeywordStart() bool {
	return p.peekTok.Type == lexer.TokenIdentifier && p.peekTok2.Type == lexer.TokenColon
}

// Parse parses the source code and returns an AST.
//
// This is the main entry point for parsing. It processes all statements
// in the program until reaching EOF (end of file).
//
// Process:
//   1. Create a Program node (the AST root)
//   2. Parse statements one by one until EOF
//   3. Add each statement to the Program's statement list
//   4. Return the completed AST or error if parsing failed
//
// Example:
//
//   Source:
//     | x |
//     x := 5.
//     x + 3.
//
//   AST:
//     Program{
//       Statements: [
//         VariableDeclaration{Names: ["x"]},
//         ExpressionStatement{Assignment{Name: "x", Value: IntegerLiteral{5}}},
//         ExpressionStatement{MessageSend{Receiver: Identifier("x"), Selector: "+", Args: [IntegerLiteral{3}]}}
//       ]
//     }
//
// Error Handling:
//   If any syntax errors were encountered, they are returned as a single
//   error containing all error messages. The AST is still returned (possibly
//   incomplete) to allow for error recovery and reporting.
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Parse statements until we hit EOF
	for p.curTok.Type != lexer.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// Move to the next token for the next iteration
		p.nextToken()
	}

	// If there were any parsing errors, return them
	if len(p.errors) > 0 {
		return program, fmt.Errorf("parser errors: %v", p.errors)
	}

	return program, nil
}

// parseStatement parses a single statement.
//
// Statements are the top-level constructs in smog. This function determines
// what kind of statement we're looking at and delegates to the appropriate
// parsing function.
//
// Statement Types:
//
//   1. Variable Declaration: | x y z |
//      Recognized by: curTok is TokenPipe
//      Parsed by: parseVariableDeclaration()
//
//   2. Return Statement: ^expression
//      Recognized by: curTok is TokenCaret
//      Parsed by: parseReturnStatement()
//
//   3. Expression Statement: any expression followed by optional period
//      Recognized by: anything else
//      Parsed by: parseExpression() wrapped in ExpressionStatement
//
// Example flows:
//
//   "| x |" -> curTok is TokenPipe -> parseVariableDeclaration()
//   "^42" -> curTok is TokenCaret -> parseReturnStatement()
//   "x := 5." -> curTok is TokenIdentifier -> parseExpression() -> Assignment
//   "3 + 4." -> curTok is TokenInteger -> parseExpression() -> MessageSend
func (p *Parser) parseStatement() ast.Statement {
	// Check for variable declarations (start with |)
	if p.curTok.Type == lexer.TokenPipe {
		return p.parseVariableDeclaration()
	}

	// Check for return statements (start with ^)
	if p.curTok.Type == lexer.TokenCaret {
		return p.parseReturnStatement()
	}

	// Check for class definitions (Identifier subclass: #ClassName [...])
	// We need to check if it's specifically: identifier "subclass" ":"
	if p.isClassDefinition() {
		return p.parseClass()
	}

	// Otherwise, treat it as an expression statement
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}

	stmt := &ast.ExpressionStatement{Expression: expr}

	// Skip optional period at end of statement
	// The period is a statement terminator but is optional at EOF
	if p.peekTok.Type == lexer.TokenPeriod {
		p.nextToken()
	}

	return stmt
}

// parseVariableDeclaration parses a variable declaration.
//
// Syntax: | varName1 varName2 ... |
//
// Variable declarations introduce local variables in the current scope.
// They consist of:
//   - Opening pipe: |
//   - Zero or more identifiers (variable names)
//   - Closing pipe: |
//
// Example:
//   | x y sum |
//
// Process:
//   1. Skip the opening | (already verified by caller)
//   2. Collect all identifier names
//   3. Expect closing |
//   4. Return VariableDeclaration with the collected names
//
// The variables are initially nil and must be assigned before use.
func (p *Parser) parseVariableDeclaration() ast.Statement {
	// Skip opening pipe (curTok is TokenPipe)
	p.nextToken()

	// Collect all variable names
	var names []string
	for p.curTok.Type == lexer.TokenIdentifier {
		names = append(names, p.curTok.Literal)
		p.nextToken()
	}

	// Expect closing pipe
	if p.curTok.Type != lexer.TokenPipe {
		p.addError("expected closing | in variable declaration")
		return nil
	}

	return &ast.VariableDeclaration{Names: names}
}

// parseExpression parses an expression.
//
// Expressions are constructs that evaluate to values. This function handles
// the top-level expression parsing and delegates to more specific parsers.
//
// Expression Types (by precedence):
//
//   1. Assignment: identifier := value
//      Special case - handled here by lookahead
//
//   2. Message Send: receiver message
//      Handled by parseMessageSend()
//
// The parser uses lookahead to distinguish assignments from other expressions.
// If we see "identifier :=", it's an assignment. Otherwise, we parse a
// message send (which might just be a primary expression with no message).
//
// Example decision trees:
//
//   "x := 5"
//     curTok=IDENT("x"), peekTok=ASSIGN
//     -> parseAssignment()
//
//   "x + 5"
//     curTok=IDENT("x"), peekTok=PLUS
//     -> parseMessageSend() -> binary message
//
//   "42"
//     curTok=INTEGER(42), peekTok=PERIOD
//     -> parseMessageSend() -> just primary expression
func (p *Parser) parseExpression() ast.Expression {
	// Check for assignment by looking ahead
	// Assignment syntax: identifier := expression
	if p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenAssign {
		return p.parseAssignment()
	}

	// Otherwise, parse as a message send (or just a primary expression)
	return p.parseMessageSend()
}

// parseAssignment parses an assignment expression.
//
// Syntax: variableName := value
//
// Assignments bind a value to a variable. The value can be any expression.
// Assignments are themselves expressions and return the assigned value.
//
// Process:
//   1. Extract the variable name from curTok
//   2. Consume the := operator
//   3. Parse the value expression (recursive - can be anything)
//   4. Return Assignment node
//
// Example:
//   x := 10
//     -> Assignment{Name: "x", Value: IntegerLiteral{10}}
//
//   y := x + 5
//     -> Assignment{Name: "y", Value: MessageSend{...}}
//
// Note: The caller has already verified curTok is IDENT and peekTok is ASSIGN.
func (p *Parser) parseAssignment() ast.Expression {
	// Get the variable name
	name := p.curTok.Literal
	p.nextToken() // consume identifier

	// Verify := operator (should always be true given caller's check)
	if p.curTok.Type != lexer.TokenAssign {
		p.addError("expected := in assignment")
		return nil
	}
	p.nextToken() // consume :=

	// Parse the value expression
	// This can be any expression, including another assignment: x := y := 5
	value := p.parseMessageSend()
	if value == nil {
		return nil
	}

	return &ast.Assignment{
		Name:  name,
		Value: value,
	}
}

// parseMessageSend parses a message send expression with proper Smalltalk precedence.
//
// Message sending is the fundamental operation in smog. All computation
// happens by sending messages to objects.
//
// Smalltalk Message Precedence (from highest to lowest):
//   1. Unary messages: receiver selector
//   2. Binary messages: receiver op argument
//   3. Keyword messages: receiver key: arg
//
// Within each level, messages are evaluated left-to-right.
//
// Examples demonstrating precedence:
//   arr size + 1        -> (arr size) + 1         (unary before binary)
//   3 + 4 * 2          -> (3 + 4) * 2             (binary left-to-right, no operator precedence)
//   arr at: i + 1      -> arr at: (i + 1)         (binary in keyword argument)
//   x sqrt negated     -> (x sqrt) negated        (unary chains left-to-right)
//
// This implementation properly handles the precedence hierarchy by
// having each precedence level call the next higher level for its components.
func (p *Parser) parseMessageSend() ast.Expression {
	// Check for super message send
	if p.curTok.Type == lexer.TokenSuper {
		return p.parseSuperMessageSend()
	}
	
	// Start with keyword messages (lowest precedence)
	// Keyword messages will call binary messages for their receiver and arguments
	return p.parseKeywordMessage()
}

// parseKeywordMessage parses keyword messages (lowest precedence).
//
// Syntax: receiver keyword1: arg1 keyword2: arg2 ...
//
// Examples:
//   array at: 1
//   array at: 1 put: 'value'
//   point x: 10 y: 20
//
// The receiver and arguments are parsed as binary messages (next higher precedence).
func (p *Parser) parseKeywordMessage() ast.Expression {
	// Parse receiver as a binary message (which will handle unary messages too)
	receiver := p.parseBinaryMessage()
	if receiver == nil {
		return nil
	}
	
	// Check if this is followed by a keyword message
	// Use the helper to check for identifier followed by colon
	if !p.peekIsKeywordStart() {
		// No keyword message, but might still have a cascade
		// Check if the receiver is a message send and if so, check for cascade
		return p.checkForCascade(receiver)
	}
	
	// It's a keyword message - parse all keyword parts
	var selector string
	var args []ast.Expression
	
	for p.peekIsKeywordStart() {
		p.nextToken() // move to keyword identifier
		selector += p.curTok.Literal + ":"
		p.nextToken() // consume colon
		
		// Parse argument as binary message (can contain unary and binary messages)
		p.nextToken() // move to argument position
		arg := p.parseBinaryMessage()
		if arg == nil {
			p.addError("expected argument after keyword")
			return nil
		}
		args = append(args, arg)
	}
	
	msgSend := &ast.MessageSend{
		Receiver: receiver,
		Selector: selector,
		Args:     args,
	}
	
	// Check for cascade after this message
	return p.checkForCascade(msgSend)
}

// parseBinaryMessage parses binary messages (middle precedence).
//
// Syntax: receiver binaryOp argument
//
// Binary operators: + - * / % < > <= >= = ~=
//
// Binary messages are left-associative with no operator precedence:
//   3 + 4 * 2  means  (3 + 4) * 2 = 14  (not 3 + 8 = 11)
//   10 - 5 + 3 means  (10 - 5) + 3 = 8
//
// The receiver and arguments are parsed as unary messages (next higher precedence).
//
// Examples:
//   3 + 4              -> MessageSend{Receiver: 3, Selector: "+", Args: [4]}
//   arr size + 1       -> MessageSend{Receiver: (arr size), Selector: "+", Args: [1]}
//   3 + 4 * 2          -> MessageSend{Receiver: (3+4), Selector: "*", Args: [2]}
func (p *Parser) parseBinaryMessage() ast.Expression {
	// Parse receiver as unary messages (which will handle primary too)
	receiver := p.parseUnaryMessage()
	if receiver == nil {
		return nil
	}
	
	// Chain binary messages (left-to-right)
	for p.isBinaryOperator(p.peekTok.Type) {
		p.nextToken() // advance to operator
		operator := p.curTok.Literal
		
		// Parse argument as unary message
		p.nextToken() // move to argument
		arg := p.parseUnaryMessage()
		if arg == nil {
			p.addError("expected argument after binary operator")
			return nil
		}
		
		// Build message send with current receiver
		receiver = &ast.MessageSend{
			Receiver: receiver,
			Selector: operator,
			Args:     []ast.Expression{arg},
		}
	}
	
	return receiver
}

// parseUnaryMessage parses unary messages (highest precedence).
//
// Syntax: receiver selector1 selector2 ...
//
// Unary messages are chained left-to-right:
//   x sqrt floor  means  (x sqrt) floor
//   arr size negated means (arr size) negated
//
// The receiver is parsed as a primary expression.
//
// Examples:
//   x println          -> MessageSend{Receiver: x, Selector: "println"}
//   arr size           -> MessageSend{Receiver: arr, Selector: "size"}
//   x sqrt floor       -> MessageSend{Receiver: (x sqrt), Selector: "floor"}
func (p *Parser) parseUnaryMessage() ast.Expression {
	// Parse the primary expression (literals, identifiers, blocks, etc.)
	receiver := p.parsePrimaryExpression()
	if receiver == nil {
		return nil
	}
	
	// Chain unary messages (left-to-right)
	// Only consume identifiers that are NOT followed by colons (which would be keyword messages)
	for p.peekTok.Type == lexer.TokenIdentifier && !p.peekIsKeywordStart() {
		p.nextToken() // move to the unary selector
		selector := p.curTok.Literal
		receiver = &ast.MessageSend{
			Receiver: receiver,
			Selector: selector,
			Args:     []ast.Expression{},
		}
	}
	
	return receiver
}

// checkForCascade checks if there's a cascade (;) after the initial expression
// and if so, parses the cascade.
//
// Syntax: receiver message1; message2; message3
//
// The receiver is evaluated once, and each message is sent to the same receiver.
// The cascade returns the receiver itself (not the result of the last message).
func (p *Parser) checkForCascade(expr ast.Expression) ast.Expression {
	// If the expression is not a message send, it can't be cascaded
	firstMsg, isMessageSend := expr.(*ast.MessageSend)
	if !isMessageSend {
		return expr
	}
	
	// Check if there's a semicolon indicating a cascade
	if p.peekTok.Type != lexer.TokenSemicolon {
		return expr
	}
	
	// We have a cascade! Build a CascadeExpression
	receiver := firstMsg.Receiver
	messages := []ast.MessageSend{*firstMsg}
	
	// Parse additional messages separated by semicolons
	for p.peekTok.Type == lexer.TokenSemicolon {
		p.nextToken() // consume the semicolon
		p.nextToken() // move to the message selector
		
		// Parse the next message (without the receiver)
		msg := p.parseMessageWithoutReceiver()
		if msg != nil {
			messages = append(messages, *msg)
		}
	}
	
	return &ast.CascadeExpression{
		Receiver: receiver,
		Messages: messages,
	}
}

// parseMessageWithoutReceiver parses a message selector and arguments
// without a receiver (used in cascades).
//
// Returns a MessageSend with nil Receiver.
// This needs to handle all three types of messages: unary, binary, and keyword.
func (p *Parser) parseMessageWithoutReceiver() *ast.MessageSend {
	// Check for keyword message (identifier followed by colon)
	if p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenColon {
		var selector string
		var args []ast.Expression
		
		// Parse keyword parts
		for p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenColon {
			selector += p.curTok.Literal + ":"
			p.nextToken() // consume colon
			
			// Parse argument as binary message (can include unary and binary)
			p.nextToken()
			arg := p.parseBinaryMessage()
			if arg == nil {
				p.addError("expected argument after keyword in cascade")
				return nil
			}
			args = append(args, arg)
			
			// Check for next keyword part using the helper
			if !p.peekIsKeywordStart() {
				break
			}
			p.nextToken() // move to next keyword identifier
		}
		
		return &ast.MessageSend{
			Receiver: nil,
			Selector: selector,
			Args:     args,
		}
	} else if p.isBinaryOperator(p.curTok.Type) {
		// Binary message - parse right side as unary message
		operator := p.curTok.Literal
		p.nextToken()
		arg := p.parseUnaryMessage()
		if arg == nil {
			p.addError("expected argument after binary operator in cascade")
			return nil
		}
		
		return &ast.MessageSend{
			Receiver: nil,
			Selector: operator,
			Args:     []ast.Expression{arg},
		}
	} else if p.curTok.Type == lexer.TokenIdentifier {
		// Unary message
		selector := p.curTok.Literal
		return &ast.MessageSend{
			Receiver: nil,
			Selector: selector,
			Args:     []ast.Expression{},
		}
	}
	
	p.addError("expected message selector in cascade")
	return nil
}

// parseSuperMessageSend parses a super message send.
//
// Syntax: super selector
//        or: super keyword: arg
//        or: super binaryOp arg
//
// Super sends start method lookup in the superclass of the current class.
// They're used to call inherited methods that have been overridden.
//
// Process:
//   1. Verify we're on the 'super' keyword
//   2. Parse the message selector and arguments with proper precedence
//   3. Return MessageSend with IsSuper flag set
//
// Examples:
//   super initialize
//     -> MessageSend{Receiver: nil, Selector: "initialize", Args: [], IsSuper: true}
//
//   super at: index
//     -> MessageSend{Receiver: nil, Selector: "at:", Args: [index], IsSuper: true}
//
//   super + other
//     -> MessageSend{Receiver: nil, Selector: "+", Args: [other], IsSuper: true}
func (p *Parser) parseSuperMessageSend() ast.Expression {
	// curTok is TokenSuper
	p.nextToken() // move to the message selector
	
	// Check if it's a keyword message (identifier followed by colon)
	if p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenColon {
		var selector string
		var args []ast.Expression
		
		// Parse keyword parts
		for p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenColon {
			selector += p.curTok.Literal + ":"
			p.nextToken() // consume colon
			
			// Parse argument as binary message
			p.nextToken()
			arg := p.parseBinaryMessage()
			if arg == nil {
				p.addError("expected argument after keyword in super send")
				return nil
			}
			args = append(args, arg)
			
			// Check for next keyword part using helper
			if !p.peekIsKeywordStart() {
				break
			}
			p.nextToken() // move to next keyword identifier
		}
		
		return &ast.MessageSend{
			Receiver: nil, // receiver is implicit (self)
			Selector: selector,
			Args:     args,
			IsSuper:  true,
		}
	} else if p.isBinaryOperator(p.curTok.Type) {
		// Binary message
		operator := p.curTok.Literal
		p.nextToken()
		arg := p.parseUnaryMessage()
		if arg == nil {
			p.addError("expected argument after binary operator in super send")
			return nil
		}
		
		return &ast.MessageSend{
			Receiver: nil, // receiver is implicit (self)
			Selector: operator,
			Args:     []ast.Expression{arg},
			IsSuper:  true,
		}
	} else if p.curTok.Type == lexer.TokenIdentifier {
		// Unary message
		selector := p.curTok.Literal
		return &ast.MessageSend{
			Receiver: nil, // receiver is implicit (self)
			Selector: selector,
			Args:     []ast.Expression{},
			IsSuper:  true,
		}
	}
	
	p.addError("expected message selector after super")
	return nil
}

// isBinaryOperator checks if a token type represents a binary operator.
//
// Binary operators are special message selectors that appear between
// the receiver and argument (infix notation).
//
// Supported binary operators:
//   Arithmetic: + - * / %
//   Comparison: < > <= >= = ~=
//
// Returns true if the token type is one of these operators.
func (p *Parser) isBinaryOperator(tt lexer.TokenType) bool {
	return tt == lexer.TokenPlus ||
		tt == lexer.TokenMinus ||
		tt == lexer.TokenStar ||
		tt == lexer.TokenSlash ||
		tt == lexer.TokenPercent ||
		tt == lexer.TokenLess ||
		tt == lexer.TokenGreater ||
		tt == lexer.TokenLessEq ||
		tt == lexer.TokenGreaterEq ||
		tt == lexer.TokenEqual ||
		tt == lexer.TokenNotEqual
}

// parsePrimaryExpression parses a primary expression (literals and identifiers).
//
// Primary expressions are the atomic building blocks of expressions.
// They don't contain any operators or sub-expressions - they're just values.
//
// Primary Expression Types:
//   - Integer literals: 42, 0, -5
//   - Float literals: 3.14, 0.5
//   - String literals: 'Hello'
//   - Boolean literals: true, false
//   - Nil literal: nil
//   - Identifiers: variableName, x, count
//   - Block literals: [ ... ], [ :x | ... ]
//   - Array literals: #(1 2 3)
//
// This function dispatches to specific parsing functions based on the
// current token type.
//
// Example mappings:
//   TokenInteger -> parseIntegerLiteral() -> IntegerLiteral{Value: 42}
//   TokenString -> parseStringLiteral() -> StringLiteral{Value: "Hello"}
//   TokenIdentifier -> Identifier{Name: "x"}
//   TokenLBracket -> parseBlockLiteral() -> BlockLiteral{...}
func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.curTok.Type {
	case lexer.TokenInteger:
		return p.parseIntegerLiteral()
	case lexer.TokenFloat:
		return p.parseFloatLiteral()
	case lexer.TokenString:
		return p.parseStringLiteral()
	case lexer.TokenTrue:
		return &ast.BooleanLiteral{Value: true}
	case lexer.TokenFalse:
		return &ast.BooleanLiteral{Value: false}
	case lexer.TokenNil:
		return &ast.NilLiteral{}
	case lexer.TokenSelf:
		// self is represented as a special identifier
		return &ast.Identifier{Name: "self"}
	case lexer.TokenIdentifier:
		return &ast.Identifier{Name: p.curTok.Literal}
	case lexer.TokenLBracket:
		return p.parseBlockLiteral()
	case lexer.TokenHashLParen:
		// Array literal #(...)
		return p.parseArrayLiteral()
	case lexer.TokenHashLBrace:
		// Dictionary literal #{...}
		return p.parseDictionaryLiteral()
	case lexer.TokenLParen:
		// Parenthesized expression (...)
		return p.parseParenthesizedExpression()
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.curTok.Type))
		return nil
	}
}

// parseIntegerLiteral parses an integer literal.
//
// Converts the token's string representation to an int64 value.
//
// Example:
//   Token{Type: TokenInteger, Literal: "42"}
//     -> IntegerLiteral{Value: 42}
//
// Error handling:
//   If the string can't be parsed as an integer (shouldn't happen if
//   the lexer is correct), an error is recorded.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curTok.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Value: value}
}

// parseFloatLiteral parses a floating-point literal.
//
// Converts the token's string representation to a float64 value.
//
// Example:
//   Token{Type: TokenFloat, Literal: "3.14"}
//     -> FloatLiteral{Value: 3.14}
//
// Error handling:
//   If the string can't be parsed as a float, an error is recorded.
func (p *Parser) parseFloatLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curTok.Literal))
		return nil
	}
	return &ast.FloatLiteral{Value: value}
}

// parseStringLiteral parses a string literal.
//
// The lexer has already removed the quotes, so we just extract the value.
//
// Example:
//   Token{Type: TokenString, Literal: "Hello"}
//     -> StringLiteral{Value: "Hello"}
//
// Note: The token's Literal field contains the string without quotes.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curTok.Literal}
}

// addError adds an error message to the error list.
//
// The parser accumulates errors rather than stopping at the first one.
// This allows reporting multiple syntax errors in a single pass.
//
// Parameters:
//   - msg: A human-readable error message
//
// Example:
//   p.addError("expected closing | in variable declaration")
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// parseBlockLiteral parses a block literal.
//
// Syntax: [ statements... ]
//        or: [ :param1 :param2 ... | statements... ]
//
// Blocks are closures that can capture variables from their environment.
//
// Process:
//   1. Skip the opening [ (already verified by caller)
//   2. Check for parameters (start with :)
//   3. If parameters exist, collect them until |
//   4. Parse statements until closing ]
//   5. Return BlockLiteral node
//
// Examples:
//   [ 'Hello' println ]
//     -> BlockLiteral{Parameters: [], Body: [println statement]}
//
//   [ :x | x * 2 ]
//     -> BlockLiteral{Parameters: ["x"], Body: [x * 2 statement]}
//
//   [ :x :y | x + y ]
//     -> BlockLiteral{Parameters: ["x", "y"], Body: [x + y statement]}
func (p *Parser) parseBlockLiteral() ast.Expression {
	// curTok is [, move to next
	p.nextToken()

	var parameters []string

	// Check for parameters (start with colon)
	if p.curTok.Type == lexer.TokenColon {
		// Parse parameters
		for p.curTok.Type == lexer.TokenColon {
			p.nextToken() // skip colon
			if p.curTok.Type != lexer.TokenIdentifier {
				p.addError("expected parameter name after :")
				return nil
			}
			parameters = append(parameters, p.curTok.Literal)
			p.nextToken() // move past parameter name
		}

		// Expect pipe after parameters
		if p.curTok.Type != lexer.TokenPipe {
			p.addError("expected | after block parameters")
			return nil
		}
		p.nextToken() // skip pipe
	}

	// Parse block body (statements until ])
	var body []ast.Statement
	for p.curTok.Type != lexer.TokenRBracket && p.curTok.Type != lexer.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		// If we're at a period, skip it and continue
		if p.curTok.Type == lexer.TokenPeriod {
			p.nextToken()
		} else if p.curTok.Type != lexer.TokenRBracket {
			// If not at ] and not at period, move forward
			p.nextToken()
		}
	}

	// Expect closing ]
	if p.curTok.Type != lexer.TokenRBracket {
		p.addError("expected ] to close block")
		return nil
	}

	return &ast.BlockLiteral{
		Parameters: parameters,
		Body:       body,
	}
}

// parseReturnStatement parses a return statement.
//
// Syntax: ^expression
//
// Return statements exit from methods, returning a value.
//
// Example:
//   ^count
//     -> ReturnStatement{Value: Identifier("count")}
//
//   ^x + y
//     -> ReturnStatement{Value: MessageSend{...}}
func (p *Parser) parseReturnStatement() ast.Statement {
	// curTok is ^, move to the expression
	p.nextToken()

	// Parse the return value expression
	value := p.parseExpression()
	if value == nil {
		p.addError("expected expression after ^")
		return nil
	}

	return &ast.ReturnStatement{Value: value}
}

// parseArrayLiteral parses an array literal.
//
// Syntax: #(element1 element2 ...)
//
// Array literals create array objects with the specified elements.
//
// Example:
//   #(1 2 3 4 5)
//     -> ArrayLiteral{Elements: [1, 2, 3, 4, 5]}
func (p *Parser) parseArrayLiteral() ast.Expression {
	// curTok is #(
	p.nextToken() // move past #(

	var elements []ast.Expression

	// Parse elements until )
	for p.curTok.Type != lexer.TokenRParen && p.curTok.Type != lexer.TokenEOF {
		elem := p.parsePrimaryExpression()
		if elem != nil {
			elements = append(elements, elem)
		}
		p.nextToken()
	}

	// Expect closing )
	if p.curTok.Type != lexer.TokenRParen {
		p.addError("expected ) to close array literal")
		return nil
	}

	return &ast.ArrayLiteral{Elements: elements}
}

// parseDictionaryLiteral parses a dictionary literal.
//
// Syntax: #{key1 -> value1. key2 -> value2. ...}
//
// Dictionary literals create dictionary objects with the specified key-value pairs.
// Each pair consists of a key expression, an arrow (->), and a value expression.
// Pairs are separated by periods.
//
// Example:
//   #{'name' -> 'Alice'. 'age' -> 30}
//     -> DictionaryLiteral{Pairs: [{'name', 'Alice'}, {'age', 30}]}
func (p *Parser) parseDictionaryLiteral() ast.Expression {
	// curTok is #{
	p.nextToken() // move past #{

	var pairs []ast.DictionaryPair

	// Parse key-value pairs until }
	for p.curTok.Type != lexer.TokenRBrace && p.curTok.Type != lexer.TokenEOF {
		// Parse key
		key := p.parsePrimaryExpression()
		if key == nil {
			p.addError("expected key in dictionary literal")
			return nil
		}
		
		p.nextToken()
		
		// Expect arrow
		if p.curTok.Type != lexer.TokenArrow {
			p.addError("expected -> after dictionary key")
			return nil
		}
		
		p.nextToken() // move past ->
		
		// Parse value
		value := p.parsePrimaryExpression()
		if value == nil {
			p.addError("expected value in dictionary literal")
			return nil
		}
		
		pairs = append(pairs, ast.DictionaryPair{Key: key, Value: value})
		
		p.nextToken()
		
		// Skip optional period between pairs
		if p.curTok.Type == lexer.TokenPeriod {
			p.nextToken()
		}
	}

	// Expect closing }
	if p.curTok.Type != lexer.TokenRBrace {
		p.addError("expected } to close dictionary literal")
		return nil
	}

	return &ast.DictionaryLiteral{Pairs: pairs}
}

// parseParenthesizedExpression parses an expression within parentheses.
//
// Syntax: (expression)
//
// Parentheses are used for grouping and controlling evaluation order.
// They override the normal precedence rules.
//
// Example:
//   (x + y) * z
//   Point x: (a + b) y: (c + d)
//   (3 + 4) sqrt
func (p *Parser) parseParenthesizedExpression() ast.Expression {
	// curTok is '('
	p.nextToken() // move past '('
	
	// Parse the full expression inside (starting with lowest precedence - keyword messages)
	expr := p.parseKeywordMessage()
	if expr == nil {
		return nil
	}
	
	// Expect closing ')'
	p.nextToken()
	if p.curTok.Type != lexer.TokenRParen {
		p.addError("expected ')' to close parenthesized expression")
		return nil
	}
	
	return expr
}

// Errors returns the list of accumulated parsing errors.
//
// This can be called after Parse() to get detailed error information
// if parsing failed.
//
// Returns:
//   - A slice of error message strings (empty if no errors)
func (p *Parser) Errors() []string {
	return p.errors
}

// isClassDefinition checks if the current position is at the start of a class definition.
//
// A class definition has the pattern: Identifier "subclass" ":" ...
// We check if curTok is identifier and peekTok is specifically "subclass".
func (p *Parser) isClassDefinition() bool {
	return p.curTok.Type == lexer.TokenIdentifier &&
		p.peekTok.Type == lexer.TokenIdentifier &&
		p.peekTok.Literal == "subclass"
}

// parseClass parses a class definition.
//
// Syntax: SuperClass subclass: #ClassName [
//           | instanceVar1 instanceVar2 |
//           <| classVar1 classVar2 |>
//           method1 [ body ]
//           <classMethod [ body ]>
//         ]
//
// Process:
//   1. Extract superclass name (already at identifier)
//   2. Verify "subclass:" keyword
//   3. Parse class name (symbol starting with #)
//   4. Parse class body within brackets [...]
//   5. Within body, parse instance variables, class variables, and methods
//
// Example:
//   Object subclass: #Counter [
//       | count |
//       initialize [ count := 0. ]
//   ]
func (p *Parser) parseClass() *ast.Class {
	// curTok should be the superclass identifier
	if p.curTok.Type != lexer.TokenIdentifier {
		p.addError("expected superclass identifier")
		return nil
	}
	superClass := p.curTok.Literal
	
	// Move to "subclass" keyword
	p.nextToken()
	if p.curTok.Type != lexer.TokenIdentifier || p.curTok.Literal != "subclass" {
		p.addError("expected 'subclass' keyword")
		return nil
	}
	
	// Expect colon after "subclass"
	p.nextToken()
	if p.curTok.Type != lexer.TokenColon {
		p.addError("expected ':' after 'subclass'")
		return nil
	}
	
	// Move to class name (should be a symbol like #Counter)
	p.nextToken()
	if p.curTok.Type != lexer.TokenHash {
		p.addError("expected '#' before class name")
		return nil
	}
	
	// Get the class name after #
	p.nextToken()
	if p.curTok.Type != lexer.TokenIdentifier {
		p.addError("expected class name after '#'")
		return nil
	}
	className := p.curTok.Literal
	
	// Expect opening bracket [
	p.nextToken()
	if p.curTok.Type != lexer.TokenLBracket {
		p.addError("expected '[' to start class body")
		return nil
	}
	
	// Parse class body
	class := &ast.Class{
		Name:           className,
		SuperClass:     superClass,
		Fields:         []string{},
		ClassVariables: []string{},
		Methods:        []*ast.Method{},
		ClassMethods:   []*ast.Method{},
	}
	
	p.nextToken() // move into the class body
	
	// Parse instance variables if present (| var1 var2 |)
	if p.curTok.Type == lexer.TokenPipe {
		p.nextToken() // skip opening |
		for p.curTok.Type == lexer.TokenIdentifier {
			class.Fields = append(class.Fields, p.curTok.Literal)
			p.nextToken()
		}
		if p.curTok.Type != lexer.TokenPipe {
			p.addError("expected '|' to close instance variables")
			return nil
		}
		p.nextToken() // skip closing |
	}
	
	// Parse class variables if present (<| classVar1 classVar2 |>)
	if p.curTok.Type == lexer.TokenLess {
		// Check if next is pipe
		if p.peekTok.Type == lexer.TokenPipe {
			p.nextToken() // skip <
			p.nextToken() // skip |
			for p.curTok.Type == lexer.TokenIdentifier {
				class.ClassVariables = append(class.ClassVariables, p.curTok.Literal)
				p.nextToken()
			}
			if p.curTok.Type != lexer.TokenPipe {
				p.addError("expected '|' to close class variables")
				return nil
			}
			p.nextToken() // skip |
			if p.curTok.Type != lexer.TokenGreater {
				p.addError("expected '>' to close class variables")
				return nil
			}
			p.nextToken() // skip >
		}
	}
	
	// Parse methods until we hit the closing bracket
	for p.curTok.Type != lexer.TokenRBracket && p.curTok.Type != lexer.TokenEOF {
		// Check if this is a class method (starts with <)
		isClassMethod := false
		if p.curTok.Type == lexer.TokenLess {
			isClassMethod = true
			// Don't consume the < yet, let parseMethod handle it
		}
		
		method := p.parseMethod()
		if method != nil {
			if isClassMethod {
				class.ClassMethods = append(class.ClassMethods, method)
			} else {
				class.Methods = append(class.Methods, method)
			}
		}
	}
	
	// Expect closing bracket ]
	if p.curTok.Type != lexer.TokenRBracket {
		p.addError("expected ']' to close class body")
		return nil
	}
	
	return class
}

// parseMethod parses a method definition within a class.
//
// Syntax: methodSelector [ body ]
//        or: keyword: param [ body ]
//        or: <classMethod [ body ]>
//
// Returns a Method with name, parameters, and body.
func (p *Parser) parseMethod() *ast.Method {
	// Check for class method (starts with <)
	isClassMethod := false
	if p.curTok.Type == lexer.TokenLess {
		isClassMethod = true
		p.nextToken() // skip <
	}
	
	// Parse method selector and parameters
	var selector string
	var params []string
	
	// Check what kind of method selector we have
	if p.curTok.Type == lexer.TokenIdentifier {
		// Could be unary or keyword method
		if p.peekTok.Type == lexer.TokenColon {
			// Keyword method - parse keyword parts
			for p.curTok.Type == lexer.TokenIdentifier && p.peekTok.Type == lexer.TokenColon {
				selector += p.curTok.Literal + ":"
				p.nextToken() // skip identifier
				p.nextToken() // skip colon
				
				// Get parameter name
				if p.curTok.Type != lexer.TokenIdentifier {
					p.addError("expected parameter name after ':'")
					return nil
				}
				params = append(params, p.curTok.Literal)
				p.nextToken()
			}
		} else {
			// Unary method
			selector = p.curTok.Literal
			p.nextToken()
		}
	} else if p.isBinaryOperator(p.curTok.Type) {
		// Binary method (e.g., +, -, etc.)
		selector = p.curTok.Literal
		p.nextToken()
		
		// Binary methods have one parameter
		if p.curTok.Type != lexer.TokenIdentifier {
			p.addError("expected parameter name for binary method")
			return nil
		}
		params = append(params, p.curTok.Literal)
		p.nextToken()
	} else {
		p.addError("expected method selector")
		return nil
	}
	
	// Expect opening bracket for method body
	if p.curTok.Type != lexer.TokenLBracket {
		p.addError("expected '[' to start method body")
		return nil
	}
	p.nextToken() // skip [
	
	// Parse method body (statements until ])
	var body []ast.Statement
	for p.curTok.Type != lexer.TokenRBracket && p.curTok.Type != lexer.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		p.nextToken()
	}
	
	// Expect closing bracket
	if p.curTok.Type != lexer.TokenRBracket {
		p.addError("expected ']' to close method body")
		return nil
	}
	p.nextToken() // skip ]
	
	// If class method, expect closing >
	if isClassMethod {
		if p.curTok.Type != lexer.TokenGreater {
			p.addError("expected '>' to close class method")
			return nil
		}
		p.nextToken() // skip >
	}
	
	method := &ast.Method{
		Name:       selector,
		Parameters: params,
		Body:       body,
	}
	
	// Note: We don't distinguish class methods from instance methods in the AST yet
	// This would need to be added to the Method struct or handled separately
	// For now, all methods go into the Methods slice
	
	return method
}
