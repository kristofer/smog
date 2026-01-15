# Smog Module System Implementation Plan

Version 0.1.0 (Draft)

## Overview

This document outlines the implementation plan for adding a module system to Smog, as specified in [MODULE_SYSTEM_SPEC.md](../spec/MODULE_SYSTEM_SPEC.md). The implementation will be done in phases to minimize risk and ensure each phase delivers working functionality.

## Implementation Phases

### Phase 1: Foundation (v0.6.0-alpha)
**Goal**: Basic module declaration and loading

#### 1.1 Lexer Extensions
- [ ] Add recognition of module declaration comments (`"! ... !"`)
- [ ] Add token types for module metadata
- [ ] Preserve backward compatibility with existing comment handling

**Files to Modify**:
- `pkg/lexer/lexer.go`
- `pkg/lexer/token.go`

**Estimated Effort**: 2-3 days

#### 1.2 Parser Extensions
- [ ] Parse module declaration header
- [ ] Parse import statements
- [ ] Create AST nodes for module metadata
- [ ] Extract module name and description

**Files to Modify**:
- `pkg/parser/parser.go`
- `pkg/ast/ast.go` (new node types)

**New AST Node Types**:
```go
type ModuleDeclaration struct {
    Name        string
    Description string
    Position    Position
}

type ImportStatement struct {
    ModulePath string
    Alias      string
    Position   Position
}

type Module struct {
    Declaration *ModuleDeclaration
    Imports     []*ImportStatement
    Classes     []*ClassDefinition
    InitBlock   *BlockExpression
}
```

**Estimated Effort**: 3-4 days

#### 1.3 Module Registry
- [ ] Create module registry to track loaded modules
- [ ] Implement module metadata storage
- [ ] Add module name validation
- [ ] Implement module lookup by name

**New Files**:
- `pkg/module/registry.go`
- `pkg/module/module.go`
- `pkg/module/metadata.go`

**Estimated Effort**: 2-3 days

#### 1.4 Basic Module Loading
- [ ] Implement file-based module loader
- [ ] Add module search path resolution
- [ ] Support both .smog and .sg files
- [ ] Detect module declaration in files

**New Files**:
- `pkg/module/loader.go`
- `pkg/module/resolver.go`

**Estimated Effort**: 4-5 days

#### 1.5 Testing
- [ ] Unit tests for module declaration parsing
- [ ] Unit tests for module registry
- [ ] Unit tests for module loading
- [ ] Integration test: load simple module

**Estimated Effort**: 2-3 days

**Total Phase 1**: ~15-20 days

---

### Phase 2: Import Resolution (v0.6.0-beta)
**Goal**: Import statements work and resolve class names

#### 2.1 Namespace Management
- [ ] Implement namespace hierarchy
- [ ] Add namespace resolution logic
- [ ] Support fully qualified names (Package.Module.Class)
- [ ] Handle default namespace for backward compatibility

**New Files**:
- `pkg/module/namespace.go`

**Estimated Effort**: 3-4 days

#### 2.2 Import Processing
- [ ] Resolve import statements to module files
- [ ] Load imported modules automatically
- [ ] Add imported classes to local namespace
- [ ] Support wildcard imports (Package.*)

**Files to Modify**:
- `pkg/module/loader.go`
- `pkg/module/resolver.go`

**Estimated Effort**: 4-5 days

#### 2.3 Name Resolution in Compiler
- [ ] Modify compiler to check imports during name resolution
- [ ] Support short names (after import) and qualified names
- [ ] Add error reporting for unresolved names
- [ ] Maintain backward compatibility with global names

**Files to Modify**:
- `pkg/compiler/compiler.go`
- `pkg/compiler/resolver.go` (new file for name resolution)

**Estimated Effort**: 5-6 days

#### 2.4 Circular Dependency Detection
- [ ] Track module loading stack
- [ ] Detect circular imports
- [ ] Provide helpful error messages
- [ ] Consider lazy resolution strategy

**Files to Modify**:
- `pkg/module/loader.go`

**Estimated Effort**: 2-3 days

#### 2.5 Testing
- [ ] Test basic imports
- [ ] Test qualified names
- [ ] Test circular dependency detection
- [ ] Test namespace isolation
- [ ] Integration tests with multiple modules

**Estimated Effort**: 3-4 days

**Total Phase 2**: ~17-22 days

---

### Phase 3: Module Initialization (v0.6.0-rc1)
**Goal**: Module initialization blocks work correctly

#### 3.1 Initialization Block Support
- [ ] Parse initialization blocks (`"! init !"`)
- [ ] Execute initialization on first module load
- [ ] Ensure initialization runs only once
- [ ] Support module-level variables

**Files to Modify**:
- `pkg/parser/parser.go`
- `pkg/ast/ast.go`
- `pkg/module/loader.go`

**Estimated Effort**: 3-4 days

#### 3.2 Initialization Ordering
- [ ] Topologically sort modules by dependencies
- [ ] Initialize modules in dependency order
- [ ] Handle initialization errors gracefully
- [ ] Report initialization cycle errors

**Files to Modify**:
- `pkg/module/loader.go`
- `pkg/module/initializer.go` (new)

**Estimated Effort**: 4-5 days

#### 3.3 Module-Level State
- [ ] Support module-level variables (similar to class variables)
- [ ] Implement module variable access from classes
- [ ] Add module variable bytecode instructions if needed

**Files to Modify**:
- `pkg/compiler/compiler.go`
- `pkg/vm/vm.go` (if new opcodes needed)
- `pkg/bytecode/opcodes.go` (if new opcodes needed)

**Estimated Effort**: 3-4 days

#### 3.4 Testing
- [ ] Test initialization block execution
- [ ] Test initialization ordering
- [ ] Test module variables
- [ ] Integration tests with complex initialization

**Estimated Effort**: 2-3 days

**Total Phase 3**: ~12-16 days

---

### Phase 4: Bytecode Format Extension (v0.6.0-rc2)
**Goal**: Compiled modules (.sg files) include module metadata

#### 4.1 Bytecode Format Updates
- [ ] Extend .sg format to include module metadata
- [ ] Add module name, version, dependencies to header
- [ ] Include import list in bytecode
- [ ] Include initialization bytecode

**Files to Modify**:
- `pkg/bytecode/serialization.go`
- `pkg/bytecode/format.go`

**Estimated Effort**: 3-4 days

#### 4.2 Compilation Updates
- [ ] Compile module metadata into bytecode
- [ ] Include imports in compiled output
- [ ] Compile initialization blocks
- [ ] Update bytecode version number

**Files to Modify**:
- `pkg/compiler/compiler.go`
- `cmd/smog/compile.go`

**Estimated Effort**: 2-3 days

#### 4.3 Loading Updates
- [ ] Load module metadata from .sg files
- [ ] Resolve imports when loading bytecode
- [ ] Execute initialization from bytecode
- [ ] Handle version compatibility

**Files to Modify**:
- `pkg/bytecode/deserialization.go`
- `pkg/module/loader.go`

**Estimated Effort**: 3-4 days

#### 4.4 Testing
- [ ] Test module compilation to .sg
- [ ] Test loading compiled modules
- [ ] Test initialization from bytecode
- [ ] Integration tests: compile and run modular programs

**Estimated Effort**: 2-3 days

**Total Phase 4**: ~10-14 days

---

### Phase 5: Standard Library Reorganization (v0.6.1)
**Goal**: Organize standard library into modules

#### 5.1 Core Module
- [ ] Create Core.Object module
- [ ] Create Core.Class module
- [ ] Create Core.Boolean module
- [ ] Create Core.Nil module

**New Files**:
- `stdlib/Core/Object.smog`
- `stdlib/Core/Class.smog`
- `stdlib/Core/Boolean.smog`
- `stdlib/Core/Nil.smog`

**Estimated Effort**: 2-3 days

#### 5.2 Collections Module
- [ ] Create Collections.Array module
- [ ] Create Collections.ArrayList module (new)
- [ ] Create Collections.HashMap module (new)
- [ ] Create Collections.LinkedList module (new)

**New Files**:
- `stdlib/Collections/Array.smog`
- `stdlib/Collections/ArrayList.smog`
- `stdlib/Collections/HashMap.smog`
- `stdlib/Collections/LinkedList.smog`

**Estimated Effort**: 5-6 days

#### 5.3 Math Module
- [ ] Create Math.Integer module
- [ ] Create Math.Double module
- [ ] Create Math.Geometry module (Point, Rectangle, etc.)

**New Files**:
- `stdlib/Math/Integer.smog`
- `stdlib/Math/Double.smog`
- `stdlib/Math/Geometry.smog`

**Estimated Effort**: 3-4 days

#### 5.4 Blocks Module
- [ ] Create Blocks.Block module
- [ ] Create Blocks.Continuation module (placeholder)

**New Files**:
- `stdlib/Blocks/Block.smog`
- `stdlib/Blocks/Continuation.smog`

**Estimated Effort**: 1-2 days

#### 5.5 Backward Compatibility Layer
- [ ] Create compatibility shim for old code
- [ ] Auto-import Core modules by default
- [ ] Provide migration guide

**New Files**:
- `stdlib/Compat/prelude.smog`
- `docs/MODULE_MIGRATION_GUIDE.md`

**Estimated Effort**: 2-3 days

#### 5.6 Testing
- [ ] Test each stdlib module independently
- [ ] Test backward compatibility
- [ ] Update all examples to use new stdlib organization
- [ ] Performance testing

**Estimated Effort**: 4-5 days

**Total Phase 5**: ~17-23 days

---

### Phase 6: Advanced Features (v0.6.2)
**Goal**: Aliased imports and advanced features

#### 6.1 Import Aliases
- [ ] Parse import aliases (`as:` syntax)
- [ ] Update namespace to support aliases
- [ ] Test aliased imports
- [ ] Document alias usage

**Files to Modify**:
- `pkg/parser/parser.go`
- `pkg/module/namespace.go`

**Estimated Effort**: 2-3 days

#### 6.2 Package Wildcards
- [ ] Implement wildcard import resolution (`Package.*`)
- [ ] Load all modules in package directory
- [ ] Handle package initialization order
- [ ] Test wildcard imports

**Files to Modify**:
- `pkg/module/resolver.go`
- `pkg/module/loader.go`

**Estimated Effort**: 3-4 days

#### 6.3 Module Metadata
- [ ] Support author, version, license in module declaration
- [ ] Store metadata in registry
- [ ] Expose metadata via reflection
- [ ] Tool to query module metadata

**Files to Modify**:
- `pkg/parser/parser.go`
- `pkg/module/metadata.go`
- `cmd/smog/module_info.go` (new command)

**Estimated Effort**: 2-3 days

#### 6.4 SMOG_PATH Environment Variable
- [ ] Implement SMOG_PATH parsing
- [ ] Add to module search path
- [ ] Document usage
- [ ] Test with custom paths

**Files to Modify**:
- `pkg/module/resolver.go`

**Estimated Effort**: 1-2 days

#### 6.5 Testing
- [ ] Test all advanced features
- [ ] Integration tests
- [ ] Documentation with examples

**Estimated Effort**: 2-3 days

**Total Phase 6**: ~10-15 days

---

### Phase 7: Tooling (v0.6.3)
**Goal**: Developer tools for working with modules

#### 7.1 Module Browser
- [ ] CLI tool to list available modules
- [ ] Show module metadata
- [ ] Display module dependencies
- [ ] Search functionality

**New Files**:
- `cmd/smog/modules.go`

**Estimated Effort**: 3-4 days

#### 7.2 Dependency Analyzer
- [ ] Analyze and display dependency tree
- [ ] Detect unused imports
- [ ] Suggest import optimizations
- [ ] Export dependency graph

**New Files**:
- `cmd/smog/deps.go`
- `pkg/module/analyzer.go`

**Estimated Effort**: 3-4 days

#### 7.3 Package Creator
- [ ] Scaffold new package structure
- [ ] Generate boilerplate module files
- [ ] Create package documentation template
- [ ] Interactive package creation wizard

**New Files**:
- `cmd/smog/new_package.go`

**Estimated Effort**: 2-3 days

#### 7.4 Documentation Generator
- [ ] Extract module metadata for docs
- [ ] Generate API documentation from modules
- [ ] Create package index
- [ ] HTML/Markdown output

**New Files**:
- `cmd/smog/doc.go`
- `pkg/doc/generator.go`

**Estimated Effort**: 4-5 days

#### 7.5 Testing
- [ ] Test all tools
- [ ] User acceptance testing
- [ ] Documentation

**Estimated Effort**: 2-3 days

**Total Phase 7**: ~14-19 days

---

## Timeline Summary

| Phase | Description | Estimated Duration | Target Version |
|-------|-------------|-------------------|----------------|
| 1 | Foundation | 15-20 days | v0.6.0-alpha |
| 2 | Import Resolution | 17-22 days | v0.6.0-beta |
| 3 | Module Initialization | 12-16 days | v0.6.0-rc1 |
| 4 | Bytecode Format | 10-14 days | v0.6.0-rc2 |
| 5 | Stdlib Reorganization | 17-23 days | v0.6.1 |
| 6 | Advanced Features | 10-15 days | v0.6.2 |
| 7 | Tooling | 14-19 days | v0.6.3 |

**Total Estimated Time**: 95-129 days (~3-4 months)

## Risk Assessment

### High Risk Items

1. **Backward Compatibility**: Ensuring existing code continues to work
   - Mitigation: Extensive testing, compatibility mode
   
2. **Circular Dependencies**: Complex to detect and resolve
   - Mitigation: Start with simple detection, fail fast
   
3. **Performance Impact**: Module loading overhead
   - Mitigation: Bytecode caching, lazy loading

### Medium Risk Items

1. **Namespace Complexity**: Managing namespaces across imports
   - Mitigation: Clear specification, good error messages
   
2. **Initialization Order**: Dependencies between module initializations
   - Mitigation: Topological sorting, clear error reporting

### Low Risk Items

1. **Tool Development**: Independent of core implementation
   - Mitigation: Can be delayed if needed

## Success Criteria

### Phase Completion Criteria

Each phase must meet these criteria before moving to the next:

1. **All planned features implemented**
2. **Unit tests passing** (>90% coverage for new code)
3. **Integration tests passing**
4. **Documentation complete**
5. **Code review completed**
6. **No critical bugs**

### Overall Success Criteria

1. **Backward Compatibility**: All existing examples and tests still work
2. **Performance**: No more than 10% overhead for non-modular code
3. **Usability**: Clear error messages, good documentation
4. **Completeness**: Can organize stdlib into modules
5. **Tools**: Developer tools make module development easy

## Testing Strategy

### Unit Testing

- Test each component in isolation
- Mock dependencies
- Test edge cases and error conditions
- Aim for >90% code coverage

### Integration Testing

- Test module loading and imports end-to-end
- Test with multiple modules
- Test initialization ordering
- Test bytecode compilation and loading

### Compatibility Testing

- Run all existing tests without modification
- Run all examples without modification
- Verify no performance regression on non-modular code

### Performance Testing

- Benchmark module loading times
- Compare .smog vs .sg loading
- Test with large numbers of modules
- Measure memory usage

## Migration Path for Existing Code

### Phase 1: No Changes Required
All existing code works without modification.

### Phase 2: Optional Module Declarations
Developers can optionally add module declarations to new code.

### Phase 3: Standard Library Migration
Standard library is reorganized, with backward compatibility layer.

### Phase 4: Full Migration
Projects can be fully modularized when ready.

## Documentation Plan

### For Users

- [ ] Module system tutorial
- [ ] Import statement guide
- [ ] Best practices for organizing code
- [ ] Migration guide from non-modular code
- [ ] Common patterns and examples

### For Developers

- [ ] Module system architecture document
- [ ] API documentation for module package
- [ ] Bytecode format changes
- [ ] Testing guide for modules

## Open Questions

1. **Default Imports**: Should Core modules be auto-imported?
2. **Module Reloading**: Support during development?
3. **Private Classes**: Need visibility control within modules?
4. **Module Versions**: Semantic versioning support?
5. **Remote Modules**: Future support for downloading modules?

## Dependencies

### Required Before Starting

- [x] Advanced class features complete (v0.5.0)
- [x] Bytecode format stable
- [x] Comprehensive test suite

### Blocking Issues

None identified at this time.

## References

- [Module System Specification](../spec/MODULE_SYSTEM_SPEC.md)
- [Language Specification](../spec/LANGUAGE_SPEC.md)
- [Bytecode Format Guide](../BYTECODE_FORMAT.md)
- [Development Roadmap](ROADMAP.md)

## Appendix: Code Structure Changes

### New Packages

```
pkg/module/
  ├── module.go           - Core module type
  ├── registry.go         - Module registry
  ├── loader.go           - Module loading
  ├── resolver.go         - Path and name resolution
  ├── namespace.go        - Namespace management
  ├── metadata.go         - Module metadata
  ├── initializer.go      - Module initialization
  └── analyzer.go         - Dependency analysis

stdlib/
  ├── Core/               - Core classes
  ├── Collections/        - Collection classes
  ├── Math/              - Math classes
  ├── Blocks/            - Block and continuation
  └── Compat/            - Backward compatibility
```

### Modified Packages

```
pkg/lexer/              - Module declaration tokens
pkg/parser/             - Module declaration parsing
pkg/ast/                - Module AST nodes
pkg/compiler/           - Module-aware name resolution
pkg/bytecode/           - Module metadata serialization
cmd/smog/               - New module-related commands
```

## Conclusion

This implementation plan provides a structured approach to adding a module system to Smog. By breaking the work into phases, we can:

1. Deliver value incrementally
2. Test thoroughly at each phase
3. Maintain backward compatibility
4. Adjust based on feedback
5. Minimize risk

The estimated timeline of 3-4 months is realistic for a comprehensive module system that will enable Smog to scale from simple scripts to large, well-organized programs.
