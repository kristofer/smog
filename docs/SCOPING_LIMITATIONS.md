# Variable Scoping Limitations (Pre-Lexical Scoping)

## Current Implementation Issues

The current smog implementation has significant limitations with variable scoping that prevent proper lexical scoping behavior.

### Problem 1: All Variables Must Be Declared at Top

**Broken:**
```smog
| x |
x := 1.

[ x println ] value.

| y |  " ERROR: Cannot declare variables after blocks
y := 2.
```

**Workaround:**
```smog
| x y |  " All variables declared at top
x := 1.

[ x println ] value.

y := 2.
```

### Problem 2: Blocks Cannot Declare Local Variables

**Broken:**
```smog
| x |
x := 10.

[ :y |
  | temp |  " ERROR: Blocks cannot have local variables
  temp := y * 2.
  x + temp
] value: 5.
```

**Workaround:**
```smog
| x temp |  " Declare temp at outer scope
x := 10.

[ :y |
  temp := y * 2.  " Use outer scope variable
  x + temp
] value: 5.
```

### Problem 3: Variable Index Conflicts

When a block is created, it copies the parent's symbol table. New variables declared after the block conflict with the block's captured variable indices, causing "local variable index out of bounds" errors.

## Root Cause

The compiler uses a **flat variable namespace** where:

1. Parent scope variables: indices 0..N
2. Block parameters: indices N+1..M
3. New variables after block: indices 0..K (CONFLICT!)

This creates index collisions because there's no distinction between:
- Variables in my own scope
- Variables captured from parent scope

## Why This Matters

This breaks fundamental programming patterns:

```smog
" Pattern: Progressive variable introduction
| numbers |
numbers := #(1 2 3).

" Process the numbers
numbers do: [ :each | each println ].

" ERROR: Cannot declare more variables here
| sum |
sum := 0.
```

## Temporary Workaround

Until proper lexical scoping is implemented:

✅ **DO:**
- Declare all variables in a single block at the top of each scope
- Avoid local variables inside blocks (use parameters only)

❌ **DON'T:**
- Declare variables after using blocks
- Use multiple variable declaration blocks
- Declare local variables inside blocks

## Next Steps

This limitation will be removed by implementing proper **lexical scoping** with:
- Environment chains for captured variables
- Separate opcodes for local vs captured variable access
- Support for block-local temporaries
- Proper closure semantics

See `docs/LEXICAL_SCOPING.md` for the implementation plan.
