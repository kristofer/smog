# Smog REPL Guide

The Smog REPL (Read-Eval-Print Loop) provides an interactive environment for experimenting with the Smog language, testing code snippets, and learning the language.

## Starting the REPL

There are two ways to start the REPL:

```bash
# Method 1: Run smog with no arguments
./bin/smog

# Method 2: Explicitly specify repl
./bin/smog repl
```

## Features

### Interactive Evaluation

The REPL evaluates expressions and statements as you type them:

```
smog> 3 + 4.
smog> 'Hello, World!' println.
```

### Persistent State

Variables and values persist across statements within a REPL session:

```
smog> | x y |
smog> x := 42.
smog> y := 10.
smog> x + y.
```

### Multi-line Input

The REPL supports multi-line input. Continue typing on new lines, and the REPL will execute when it sees a complete statement (ending with a period):

```
smog> | counter |
....> counter := 0.
smog> counter.
```

### Error Recovery

Errors don't crash the REPL - you can continue working after fixing mistakes:

```
smog> xyz.
Runtime error: undefined identifier: xyz
smog> | xyz |
smog> xyz := 5.
smog> xyz.
```

## Special Commands

The REPL supports several special commands that start with a colon:

- `:help` - Display help information
- `:quit` - Exit the REPL
- `:exit` - Exit the REPL (alternative to :quit)

Example:
```
smog> :help
smog REPL Help
...
smog> :quit
Goodbye!
```

## Example Session

Here's a complete example session showing the REPL in action:

```
$ ./bin/smog
smog REPL v0.4.0
Type ':help' for help, ':quit' or ':exit' to exit

smog> | x |
smog> x := 10.
smog> x + 5.
smog> | arr |
smog> arr := #(1 2 3).
smog> | dict |
smog> dict := #{'name' -> 'Alice'. 'age' -> 30}.
smog> :quit
Goodbye!
```

## Tips

1. **Variable Declaration**: Always declare variables before using them with `| varName |`
2. **Statement Termination**: End statements with a period (`.`) to execute them
3. **Exploring**: Use the REPL to test language features and experiment with syntax
4. **Learning**: Try examples from the documentation interactively

## Limitations

- No way to view current variable values (yet)
- No command history or editing (depends on terminal)
- No way to clear variables or reset state (restart the REPL)
- Limited error messages compared to file execution

## Use Cases

### Learning Smog

Perfect for beginners learning the language syntax:

```
smog> true ifTrue: [ 'yes' println ].
smog> false ifFalse: [ 'no' println ].
```

### Quick Testing

Test small code snippets without creating files:

```
smog> | point |
smog> point x: 10; y: 20; z: 30.
```

### Interactive Development

Develop and test small pieces of functionality:

```
smog> | factorial |
smog> factorial := [ :n | n <= 1 ifTrue: [ 1 ] ifFalse: [ n * factorial value: n - 1 ] ].
smog> factorial value: 5.
```

## Troubleshooting

### REPL won't start
- Make sure you've built the binary: `go build -o bin/smog ./cmd/smog`
- Check that the binary is executable

### Can't exit the REPL
- Use `:quit` or `:exit` command
- Use Ctrl+C to force quit
- Use Ctrl+D (EOF) on Unix systems

### Errors persist
- The VM state persists across statements
- Variables can't be redeclared
- Restart the REPL to get a clean state

## Future Enhancements

Planned improvements to the REPL:

- Command history and line editing
- Ability to inspect variable values
- Pretty-printing of results
- Multi-line editing support
- Tab completion for identifiers
- Save/load REPL sessions
- Reset command to clear state
