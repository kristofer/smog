# Smog Learning Guide: A Beginner's Mental Model

## Introduction

Welcome to programming with Smog! This guide is designed for AI-era beginners who want to understand not just *how* to write code, but *why* programming systems work the way they do. We'll build a clear mental model of Smog from the ground up.

## What is a Programming Language?

Think of a programming language as a way to communicate instructions to a computer. But there's a journey your code takes before the computer can execute it:

```
Your Ideas → Smog Code → Computer Understanding → Results
```

This guide explains every step of that journey.

## The Big Picture: How Smog Works

### The Four-Stage Pipeline

Smog transforms your code through four stages:

```
┌──────────┐    ┌────────┐    ┌──────────┐    ┌────────┐
│  Source  │ →  │ Tokens │ →  │   AST    │ →  │Bytecode│ → Results
│   Code   │    │ Stream │    │  (Tree)  │    │        │
└──────────┘    └────────┘    └──────────┘    └────────┘
   Lexer          Parser        Compiler          VM
```

Let's understand each stage with a simple example:

**Your Code:**
```smog
x := 10 + 5.
```

**Stage 1: Lexer (Breaking into Pieces)**
```
[identifier('x'), assign(':='), number('10'), plus('+'), number('5'), period('.')]
```
*Like breaking a sentence into words*

**Stage 2: Parser (Understanding Structure)**
```
Assignment
├── target: x
└── value: Add
    ├── left: 10
    └── right: 5
```
*Like understanding grammar in a sentence*

**Stage 3: Compiler (Creating Instructions)**
```
PUSH 10
PUSH 5
ADD
STORE x
```
*Like translating to a language the computer understands*

**Stage 4: VM (Execution)**
```
Stack: [] → [10] → [10, 5] → [15]
Store 15 in x
```
*Like actually performing the actions*

## Mental Model #1: Everything is an Object

In Smog, **everything** is an object. This is simpler than it sounds!

### What's an Object?

An object is:
1. **Data** (information it stores)
2. **Behavior** (what it can do)
3. **Identity** (it's unique)

**Examples:**

```smog
" The number 5 is an object "
5           " Data: the value 5 "
            " Behavior: can add, subtract, etc. "
            " Identity: this specific 5 "

" A string is an object "
'hello'     " Data: the characters "
            " Behavior: can get length, uppercase, etc. "
            " Identity: this specific string "

" Even true/false are objects! "
true        " Data: truth value "
            " Behavior: can make decisions "
            " Identity: THE true object "
```

### Why Does This Matter?

Because in Smog, you never work directly with "raw data." You always work with objects by sending them messages.

## Mental Model #2: Communication Through Messages

In Smog, you don't "call functions" or "use operators." You **send messages to objects**.

### The Message Sending Pattern

```
receiver message
```

**Examples:**

```smog
" Ask an array for its size "
#(1 2 3) size         " Array receives 'size' message "

" Ask a number to add another number "
10 + 5                " 10 receives '+' message with argument 5 "

" Ask a string to print itself "
'hello' println       " String receives 'println' message "
```

### Three Types of Messages

**1. Unary Messages** (no arguments - just the message name)
```smog
array size
object class
counter increment
```

**2. Binary Messages** (one argument - operator-like)
```smog
3 + 4        " Send + with argument 4 to receiver 3 "
10 - 2       " Send - with argument 2 to receiver 10 "
5 < 10       " Send < with argument 10 to receiver 5 "
```

**3. Keyword Messages** (multiple arguments - descriptive names)
```smog
array at: 1 put: 'value'          " Two arguments: index and value "
point x: 10 y: 20                 " Two arguments: x and y "
string copyFrom: 0 to: 5          " Two arguments: start and end "
```

### Message Precedence (Order of Operations)

Like math has order (multiplication before addition), messages have order:

1. **Unary** (highest priority)
2. **Binary**
3. **Keyword** (lowest priority)

**Example:**
```smog
array size + 1 * 2
```

**Evaluation order:**
1. `array size` (unary - get size, let's say 5)
2. `5 + 1` (binary - equals 6)
3. `6 * 2` (binary - equals 12)

Use parentheses to change order:
```smog
array size + (1 * 2)    " size + 2 "
(array size + 1) * 2    " (size + 1) * 2 "
```

## Mental Model #3: Blocks Are Frozen Code

Blocks are pieces of code you can save and execute later. Think of them as recipes you write down and use when needed.

### Creating a Block

```smog
[ 'Hello!' println ]
```

This creates a block but **doesn't execute it** yet. It's like writing a recipe without cooking.

### Executing a Block

```smog
| greet |
greet := [ 'Hello!' println ].
greet value.                      " Now it executes "
```

### Blocks with Parameters

Blocks can accept inputs (parameters):

```smog
| square |
square := [ :x | x * x ].         " :x is the parameter "
square value: 5.                  " Pass 5 as x, result is 25 "
```

**Think of it like:**
- Recipe with ingredients: "Take ingredient :x, multiply it by itself"
- When you cook: "Use 5 as the ingredient"

### Why Blocks Matter

Blocks enable:
1. **Reusable code** - Write once, use many times
2. **Control flow** - if/then/else, loops
3. **Higher-order programming** - Pass code as data

**Example - Control Flow:**
```smog
x > 0 ifTrue: [ 'positive' println ].
```

Here's what happens:
1. `x > 0` sends `>` message, returns `true` or `false`
2. `ifTrue:` is sent to that true/false object
3. The block `[ 'positive' println ]` is passed as an argument
4. If true, the object executes the block
5. If false, the object ignores the block

## Mental Model #4: Classes Are Object Factories

Classes are blueprints for creating objects.

### Defining a Class

```smog
Object subclass: #Counter [
    | count |              " Instance variable "
    
    initialize [           " Constructor "
        count := 0.
    ]
    
    increment [            " Method "
        count := count + 1.
    ]
    
    value [                " Getter "
        ^count             " ^ means return "
    ]
]
```

**Mental Model:**
- `Object subclass:` - "Create a new kind of object"
- `| count |` - "Each instance will have its own count"
- `initialize` - "What to do when creating an instance"
- Methods - "What instances can do"

### Using a Class

```smog
| counter1 counter2 |
counter1 := Counter new.    " Create first counter "
counter1 initialize.        " Set it up "

counter2 := Counter new.    " Create second counter "
counter2 initialize.        " Set it up separately "

counter1 increment.         " Only affects counter1 "
counter1 value println.     " Prints: 1 "
counter2 value println.     " Prints: 0 (independent) "
```

## Mental Model #5: The Stack Machine

Understanding how code executes helps you write better code.

### The Stack Concept

Think of a stack like a stack of plates:
- **Push**: Add plate on top
- **Pop**: Remove plate from top
- You can only access the top plate

### Example Execution

**Code:**
```smog
3 + 4 * 5
```

**Stack Evolution:**
```
Step 1: PUSH 3        Stack: [3]
Step 2: PUSH 4        Stack: [3, 4]
Step 3: PUSH 5        Stack: [3, 4, 5]
Step 4: MULTIPLY      Stack: [3, 20]      (pop 5 and 4, push 20)
Step 5: ADD           Stack: [23]         (pop 20 and 3, push 23)
```

**Why This Matters:**
- Helps you understand evaluation order
- Makes debugging easier
- Explains why some operations are faster than others

## Mental Model #6: Variables Are Labels

Variables don't "contain" values - they're labels pointing to objects.

```smog
| x y |
x := 10.
y := x.
x := 20.

x println.    " Prints: 20 "
y println.    " Prints: 10 "
```

**What Happened:**
1. `x := 10` - x points to object 10
2. `y := x` - y points to same object 10
3. `x := 20` - x now points to object 20
4. y still points to object 10

**Mental Model:**
```
Step 1:  x ──→ [10]
         y

Step 2:  x ──→ [10]
         y ──→ [10]

Step 3:  x ──→ [20]
         y ──→ [10]
```

## Learning Path: From Beginner to Advanced

### Level 1: Basics (Start Here)

**What to Learn:**
- Variables and assignment
- Basic message sending
- Simple blocks
- Arrays

**Practice Project:** Calculator
```smog
Object subclass: #Calculator [
    add: a to: b [
        ^a + b
    ]
    
    subtract: a from: b [
        ^b - a
    ]
]
```

### Level 2: Control Flow

**What to Learn:**
- Conditionals (ifTrue:, ifFalse:)
- Loops (timesRepeat:, do:)
- Collection iteration

**Practice Project:** Number Guesser
```smog
Object subclass: #GuessingGame [
    | secret guesses |
    
    initialize [
        secret := 42.
        guesses := 0.
    ]
    
    guess: number [
        guesses := guesses + 1.
        number = secret ifTrue: [ 
            ^'Correct! Guesses: ' + guesses asString
        ].
        number < secret ifTrue: [ ^'Too low' ].
        ^'Too high'
    ]
]
```

### Level 3: Object-Oriented Design

**What to Learn:**
- Class design
- Instance variables
- Methods and encapsulation
- Inheritance

**Practice Project:** Shape Hierarchy
```smog
Object subclass: #Shape [
    area [
        self subclassResponsibility.
    ]
]

Shape subclass: #Circle [
    | radius |
    
    radius: r [
        radius := r.
    ]
    
    area [
        ^3.14159 * radius * radius
    ]
]
```

### Level 4: Advanced Patterns

**What to Learn:**
- Higher-order functions
- Closures and captured variables
- Design patterns
- Algorithms

**Practice Project:** Collection Library
```smog
Object subclass: #MyCollection [
    map: transformBlock [
        " Transform each element "
    ]
    
    filter: predicateBlock [
        " Keep only elements matching predicate "
    ]
    
    reduce: binaryBlock [
        " Combine all elements "
    ]
]
```

## Common Beginner Mistakes

### Mistake 1: Forgetting the Period

```smog
" Wrong "
x := 10
y := 20

" Right "
x := 10.
y := 20.
```

### Mistake 2: Wrong Message Precedence

```smog
" Wrong - parses as: 10 + (array size) "
10 + array size

" Right - parses as: (10 + (array size)) "
10 + (array size)

" Or use parentheses to be clear "
(array size) + 10
```

### Mistake 3: Forgetting to Initialize

```smog
" Wrong "
| counter |
counter := Counter new.
counter increment.    " count is nil! "

" Right "
| counter |
counter := Counter new.
counter initialize.   " Now count is 0 "
counter increment.
```

### Mistake 4: Confusing := and =

```smog
" := is assignment "
x := 10.

" = is comparison "
x = 10 ifTrue: [ 'yes' println ].
```

## Understanding Error Messages

### "Stack Underflow"

**What it means:** Tried to pop from empty stack

**Common cause:** Mismatched operations

```smog
" Wrong - tries to add but only one value on stack "
+ 5

" Right "
10 + 5
```

### "Undefined Variable"

**What it means:** Variable not declared

**Fix:** Declare in variable list

```smog
" Wrong "
x := 10.    " x not declared "

" Right "
| x |
x := 10.
```

### "Message Not Understood"

**What it means:** Object doesn't have that method

**Common cause:** Typo or wrong object type

```smog
" Wrong "
5 length    " Numbers don't have length "

" Right "
'hello' length    " Strings do "
```

## Study Strategy

### 1. Read Code Daily

Start with small examples in `examples/` directory:
- `hello.smog` - Basic I/O
- `factorial.smog` - Recursion
- `counter.smog` - Classes

### 2. Type Out Examples

Don't copy-paste! Type each example:
- Builds muscle memory
- Catches details you'd miss
- Helps internalize syntax

### 3. Modify Examples

After typing, make changes:
- What if I change this number?
- What if I add another method?
- Can I combine two examples?

### 4. Build Small Projects

Start simple:
1. **Week 1:** Calculator (add, subtract, multiply, divide)
2. **Week 2:** Todo list (add, remove, list items)
3. **Week 3:** Simple game (number guessing)
4. **Week 4:** Data structure (stack or queue)

### 5. Read Documentation Progressively

Don't read everything at once! Follow this order:

1. This guide (you're here!)
2. [User's Guide](USERS_GUIDE.md) - Practical examples
3. [Language Spec](spec/LANGUAGE_SPEC.md) - Complete syntax
4. Technical docs (Lexer, Parser, etc.) - Only when curious

## Debugging Strategies

### 1. Print Debugging

Add println everywhere to see what's happening:

```smog
| x y |
x := 10.
'x is: ' print.
x println.           " See value of x "

y := x + 5.
'y is: ' print.
y println.           " See value of y "
```

### 2. Simplify

If code doesn't work, make it simpler:

```smog
" Complex - hard to debug "
result := ((array at: index) + offset) * multiplier / divisor.

" Simplified - see each step "
| element adjusted scaled |
element := array at: index.
'element: ' print. element println.

adjusted := element + offset.
'adjusted: ' print. adjusted println.

scaled := adjusted * multiplier.
'scaled: ' print. scaled println.

result := scaled / divisor.
'result: ' print. result println.
```

### 3. Check Types

Make sure objects are what you think:

```smog
| value |
value := 'hello'.
value class println.    " Prints: String "

value := 42.
value class println.    " Prints: Integer "
```

## Conceptual Analogies

### Programming Concepts as Everyday Things

**Variables** = Labels on boxes
- The box contains the object
- The label helps you find it
- You can relabel boxes

**Classes** = Cookie cutters
- Define the shape
- Make many cookies (objects)
- Each cookie is independent

**Methods** = Recipes in a cookbook
- Cookbook (class) has many recipes (methods)
- Each recipe has steps
- You can follow the same recipe many times

**Blocks** = Instruction cards
- Write instructions on a card
- Hand card to someone
- They follow it when ready

**Messages** = Asking favors
- You ask someone (object) to do something (method)
- They do it and maybe give you something back (return value)

## Next Steps

Now that you have a mental model:

1. **Try the examples** in [User's Guide](USERS_GUIDE.md)
2. **Build something small** - a calculator or counter
3. **Read other people's code** - learn from examples
4. **Experiment** - change things and see what happens
5. **Ask questions** - why does this work? What if I change it?

## Resources for Deeper Learning

**Within Smog:**
- [User's Guide](USERS_GUIDE.md) - Practical programming
- [Language Spec](spec/LANGUAGE_SPEC.md) - Complete reference
- [Examples](../examples/) - Working code to study

**Understanding the System:**
- [Lexer Guide](LEXER.md) - How code becomes tokens
- [Parser Guide](PARSER.md) - How tokens become structure
- [Compiler Guide](COMPILER.md) - How structure becomes bytecode
- [VM Guide](VM_DEEP_DIVE.md) - How bytecode executes

**Smalltalk Resources** (Smog's inspiration):
- Smalltalk-80 Blue Book
- Pharo by Example
- Squeak tutorials

## Final Thoughts

Programming is a skill that develops over time. You'll make mistakes - that's how you learn! The key is to:

1. **Start simple** - master basics before complexity
2. **Practice regularly** - a little daily beats cramming
3. **Read and write** - both are essential
4. **Experiment** - break things and fix them
5. **Be patient** - everyone was a beginner once

Welcome to the world of Smog programming. Enjoy the journey!
