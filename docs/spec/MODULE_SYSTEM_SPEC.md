# Smog Module System Specification

Version 0.1.0 (Draft)

## Overview

This specification defines a module system for Smog that enables code organization, namespace management, and multi-file programs. The design is inspired by Python's simple import system and Go's package model, while maintaining Smalltalk's object-oriented philosophy.

## Design Goals

1. **Simplicity**: Easy to understand and use, like Python's import system
2. **Explicit**: Clear declaration of dependencies and exports
3. **Namespace Safety**: Prevent name collisions in larger programs
4. **Smalltalk Philosophy**: Maintain consistency with Smalltalk's message-passing model
5. **File Organization**: Support multi-file programs naturally
6. **Compatibility**: Work seamlessly with existing bytecode (.sg) system

## Core Concepts

### Module

A **module** is a single .smog file that contains class definitions and optional executable code. Each module has:
- A unique name (typically matching the filename)
- A set of classes it defines
- Optional imports of other modules
- An optional module-level initialization block

### Package

A **package** is a collection of related modules organized in a directory. Packages provide:
- Logical grouping of related modules
- Namespace hierarchy
- Reusable library structure

### Namespace

A **namespace** maps module/package names to their exported classes and objects, preventing naming conflicts.

## Module Declaration Syntax

### Module Header (Optional but Recommended)

Every Smog file can optionally begin with a module declaration that specifies its name and namespace:

```smog
"! module: Collections.ArrayList
   This module provides an ArrayList implementation.
!"

Object subclass: #ArrayList [
    | elements size |
    " ... implementation ... "
]
```

The module declaration is a special comment starting with `"!` and ending with `!"`. It contains:
- `module:` - The fully qualified module name (package.module format)
- Optional description and metadata

**Alternative compact syntax:**

```smog
"! module: Math.Geometry !"

Object subclass: #Point [
    " ... "
]
```

### Module Name Rules

- Module names use dot notation for hierarchy: `Package.Subpackage.ModuleName`
- The last component should match the filename (without .smog extension)
- Examples:
  - `Collections.ArrayList` → `ArrayList.smog` in `Collections/` directory
  - `Math.Geometry` → `Geometry.smog` in `Math/` directory
  - `Core.Object` → `Object.smog` in `Core/` directory

### Files Without Module Declaration

Files without a module declaration are considered to be in the **default namespace** and their classes are globally accessible. This maintains backward compatibility with existing code.

## Import Syntax

### Basic Import

Import all exported classes from a module:

```smog
"! import: Collections.ArrayList !"

| list |
list := ArrayList new.
list add: 1.
list add: 2.
```

### Multiple Imports

Import multiple modules:

```smog
"! import: Collections.ArrayList !"
"! import: Collections.HashMap !"
"! import: Math.Geometry !"

| list map point |
list := ArrayList new.
map := HashMap new.
point := Point x: 10 y: 20.
```

### Package Import

Import all modules from a package:

```smog
"! import: Collections.* !"

" Now ArrayList, HashMap, LinkedList, etc. are all available "
| list map |
list := ArrayList new.
map := HashMap new.
```

### Aliased Import

Import with an alias to avoid naming conflicts:

```smog
"! import: Graphics.Point as: GPoint !"
"! import: Math.Point as: MPoint !"

| graphicsPoint mathPoint |
graphicsPoint := GPoint x: 100 y: 200.
mathPoint := MPoint x: 1.5 y: 2.7.
```

## Module Exports

### Automatic Exports

By default, all classes defined in a module are exported and available to importers. This follows Smalltalk's open philosophy.

### Selective Export (Future Enhancement)

For more control, modules can explicitly declare exports:

```smog
"! module: Collections.Internal
   export: ArrayList, HashMap
!"

Object subclass: #ArrayList [
    " ... public API ... "
]

Object subclass: #InternalHelper [
    " ... not exported, private to module ... "
]
```

## Module Resolution

### Search Path

When resolving imports, the Smog runtime searches in this order:

1. **Standard Library Path**: Built-in modules (e.g., `Core.Object`, `Collections.*`)
2. **Project Path**: Modules in the current project directory
3. **SMOG_PATH**: Additional paths from the SMOG_PATH environment variable

### File Resolution

For an import like `"! import: Collections.ArrayList !"`, the runtime searches for:

1. `Collections/ArrayList.smog` or `Collections/ArrayList.sg`
2. `Collections.ArrayList.smog` or `Collections.ArrayList.sg`
3. In each directory of the search path

### Bytecode Priority

If both `.smog` and `.sg` files exist, the `.sg` file is preferred for faster loading (unless it's older than the `.smog` file).

## Module Initialization

### Initialization Block

Modules can have initialization code that runs once when first imported:

```smog
"! module: Config.Settings !"

" Module-level variables (class variables of a pseudo-class) "
| defaultTimeout maxRetries |

"! init !"
defaultTimeout := 30.
maxRetries := 3.
'Settings module initialized' println.

Object subclass: #Settings [
    timeout [
        ^defaultTimeout
    ]
    
    retries [
        ^maxRetries
    ]
]
```

The `"! init !"` marker indicates code that runs during module initialization.

## Namespace Access

### Fully Qualified Names

Classes can always be referenced by their fully qualified name:

```smog
| point |
point := Math.Geometry.Point x: 5 y: 10.
```

### After Import

Once imported, classes can be used by their short name:

```smog
"! import: Math.Geometry !"

| point |
point := Point x: 5 y: 10.  " Short name "
```

## Standard Library Organization

The Smog standard library will be organized into packages:

```
Core/
  Object.smog          - Root object class
  Class.smog           - Metaclass
  Boolean.smog         - Boolean classes
  Nil.smog            - Nil class

Collections/
  Array.smog          - Array class
  ArrayList.smog      - Dynamic array
  HashMap.smog        - Hash map
  LinkedList.smog     - Linked list
  Set.smog            - Set collection

Math/
  Integer.smog        - Integer class
  Double.smog         - Float class
  Geometry.smog       - Point, Rectangle, etc.

IO/
  File.smog           - File operations
  Stream.smog         - Stream abstraction
  Console.smog        - Console I/O

Blocks/
  Block.smog          - Block/closure class
  Continuation.smog   - First-class continuations
```

## Compatibility with Current Code

### Backward Compatibility

To maintain compatibility with existing Smog code:

1. Files without module declarations work as before (global namespace)
2. All current examples and code continue to work unchanged
3. Module system is purely additive - no breaking changes

### Migration Path

Existing code can be gradually migrated:

1. **Phase 1**: Continue using code without module declarations
2. **Phase 2**: Add module declarations to new files
3. **Phase 3**: Organize code into packages
4. **Phase 4**: Use imports to manage dependencies

## Examples

### Example 1: Simple Module

**File: Math/Factorial.smog**
```smog
"! module: Math.Factorial !"

Object subclass: #Factorial [
    compute: n [
        n <= 1 ifTrue: [ ^1 ].
        ^n * (self compute: (n - 1))
    ]
]
```

**File: main.smog**
```smog
"! import: Math.Factorial !"

| calc result |
calc := Factorial new.
result := calc compute: 5.
result println.  " Prints: 120 "
```

### Example 2: Package with Multiple Modules

**File: Collections/ArrayList.smog**
```smog
"! module: Collections.ArrayList !"

Object subclass: #ArrayList [
    | elements capacity |
    
    initialize [
        capacity := 10.
        elements := Array new: capacity.
    ]
    
    add: element [
        " ... implementation ... "
    ]
]
```

**File: Collections/HashMap.smog**
```smog
"! module: Collections.HashMap !"

Object subclass: #HashMap [
    | buckets |
    
    initialize [
        buckets := Array new: 16.
    ]
    
    at: key put: value [
        " ... implementation ... "
    ]
]
```

**File: app.smog**
```smog
"! import: Collections.ArrayList !"
"! import: Collections.HashMap !"

| list map |
list := ArrayList new initialize.
map := HashMap new initialize.

list add: 'item1'.
map at: 'key' put: 'value'.
```

### Example 3: Avoiding Name Conflicts

**File: Graphics/Point.smog**
```smog
"! module: Graphics.Point !"

Object subclass: #Point [
    | x y color |
    
    x: xVal y: yVal color: c [
        x := xVal.
        y := yVal.
        color := c.
    ]
]
```

**File: Math/Point.smog**
```smog
"! module: Math.Point !"

Object subclass: #Point [
    | x y |
    
    x: xVal y: yVal [
        x := xVal.
        y := yVal.
    ]
    
    distance [
        ^((x * x) + (y * y)) sqrt
    ]
]
```

**File: app.smog**
```smog
"! import: Graphics.Point as: GPoint !"
"! import: Math.Point as: MPoint !"

| screenPoint mathPoint |
screenPoint := GPoint x: 100 y: 200 color: 'red'.
mathPoint := MPoint x: 3 y: 4.
mathPoint distance println.  " Prints: 5 "
```

## Implementation Considerations

### Compiler Changes

1. **Lexer**: Recognize module declaration comments (`"! ... !"`)
2. **Parser**: Parse module declarations and import statements
3. **Module Registry**: Track loaded modules and their exports
4. **Name Resolution**: Resolve class names considering imports and namespaces

### Runtime Changes

1. **Module Loader**: Load and initialize modules on demand
2. **Namespace Management**: Maintain namespace hierarchy
3. **Circular Dependency Detection**: Detect and handle circular imports
4. **Caching**: Cache loaded modules to avoid reloading

### Bytecode Format

Extend the .sg format to include:
- Module metadata (name, version, dependencies)
- Import list
- Export list
- Initialization code

### Tools

1. **Module Browser**: Tool to explore available modules and packages
2. **Dependency Analyzer**: Show module dependencies
3. **Package Creator**: Scaffold for new packages
4. **Documentation Generator**: Generate docs from module metadata

## Comparison with Other Languages

### Python-like Simplicity

```python
# Python
import collections
from math import sqrt

# Smog equivalent
"! import: Collections.* !"
"! import: Math.sqrt !"
```

### Go-like Package Structure

```go
// Go package structure
package collections

import "fmt"

// Smog equivalent structure
"! module: Collections.ArrayList !"
"! import: IO.Console !"
```

### Smalltalk Philosophy

Unlike typical module systems, Smog maintains Smalltalk's runtime flexibility:
- Modules can be loaded at runtime
- Classes can be modified after loading
- Message sends work across module boundaries
- No compile-time linking required

## Future Enhancements

### Version Management

```smog
"! module: Collections.ArrayList
   version: 1.2.0
   requires: Core.Object >= 1.0.0
!"
```

### Private Exports

```smog
"! module: Collections.Internal
   export: public ArrayList, HashMap
   private: InternalHelper
!"
```

### Conditional Imports

```smog
"! import: Platform.Unix if: Platform unix !"
"! import: Platform.Windows if: Platform windows !"
```

### Module Metadata

```smog
"! module: Collections.ArrayList
   author: Kristofer
   license: MIT
   description: Dynamic array implementation
   tags: collection, array, list
!"
```

## Open Questions for Discussion

1. **Initialization Order**: How to handle initialization dependencies between modules?
2. **Circular Dependencies**: Allow with lazy resolution, or prohibit?
3. **Module Reloading**: Support runtime module reloading for development?
4. **Visibility Modifiers**: Need for private/protected classes within modules?
5. **Module Variables**: Support module-level variables (like class variables)?

## Summary

This module system proposal provides:

- **Simple syntax** inspired by Python and Go
- **Namespace management** to prevent conflicts
- **Backward compatibility** with existing code
- **Natural organization** for multi-file programs
- **Foundation for package ecosystem**
- **Consistent with Smalltalk** message-passing philosophy

The design balances simplicity with power, allowing developers to organize code naturally while maintaining Smog's core object-oriented principles.
