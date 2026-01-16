# Standard Library Examples - Quick Start Guide

This directory contains working examples demonstrating the Smog standard library.

## Running Examples

```bash
# Build the smog interpreter first
go build -o bin/smog ./cmd/smog

# Run any example
./bin/smog examples/stdlib/set_example.smog
./bin/smog examples/stdlib/math_example.smog
./bin/smog examples/stdlib/ordered_collection_example.smog
./bin/smog examples/stdlib/comprehensive_example.smog
```

## Available Examples

### set_example.smog
Demonstrates the **Set** collection class for managing unique elements.

**Features shown:**
- Creating sets
- Adding elements (duplicates ignored)
- Testing membership with `includes:`
- Set operations: `union:`, `intersection:`
- Iterating with `do:`

**Output highlights:**
- Shows how duplicates are automatically filtered
- Demonstrates set union (combining sets)
- Demonstrates set intersection (finding common elements)

### math_example.smog
Demonstrates the **Math** utility class for mathematical operations.

**Features shown:**
- Mathematical constants (pi, e)
- Basic operations (abs, max, min)
- Powers and square roots
- Factorial calculations
- Fibonacci sequence generation
- Greatest common divisor (GCD)

**Output highlights:**
- Computes factorials up to 10!
- Generates first 15 Fibonacci numbers
- Calculates square roots using Newton's method

### ordered_collection_example.smog
Demonstrates the **OrderedCollection** class for flexible list operations.

**Features shown:**
- Creating and populating collections
- Accessing elements (first, last, at:)
- Transforming with `collect:`
- Filtering with `select:`
- Finding with `detect:`
- Testing with `anySatisfy:` and `allSatisfy:`

**Output highlights:**
- Doubles each number using collect
- Filters even and odd numbers
- Detects first number greater than 5
- Chains operations for complex queries

### comprehensive_example.smog
Demonstrates multiple stdlib classes working together to analyze numbers.

**Features shown:**
- Combines OrderedCollection, Set, and Math classes
- Creates a NumberAnalyzer class that uses stdlib
- Statistical calculations (sum, mean, max, min)
- Filtering (evens, odds, positives)
- Transformations (squares)
- Finding unique values

**Output highlights:**
- Shows how stdlib classes compose naturally
- Demonstrates real-world use case
- Illustrates object-oriented design with stdlib

## Learning Path

1. **Start with set_example.smog** - Learn about unique collections
2. **Try math_example.smog** - Explore mathematical utilities
3. **Explore ordered_collection_example.smog** - Master list operations
4. **Study comprehensive_example.smog** - See how it all fits together

## Tips

- All examples include the class definitions inline, so they're self-contained
- Once the module system is available (v0.6.0), you'll be able to import these classes
- The examples demonstrate clean, readable Smog code following best practices
- Each example builds on concepts from previous ones

## Next Steps

After running these examples:
- Read the [Standard Library INDEX](../../stdlib/INDEX.md) for complete API reference
- Check out [Standard Library README](../../stdlib/README.md) for design philosophy
- Explore the stdlib source files in the `stdlib/` directory
- Try modifying the examples to experiment with the APIs

## Contributing Examples

When adding new stdlib examples:
1. Include the class definition inline (until module system exists)
2. Add clear comments explaining what's being demonstrated
3. Show realistic use cases, not just trivial demos
4. Include error cases where appropriate
5. Add output that helps verify correctness
6. Consider adding a test in `test/stdlib_test.go`

Happy coding with Smog's standard library!
