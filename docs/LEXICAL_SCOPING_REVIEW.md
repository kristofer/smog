# Lexical Scoping Implementation Review

## Executive Summary

The lexical scoping implementation in Smog is **partially complete**. The foundational work has been done (bytecode opcodes defined, compiler struct updated), but the full implementation is not yet finished. The current state allows the system to compile and run, but true lexical scoping with proper variable capture is not yet functional.

## Current Status Assessment

### ✅ Completed Components

1. **Bytecode Opcodes (pkg/bytecode/bytecode.go)**
   - ✅ `OpLoadCaptured` - defined and documented
   - ✅ `OpStoreCaptured` - defined and documented  
   - ✅ `OpMakeClosureWithEnv` - defined and documented
   - ✅ `CapturedVar` structure - properly defined
   - ✅ `Bytecode` structure extended with `CapturedVars` and `LocalCount` fields
   - ✅ String representation for new opcodes

2. **Compiler Structure Updates (pkg/compiler/compiler.go)**
   - ✅ `localVars []string` field added
   - ✅ `capturedVars []bytecode.CapturedVar` field added
   - ✅ `parent *Compiler` field added
   - ✅ `New()` function updated to initialize new fields
   - ✅ Variable declaration updated to use `localVars` slice
   - ✅ `findLocalVar()` helper function added
   - ✅ Variable resolution updated to use `findLocalVar()`

3. **Documentation**
   - ✅ `docs/LEXICAL_SCOPING.md` - comprehensive design document
   - ✅ `docs/SCOPING_LIMITATIONS.md` - current limitations documented
   - ✅ `docs/LEXICAL_SCOPING_IMPLEMENTATION_STATUS.md` - status tracking

4. **Parser Error Reporting (Recently Completed)**
   - ✅ Enhanced error messages with line/column information
   - ✅ Source context showing in error output
   - ✅ Visual pointer (^) showing exact error location
   - ✅ Helpful suggestions for common errors

### ❌ Not Implemented

1. **Compiler - Variable Resolution**
   - ❌ No `resolveVariable()` method for lexical scope chain traversal
   - ❌ No distinction between local and captured variable access
   - ❌ Block compilation still uses flat variable copying (not environment chains)
   - ❌ No support for multi-level variable capture
   - ❌ Compiler doesn't set `CapturedVars` or `LocalCount` in returned bytecode

2. **VM - New Opcode Handlers**
   - ❌ `OpLoadCaptured` handler not implemented
   - ❌ `OpStoreCaptured` handler not implemented
   - ❌ `OpMakeClosureWithEnv` handler not implemented
   - ❌ No `Closure` struct with environment
   - ❌ Block execution doesn't use captured environments

3. **Parser - Block-Local Variables**
   - ❌ Cannot declare variables inside blocks
   - Current: `[ :param | body ]`
   - Needed: `[ :param | | local1 local2 | body ]`

4. **Testing**
   - ❌ No lexical scoping specific tests
   - ❌ No tests for variable capture
   - ❌ No tests for multi-level capture
   - ❌ No tests for block-local variables

## Critical Issues

### 1. Compiler Block Handling is Still Flat

**Current Implementation (lines 662-690 in compiler.go):**
```go
// Copy parent's local variables to support closures
blockCompiler.localVars = append([]string{}, c.localVars...)
blockCompiler.localCount = c.localCount
```

**Problem:** This creates a flat copy of variables, not true lexical scoping with environment chains.

**Impact:** 
- Cannot distinguish between local and captured variables
- Multiple variable declarations after blocks cause index conflicts
- True closures don't work (variables aren't properly captured)

### 2. VM Doesn't Implement New Opcodes

The VM has no handlers for:
- `OpLoadCaptured`
- `OpStoreCaptured`
- `OpMakeClosureWithEnv`

**Impact:** If the compiler emitted these opcodes, the VM would crash.

### 3. No Environment Chain

**Current Block Structure (in VM):**
```go
type Block struct {
    Code             *bytecode.Bytecode
    ParentLocalCount int
    ParameterCount   int
}
```

**Missing:**
- No captured environment storage
- No link to parent closure
- No way to resolve captured variables at runtime

## Comparison with Requirements

### From Issue Comments

1. **Multiple Variable Declaration Blocks**
   - Requirement: "compiler doesn't handle variable scoping correctly when new variables are declared after blocks"
   - Status: ❌ NOT FIXED - Still has this limitation
   - Note in docs: Language should require single variable declaration block

2. **Review Lexical Scoping Implementation**
   - Status: ✅ COMPLETED (this document)

3. **Reference SOM Implementation**
   - Status: ⚠️ PARTIAL - Design document references SOM concepts
   - Action Needed: Could benefit from deeper analysis of SOM's actual code

### From Design Document (docs/LEXICAL_SCOPING.md)

The design document outlines a comprehensive implementation that would:
- ✅ Support multiple variable declaration blocks
- ✅ Allow block-local temporaries  
- ✅ Properly capture variables from parent scopes
- ✅ Support multi-level capture (nested blocks)
- ✅ Implement proper closure semantics

**Reality:** None of the runtime behavior is implemented yet.

## SOM Reference Analysis

The issue mentions referencing SOM's implementation at:
https://github.com/smarr/SOM/tree/master/SomSom/src

### Key SOM Patterns We Should Study

1. **Lexical Context Chain**
   - SOM maintains a chain of lexical contexts
   - Each context knows its parent
   - Variable lookup traverses the chain

2. **Variable Resolution**
   - First check local scope
   - Then check each parent scope in order
   - Cache the resolution for performance

3. **Closure Creation**
   - Capture only variables actually used (not all parent variables)
   - Store captured values in the closure object
   - Reference by index in captured array

4. **Block Compilation**
   - Blocks compile in the context of their parent
   - Track which parent variables are accessed
   - Generate different opcodes for local vs. captured access

## Recommendations

### Short-term (Current Issue)

1. **Document the Limitation ✅ (Already Done)**
   - Update language spec to require single variable declaration block
   - Add to user guide with clear examples
   - Ensure all examples follow this pattern

2. **Add Parser Validation**
   - Add a check to enforce single variable declaration block
   - Give helpful error message when violated
   - This prevents confusing runtime errors

3. **Update Examples**
   - Audit all examples to ensure single declaration block
   - Add comments explaining the limitation

### Long-term (Full Lexical Scoping)

To complete the lexical scoping implementation:

1. **Compiler Changes (2-3 days)**
   - Implement `resolveVariable()` method with environment chain traversal
   - Update `compileBlockLiteral()` to build environment capture list
   - Emit `OpMakeClosureWithEnv` with captured variable indices
   - Set `CapturedVars` and `LocalCount` in returned bytecode

2. **VM Changes (1-2 days)**
   - Create `Closure` struct with environment array
   - Implement `OpLoadCaptured` handler
   - Implement `OpStoreCaptured` handler
   - Implement `OpMakeClosureWithEnv` handler
   - Update block invocation to use closure environments

3. **Parser Changes (1 day)**
   - Add support for block-local variable declarations
   - Syntax: `[ :param | | local1 local2 | body ]`

4. **Testing (2-3 days)**
   - Unit tests for each opcode
   - Integration tests for variable capture scenarios
   - Edge case testing (shadowing, multi-level capture)
   - Performance benchmarks

5. **Documentation Updates (1 day)**
   - Update all affected documentation
   - Create examples demonstrating new capabilities
   - Update migration guide

**Total Estimated Effort:** 7-10 days for a full implementation

## Interim Solution

For now, the language works correctly with these constraints:

1. **Single Variable Declaration Block**
   ```smog
   | x y z |  " Declare all variables at the top
   x := 1.
   [ y := 2 ] value.  " This works
   " | w | <- This would cause errors (don't do it)
   ```

2. **No Block-Local Variables**
   ```smog
   | x temp |  " Declare temp at outer scope
   x := 10.
   [ :y |
     temp := y * 2.  " Use outer scope variable
     x + temp
   ] value: 5.
   ```

3. **Blocks Can Still Capture Variables**
   - The current flat copying does allow basic closure behavior
   - Just can't declare new variables after creating blocks

## Conclusion

The lexical scoping implementation has solid **design and infrastructure** but lacks **runtime implementation**. The current system is in a working state with documented limitations. 

**Key Finding:** The limitation on multiple variable declaration blocks is a direct consequence of the incomplete lexical scoping implementation, not a language design choice.

**Recommendation:** 
1. For this issue: Document the limitation clearly and enforce it with parser validation
2. For future: Complete the lexical scoping implementation to remove the limitation

The current approach of documenting the limitation and adjusting examples is the right pragmatic solution while the full implementation is pending.
