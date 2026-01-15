# Smog Debugger Guide

The Smog debugger provides interactive debugging capabilities for Smog programs, allowing you to step through code execution, inspect variables, and set breakpoints.

## Starting the Debugger

To run a Smog program with the debugger enabled:

```bash
smog debug myprogram.smog
```

The debugger will start in step mode, pausing before the first instruction.

## Debugger Commands

When the debugger pauses, you'll see a `debug>` prompt. The following commands are available:

### Execution Control

- **`help` or `h` or `?`** - Show available commands
- **`continue` or `c`** - Continue execution until next breakpoint or completion
- **`step` or `s`** - Enable step mode (pause after each instruction)
- **`next` or `n`** - Execute the next instruction and pause
- **`quit` or `q`** - Quit debugging (aborts program execution)

### Inspection Commands

- **`stack` or `st`** - Display the current VM stack (values being computed)
- **`locals` or `l`** - Show local variables and their values
- **`globals` or `g`** - Show global variables and their values
- **`callstack` or `cs`** - Show the call stack (function/method call chain)
- **`instruction` or `i`** - Show the current instruction being executed

### Breakpoints

- **`breakpoint <n>` or `b <n>`** - Add a breakpoint at instruction number n
- **`delete <n>` or `d <n>`** - Remove the breakpoint at instruction number n
- **`list` or `ls`** - List all instructions with breakpoint markers

## Example Debugging Session

Here's an example of using the debugger:

```smog
| x y result |
x := 10.
y := 5.
result := x + y.
result println.
```

Running with the debugger:

```bash
$ smog debug example.smog

=== Smog Debugger ===
Type 'help' at the debug prompt for available commands
Starting in step mode...

=== Debugger Paused ===
     0: PUSH 0

debug> help
Debugger Commands:
  help, h, ?           Show this help
  continue, c          Continue execution
  step, s              Enable step mode (pause after each instruction)
  next, n              Execute next instruction
  stack, st            Show VM stack
  locals, l            Show local variables
  globals, g           Show global variables
  callstack, cs        Show call stack
  instruction, i       Show current instruction
  breakpoint <n>, b    Add breakpoint at instruction n
  delete <n>, d        Remove breakpoint at instruction n
  list, ls             List all instructions
  quit, q              Quit debugging (abort execution)

debug> list
->    0: PUSH 0
      1: STORE_LOCAL 0
      2: PUSH 1
      3: STORE_LOCAL 1
      4: LOAD_LOCAL 0
      5: LOAD_LOCAL 1
      6: SEND selector=2 args=1 (+)
      7: STORE_LOCAL 2
      8: LOAD_LOCAL 2
      9: SEND selector=3 args=0 (println)
     10: POP
     11: RETURN 

debug> b 6
Breakpoint added at instruction 6

debug> c
=== Debugger Paused ===
  *   6: SEND selector=2 args=1 (+)

debug> stack
Stack (top to bottom):
  [1] 5 (int64)
  [0] 10 (int64)

debug> n
=== Debugger Paused ===
     7: STORE_LOCAL 2

debug> stack
Stack (top to bottom):
  [0] 15 (int64)

debug> locals
Local variables:
  [0] 10 (int64)
  [1] 5 (int64)
  [2] 15 (int64)

debug> c
15

Program completed successfully
```

## Understanding the Instruction Listing

The `list` command shows all bytecode instructions. The format is:

```
[marker] [position]: [opcode] [operands]
```

Where:
- `marker` is:
  - `->` for the current instruction
  - `*` for a breakpoint
  - Empty for normal instructions
- `position` is the instruction number (for setting breakpoints)
- `opcode` is the bytecode operation (PUSH, SEND, etc.)
- `operands` are instruction arguments

## Stack Inspection

The VM uses a stack-based execution model. The `stack` command shows:
- Values on the stack (top to bottom)
- The index of each value
- The type of each value

This is useful for understanding how expressions are evaluated.

## Variable Inspection

Variables come in two types:

1. **Local variables** (`locals`) - Declared with `| var |` syntax, function parameters
2. **Global variables** (`globals`) - Variables assigned without declaration

The debugger shows the slot number, current value, and type for each variable.

## Call Stack

The `callstack` command shows the chain of function/method calls:
- Method or function name
- Message selector (if applicable)
- Instruction pointer at time of call

This is particularly useful for understanding execution flow in programs with nested calls.

## Tips for Effective Debugging

1. **Use step mode initially** to understand program flow
2. **Set breakpoints** at interesting locations, then use `continue`
3. **Check the stack** before and after operations to understand computations
4. **Use locals and globals** to verify variable states
5. **List instructions** to get oriented in the code

## Limitations

- The debugger works with bytecode, not source code (no source-level debugging yet)
- Instruction numbers are used instead of line numbers
- Some optimizations may make instruction flow non-obvious

## Advanced Usage

### Debugging Class Methods

When debugging programs with classes, use `callstack` to understand which method you're in:

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
```

Set breakpoints in specific methods by first running with `step` mode to find instruction positions.

### Debugging Blocks and Closures

Blocks create nested execution contexts. Use the call stack to understand which block is executing:

```smog
| block |
block := [ :x | x * 2 ].
block value: 5.
```

The debugger will show when you're inside the block's bytecode.

## Related Documentation

- [VM Deep Dive](VM_DEEP_DIVE.md) - Understanding bytecode execution
- [Bytecode Format](BYTECODE_FORMAT.md) - Bytecode instruction reference
- [REPL Guide](REPL.md) - Interactive programming without debugging

## Future Enhancements

Planned improvements for the debugger:

- Source-level debugging (line numbers instead of instruction positions)
- Conditional breakpoints
- Watchpoints (break when variable changes)
- Reverse execution
- Save/load debugging sessions
