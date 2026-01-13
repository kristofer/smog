# Smog Parser Documentation

## Overview

The parser is responsible for transforming a stream of tokens from the lexer into an Abstract Syntax Tree (AST). It enforces the syntactic rules of the Smog language and creates a structured representation that the compiler can process.

## Purpose and Role

The parser is the second stage in the Smog compilation pipeline:

```
Source Code → Lexer → Tokens → **PARSER** → AST → Compiler → Bytecode
```

Its primary responsibilities are:
1. **Syntax Analysis**: Verify token sequences follow Smog grammar rules
2. **AST Construction**: Build a tree representation of the program structure
3. **Error Reporting**: Provide clear, actionable syntax error messages
4. **Precedence Handling**: Correctly parse message precedence (unary > binary > keyword)

## Key Concepts

### 1. Abstract Syntax Tree (AST)

The AST is a hierarchical representation of program structure:

**Source Code:**
```smog
x := 3 + 4.
```

**AST Structure:**
```
AssignmentNode
├── target: IdentifierNode("x")
└── value: MessageSendNode
    ├── receiver: IntegerNode(3)
    ├── selector: "+"
    └── arguments: [IntegerNode(4)]
```

### 2. Message Precedence

Smog follows Smalltalk's message precedence rules:

**Precedence Levels (highest to lowest):**
1. **Unary messages**: `object message`
2. **Binary messages**: `object + argument`
3. **Keyword messages**: `object at: 1 put: value`

**Example:**
```smog
array size + 1 * 2
```

**Parsing Order:**
1. `array size` (unary)
2. Result `+ 1` (binary)
3. Result `* 2` (binary, left-to-right)

**AST:**
```
MessageSend(*)
├── receiver: MessageSend(+)
│   ├── receiver: MessageSend(size)
│   │   └── receiver: Identifier(array)
│   └── arguments: [Integer(1)]
└── arguments: [Integer(2)]
```

### 3. Recursive Descent Parsing

The parser uses recursive descent, where each grammar rule has a corresponding parsing function:

```
expression → assignment | messageExpression
assignment → identifier ':=' expression
messageExpression → unaryExpression | binaryExpression | keywordExpression
```

**Implementation:**
```go
func (p *Parser) parseExpression() ast.Node {
    if p.peekToken() == TokenAssign {
        return p.parseAssignment()
    }
    return p.parseMessageExpression()
}
```

## Parsing Process

### Step-by-Step Example

**Input Source:**
```smog
| x |
x := 10 + 5.
x println.
```

**Step 1: Lexer Produces Tokens**
```
[TokenPipe, TokenIdentifier("x"), TokenPipe, 
 TokenIdentifier("x"), TokenAssign, TokenInteger("10"), 
 TokenPlus, TokenInteger("5"), TokenPeriod,
 TokenIdentifier("x"), TokenIdentifier("println"), TokenPeriod]
```

**Step 2: Parse Variable Declaration**
```
parseVariableDeclaration()
  → LocalVariables: ["x"]
```

**Step 3: Parse First Statement**
```
parseStatement()
  → parseAssignment()
    ├── target: Identifier("x")
    └── value: parseBinaryExpression()
        ├── receiver: Integer(10)
        ├── operator: "+"
        └── argument: Integer(5)
```

**Step 4: Parse Second Statement**
```
parseStatement()
  → parseMessageExpression()
    → parseUnaryMessage()
      ├── receiver: Identifier("x")
      └── selector: "println"
```

**Final AST:**
```
Program
├── locals: ["x"]
└── statements:
    ├── Assignment
    │   ├── target: Identifier("x")
    │   └── value: BinaryMessage("+", Integer(10), Integer(5))
    └── UnaryMessage("println", Identifier("x"))
```

## Grammar Rules

The complete Smog grammar (simplified):

```
program        → variableDecl? statement* EOF
variableDecl   → '|' identifier* '|'
statement      → expression '.'
expression     → assignment | messageExpr
assignment     → identifier ':=' expression
messageExpr    → keywordExpr | binaryExpr | unaryExpr | primary
keywordExpr    → binaryExpr (keyword binaryExpr)+
binaryExpr     → unaryExpr (binaryOp unaryExpr)*
unaryExpr      → primary unaryMsg*
primary        → literal | identifier | block | '(' expression ')'
literal        → integer | float | string | symbol | array | boolean | nil
block          → '[' blockParams? statement* ']'
blockParams    → (':' identifier)+ '|'
array          → '#(' literal* ')'
```

## Parsing Different Constructs

### 1. Literals

**Numbers:**
```smog
42          → IntegerNode(42)
3.14        → FloatNode(3.14)
-17         → IntegerNode(-17)
```

**Strings:**
```smog
'Hello'     → StringNode("Hello")
```

**Arrays:**
```smog
#(1 2 3)    → ArrayNode([Integer(1), Integer(2), Integer(3)])
```

### 2. Variables and Assignment

**Local Variable Declaration:**
```smog
| x y z |
```

Parser creates symbol table entries for x, y, z.

**Assignment:**
```smog
x := 42.
```

**AST:**
```
Assignment
├── target: Identifier("x")
└── value: Integer(42)
```

### 3. Message Sends

**Unary Message:**
```smog
object method
```

**AST:**
```
UnaryMessage
├── receiver: Identifier("object")
└── selector: "method"
```

**Binary Message:**
```smog
3 + 4
```

**AST:**
```
BinaryMessage
├── receiver: Integer(3)
├── selector: "+"
└── argument: Integer(4)
```

**Keyword Message:**
```smog
array at: 1 put: 'value'
```

**AST:**
```
KeywordMessage
├── receiver: Identifier("array")
├── selector: "at:put:"
└── arguments: [Integer(1), String("value")]
```

### 4. Blocks (Closures)

**Simple Block:**
```smog
[ 'Hello' println ]
```

**AST:**
```
Block
├── parameters: []
└── body: [UnaryMessage("println", String("Hello"))]
```

**Block with Parameters:**
```smog
[ :x :y | x + y ]
```

**AST:**
```
Block
├── parameters: ["x", "y"]
└── body: [BinaryMessage("+", Identifier("x"), Identifier("y"))]
```

### 5. Class Definitions

**Class Syntax:**
```smog
Object subclass: #Counter [
    | count |
    
    initialize [
        count := 0.
    ]
    
    increment [
        count := count + 1.
    ]
]
```

**AST:**
```
ClassDefinition
├── name: "Counter"
├── superclass: "Object"
├── instanceVars: ["count"]
└── methods:
    ├── Method
    │   ├── selector: "initialize"
    │   ├── parameters: []
    │   └── body: [Assignment(Identifier("count"), Integer(0))]
    └── Method
        ├── selector: "increment"
        ├── parameters: []
        └── body: [Assignment(Identifier("count"), 
                   BinaryMessage("+", Identifier("count"), Integer(1)))]
```

## Error Handling

The parser provides detailed error messages with location information:

### Syntax Errors

**Missing Period:**
```smog
x := 10
y := 20.
```

**Error:**
```
Parse error at line 2, column 1: expected '.', got 'y'
```

**Unbalanced Brackets:**
```smog
[ x + 1
```

**Error:**
```
Parse error at line 1, column 8: expected ']', got EOF
```

**Invalid Assignment Target:**
```smog
5 := 10.
```

**Error:**
```
Parse error at line 1, column 3: invalid assignment target
```

### Recovery Strategies

The parser attempts to recover from errors to find additional issues:

1. **Synchronization points**: Period (`.`) marks statement boundaries
2. **Skip tokens**: Discard tokens until synchronization point
3. **Continue parsing**: Find more errors in single pass

## Parsing Algorithm Details

### Operator Precedence Climbing

For binary operators with equal precedence (left-to-right):

```smog
a + b - c * d
```

**Algorithm:**
1. Parse `a` (primary)
2. See `+` (binary operator)
3. Parse right side with precedence check
4. Continue while operators have equal/higher precedence

**Result:** `((a + b) - c) * d` with correct precedence

### Look-Ahead

Parser uses one-token lookahead to make decisions:

```go
func (p *Parser) parseExpression() ast.Node {
    // Look ahead to decide path
    next := p.peekToken()
    
    if next == TokenAssign {
        return p.parseAssignment()
    } else if next == TokenColon {
        return p.parseKeywordMessage()
    }
    // ... etc
}
```

## Parser API

### Creating a Parser

```go
import "github.com/kristofer/smog/pkg/parser"

// From source string
p := parser.New(sourceCode)

// From file
content, _ := os.ReadFile("program.smog")
p := parser.New(string(content))
```

### Parsing to AST

```go
ast, err := p.Parse()
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
    return
}

// AST is now ready for compilation
```

### Accessing Parse Results

```go
// Visit nodes
for _, stmt := range ast.Statements {
    switch node := stmt.(type) {
    case *ast.Assignment:
        fmt.Printf("Assignment to %s\n", node.Target)
    case *ast.MessageSend:
        fmt.Printf("Message send: %s\n", node.Selector)
    }
}
```

## Testing the Parser

Example test structure:

```go
func TestParseAssignment(t *testing.T) {
    input := "x := 42."
    
    p := parser.New(input)
    ast, err := p.Parse()
    
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    // Verify AST structure
    if len(ast.Statements) != 1 {
        t.Errorf("Expected 1 statement, got %d", len(ast.Statements))
    }
    
    assign, ok := ast.Statements[0].(*ast.Assignment)
    if !ok {
        t.Errorf("Expected Assignment node")
    }
    
    if assign.Target.Name != "x" {
        t.Errorf("Expected target 'x', got '%s'", assign.Target.Name)
    }
}
```

## Common Parsing Patterns

### Parsing Lists

Many constructs involve lists (parameters, arguments, array elements):

```go
func (p *Parser) parseList(terminator TokenType) []ast.Node {
    nodes := []ast.Node{}
    
    for !p.match(terminator) && !p.isAtEnd() {
        nodes = append(nodes, p.parseExpression())
        
        if !p.match(terminator) {
            p.consume(TokenComma, "expected ','")
        }
    }
    
    return nodes
}
```

### Parsing Sequences

Statements separated by periods:

```go
func (p *Parser) parseStatements() []ast.Node {
    statements := []ast.Node{}
    
    for !p.isAtEnd() {
        stmt := p.parseStatement()
        statements = append(statements, stmt)
        
        p.consume(TokenPeriod, "expected '.' after statement")
    }
    
    return statements
}
```

## Best Practices

1. **Fail fast**: Report errors immediately when detected
2. **Clear error messages**: Include line/column and expected vs. actual
3. **Maintain invariants**: Ensure AST nodes are always valid
4. **Use helper methods**: Keep parsing functions focused and small
5. **Test edge cases**: Empty programs, deeply nested expressions, etc.

## Performance Considerations

- **Single pass**: Parser traverses tokens once
- **No backtracking**: Predictive parsing (LL(1) grammar)
- **Lazy evaluation**: Nodes created only when needed
- **Memory efficiency**: AST nodes are lightweight

## Debugging Tips

### Print Token Stream

```go
// Before parsing
tokens := lexer.Tokenize(input)
for _, tok := range tokens {
    fmt.Printf("%s: %s\n", tok.Type, tok.Literal)
}
```

### Print AST

```go
// After parsing
ast, _ := parser.Parse()
fmt.Printf("%+v\n", ast) // Deep print
```

### Enable Debug Mode

```go
parser.SetDebugMode(true)
// Parser will print each production rule as it's matched
```

## Related Documentation

- [Lexer Documentation](LEXER.md) - How tokens are created
- [AST Reference](../pkg/ast/) - AST node types
- [Compiler Documentation](COMPILER.md) - How AST is compiled
- [Language Specification](spec/LANGUAGE_SPEC.md) - Complete grammar

## Summary

The Smog parser transforms token streams into abstract syntax trees by applying grammar rules through recursive descent parsing. It handles message precedence, builds structured representations of code, and provides helpful error messages. Understanding the parser is essential for working with Smog's syntax and adding new language features.
