# Smog User's Guide

## Introduction

Welcome to Smog! This guide will teach you how to use the Smog programming language for common programming tasks. Whether you're learning programming for the first time or coming from another language, this guide provides practical examples and patterns you can use right away.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Concepts](#basic-concepts)
3. [Built-in Methods Reference](#built-in-methods-reference)
4. [Data Structures](#data-structures)
5. [Sorting and Searching](#sorting-and-searching)
6. [Object-Oriented Patterns](#object-oriented-patterns)
7. [Common Algorithms](#common-algorithms)
8. [Best Practices](#best-practices)

## Getting Started

### Your First Program

Create a file called `hello.smog`:

```smog
'Hello, World!' println.
```

Run it:
```bash
./bin/smog hello.smog
```

Output:
```
Hello, World!
```

### Compiling for Faster Execution

For frequently-run programs, compile to bytecode:

```bash
# Compile to .sg bytecode file
./bin/smog compile hello.smog hello.sg

# Run the compiled bytecode (faster startup)
./bin/smog hello.sg

# Inspect the bytecode
./bin/smog disassemble hello.sg
```

**Benefits of .sg files:**
- 5-50x faster startup time
- Distribute programs without source code
- Foundation for module systems

See the [Bytecode Format Guide](BYTECODE_FORMAT.md) for details.

### Working with Variables

```smog
" Declare local variables "
| name age city |

name := 'Alice'.
age := 30.
city := 'San Francisco'.

name println.
age println.
city println.
```

**Key Points:**
- Variables declared between pipes: `| var1 var2 |`
- Assignment uses `:=`
- Statements end with period `.`
- **All variables must be declared in a single block at the beginning of each scope**

**Important Scoping Limitation:**

Due to the current implementation, you must declare all variables you will use in a single declaration block at the top of each scope. You cannot add new variable declarations after creating blocks or executing statements.

```smog
" ✅ CORRECT: All variables declared together "
| x y result temp |
x := 10.
y := 20.
[ temp := x * 2 ] value.
result := temp + y.

" ❌ INCORRECT: Multiple declaration blocks "
| x y |
x := 10.
[ x + 5 ] value.
| result |  " ERROR: Cannot declare here
result := x.
```

**Best Practice:** Plan ahead and declare all variables you will need at the beginning of your code.

## Basic Concepts

### 1. Everything is an Object

In Smog, even numbers and booleans are objects:

```smog
" Numbers are objects "
5 class println.          " Prints: Integer "

" Booleans are objects "
true class println.       " Prints: True "

" You send messages to objects "
10 + 5 println.           " Message + sent to 10 "
'hello' length println.   " Message length sent to 'hello' "
```

### 2. Message Sending

All computation happens through sending messages to objects:

**Unary Messages** (no arguments):
```smog
| array |
array := #(1 2 3 4 5).
array size println.       " Prints: 5 "
```

**Binary Messages** (one argument, operator-like):
```smog
10 + 5 println.          " Addition "
20 - 3 println.          " Subtraction "
4 * 6 println.           " Multiplication "
10 / 2 println.          " Division "
10 \\ 3 println.         " Modulo (remainder): 1 "
5 < 10 println.          " Comparison "
```

**Note on Operators:**
- `+`, `-`, `*`, `/` - Standard arithmetic operators
- `\\` - Modulo operator (returns remainder after division)
  - Example: `10 \\ 3` returns `1` (10 ÷ 3 = 3 remainder 1)
  - Useful for checking even/odd: `(n \\ 2) = 0` means n is even
- `<`, `>`, `<=`, `>=` - Comparison operators
- `=` - Equality check
- `~=` - Not equal

**Keyword Messages** (multiple parts):
```smog
| array |
array := #(1 2 3 4 5).
array at: 1 put: 10.     " Set element at index 1 "
array at: 1 println.     " Get element at index 1 "
```

### 3. Blocks (Anonymous Functions)

Blocks are reusable pieces of code:

```smog
" Simple block "
| greet |
greet := [ 'Hello!' println ].
greet value.

" Block with parameters "
| double |
double := [ :x | x * 2 ].
(double value: 5) println.    " Prints: 10 "

" Block with multiple parameters "
| add |
add := [ :x :y | x + y ].
(add value: 3 value: 7) println.  " Prints: 10 "
```

### 4. Control Flow

Control flow uses blocks and message passing:

```smog
| x |
x := 10.

" Conditional "
x > 0 ifTrue: [ 'positive' println ].
x < 0 ifFalse: [ 'not negative' println ].

" if-then-else "
x > 0
    ifTrue: [ 'positive' println ]
    ifFalse: [ 'not positive' println ].

" Loops "
5 timesRepeat: [ 'hello' println ].

#(1 2 3) do: [ :each |
    each println.
].
```

## Built-in Methods Reference

Smog provides built-in methods for core types. All computation happens by sending messages to objects.

### Boolean Methods

Booleans (`true` and `false`) respond to conditional messages:

#### `ifTrue: aBlock`
Execute the block if the boolean is true.
```smog
true ifTrue: [ 'This executes' println ].
false ifTrue: [ 'This does not execute' println ].
```

#### `ifFalse: aBlock`
Execute the block if the boolean is false.
```smog
false ifFalse: [ 'This executes' println ].
true ifFalse: [ 'This does not execute' println ].
```

#### `ifTrue: trueBlock ifFalse: falseBlock`
Execute the first block if true, otherwise execute the second block.
```smog
| x result |
x := 10.
result := (x > 5)
    ifTrue: [ 'greater' ]
    ifFalse: [ 'not greater' ].
result println.  " Prints: greater "
```

### Integer Methods

Integers support arithmetic, comparison, and iteration messages:

#### Arithmetic Operations
- `+ other` - Addition
- `- other` - Subtraction
- `* other` - Multiplication
- `/ other` - Division (integer division)
- `\\ other` - Modulo (remainder)

```smog
10 + 5 println.   " Prints: 15 "
10 - 3 println.   " Prints: 7 "
10 * 2 println.   " Prints: 20 "
10 / 3 println.   " Prints: 3 "
10 \\ 3 println.  " Prints: 1 "
```

#### Comparison Operations
- `< other` - Less than
- `> other` - Greater than
- `<= other` - Less than or equal
- `>= other` - Greater than or equal
- `= other` - Equal to
- `~= other` - Not equal to

```smog
5 < 10 println.   " Prints: true "
5 > 10 println.   " Prints: false "
5 = 5 println.    " Prints: true "
5 ~= 3 println.   " Prints: true "
```

#### Iteration Methods

#### `timesRepeat: aBlock`
Execute the block N times.
```smog
3 timesRepeat: [ 'Hello' println ].
" Prints Hello three times "
```

### String Methods

Strings support printing and comparison:

#### `println`
Print the string followed by a newline.
```smog
'Hello, World!' println.
```

#### `print`
Print the string without a newline.
```smog
'Name: ' print.
'Alice' println.
" Prints: Name: Alice "
```

#### Comparison
Strings support `=` and `~=` for equality testing.
```smog
'hello' = 'hello' println.  " Prints: true "
'hello' = 'world' println.  " Prints: false "
```

### Array Methods

Arrays are ordered collections of elements:

#### `size`
Return the number of elements in the array.
```smog
| arr |
arr := #(1 2 3 4 5).
arr size println.  " Prints: 5 "
```

#### `at: index`
Get the element at the given index (1-based indexing).
```smog
| arr elem |
arr := #(10 20 30).
elem := arr at: 1.
elem println.  " Prints: 10 "

elem := arr at: 2.
elem println.  " Prints: 20 "
```

**Note:** Arrays in Smog use 1-based indexing (like Smalltalk), not 0-based indexing.

#### `at: index put: value`
Set the element at the given index (1-based indexing).
```smog
| arr |
arr := #(1 2 3).
arr at: 2 put: 99.
arr at: 2 println.  " Prints: 99 "
```

#### `do: aBlock`
Iterate over each element, executing the block with each element as a parameter.
```smog
| numbers |
numbers := #(1 2 3 4 5).
numbers do: [ :each |
    each println.
].
" Prints: 1 2 3 4 5 (each on a new line) "
```

#### Advanced Array Methods

The following methods are commonly implemented in user code (not built-in, but shown as patterns):

##### `collect: transformBlock`
Transform each element and return a new array (also known as "map").
```smog
| numbers doubled |
numbers := #(1 2 3 4 5).
doubled := numbers collect: [ :each | each * 2 ].
doubled do: [ :each | each println ].
" Prints: 2 4 6 8 10 "
```

##### `select: predicateBlock`
Filter elements that satisfy a condition (also known as "filter").
```smog
| numbers evens |
numbers := #(1 2 3 4 5 6).
evens := numbers select: [ :each | (each \\ 2) = 0 ].
evens do: [ :each | each println ].
" Prints: 2 4 6 "
```

##### `inject: initialValue into: binaryBlock`
Reduce the array to a single value (also known as "fold" or "reduce").
```smog
| numbers sum |
numbers := #(1 2 3 4 5).
sum := numbers inject: 0 into: [ :acc :each | acc + each ].
sum println.  " Prints: 15 "
```

**Note:** The `collect:`, `select:`, and `inject:into:` methods are patterns you implement in your own classes, not built-in VM operations. See the [Data Structures](#data-structures) section for examples.

### Block Methods

Blocks (closures/anonymous functions) respond to value messages:

#### `value`
Execute a block with no parameters.
```smog
| greet |
greet := [ 'Hello!' println ].
greet value.  " Prints: Hello! "
```

#### `value: arg`
Execute a block with one parameter.
```smog
| double |
double := [ :x | x * 2 ].
(double value: 5) println.  " Prints: 10 "
```

#### `value: arg1 value: arg2`
Execute a block with two parameters.
```smog
| add |
add := [ :x :y | x + y ].
(add value: 3 value: 7) println.  " Prints: 10 "
```

Blocks with more parameters follow the same pattern: `value: a value: b value: c` etc.

#### `whileTrue: aBlock`
Execute the receiver block, and while it returns true, execute the argument block.
```smog
i := 1.
[i <= 5] whileTrue: [
    i println.
    i := i + 1.
].
" Prints: 1 2 3 4 5 "
```

**Note:** Due to current closure limitations, variables accessed in loops should be global (not declared with `| var |`). Blocks can access and modify global variables but not local variables from the enclosing scope.

#### `whileFalse: aBlock`
Execute the receiver block, and while it returns false, execute the argument block.
```smog
i := 1.
[i > 5] whileFalse: [
    i println.
    i := i + 1.
].
" Prints: 1 2 3 4 5 "
```

### Class Methods

All classes respond to:

#### `new`
Create a new instance of the class.
```smog
Object subclass: #Person [
    | name |
].

| person |
person := Person new.
```

### Object Methods

All objects inherit from `Object` and respond to:

#### `class`
Return the class of the object.
```smog
5 class println.         " Prints: Integer "
'hello' class println.   " Prints: String "
true class println.      " Prints: True "
```

#### `println`
Print the object followed by a newline.
```smog
42 println.
'text' println.
true println.
```

#### `print`
Print the object without a newline.
```smog
'Answer: ' print.
42 println.
" Prints: Answer: 42 "
```

### Method Lookup and User-Defined Classes

When you define your own classes, you can add methods that override or extend the built-in behavior:

```smog
Object subclass: #Point [
    | x y |

    " Custom initialization "
    x: xVal y: yVal [
        x := xVal.
        y := yVal.
    ]

    " Custom println "
    println [
        '(' print.
        x print.
        ', ' print.
        y print.
        ')' println.
    ]
]

| p |
p := Point new.
p x: 10 y: 20.
p println.  " Prints: (10, 20) "
```

## Data Structures

### Arrays

Arrays store ordered collections of elements:

**Creating Arrays:**
```smog
" Array literals "
| numbers names mixed |
numbers := #(1 2 3 4 5).
names := #('Alice' 'Bob' 'Charlie').
mixed := #(1 'hello' true 3.14).
```

**Accessing Elements:**
```smog
| array |
array := #(10 20 30 40 50).

" Get element "
(array at: 0) println.    " Prints: 10 (first element) "
(array at: 2) println.    " Prints: 30 (third element) "

" Set element "
array at: 1 put: 99.
(array at: 1) println.    " Prints: 99 "

" Array size "
array size println.       " Prints: 5 "
```

**Iterating Arrays:**
```smog
| numbers |
numbers := #(1 2 3 4 5).

" Print each element "
numbers do: [ :each |
    each println.
].

" Sum all elements "
| sum |
sum := 0.
numbers do: [ :each |
    sum := sum + each.
].
sum println.  " Prints: 15 "
```

**Array Operations:**
```smog
| numbers |
numbers := #(1 2 3 4 5).

" Transform (map) "
| doubled |
doubled := numbers collect: [ :each |
    each * 2
].
doubled do: [ :each | each println ].  " Prints: 2 4 6 8 10 "

" Filter (select) "
| evens |
evens := numbers select: [ :each |
    (each \\ 2) = 0
].
evens do: [ :each | each println ].  " Prints: 2 4 "

" Reduce (inject) "
| product |
product := numbers inject: 1 into: [ :acc :each |
    acc * each
].
product println.  " Prints: 120 (1*2*3*4*5) "
```

### Building a Stack

```smog
Object subclass: #Stack [
    | items |
    
    initialize [
        items := #().
    ]
    
    push: item [
        items := items copyWith: item.
    ]
    
    pop [
        | item |
        items size = 0 ifTrue: [ ^nil ].
        item := items at: (items size - 1).
        items := items copyFrom: 0 to: (items size - 2).
        ^item
    ]
    
    peek [
        items size = 0 ifTrue: [ ^nil ].
        ^items at: (items size - 1)
    ]
    
    isEmpty [
        ^items size = 0
    ]
    
    size [
        ^items size
    ]
]

" Using the stack "
| stack |
stack := Stack new.
stack initialize.

stack push: 10.
stack push: 20.
stack push: 30.

stack size println.      " Prints: 3 "
stack peek println.      " Prints: 30 "
stack pop println.       " Prints: 30 "
stack size println.      " Prints: 2 "
```

### Building a Queue

```smog
Object subclass: #Queue [
    | items |
    
    initialize [
        items := #().
    ]
    
    enqueue: item [
        items := items copyWith: item.
    ]
    
    dequeue [
        | item |
        items size = 0 ifTrue: [ ^nil ].
        item := items at: 0.
        items := items copyFrom: 1 to: (items size - 1).
        ^item
    ]
    
    isEmpty [
        ^items size = 0
    ]
    
    size [
        ^items size
    ]
]

" Using the queue "
| queue |
queue := Queue new.
queue initialize.

queue enqueue: 'first'.
queue enqueue: 'second'.
queue enqueue: 'third'.

queue dequeue println.   " Prints: first "
queue dequeue println.   " Prints: second "
queue size println.      " Prints: 1 "
```

## Sorting and Searching

### Bubble Sort

```smog
Object subclass: #Sorter [
    bubbleSort: array [
        | n swapped temp |
        n := array size.
        
        [ true ] whileTrue: [
            swapped := false.
            
            1 to: (n - 1) do: [ :i |
                (array at: (i - 1)) > (array at: i) ifTrue: [
                    temp := array at: (i - 1).
                    array at: (i - 1) put: (array at: i).
                    array at: i put: temp.
                    swapped := true.
                ]
            ].
            
            swapped ifFalse: [ ^array ].
        ].
    ]
]

" Usage "
| sorter numbers |
sorter := Sorter new.
numbers := #(64 34 25 12 22 11 90).

'Before sorting:' println.
numbers do: [ :each | each println ].

sorter bubbleSort: numbers.

'After sorting:' println.
numbers do: [ :each | each println ].
```

### Quick Sort

```smog
Object subclass: #QuickSorter [
    sort: array [
        ^self quickSort: array from: 0 to: (array size - 1)
    ]
    
    quickSort: array from: low to: high [
        | pi |
        low < high ifTrue: [
            pi := self partition: array from: low to: high.
            self quickSort: array from: low to: (pi - 1).
            self quickSort: array from: (pi + 1) to: high.
        ].
        ^array
    ]
    
    partition: array from: low to: high [
        | pivot i j temp |
        pivot := array at: high.
        i := low - 1.
        
        low to: (high - 1) do: [ :j |
            (array at: j) <= pivot ifTrue: [
                i := i + 1.
                temp := array at: i.
                array at: i put: (array at: j).
                array at: j put: temp.
            ]
        ].
        
        temp := array at: (i + 1).
        array at: (i + 1) put: (array at: high).
        array at: high put: temp.
        
        ^i + 1
    ]
]
```

### Binary Search

```smog
Object subclass: #Searcher [
    binarySearch: array for: target [
        ^self search: array for: target from: 0 to: (array size - 1)
    ]
    
    search: array for: target from: low to: high [
        | mid midValue |
        
        low > high ifTrue: [ ^-1 ].
        
        mid := (low + high) // 2.
        midValue := array at: mid.
        
        midValue = target ifTrue: [ ^mid ].
        midValue > target ifTrue: [ 
            ^self search: array for: target from: low to: (mid - 1)
        ].
        ^self search: array for: target from: (mid + 1) to: high
    ]
]

" Usage "
| searcher sortedArray index |
searcher := Searcher new.
sortedArray := #(1 3 5 7 9 11 13 15 17 19).

index := searcher binarySearch: sortedArray for: 7.
index println.  " Prints: 3 "

index := searcher binarySearch: sortedArray for: 8.
index println.  " Prints: -1 (not found) "
```

### Linear Search

```smog
Object subclass: #LinearSearcher [
    search: array for: target [
        | index |
        index := 0.
        
        array do: [ :each |
            each = target ifTrue: [ ^index ].
            index := index + 1.
        ].
        
        ^-1  " Not found "
    ]
]
```

## Object-Oriented Patterns

### 1. Encapsulation (Data Hiding)

```smog
Object subclass: #BankAccount [
    | balance |
    
    initialize [
        balance := 0.
    ]
    
    deposit: amount [
        amount > 0 ifTrue: [
            balance := balance + amount.
        ].
    ]
    
    withdraw: amount [
        amount > 0 ifTrue: [
            balance >= amount ifTrue: [
                balance := balance - amount.
                ^true
            ].
        ].
        ^false
    ]
    
    getBalance [
        ^balance
    ]
]

" Usage "
| account |
account := BankAccount new.
account initialize.

account deposit: 100.
account deposit: 50.
account withdraw: 30.

account getBalance println.  " Prints: 120 "
```

### 2. Inheritance

```smog
Object subclass: #Animal [
    | name |
    
    name: aName [
        name := aName.
    ]
    
    getName [
        ^name
    ]
    
    speak [
        'Some generic sound' println.
    ]
]

Animal subclass: #Dog [
    speak [
        'Woof!' println.
    ]
]

Animal subclass: #Cat [
    speak [
        'Meow!' println.
    ]
]

" Usage "
| dog cat |
dog := Dog new.
dog name: 'Buddy'.
dog getName println.  " Prints: Buddy "
dog speak.            " Prints: Woof! "

cat := Cat new.
cat name: 'Whiskers'.
cat getName println.  " Prints: Whiskers "
cat speak.            " Prints: Meow! "
```

### 3. Polymorphism

```smog
Object subclass: #Shape [
    area [
        self subclassResponsibility.
    ]
    
    perimeter [
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
    
    perimeter [
        ^2 * 3.14159 * radius
    ]
]

Shape subclass: #Rectangle [
    | width height |
    
    width: w height: h [
        width := w.
        height := h.
    ]
    
    area [
        ^width * height
    ]
    
    perimeter [
        ^2 * (width + height)
    ]
]

" Usage - polymorphic behavior "
| shapes |
shapes := #().

| circle |
circle := Circle new.
circle radius: 5.
shapes := shapes copyWith: circle.

| rect |
rect := Rectangle new.
rect width: 4 height: 6.
shapes := shapes copyWith: rect.

" Calculate total area "
| totalArea |
totalArea := 0.
shapes do: [ :shape |
    totalArea := totalArea + shape area.
].
totalArea println.
```

### 4. Composition

```smog
Object subclass: #Engine [
    start [
        'Engine started' println.
    ]
    
    stop [
        'Engine stopped' println.
    ]
]

Object subclass: #Car [
    | engine model |
    
    initialize [
        engine := Engine new.
        model := 'Generic Car'.
    ]
    
    setModel: m [
        model := m.
    ]
    
    start [
        model print.
        ' is starting...' println.
        engine start.
    ]
    
    stop [
        model print.
        ' is stopping...' println.
        engine stop.
    ]
]

" Usage "
| car |
car := Car new.
car initialize.
car setModel: 'Tesla Model S'.
car start.
car stop.
```

## Common Algorithms

### Factorial (Recursive)

```smog
Object subclass: #Math [
    factorial: n [
        n <= 1 ifTrue: [ ^1 ].
        ^n * (self factorial: (n - 1))
    ]
]

" Usage "
| math |
math := Math new.
(math factorial: 5) println.  " Prints: 120 "
```

### Fibonacci Sequence

```smog
Object subclass: #Fibonacci [
    " Recursive (simple but slow) "
    fibRecursive: n [
        n <= 1 ifTrue: [ ^n ].
        ^(self fibRecursive: (n - 1)) + (self fibRecursive: (n - 2))
    ]
    
    " Iterative (faster) "
    fib: n [
        | a b temp i |
        n <= 1 ifTrue: [ ^n ].
        
        a := 0.
        b := 1.
        
        2 to: n do: [ :i |
            temp := a + b.
            a := b.
            b := temp.
        ].
        
        ^b
    ]
]

" Usage "
| fib |
fib := Fibonacci new.
(fib fib: 10) println.  " Prints: 55 "
```

### Greatest Common Divisor (GCD)

```smog
Object subclass: #Math [
    gcd: a and: b [
        b = 0 ifTrue: [ ^a ].
        ^self gcd: b and: (a \\ b)
    ]
]

" Usage "
| math |
math := Math new.
(math gcd: 48 and: 18) println.  " Prints: 6 "
```

### Prime Number Checker

```smog
Object subclass: #PrimeChecker [
    isPrime: n [
        | i |
        n <= 1 ifTrue: [ ^false ].
        n = 2 ifTrue: [ ^true ].
        (n \\ 2) = 0 ifTrue: [ ^false ].
        
        i := 3.
        [ i * i <= n ] whileTrue: [
            (n \\ i) = 0 ifTrue: [ ^false ].
            i := i + 2.
        ].
        
        ^true
    ]
    
    primesUpTo: limit [
        | primes i |
        primes := #().
        
        2 to: limit do: [ :i |
            (self isPrime: i) ifTrue: [
                primes := primes copyWith: i.
            ]
        ].
        
        ^primes
    ]
]

" Usage "
| checker |
checker := PrimeChecker new.
(checker isPrime: 17) println.  " Prints: true "

| primes |
primes := checker primesUpTo: 20.
primes do: [ :each | each println ].  " Prints: 2 3 5 7 11 13 17 19 "
```

## Best Practices

### 1. Use Meaningful Names

```smog
" Good "
| customerName totalPrice isValid |

" Bad "
| cn tp iv |
```

### 2. Keep Methods Small

```smog
" Good - focused, single responsibility "
Object subclass: #Calculator [
    add: a to: b [
        ^a + b
    ]
    
    subtract: a from: b [
        ^b - a
    ]
]

" Bad - doing too much "
Object subclass: #Calculator [
    compute: a and: b operation: op [
        " Complex logic handling many operations "
    ]
]
```

### 3. Use Blocks for Abstraction

```smog
" Reusable higher-order function "
Object subclass: #Collection [
    map: transformBlock [
        | result |
        result := #().
        self do: [ :each |
            result := result copyWith: (transformBlock value: each).
        ].
        ^result
    ]
]
```

### 4. Initialize Objects Properly

```smog
Object subclass: #Person [
    | name age |
    
    initialize [
        name := 'Unknown'.
        age := 0.
    ]
    
    name: n age: a [
        name := n.
        age := a.
    ]
]

" Always initialize "
| person |
person := Person new.
person initialize.  " Important! "
```

### 5. Handle Edge Cases

```smog
Object subclass: #Divider [
    divide: a by: b [
        b = 0 ifTrue: [
            'Error: Division by zero' println.
            ^nil
        ].
        ^a / b
    ]
]
```

## Tips for Learning Smog

1. **Think in Messages**: Everything is sending messages to objects
2. **Experiment with Blocks**: They're powerful for abstraction
3. **Start Simple**: Begin with basic programs and build up
4. **Read Examples**: Study the examples in the `examples/` directory
5. **Practice OOP**: Smog is pure object-oriented - embrace it!

## Standard Library

Smog includes a growing standard library with common data structures and utilities. The library is organized into modules:

### Collections
- **Set** - Unordered collection of unique elements
- **OrderedCollection** - Growable, ordered list
- **Bag** - Multiset that tracks element occurrences

### Core Utilities
- **Math** - Mathematical functions (sqrt, factorial, fibonacci, gcd, etc.)
- **Stream** - Sequential data access (ReadStream, WriteStream)

### I/O, Crypto, and Compression
- **HTTP** - HTTP client (interface ready, requires VM primitives)
- **AES** - AES-256 encryption (interface ready, requires VM primitives)
- **ZIP/GZIP** - Compression (interface ready, requires VM primitives)

For detailed documentation and examples, see:
- [Standard Library README](../stdlib/README.md)
- [Standard Library Index](../stdlib/INDEX.md)
- [Standard Library Examples](../examples/stdlib/)

## Next Steps

- Explore the [Standard Library](../stdlib/README.md) for common utilities
- Read the [Language Specification](../docs/spec/LANGUAGE_SPEC.md) for complete syntax reference
- Read the [Learning Guide](LEARNING_GUIDE.md) for conceptual understanding
- Check out [Example Programs](../examples/) for more inspiration
- Study the [Architecture](design/ARCHITECTURE.md) to understand how Smog works internally

Happy coding in Smog!
