# Smog Language Specification

Version 0.1.0

## Overview

Smog is a simple object-oriented language inspired by Smalltalk and the Simple Object Machine (SOM). It emphasizes simplicity, message passing, and everything-is-an-object philosophy.

## Language Principles

1. **Everything is an object** - Even classes and primitive values are objects
2. **Message passing** - All computation happens through sending messages to objects
3. **Simple syntax** - Minimal syntax inspired by Smalltalk
4. **Bytecode compilation** - Source code is compiled to bytecode for efficient execution

## Syntax

### Comments

```smog
" This is a comment "
" Comments are enclosed in double quotes
  and can span multiple lines "
```

### Literals

#### Numbers
```smog
42          " Integer "
3.14        " Float "
-17         " Negative number "
```

#### Strings
```smog
'Hello, World!'
'Multiple words in a string'
```

#### Booleans and Nil
```smog
true
false
nil
```

#### Arrays
```smog
#(1 2 3 4 5)
#('hello' 'world')
```

### Variables

#### Local Variables

Local variables are declared using pipe symbols (`|`) and must be declared in a **single declaration block** at the beginning of a scope (method, block, or top-level).

```smog
| x y |
x := 10.
y := 20.
```

**Important Scoping Rule:** All local variables in a scope must be declared together in one declaration block. You cannot add new variable declarations after executing statements or creating blocks.

**✅ Correct:**
```smog
| x y result |  " All variables declared together
x := 10.
y := 20.
result := x + y.
```

**❌ Incorrect:**
```smog
| x y |
x := 10.
[ y := 20 ] value.
| result |  " ERROR: Cannot declare variables after blocks
result := x + y.
```

This limitation is due to the current variable scoping implementation. As a workaround, declare all variables you will need at the beginning of the scope.

#### Instance Variables
Defined in class declarations and accessed directly in methods.

### Classes

Class definitions follow the Smalltalk style:

```smog
Object subclass: #Point [
    | x y |
    
    "Constructor"
    x: xValue y: yValue [
        x := xValue.
        y := yValue.
    ]
    
    "Accessors"
    x [
        ^x
    ]
    
    y [
        ^y
    ]
    
    "Methods"
    + aPoint [
        ^Point x: (x + aPoint x) y: (y + aPoint y)
    ]
    
    printOn: stream [
        stream nextPutAll: x asString.
        stream nextPut: '@'.
        stream nextPutAll: y asString.
    ]
]
```

### Messages

#### Unary Messages
```smog
object message
myArray size
```

#### Binary Messages
```smog
3 + 4
x < y
a & b
```

#### Keyword Messages
```smog
point x: 10 y: 20
array at: 1 put: 'value'
```

#### Message Precedence
1. Unary messages (highest)
2. Binary messages
3. Keyword messages (lowest)

Parentheses can be used to control precedence:
```smog
(3 + 4) * 5
array at: (index + 1)
```

### Blocks (Closures)

Blocks are anonymous functions:

```smog
[ "empty block" ]
[ :x | x + 1 ]
[ :x :y | x + y ]
```

Blocks are evaluated with the `value` message:
```smog
| block result |
block := [ :x | x * 2 ].
result := block value: 5.  " result is 10 "
```

### Control Structures

Control structures are implemented as messages to booleans and blocks:

#### Conditionals
```smog
x > 0 ifTrue: [ 'positive' ].
x > 0 ifFalse: [ 'not positive' ].
x > 0 
    ifTrue: [ 'positive' ]
    ifFalse: [ 'not positive' ].
```

#### Loops
```smog
"While loop"
[ x < 10 ] whileTrue: [ 
    x := x + 1.
].

"Times loop"
5 timesRepeat: [ 
    'hello' println.
].

"Collection iteration"
#(1 2 3 4 5) do: [ :each |
    each println.
].
```

### Return Statements

Methods return the last expression by default. Explicit return uses `^`:

```smog
factorial: n [
    n <= 1 ifTrue: [ ^1 ].
    ^n * (self factorial: (n - 1))
]
```

#### Non-Local Returns

In Smalltalk-style languages, return statements (`^`) in blocks perform **non-local returns** - they return from the method that created the block, not just from the block itself. This is a fundamental feature for control flow.

```smog
findFirst: predicate [
    " Returns the first element matching the predicate, or nil "
    self do: [ :each |
        (predicate value: each) ifTrue: [
            ^each    " Returns from findFirst:, not just from the ifTrue: block "
        ]
    ].
    ^nil
]
```

In this example, when `^each` executes inside the nested blocks (`do:` → `ifTrue:`), it returns `each` from the `findFirst:` method, skipping the rest of the iteration and the `^nil` at the end.

**Key Points:**
- `^expression` in a method returns from that method
- `^expression` in a block returns from the method that created the block
- This makes control flow constructs like `ifTrue:`, `ifFalse:`, `whileTrue:` work naturally
- Blocks without explicit returns yield the value of their last expression

**Example with nested blocks:**
```smog
method [
    (true) ifTrue: [
        (true) ifTrue: [
            ^42    " Returns 42 from method, not from either ifTrue: block "
        ]
    ].
    'This will not execute' println.
    ^99
]
```

The method returns `42`, not `99`, because the non-local return exits the entire method.

## Standard Library

### Core Classes

- **Object** - Root of the class hierarchy
- **Class** - Metaclass for all classes
- **Boolean** - Abstract class for true/false
- **True** - The singleton true object
- **False** - The singleton false object
- **Integer** - Integer numbers
- **Double** - Floating point numbers
- **String** - Character strings
- **Array** - Fixed-size arrays
- **Block** - Closures/anonymous functions
- **Nil** - The singleton nil object

### Common Methods

#### Object
- `class` - Returns the receiver's class
- `== anObject` - Identity comparison
- `= anObject` - Equality comparison
- `println` - Print the object followed by newline
- `asString` - Convert to string

#### Integer/Double
- `+ aNumber` - Addition
- `- aNumber` - Subtraction
- `* aNumber` - Multiplication
- `/ aNumber` - Division
- `< aNumber` - Less than
- `> aNumber` - Greater than
- `<= aNumber` - Less than or equal
- `>= aNumber` - Greater than or equal

#### String
- `length` - String length
- `at: index` - Character at index
- `, aString` - Concatenation

#### Array
- `size` - Array size
- `at: index` - Element at index
- `at: index put: value` - Set element at index
- `do: aBlock` - Iterate over elements

#### Block
- `value` - Evaluate with no arguments
- `value: arg` - Evaluate with one argument
- `value: arg1 value: arg2` - Evaluate with two arguments

## File Structure

A smog source file (`.smog`) contains class definitions and optionally a main execution block:

```smog
Object subclass: #Hello [
    greet [
        'Hello, Smog!' println.
    ]
]

"Main execution"
Hello new greet.
```

### Bytecode Files (.sg)

Smog source files can be compiled to binary bytecode files with the `.sg` extension:

```bash
# Compile source to bytecode
smog compile hello.smog hello.sg

# Run bytecode directly
smog hello.sg
```

**Bytecode File Format:**
- Magic number: "SMOG" (0x534D4F47)
- Version tracking for compatibility
- Binary encoding of instructions and constants
- Includes all class definitions and methods

**Use cases:**
- Faster program loading (5-50x improvement)
- Distribution without source code
- Multi-file programs with pre-compiled modules
- Production deployment

See the [Bytecode Format Guide](../BYTECODE_FORMAT.md) for complete specification.

## Examples

### Hello World
```smog
'Hello, World!' println.
```

### Factorial
```smog
Object subclass: #Math [
    factorial: n [
        n <= 1 ifTrue: [ ^1 ].
        ^n * (self factorial: (n - 1))
    ]
]

Math new factorial: 5.  " Returns 120 "
```

### Counter Class
```smog
Object subclass: #Counter [
    | count |
    
    initialize [
        count := 0.
    ]
    
    increment [
        count := count + 1.
    ]
    
    value [
        ^count
    ]
]

| counter |
counter := Counter new.
counter initialize.
counter increment.
counter increment.
counter value println.  " Prints: 2 "
```

## Implementation Notes

### Compilation Phases

1. **Lexing** - Source text → Tokens
2. **Parsing** - Tokens → AST
3. **Compilation** - AST → Bytecode
4. **Optional: Serialization** - Bytecode → .sg file
5. **Execution** - Bytecode (from source or .sg) → VM execution

### Bytecode Format

The bytecode is a sequence of instructions operating on a stack-based virtual machine. Bytecode can be:
- **Executed directly** from memory after compilation
- **Saved to .sg files** for faster loading
- **Loaded from .sg files** and executed without recompilation

See the [Bytecode Format documentation](../BYTECODE_FORMAT.md) for details on:
- Binary file structure
- Instruction encoding
- Constant pool format
- Class and method serialization

## Future Considerations

- Module system
- Exception handling
- Concurrency primitives
- Foreign function interface
- Optimizing JIT compiler
- Garbage collector enhancements
