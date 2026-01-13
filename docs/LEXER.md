# Smog Lexer Documentation

## Overview

The lexer (also called scanner or tokenizer) is the first stage of the Smog compilation pipeline. It transforms raw source code text into a stream of tokens that the parser can process.

## Purpose and Role

The lexer is the entry point of the compilation pipeline:

```
Source Code → **LEXER** → Tokens → Parser → AST → Compiler → Bytecode
```

Its primary responsibilities are:
1. **Character Recognition**: Read source code character by character
2. **Token Classification**: Group characters into meaningful tokens
3. **Whitespace Handling**: Skip spaces, tabs, and newlines
4. **Comment Removal**: Ignore content between double quotes
5. **Error Detection**: Report illegal characters and malformed tokens

## Key Concepts

### 1. What is a Token?

A token is the smallest meaningful unit in the language:

**Source Code:**
```smog
x := 42.
```

**Tokens:**
```
1. TokenIdentifier("x")
2. TokenAssign(":=")
3. TokenInteger("42")
4. TokenPeriod(".")
```

Each token has:
- **Type**: What kind of token (identifier, number, operator, etc.)
- **Literal**: The actual text from source
- **Line**: Line number in source (for error messages)
- **Column**: Column number in source (for error messages)

### 2. Token Types

Smog has several categories of tokens:

**Literals:**
```smog
42          → TokenInteger
3.14        → TokenFloat
'hello'     → TokenString
#symbol     → TokenSymbol
true/false  → TokenTrue/TokenFalse
nil         → TokenNil
```

**Keywords:**
```smog
self        → TokenSelf
super       → TokenSuper
true        → TokenTrue
false       → TokenFalse
nil         → TokenNil
```

**Delimiters:**
```smog
.           → TokenPeriod
|           → TokenPipe
:           → TokenColon
:=          → TokenAssign
^           → TokenCaret (return)
()          → TokenLParen, TokenRParen
[]          → TokenLBracket, TokenRBracket
```

**Operators (Binary Messages):**
```smog
+           → TokenPlus
-           → TokenMinus
*           → TokenStar
/           → TokenSlash
<           → TokenLess
>           → TokenGreater
=           → TokenEqual
~=          → TokenNotEqual
```

**Special:**
```smog
#(          → TokenHashLParen (array literal)
;           → TokenSemicolon (cascade)
```

### 3. Lexical Rules

**Identifiers:**
- Start with lowercase letter or underscore
- Contain letters, digits, underscores
- Examples: `count`, `x`, `myVariable`, `_temp`

**Symbols:**
- Start with `#` followed by letters
- Examples: `#Counter`, `#Point`, `#at:put:`

**Numbers:**
- Integers: sequence of digits (`42`, `-17`)
- Floats: digits with decimal point (`3.14`, `-2.5`)

**Strings:**
- Enclosed in single quotes
- Examples: `'hello'`, `'Hello, World!'`

**Comments:**
- Enclosed in double quotes
- Can span multiple lines

**Example:**
```smog
" This is a comment "
```

## Lexing Process

### Step-by-Step Example

**Input Source:**
```smog
" Calculate sum "
| x |
x := 10 + 5.
```

**Step 1: Skip Comment**
```
Read: " Calculate sum "
Action: Skip (comment), continue to next character
```

**Step 2: Skip Whitespace**
```
Read: \n (newline)
Action: Skip, increment line counter
```

**Step 3: Read Pipe**
```
Read: |
Token: TokenPipe("|", line: 2, col: 1)
```

**Step 4: Skip Whitespace**
```
Read: ' ' (space)
Action: Skip
```

**Step 5: Read Identifier**
```
Read: x
Token: TokenIdentifier("x", line: 2, col: 3)
```

**Step 6: Read Pipe**
```
Read: |
Token: TokenPipe("|", line: 2, col: 5)
```

**Step 7: Read Identifier**
```
Read: x
Token: TokenIdentifier("x", line: 3, col: 1)
```

**Step 8: Read Assignment**
```
Read: :=
Token: TokenAssign(":=", line: 3, col: 3)
```

**Step 9: Read Number**
```
Read: 10
Token: TokenInteger("10", line: 3, col: 6)
```

**Step 10: Read Operator**
```
Read: +
Token: TokenPlus("+", line: 3, col: 9)
```

**Step 11: Read Number**
```
Read: 5
Token: TokenInteger("5", line: 3, col: 11)
```

**Step 12: Read Period**
```
Read: .
Token: TokenPeriod(".", line: 3, col: 12)
```

**Final Token Stream:**
```
[
  TokenPipe,
  TokenIdentifier("x"),
  TokenPipe,
  TokenIdentifier("x"),
  TokenAssign,
  TokenInteger("10"),
  TokenPlus,
  TokenInteger("5"),
  TokenPeriod,
  TokenEOF
]
```

## Character Classification

The lexer uses character classification to decide what to do:

```go
func isLetter(ch rune) bool {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
    return ch >= '0' && ch <= '9'
}

func isWhitespace(ch rune) bool {
    return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
```

## Scanning Different Constructs

### 1. Identifiers and Keywords

**Algorithm:**
1. Read first character (must be letter or underscore)
2. Read subsequent characters (letters, digits, underscores)
3. Check if identifier is a keyword
4. Return appropriate token

**Example:**
```smog
false → Check if keyword → TokenFalse
count → Not a keyword → TokenIdentifier("count")
```

**Code:**
```go
func (l *Lexer) scanIdentifier() Token {
    start := l.position
    
    for isLetter(l.peek()) || isDigit(l.peek()) {
        l.advance()
    }
    
    text := l.source[start:l.position]
    
    // Check for keywords
    switch text {
    case "true":
        return Token{Type: TokenTrue, Literal: text}
    case "false":
        return Token{Type: TokenFalse, Literal: text}
    case "nil":
        return Token{Type: TokenNil, Literal: text}
    case "self":
        return Token{Type: TokenSelf, Literal: text}
    case "super":
        return Token{Type: TokenSuper, Literal: text}
    default:
        return Token{Type: TokenIdentifier, Literal: text}
    }
}
```

### 2. Numbers

**Algorithm:**
1. Read digits
2. If see decimal point, read more digits (float)
3. Otherwise, return integer

**Examples:**
```
42   → TokenInteger("42")
3.14 → TokenFloat("3.14")
-17  → TokenMinus, TokenInteger("17")
```

**Code:**
```go
func (l *Lexer) scanNumber() Token {
    start := l.position
    
    // Read integer part
    for isDigit(l.peek()) {
        l.advance()
    }
    
    // Check for decimal point
    if l.peek() == '.' && isDigit(l.peekNext()) {
        l.advance() // consume '.'
        
        // Read fractional part
        for isDigit(l.peek()) {
            l.advance()
        }
        
        text := l.source[start:l.position]
        return Token{Type: TokenFloat, Literal: text}
    }
    
    text := l.source[start:l.position]
    return Token{Type: TokenInteger, Literal: text}
}
```

### 3. Strings

**Algorithm:**
1. Consume opening single quote
2. Read characters until closing single quote
3. Handle escape sequences (if supported)

**Example:**
```smog
'Hello, World!' → TokenString("Hello, World!")
```

**Code:**
```go
func (l *Lexer) scanString() Token {
    start := l.position
    l.advance() // skip opening quote
    
    for l.peek() != '\'' && !l.isAtEnd() {
        if l.peek() == '\n' {
            l.line++
        }
        l.advance()
    }
    
    if l.isAtEnd() {
        return Token{Type: TokenIllegal, Literal: "unterminated string"}
    }
    
    l.advance() // skip closing quote
    
    // Extract string content (without quotes)
    text := l.source[start+1:l.position-1]
    return Token{Type: TokenString, Literal: text}
}
```

### 4. Comments

**Algorithm:**
1. Consume opening double quote
2. Skip all characters until closing double quote
3. Don't create a token (comments are ignored)

**Example:**
```smog
" This is a comment " → (no token, just skip)
```

**Code:**
```go
func (l *Lexer) skipComment() {
    l.advance() // skip opening quote
    
    for l.peek() != '"' && !l.isAtEnd() {
        if l.peek() == '\n' {
            l.line++
        }
        l.advance()
    }
    
    if !l.isAtEnd() {
        l.advance() // skip closing quote
    }
}
```

### 5. Multi-Character Operators

Some operators span multiple characters:

```smog
:=  → TokenAssign
<=  → TokenLessEq
>=  → TokenGreaterEq
~=  → TokenNotEqual
#(  → TokenHashLParen
```

**Algorithm:**
1. Read first character
2. Peek at next character
3. If combination is valid operator, consume both
4. Otherwise, return single-character token

**Code:**
```go
func (l *Lexer) scanOperator() Token {
    ch := l.advance()
    
    switch ch {
    case ':':
        if l.peek() == '=' {
            l.advance()
            return Token{Type: TokenAssign, Literal: ":="}
        }
        return Token{Type: TokenColon, Literal: ":"}
    
    case '<':
        if l.peek() == '=' {
            l.advance()
            return Token{Type: TokenLessEq, Literal: "<="}
        }
        return Token{Type: TokenLess, Literal: "<"}
    
    // ... etc
    }
}
```

## Error Handling

The lexer reports several types of errors:

### Illegal Characters

```smog
x := 42@  " @ is not valid in Smog "
```

**Error:**
```
Lexical error at line 1, column 8: illegal character '@'
```

### Unterminated Strings

```smog
x := 'hello
```

**Error:**
```
Lexical error at line 1, column 6: unterminated string
```

### Invalid Numbers

```smog
x := 3.14.15
```

**Error:**
```
Lexical error at line 1, column 10: invalid number format
```

## Lexer State

The lexer maintains state as it scans:

```go
type Lexer struct {
    source   string    // Source code
    position int       // Current position
    line     int       // Current line number
    column   int       // Current column number
}
```

**State Changes:**
- `position`: Advanced after each character read
- `line`: Incremented on newline (`\n`)
- `column`: Incremented on each character, reset on newline

## Lexer API

### Creating a Lexer

```go
import "github.com/kristofer/smog/pkg/lexer"

l := lexer.New(sourceCode)
```

### Getting Tokens

**One at a time:**
```go
tok := l.NextToken()
for tok.Type != lexer.TokenEOF {
    fmt.Printf("%s: %s\n", tok.Type, tok.Literal)
    tok = l.NextToken()
}
```

**All at once:**
```go
tokens := l.TokenizeAll()
for _, tok := range tokens {
    fmt.Printf("%s: %s\n", tok.Type, tok.Literal)
}
```

## Testing the Lexer

Example test structure:

```go
func TestLexNumber(t *testing.T) {
    input := "42 3.14"
    
    l := lexer.New(input)
    
    // First token: integer
    tok1 := l.NextToken()
    if tok1.Type != lexer.TokenInteger {
        t.Errorf("Expected Integer, got %s", tok1.Type)
    }
    if tok1.Literal != "42" {
        t.Errorf("Expected '42', got '%s'", tok1.Literal)
    }
    
    // Second token: float
    tok2 := l.NextToken()
    if tok2.Type != lexer.TokenFloat {
        t.Errorf("Expected Float, got %s", tok2.Type)
    }
    if tok2.Literal != "3.14" {
        t.Errorf("Expected '3.14', got '%s'", tok2.Literal)
    }
}
```

## Performance Considerations

### Efficient String Scanning

Use slicing instead of concatenation:

```go
// Good: slice the source
text := l.source[start:end]

// Bad: build string character by character
var text string
for /* ... */ {
    text += string(ch)  // Allocates new string each time
}
```

### Single Pass

The lexer reads the source once, left to right, never backtracking.

### Minimal State

Only essential state is maintained (position, line, column).

## Common Pitfalls

### 1. Forgetting to Advance Position

```go
// Wrong: infinite loop
for l.peek() == ' ' {
    // Forgot to advance!
}

// Right:
for l.peek() == ' ' {
    l.advance()
}
```

### 2. Not Checking EOF

```go
// Wrong: can read past end
while (true) {
    ch := l.advance()
    // Process ch
}

// Right:
while (!l.isAtEnd()) {
    ch := l.advance()
    // Process ch
}
```

### 3. Incorrect Line/Column Tracking

```go
// Update line/column before advancing
if ch == '\n' {
    l.line++
    l.column = 0
}
l.position++
```

## Best Practices

1. **Handle EOF gracefully**: Always check before reading
2. **Track position accurately**: Essential for error messages
3. **Fail fast on errors**: Report illegal characters immediately
4. **Keep state minimal**: Only what's absolutely needed
5. **Use character classification**: `isDigit()`, `isLetter()`, etc.
6. **Test edge cases**: Empty input, only whitespace, unterminated strings

## Debugging Tips

### Print Character Stream

```go
for !l.isAtEnd() {
    ch := l.peek()
    fmt.Printf("Char: %c (code: %d)\n", ch, ch)
    l.advance()
}
```

### Print All Tokens

```go
for {
    tok := l.NextToken()
    fmt.Printf("%-15s %q\n", tok.Type, tok.Literal)
    if tok.Type == TokenEOF {
        break
    }
}
```

### Enable Debug Mode

```go
l.SetDebugMode(true)
// Lexer prints each character and token as processed
```

## Unicode Support

Smog supports Unicode identifiers:

```smog
| ñame café |  " Valid variable names "
ñame := 'Señor'.
```

**Implementation:**
```go
func isLetter(ch rune) bool {
    return unicode.IsLetter(ch) || ch == '_'
}
```

## Related Documentation

- [Parser Documentation](PARSER.md) - How tokens are parsed
- [Language Specification](spec/LANGUAGE_SPEC.md) - Lexical syntax rules
- [Token Types](../pkg/lexer/lexer.go) - Complete token type reference

## Summary

The Smog lexer transforms raw source text into a stream of tokens by recognizing patterns of characters. It handles whitespace, comments, and various literal types while tracking position for error reporting. Understanding the lexer is fundamental to understanding how Smog source code is processed, and it's the foundation upon which the parser and compiler build.
