# Smog Design Decisions

This document records key design decisions made during the development of smog, along with the rationale behind them.

## Language Design

### Decision: Stack-Based Bytecode VM

**Decision**: Use a stack-based bytecode virtual machine rather than register-based or direct AST interpretation.

**Rationale**:
- Simpler bytecode generation - no register allocation needed
- Smaller bytecode size due to implicit operands
- Well-understood implementation from JVM, Python, and other languages
- Easier to implement initially, can optimize later
- Stack-based VMs are easier to verify and reason about

**Tradeoffs**:
- May require more bytecode instructions than register-based
- More stack manipulation operations
- Potential performance impact (mitigated by future optimizations)

**Alternatives Considered**:
- Register-based VM (Lua-style): More complex but potentially faster
- AST interpretation: Simpler to start but slower execution

### Decision: Smalltalk-Style Syntax

**Decision**: Adopt Smalltalk's syntax including message passing, blocks, and class definitions.

**Rationale**:
- Proven syntax with decades of use
- Extremely simple and consistent grammar
- Everything-is-an-object philosophy maps naturally to syntax
- Excellent for teaching OOP concepts
- Strong alignment with SOM project

**Tradeoffs**:
- May be unfamiliar to developers from C-family languages
- Limited adoption compared to mainstream syntax

**Alternatives Considered**:
- Ruby-style syntax: More familiar but less pure
- Custom syntax: Unnecessary when Smalltalk syntax works well

### Decision: Pure Object-Oriented Model

**Decision**: Everything is an object, including classes, booleans, numbers, and blocks.

**Rationale**:
- Conceptual simplicity and consistency
- Follows Smalltalk and SOM philosophy
- Enables powerful metaprogramming
- No special cases in the language

**Tradeoffs**:
- Potential performance overhead for primitives
- More complex implementation of basic operations

**Alternatives Considered**:
- Primitive types with automatic boxing: Less pure, more complexity
- Hybrid model: Violates simplicity principle

## Implementation Design

### Decision: Go Programming Language

**Decision**: Implement smog in Go rather than C, C++, Rust, or self-hosted.

**Rationale**:
- Good balance of performance and development speed
- Excellent concurrency primitives for future features
- Strong standard library
- Memory safety without manual management
- Good cross-platform support
- Fast compilation times

**Tradeoffs**:
- Garbage collection may add latency
- Not as performant as C/C++/Rust
- Larger binary sizes

**Alternatives Considered**:
- C/C++: Maximum performance but slower development
- Rust: Great safety but steeper learning curve
- Self-hosted (smog in smog): Bootstrapping challenge, future goal

### Decision: Separate Compilation and Execution

**Decision**: Generate intermediate bytecode rather than interpreting AST directly.

**Rationale**:
- Enables bytecode caching and distribution
- Separates concerns (compilation vs execution)
- Opportunity for bytecode-level optimizations
- Bytecode can be verified independently
- Easier to add JIT compilation later

**Tradeoffs**:
- More complex implementation
- Additional compilation phase
- Bytecode format to maintain

**Alternatives Considered**:
- AST interpretation: Simpler but slower and less flexible
- Direct compilation to machine code: Too complex for initial version

### Decision: Simple Garbage Collection Initially

**Decision**: Rely on Go's garbage collector for object memory management initially.

**Rationale**:
- Simplifies initial implementation significantly
- Proven GC implementation
- Good enough performance for early versions
- Can implement custom GC later if needed

**Tradeoffs**:
- Less control over memory management
- Go GC may not be optimal for smog's allocation patterns
- Potential for GC pauses in performance-critical code

**Alternatives Considered**:
- Custom mark-and-sweep GC: More work, can add later
- Reference counting: Simpler but doesn't handle cycles
- Manual memory management: Against project goals

### Decision: Minimal Standard Library Initially

**Decision**: Start with a minimal standard library covering only essential classes.

**Rationale**:
- Faster time to working interpreter
- Allows core language features to stabilize first
- Standard library can grow based on real needs
- Keeps initial scope manageable

**Tradeoffs**:
- Limited initial functionality
- May make early examples less impressive

**Alternatives Considered**:
- Comprehensive standard library: Too much work upfront
- No standard library: Too minimal to be useful

## Project Structure

### Decision: Standard Go Project Layout

**Decision**: Organize code following Go community conventions (cmd/, pkg/, internal/).

**Rationale**:
- Familiar to Go developers
- Clear separation of public APIs and internal implementation
- Supports multiple executables if needed
- Well-supported by Go tools

**Tradeoffs**:
- May seem complex for a simple project initially
- Requires understanding of Go packaging

**Alternatives Considered**:
- Flat structure: Too simple, hard to scale
- Custom organization: Unfamiliar to Go developers

### Decision: Documentation-First Approach

**Decision**: Write comprehensive documentation before implementation.

**Rationale**:
- Forces clear thinking about design
- Provides roadmap for implementation
- Helps onboard contributors
- Enables parallel work on docs and code

**Tradeoffs**:
- Documentation may become outdated if not maintained
- Upfront time investment before coding

**Alternatives Considered**:
- Code-first approach: Faster start but less clarity
- Minimal documentation: Harder to maintain and contribute

## Testing Strategy

### Decision: Unit Tests for Components, Integration Tests for Features

**Decision**: Write unit tests for individual packages and integration tests for complete features.

**Rationale**:
- Unit tests catch component-level bugs
- Integration tests verify component interaction
- Easier to maintain than only integration tests
- Fast unit tests for quick feedback

**Tradeoffs**:
- More tests to write and maintain
- Some duplication between test levels

**Alternatives Considered**:
- Only integration tests: Slower, harder to debug
- Only unit tests: Misses integration issues

## Future Decisions to Make

The following decisions are deferred to future versions:

1. **Module System Design**: How to organize and import code across files
2. **Concurrency Model**: Actor model, CSP, or other approach
3. **Type System**: Optional static typing, gradual typing, or dynamic only
4. **FFI Design**: How to interface with Go or C code
5. **JIT Compilation**: When and how to add just-in-time compilation
6. **Package Manager**: Format, hosting, versioning strategy

## Decision Review Process

Design decisions should be:
1. Documented in this file when made
2. Revisited as we gain experience with implementation
3. Updated if we discover better alternatives
4. Discussed with contributors before major changes

## References

- [SOM Design Decisions](http://som-st.github.io/)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Smalltalk-80 Blue Book](https://rmod-files.lille.inria.fr/FreeBooks/BlueBook/Bluebook.pdf)
