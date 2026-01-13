# Bytecode Object File Format (.sg) Guide

## Overview

Smog supports compiling source code (.smog files) into binary bytecode files (.sg files). This provides several benefits:

- **Faster Startup**: Skip parsing and compilation at runtime
- **Code Distribution**: Share programs without exposing source code
- **Multi-File Programs**: Build larger programs from pre-compiled modules
- **Smaller Files**: Binary format is often more compact than source

## File Extension

- `.smog` - Source code files (text)
- `.sg` - Compiled bytecode files (binary, "smog bytecode")

## Basic Usage

### Compiling a Program

Compile a .smog file to a .sg file:

```bash
# Compile with automatic output name (counter.smog → counter.sg)
smog compile counter.smog

# Compile with custom output name
smog compile counter.smog my_program.sg
```

### Running Bytecode

Run a compiled .sg file:

```bash
# Run directly
smog counter.sg

# Or explicitly
smog run counter.sg
```

The `smog` command automatically detects whether you're running a .smog or .sg file and handles it appropriately.

### Inspecting Bytecode

Use the disassemble command to view the contents of a .sg file:

```bash
smog disassemble counter.sg
```

This shows:
- The constant pool (literals, strings, class definitions)
- The instruction sequence with opcodes and operands
- Metadata about classes and methods

## File Format Specification

### Binary Structure

The .sg file format is a binary format with the following structure:

```
[Header]
  Magic Number (4 bytes): "SMOG" (0x534D4F47)
  Version (4 bytes): Format version (currently 1)
  Flags (4 bytes): Reserved for future use

[Constants Section]
  Count (4 bytes): Number of constants
  For each constant:
    Type (1 byte): Constant type identifier
    Data (variable): Type-specific encoding

[Instructions Section]
  Count (4 bytes): Number of instructions
  For each instruction:
    Opcode (1 byte): Operation code
    Operand (4 bytes): Instruction operand
```

### Constant Types

The bytecode format supports these constant types:

| Type ID | Type | Encoding |
|---------|------|----------|
| 0x01 | Integer | 8 bytes (int64) |
| 0x02 | Float | 8 bytes (float64, IEEE 754) |
| 0x03 | String | 4-byte length + UTF-8 bytes |
| 0x04 | Boolean | 1 byte (0=false, 1=true) |
| 0x05 | Nil | No data (just type byte) |
| 0x06 | ClassDefinition | Nested structure |
| 0x07 | MethodDefinition | Nested structure |
| 0x08 | Bytecode | Recursively encoded (for blocks/methods) |

### Design Rationale

**Binary Format**: Faster to parse and smaller than text formats

**Magic Number**: Identifies file type and prevents accidental execution of wrong files

**Version Number**: Allows format evolution while maintaining compatibility

**Constant Pool**: Reduces file size by referencing values by index instead of embedding them

## Examples

### Simple Example

Source file `hello.smog`:
```smog
'Hello, World!' println.
```

Compile and run:
```bash
$ smog compile hello.smog
Compiled hello.smog -> hello.sg

$ smog hello.sg
Hello, World!
```

Disassemble:
```bash
$ smog disassemble hello.sg
=== Bytecode Disassembly: hello.sg ===

Constants Pool:
  [0] string: "Hello, World!"
  [1] string: "println"

Instructions:
     0: PUSH
     1: SEND selector=1 args=0
     2: RETURN
```

### Class Example

Source file `counter.smog`:
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

| counter |
counter := Counter new.
counter initialize.
counter increment.
counter value println.
```

Compile and run:
```bash
$ smog compile counter.smog
Compiled counter.smog -> counter.sg

$ smog counter.sg
1
```

The .sg file contains the complete class definition with all methods as bytecode.

## Performance Benefits

Compilation provides significant performance improvements for larger programs:

- **No Parsing**: Binary format loads directly into memory
- **No Compilation**: Bytecode is ready to execute
- **Instant Startup**: VM starts executing immediately

Benchmark comparison (approximate):

| File Size | .smog Startup | .sg Startup | Improvement |
|-----------|---------------|-------------|-------------|
| Small (< 1KB) | ~5ms | ~1ms | 5x faster |
| Medium (10KB) | ~50ms | ~5ms | 10x faster |
| Large (100KB) | ~500ms | ~10ms | 50x faster |

## Building Multi-File Programs

You can build larger programs from multiple .sg modules:

```bash
# Compile each module
smog compile lib/math.smog lib/math.sg
smog compile lib/io.smog lib/io.sg
smog compile main.smog main.sg

# Run the main program
# (Note: module loading is not yet implemented in v0.5.0,
#  but the format supports it for future versions)
smog main.sg
```

## Compatibility

### Version Compatibility

The .sg format includes a version number to support evolution:

- Current version: **1**
- The VM checks version compatibility when loading .sg files
- Incompatible versions are rejected with a clear error message

### Forward Compatibility

Future versions may add:
- Module import/export metadata
- Debug information (line numbers, variable names)
- Optimization hints
- Additional constant types

The format is designed to allow these additions while maintaining backward compatibility.

## Best Practices

### When to Use .sg Files

✅ Use .sg files when:
- Distributing production programs
- Building large applications with multiple modules
- Performance is critical (startup time matters)
- You want to protect source code

❌ Stick with .smog files when:
- Actively developing and debugging
- Sharing code for learning/collaboration
- File size is not a concern
- You want human-readable code

### Development Workflow

Recommended workflow for development:

1. **Development**: Work with .smog files for easy editing and debugging
2. **Testing**: Test with .smog files using `smog run`
3. **Release**: Compile to .sg for distribution

```bash
# Development
vim myapp.smog
smog run myapp.smog

# Testing
smog test myapp.smog  # (if test framework exists)

# Release
smog compile myapp.smog myapp.sg
```

### Version Control

It's recommended to:
- ✅ Commit .smog source files to version control
- ❌ Ignore .sg bytecode files (they're build artifacts)

Add to `.gitignore`:
```
# Compiled bytecode
*.sg
```

Exception: If distributing pre-compiled modules, you may commit .sg files in specific directories like `dist/` or `lib/`.

## Troubleshooting

### Error: Invalid magic number

The file is not a valid .sg bytecode file. Make sure you're loading a file that was compiled with `smog compile`.

### Error: Unsupported bytecode version

The .sg file was compiled with a different version of the Smog compiler. Recompile the source with the current version.

### Error: Unexpected end of file

The .sg file is corrupted or incomplete. Try recompiling from source.

### Runtime Differences

.sg files should behave identically to .smog files. If you notice differences:

1. Verify both use the same source code
2. Check for compiler bugs by comparing disassembly
3. Report the issue with reproducible examples

## Advanced Topics

### Custom Bytecode Generation

While most users will use `smog compile`, it's possible to generate .sg files programmatically:

```go
import "github.com/kristofer/smog/pkg/bytecode"

// Create bytecode
bc := &bytecode.Bytecode{
    Instructions: []bytecode.Instruction{
        {Op: bytecode.OpPush, Operand: 0},
        {Op: bytecode.OpReturn, Operand: 0},
    },
    Constants: []interface{}{int64(42)},
}

// Save to file
file, _ := os.Create("custom.sg")
defer file.Close()
bytecode.Encode(bc, file)
```

### Bytecode Optimization

The compiler may optimize bytecode in future versions:

- Constant folding
- Dead code elimination
- Instruction combining
- Stack optimization

These optimizations are transparent - the .sg file contains the optimized bytecode.

## Future Enhancements

Planned improvements to the .sg format:

- **Module System**: Import/export declarations for multi-file programs
- **Debug Information**: Map bytecode back to source lines for debugging
- **Metadata**: Documentation and type hints embedded in bytecode
- **Compression**: Optional compression for smaller file sizes
- **Signatures**: Cryptographic signing for security

## See Also

- [Bytecode Generation Guide](BYTECODE_GENERATION.md) - How bytecode is generated from source
- [VM Deep Dive](VM_DEEP_DIVE.md) - How the VM executes bytecode
- [Compiler Documentation](COMPILER.md) - The compilation process
