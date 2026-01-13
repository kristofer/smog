# Syntax-Only Examples

This directory contains example programs that demonstrate valid Smog syntax but cannot currently be executed because they depend on features that are not yet fully implemented.

## Why These Examples Don't Run

These examples use **class definitions** and **object instantiation**, which are parsed but not yet compiled or executed by the Smog VM. The parser can recognize the syntax, but the compiler and VM do not support classes yet.

## Examples in This Directory

### counter.smog
Demonstrates:
- Class definition with instance variables
- Multiple methods (initialize, increment, decrement, value, reset)
- Object state management
- Object instantiation with `new`

**Status**: Parses correctly, but compilation fails with "unknown statement type: *ast.Class"

### factorial.smog
Demonstrates:
- Class definition
- Recursive method calls
- Conditional execution with `ifTrue:`
- Return statements (`^`)
- Method parameters

**Status**: Parses correctly, but compilation fails with "unknown statement type: *ast.Class"

### point.smog
Demonstrates:
- Constructor methods with multiple parameters
- Accessor methods
- Binary operator methods (`+`, `-`)
- Custom string representation
- Object composition

**Status**: Parses correctly, but compilation fails with "unknown statement type: *ast.Class"

### cascade_example.smog (v0.4.0)
Demonstrates:
- Cascading message sends with `;`
- Fluent interface pattern
- Multiple operations on same receiver

**Status**: Runs but fails at runtime because required objects don't exist

### self_example.smog (v0.4.0)
Demonstrates:
- Using the `self` keyword
- Method chaining
- Internal method calls

**Status**: Runs but fails at runtime because required objects don't exist

### super_example.smog (v0.4.0)
Demonstrates:
- Class inheritance with `subclass:`
- Super method calls
- Parent class initialization
- Method overriding

**Status**: Parses correctly, but compilation fails with "unknown statement type: *ast.Class"

## When Will These Examples Work?

These examples will become executable once class definition compilation and object instantiation are implemented in the Smog compiler and VM. This is planned for a future version.

See the [roadmap](../../docs/planning/ROADMAP.md) for implementation timeline.

## Runnable Examples

For examples that can be executed now, see:
- [examples/](../) - Basic runnable examples (hello.smog, arrays.smog, blocks.smog)
- [examples/v0.2.0/](../v0.2.0/) - Version 0.2.0 feature examples
- [examples/v0.3.0/](../v0.3.0/) - Version 0.3.0 feature examples
- [examples/v0.4.0/](../v0.4.0/) - Version 0.4.0 feature examples (dictionary_example.smog)

## Contributing

If you're interested in implementing class support in Smog, please see the [contribution guidelines](../../README.md#contributing) and check out the architecture documentation.
