# Version 0.2.0 Example Programs

This directory contains example programs that demonstrate the features implemented in version 0.2.0 of the Smog language.

## Running Examples

Build the interpreter first:
```bash
go build -o bin/smog ./cmd/smog
```

Then run any example:
```bash
./bin/smog examples/v0.2.0/literals.smog
./bin/smog examples/v0.2.0/variables.smog
./bin/smog examples/v0.2.0/arithmetic.smog
./bin/smog examples/v0.2.0/comparison.smog
./bin/smog examples/v0.2.0/print.smog
./bin/smog examples/v0.2.0/complex.smog
```

## Feature Examples

### literals.smog
Tests basic integer literals. Returns 42.

### variables.smog
Demonstrates variable declarations and assignments. Declares two variables, assigns values, and adds them together. Returns 30.

### arithmetic.smog
Shows binary message sends with arithmetic operators. Tests addition. Returns 7.

### comparison.smog
Demonstrates comparison operations. Compares two numbers and returns a boolean. Returns true.

### print.smog
Shows unary message sends. Prints "Hello from v0.2.0!" to the console.

### complex.smog
Combines multiple features: variable declarations, assignments, arithmetic, and variable references. Returns 42.

## What's Supported in v0.2.0

- ✅ Literals: integers, floats, strings, booleans, nil
- ✅ Variable declarations: `| x y z |`
- ✅ Assignments: `x := 42`
- ✅ Arithmetic: `+`, `-`, `*`, `/`
- ✅ Comparisons: `<`, `>`, `<=`, `>=`, `=`, `~=`
- ✅ Unary messages: `object method`
- ✅ Binary messages: `3 + 4`
- ✅ Comments: `" This is a comment "`

## What's NOT Supported Yet

- ❌ Blocks/closures: `[ :x | x + 1 ]`
- ❌ Classes: `Object subclass: #Point`
- ❌ Methods: Method definitions in classes
- ❌ Control flow: `ifTrue:`, `ifFalse:`, `whileTrue:`
- ❌ Arrays: `#(1 2 3)`
- ❌ Keyword messages: `point x: 10 y: 20` (parser supports it but no classes/objects yet)

These features will be implemented in version 0.3.0 and later.

## Testing

For comprehensive testing instructions, see: `docs/MANUAL_TESTING_GUIDE.md`

For automated tests, run:
```bash
go test ./test -v -run TestVersion0_2_0
```
