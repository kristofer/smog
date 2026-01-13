# Smog v0.3.0 Examples

This directory contains example programs demonstrating the features added in v0.3.0.

## Features Demonstrated

### blocks_and_control_flow.smog

This comprehensive example showcases:

1. **Simple Blocks**: Basic block creation and execution with `value`
2. **Parameterized Blocks**: Blocks with one or more parameters using `value:`
3. **Conditional Execution**: 
   - `ifTrue:` - execute block if condition is true
   - `ifFalse:` - execute block if condition is false
4. **Loops**: `timesRepeat:` for repeating an action N times
5. **Array Operations**:
   - `do:` - iterate over array elements
   - `size` - get array length
   - `at:` - access element by index (1-based)
6. **Combining Features**: Using blocks with variables and arrays

## Running the Examples

From the project root:

```bash
go run ./cmd/smog examples/v0.3.0/blocks_and_control_flow.smog
```

## Expected Output

The example should demonstrate each feature with clear output showing:
- Block execution results
- Conditional logic working correctly
- Repeated actions via loops
- Array iteration and access
- Combined operations

## Notes

- Arrays use 1-based indexing (like Smalltalk)
- Blocks capture parameters but not outer scope yet (closures fully in v0.4.0)
- All control flow is implemented via message sends to objects
