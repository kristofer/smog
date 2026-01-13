# Smog v0.4.0 Example Programs

This directory contains example programs demonstrating the new features in Smog v0.4.0.

## New Features

### 1. Super Message Sends
See: `super_example.smog`

The `super` keyword allows calling parent class methods from overridden methods.

```smog
super initialize.
super accelerate.
```

### 2. Self Keyword
See: `self_example.smog`

The `self` keyword refers to the current object (receiver).

```smog
self displayTotal.
```

### 3. Cascading Messages
See: `cascade_example.smog`

Cascading allows sending multiple messages to the same receiver using semicolons.

```smog
point x: 10; y: 20; z: 30; display.
```

The cascade returns the receiver itself, not the result of the last message.

### 4. Dictionary Literals
See: `dictionary_example.smog`

Dictionary literals provide a concise syntax for creating key-value maps.

```smog
person := #{'name' -> 'Alice'. 'age' -> 30. 'city' -> 'Wonderland'}.
```

## Running Examples

**Note**: Class parsing is not yet fully implemented in the parser. These examples demonstrate the syntax and can be parsed, but full execution with class instantiation is pending.

To parse an example:
```bash
./bin/smog examples/v0.4.0/dictionary_example.smog
```

## Feature Combinations

All v0.4.0 features can be combined:

```smog
" Self with cascading "
self x: 10; y: 20; display.

" Super in a cascade "
super initialize; display.

" Dictionary in cascading "
obj data: #{'x' -> 10. 'y' -> 20}; process.
```

## Testing

Integration tests for these features are in `/test/version_0_4_0_test.go`.

Benchmarks for performance testing are in `/test/benchmark_0_4_0_test.go`.

Run tests:
```bash
go test ./test -v
```

Run benchmarks:
```bash
go test ./test -bench=.
```
