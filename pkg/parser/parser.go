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
//   - errors: Accumulated syntax errors
//
// The parser is stateful and single-use: create a new parser for each
// source file or code snippet.
type Parser struct {
	l       *lexer.Lexer    // Token source
	curTok  lexer.Token     // Current token
	peekTok lexer.Token     // Next token (lookahead)
	errors  []string        // Accumulated error messages
}

// New creates a new parser for the given source code.
//
// The parser is initialized with the first two tokens from the lexer,
// setting up the two-token lookahead window that the parser uses for
// decision making.
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

	// Read two tokens to populate curTok and peekTok.
	// After this, curTok has the first token and peekTok has the second.
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token.
//
// This moves the lookahead window forward by one token:
//   - curTok becomes the old peekTok
//   - peekTok becomes the next token from the lexer
//
// This is called after successfully processing curTok to move to the next token.
func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
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

// parseMessageSend parses a message send expression.
//
// Message sending is the fundamental operation in smog. All computation
// happens by sending messages to objects.
//
// Syntax Types:
//
//   1. Unary messages (no arguments):
//        receiver selector
//        Example: 'Hello' println
//
//   2. Binary messages (one argument, operator-like):
//        receiver binaryOp argument
//        Example: 3 + 4
//
//   3. Keyword messages (one or more named arguments):
//        receiver key1: arg1 key2: arg2 ...
//        Example: array at: 1 put: 'value'
//
// Precedence (from highest to lowest):
//   - Unary messages
//   - Binary messages
//   - Keyword messages
//
// However, this current implementation uses a simplified left-to-right
// parsing strategy that handles one message at a time.
//
// Process:
//   1. Parse the receiver (primary expression)
//   2. Check for a message following the receiver
//   3. If found, parse it (unary, binary, or keyword)
//   4. Return MessageSend or just the receiver if no message
//
// Examples:
//
//   "42" -> just IntegerLiteral{42} (no message)
//   "x println" -> MessageSend{Receiver: Identifier("x"), Selector: "println"}
//   "3 + 4" -> MessageSend{Receiver: IntegerLiteral(3), Selector: "+", Args: [IntegerLiteral(4)]}
func (p *Parser) parseMessageSend() ast.Expression {
	// Step 1: Parse the receiver (the object that will receive the message)
	receiver := p.parsePrimaryExpression()
	if receiver == nil {
		return nil
	}

	// Step 2: Check if there's a message following the receiver
	// We peek ahead to see if the next token could start a message

	// Check for unary or keyword message (starts with identifier)
	if p.peekTok.Type == lexer.TokenIdentifier {
		p.nextToken() // advance to the identifier

		// Now curTok is the identifier, check if it's followed by a colon
		if p.peekTok.Type == lexer.TokenColon {
			// It's a keyword message - parse it
			// The identifier we just read is the first keyword part
			keyword := p.curTok.Literal
			p.nextToken() // consume colon
			selector := keyword + ":"
			
			// Parse first argument (move to the argument position)
			p.nextToken()
			arg := p.parsePrimaryExpression()
			if arg == nil {
				p.addError("expected argument after keyword")
				return nil
			}
			args := []ast.Expression{arg}
			
			// Check for additional keyword parts
			// Keyword messages can have multiple parts: at: 1 put: 'x'
			for p.peekTok.Type == lexer.TokenIdentifier {
				// Save position to check if it's another keyword part
				savedCur := p.curTok
				savedPeek := p.peekTok
				p.nextToken() // advance to identifier
				
				if p.peekTok.Type == lexer.TokenColon {
					// Yes, another keyword part - add it to the selector
					keyword := p.curTok.Literal
					p.nextToken() // consume colon
					selector += keyword + ":"
					
					// Parse the argument for this keyword part
					p.nextToken()
					arg := p.parsePrimaryExpression()
					if arg == nil {
						p.addError("expected argument after keyword")
						return nil
					}
					args = append(args, arg)
				} else {
					// Not a keyword part - restore position and stop
					p.curTok = savedCur
					p.peekTok = savedPeek
					break
				}
			}
			
			return &ast.MessageSend{
				Receiver: receiver,
				Selector: selector,
				Args:     args,
			}
		} else {
			// It's a unary message (no colon after the identifier)
			// Example: 'hello' println
			selector := p.curTok.Literal
			return &ast.MessageSend{
				Receiver: receiver,
				Selector: selector,
				Args:     []ast.Expression{}, // No arguments for unary messages
			}
		}
	}

	// Check for binary message (operator between receiver and argument)
	// Binary operators: + - * / % < > <= >= = ~=
	if p.isBinaryOperator(p.peekTok.Type) {
		p.nextToken() // advance to the operator
		operator := p.curTok.Literal
		
		// Parse the argument
		p.nextToken()
		arg := p.parsePrimaryExpression()
		if arg == nil {
			return nil
		}
		
		return &ast.MessageSend{
			Receiver: receiver,
			Selector: operator,
			Args:     []ast.Expression{arg},
		}
	}

	// No message found - just return the receiver
	return receiver
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
	case lexer.TokenIdentifier:
		return &ast.Identifier{Name: p.curTok.Literal}
	case lexer.TokenLBracket:
		return p.parseBlockLiteral()
	case lexer.TokenHashLParen:
		// Array literal #(...)
		return p.parseArrayLiteral()
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
