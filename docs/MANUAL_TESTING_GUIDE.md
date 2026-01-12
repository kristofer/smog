# Manual Testing Guide for v0.2.0 Features in VS Code

This guide shows you how to manually test the v0.2.0 features of the Smog language interpreter using VS Code.

## Prerequisites

1. **Open the smog repository in VS Code**
2. **Build the interpreter**:
   ```bash
   go build -o bin/smog ./cmd/smog
   ```

## Testing Overview

Version 0.2.0 implements the following features:
- ✅ Literal expressions (integers, floats, strings, booleans, nil)
- ✅ Variable declarations and assignments
- ✅ Arithmetic operations (+, -, *, /)
- ✅ Comparison operations (<, >, <=, >=, =, ~=)
- ✅ Unary message sends (e.g., `'text' println`)
- ✅ Binary message sends (e.g., `3 + 4`)
- ✅ Keyword message sends (e.g., `point x: 10 y: 20`)

## Quick Start - Using the Integrated Terminal

### 1. Open Terminal in VS Code
- Press `` Ctrl+` `` (backtick) or go to **Terminal → New Terminal**

### 2. Build the Interpreter
```bash
go build -o bin/smog ./cmd/smog
```

### 3. Run Example Programs

Navigate to the `examples/v0.2.0/` directory to find test programs for each feature.

## Feature Testing

### Test 1: Literals

**File:** `examples/v0.2.0/literals.smog`

```smog
" Test literals - version 0.2.0 "
42
```

**Run:**
```bash
./bin/smog examples/v0.2.0/literals.smog
```

**Expected:** Program executes without error (returns integer 42 on the stack)

**What it tests:** Integer literal parsing, compilation, and execution

---

### Test 2: Variables and Assignments

**File:** `examples/v0.2.0/variables.smog`

```smog
" Test variables and assignments - version 0.2.0 "
| x y |
x := 10.
y := 20.
x + y
```

**Run:**
```bash
./bin/smog examples/v0.2.0/variables.smog
```

**Expected:** Program executes without error (returns 30 on the stack)

**What it tests:**
- Variable declarations (`| x y |`)
- Variable assignments (`x := 10`)
- Variable references (`x + y`)

---

### Test 3: Arithmetic Operations

**File:** `examples/v0.2.0/arithmetic.smog`

```smog
" Test arithmetic operations - version 0.2.0 "
3 + 4
```

**Run:**
```bash
./bin/smog examples/v0.2.0/arithmetic.smog
```

**Expected:** Program executes without error (returns 7 on the stack)

**What it tests:** Binary message sends with arithmetic operators

**Try modifying:**
- `3 - 4` (subtraction)
- `3 * 4` (multiplication)
- `12 / 3` (division)

---

### Test 4: Comparison Operations

**File:** `examples/v0.2.0/comparison.smog`

```smog
" Test comparison operations - version 0.2.0 "
| x y result |
x := 10.
y := 20.
result := x < y.
result
```

**Run:**
```bash
./bin/smog examples/v0.2.0/comparison.smog
```

**Expected:** Program executes without error (returns true on the stack)

**What it tests:** Comparison operators returning boolean values

**Try modifying:**
- `x > y` (should return false)
- `x <= y` (should return true)
- `x = y` (should return false)
- `x ~= y` (should return true)

---

### Test 5: Print Message (Unary Message Send)

**File:** `examples/v0.2.0/print.smog`

```smog
" Test print message - version 0.2.0 "
'Hello from v0.2.0!' println
```

**Run:**
```bash
./bin/smog examples/v0.2.0/print.smog
```

**Expected Output:**
```
Hello from v0.2.0!
```

**What it tests:** Unary message sends (messages with no arguments)

**Try modifying:**
- `'Hello from v0.2.0!' print` (prints without newline)
- `42 println` (prints a number)

---

### Test 6: Complex Example

**File:** `examples/v0.2.0/complex.smog`

```smog
" Complex example - version 0.2.0 "
| a b c result |
a := 15.
b := 27.
c := a + b.
result := c.
result
```

**Run:**
```bash
./bin/smog examples/v0.2.0/complex.smog
```

**Expected:** Program executes without error (returns 42 on the stack)

**What it tests:** Combination of multiple features together

---

## Testing with Custom Programs

### Using the Integrated Terminal

1. **Create a new file** in VS Code (e.g., `test.smog`)
2. **Write your program:**
   ```smog
   | x |
   x := 100.
   x
   ```
3. **Save the file** (Ctrl+S)
4. **Run in terminal:**
   ```bash
   ./bin/smog test.smog
   ```

### Using VS Code Tasks (Optional)

You can create a task to run the current file:

1. Press **Ctrl+Shift+P** and type "Tasks: Configure Task"
2. Select "Create tasks.json from template" → "Others"
3. Replace the content with:

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Run Smog File",
            "type": "shell",
            "command": "./bin/smog",
            "args": [
                "${file}"
            ],
            "problemMatcher": [],
            "group": {
                "kind": "build",
                "isDefault": true
            }
        }
    ]
}
```

4. Now you can run the current file with **Ctrl+Shift+B**

---

## Debugging Failed Tests

If a program fails, you'll see error messages like:

- **Parse error:** Syntax issue in your smog code
- **Compile error:** Problem converting AST to bytecode
- **Runtime error:** Issue during execution

### Example Error Messages

**Parse Error:**
```
Parse error: parser errors: [unexpected token: STAR]
```
→ Check your syntax, you may have an unsupported construct

**Runtime Error:**
```
Runtime error: undefined global variable: x
```
→ Make sure you declared the variable in `| x |`

---

## Running Automated Tests

To verify that v0.2.0 is working correctly, run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run only v0.2.0 integration tests
go test ./test -v -run TestVersion0_2_0
```

---

## Feature Coverage Summary

| Feature | Example File | Status |
|---------|-------------|--------|
| Integer literals | literals.smog | ✅ |
| Float literals | literals.smog | ✅ |
| String literals | print.smog | ✅ |
| Boolean literals | comparison.smog | ✅ |
| Variable declarations | variables.smog | ✅ |
| Variable assignments | variables.smog | ✅ |
| Addition (+) | arithmetic.smog | ✅ |
| Subtraction (-) | arithmetic.smog | ✅ |
| Multiplication (*) | arithmetic.smog | ✅ |
| Division (/) | arithmetic.smog | ✅ |
| Less than (<) | comparison.smog | ✅ |
| Greater than (>) | comparison.smog | ✅ |
| Equal (=) | comparison.smog | ✅ |
| Not equal (~=) | comparison.smog | ✅ |
| Unary messages | print.smog | ✅ |
| Binary messages | arithmetic.smog | ✅ |
| Complex programs | complex.smog | ✅ |

---

## Tips for Manual Testing

1. **Use the integrated terminal** - Keep it open at the bottom of VS Code for quick testing
2. **Create a scratch file** - Keep a `scratch.smog` file for quick experiments
3. **Check exit codes** - A successful program returns exit code 0
4. **Read error messages** - They indicate the exact line and type of error
5. **Start simple** - Test individual features before combining them
6. **Use comments** - Comments help document what you're testing: `" This tests X "`

---

## Next Steps

After confirming v0.2.0 works:
- Explore the existing examples in `examples/` directory
- Read `docs/IMPLEMENTATION_SUMMARY.md` for implementation details
- Check `docs/spec/LANGUAGE_SPEC.md` for language syntax reference
- Review the automated tests in `test/` for more examples

---

## Troubleshooting

### "smog: command not found"
**Solution:** Build the interpreter first:
```bash
go build -o bin/smog ./cmd/smog
```

### "no such file or directory"
**Solution:** Make sure you're in the repository root directory:
```bash
cd /path/to/smog
```

### "Parse error" when running examples
**Solution:** Make sure you're using files from `examples/v0.2.0/` - older examples may use v0.3.0+ features not yet implemented.

---

## Version History

- **v0.1.0**: Basic lexer and parser
- **v0.2.0**: Variables, assignments, arithmetic, comparisons, message sends ← *You are here*
- **v0.3.0**: Blocks, classes, methods (coming soon)
