# Smog Microcontroller Porting Plan

Version 1.0.0 (Draft)

## Executive Summary

This document analyzes two approaches for running Smog on resource-constrained microcontrollers:

1. **C-based VM Implementation** - Rewriting the VM in C with custom memory management
2. **TinyGo Compilation** - Using TinyGo to compile the existing Go VM for embedded systems

Both approaches are evaluated for feasibility, implementation effort, performance characteristics, and long-term maintainability. This analysis provides the foundation for deciding which path best serves Smog's goals for embedded deployment.

## Current State Analysis

### Existing Implementation

**Language**: Go 1.24.11  
**VM Architecture**: Stack-based bytecode interpreter  
**Memory Management**: Go's automatic garbage collection  
**Bytecode Format**: Binary .sg files with versioning  

**Key VM Components**:
- Execution stack (1024 slots)
- Local variables array (256 slots)
- Global variables map (dynamic)
- Constant pool (read-only)
- Instruction pointer and stack pointer

**Memory Characteristics**:
- Stack: ~8 KB (1024 × 8 bytes)
- Locals: ~2 KB (256 × 8 bytes)
- Globals: Variable (hash map overhead)
- Baseline: ~10-15 KB per VM instance
- Go runtime overhead: ~200 KB - 2 MB

### Target Microcontroller Constraints

**Typical Target Platforms**:
- ARM Cortex-M0/M0+: 32 KB RAM, 256 KB Flash
- ARM Cortex-M3/M4: 64-256 KB RAM, 512 KB - 2 MB Flash
- ESP32-C3: 400 KB RAM, 4 MB Flash
- RP2040: 264 KB RAM, 2 MB Flash

**Critical Constraints**:
- Limited RAM (32-400 KB typical)
- No virtual memory or swap
- Flash storage for code/bytecode
- Low power requirements
- Real-time constraints possible
- Minimal OS support (bare metal or RTOS)

## Option 1: C-based VM Implementation

### Overview

Reimplement the Smog VM in portable C, with explicit memory management suitable for embedded systems. The bytecode format and compiler remain unchanged - only the execution engine is rewritten.

### Architecture Design

#### Core Components

```c
// VM State Structure
typedef struct {
    Value stack[1024];          // Value stack
    uint16_t sp;                // Stack pointer
    Value locals[256];          // Local variables
    HashMap* globals;           // Global variables
    const Value* constants;     // Constant pool (read-only)
    const Instruction* code;    // Bytecode instructions
    uint32_t ip;                // Instruction pointer
    MemoryManager* mm;          // Memory manager
} VM;

// Value representation (tagged union)
typedef struct {
    ValueType type;
    union {
        int64_t integer;
        double floating;
        Object* object;
        const char* string;
        bool boolean;
    } as;
} Value;

// Object header for heap-allocated objects
typedef struct Object {
    ObjectType type;
    bool marked;                // GC mark bit
    struct Object* next;        // Free list / GC linked list
    // ... type-specific data follows
} Object;
```

#### Memory Management Strategy

**Arena Allocation**:
```c
typedef struct {
    uint8_t* memory;            // Pre-allocated memory pool
    size_t size;                // Total pool size
    size_t used;                // Currently allocated
    Object* objects;            // Linked list of all objects
    size_t gc_threshold;        // Trigger GC when used > threshold
} MemoryManager;
```

**Benefits**:
- Predictable allocation performance
- No fragmentation (or controlled fragmentation)
- Can allocate entire pool at startup
- Easy to implement

**Challenges**:
- Fixed memory limit
- May waste space if over-allocated
- Requires GC to reclaim space

### Garbage Collection Options

#### Option 1A: Mark-and-Sweep GC

**Algorithm**:
```c
void gc_collect(VM* vm) {
    // 1. Mark phase - trace from roots
    gc_mark_roots(vm);
    
    // 2. Sweep phase - free unmarked objects
    Object* obj = vm->mm->objects;
    Object* prev = NULL;
    
    while (obj != NULL) {
        if (!obj->marked) {
            // Unreachable - free it
            Object* unreached = obj;
            obj = obj->next;
            if (prev != NULL) {
                prev->next = obj;
            } else {
                vm->mm->objects = obj;
            }
            free_object(unreached);
        } else {
            // Reachable - unmark for next cycle
            obj->marked = false;
            prev = obj;
            obj = obj->next;
        }
    }
}

void gc_mark_value(Value value) {
    if (value.type != VALUE_OBJECT) return;
    if (value.as.object->marked) return;
    
    value.as.object->marked = true;
    
    // Mark objects referenced by this object
    mark_object_references(value.as.object);
}

void gc_mark_roots(VM* vm) {
    // Mark stack values
    for (int i = 0; i < vm->sp; i++) {
        gc_mark_value(vm->stack[i]);
    }
    
    // Mark local variables
    for (int i = 0; i < 256; i++) {
        gc_mark_value(vm->locals[i]);
    }
    
    // Mark global variables
    mark_hashmap_values(vm->globals);
}
```

**Characteristics**:
- **Pause time**: Proportional to live objects (10-100ms typical)
- **Memory overhead**: 1 bit per object (mark bit)
- **Implementation**: Simple, well-understood
- **Fragmentation**: Can be an issue over time

**Pros**:
- Simple to implement
- Handles cycles naturally
- Predictable behavior
- Low memory overhead

**Cons**:
- Stop-the-world pauses
- Pause time proportional to heap size
- Potential fragmentation
- Not real-time friendly

#### Option 1B: Reference Counting

**Implementation**:
```c
typedef struct Object {
    ObjectType type;
    uint16_t ref_count;
    struct Object* next;
    // ... type-specific data
} Object;

void retain_object(Object* obj) {
    if (obj != NULL) {
        obj->ref_count++;
    }
}

void release_object(Object* obj) {
    if (obj == NULL) return;
    
    obj->ref_count--;
    if (obj->ref_count == 0) {
        // Release all objects this one references
        release_object_references(obj);
        free_object(obj);
    }
}

void set_value(Value* dest, Value src) {
    // Release old value
    if (dest->type == VALUE_OBJECT) {
        release_object(dest->as.object);
    }
    
    // Set new value and retain
    *dest = src;
    if (src.type == VALUE_OBJECT) {
        retain_object(src.as.object);
    }
}
```

**Characteristics**:
- **Pause time**: None (incremental)
- **Memory overhead**: 2 bytes per object (ref count)
- **Implementation**: Moderate complexity
- **Fragmentation**: Can be an issue

**Pros**:
- Immediate reclamation
- No GC pauses
- Predictable timing
- Real-time friendly

**Cons**:
- Cannot handle cycles (needs cycle detector)
- Overhead on every assignment
- Memory overhead for counts
- Complex to implement correctly

#### Option 1C: Hybrid Approach (Recommended)

**Strategy**:
- Reference counting for immediate reclamation
- Periodic mark-and-sweep for cycle collection
- Trigger sweep only when needed

```c
typedef struct Object {
    ObjectType type;
    uint16_t ref_count;
    bool marked;
    bool potentially_cyclic;    // Flag for cycle detection
    struct Object* next;
} Object;

void gc_collect_hybrid(VM* vm, bool force_sweep) {
    // Fast path: reference counting handles most objects
    // Slow path: mark-and-sweep handles cycles
    
    if (force_sweep || vm->mm->potentially_cyclic_count > CYCLE_THRESHOLD) {
        // Run mark-and-sweep to collect cycles
        gc_mark_and_sweep(vm);
        vm->mm->potentially_cyclic_count = 0;
    }
}
```

**Pros**:
- Combines benefits of both approaches
- Immediate reclamation for acyclic structures
- Handles cycles correctly
- Can tune for application needs

**Cons**:
- Most complex to implement
- Still has occasional pauses
- Higher memory overhead

#### Recommendation: Mark-and-Sweep

For the initial C VM implementation, **mark-and-sweep** is recommended:

**Rationale**:
1. Simplest to implement correctly
2. Handles all object graphs (including cycles)
3. Lower memory overhead
4. Easier to debug and verify
5. Can optimize later if needed

**For Real-Time Systems**: If real-time constraints are identified later, can migrate to reference counting or hybrid approach.

### Implementation Plan

#### Phase 1: Core VM in C (4-6 weeks)

**Week 1-2: Basic VM Structure**
- [ ] Define value representation (tagged union)
- [ ] Implement stack operations (push, pop, peek)
- [ ] Implement local variable array
- [ ] Implement global variable hash map
- [ ] Create instruction decoder

**Files to Create**:
- `c-vm/src/value.h` - Value type definitions
- `c-vm/src/value.c` - Value operations
- `c-vm/src/vm.h` - VM structure and API
- `c-vm/src/vm.c` - VM implementation
- `c-vm/src/stack.c` - Stack operations

**Week 3-4: Instruction Execution**
- [ ] Implement all opcodes from bytecode spec
- [ ] Implement primitive operations (arithmetic, comparison)
- [ ] Implement message dispatch framework
- [ ] Add error handling and stack traces

**Files to Create**:
- `c-vm/src/opcodes.h` - Opcode definitions
- `c-vm/src/execute.c` - Instruction execution
- `c-vm/src/primitives.c` - Primitive operations
- `c-vm/src/errors.c` - Error handling

**Week 5-6: Bytecode Loading**
- [ ] Implement .sg file parser
- [ ] Load constant pool
- [ ] Load instructions
- [ ] Validate bytecode format

**Files to Create**:
- `c-vm/src/bytecode.h` - Bytecode format
- `c-vm/src/loader.c` - Bytecode loading
- `c-vm/src/validator.c` - Bytecode validation

#### Phase 2: Memory Management (3-4 weeks)

**Week 7-8: Arena Allocator**
- [ ] Implement memory pool allocation
- [ ] Implement object allocation
- [ ] Add memory statistics
- [ ] Add out-of-memory handling

**Files to Create**:
- `c-vm/src/memory.h` - Memory manager API
- `c-vm/src/memory.c` - Arena allocator
- `c-vm/src/object.h` - Object structures
- `c-vm/src/object.c` - Object operations

**Week 9-10: Garbage Collector**
- [ ] Implement mark phase
- [ ] Implement sweep phase
- [ ] Add GC triggering logic
- [ ] Add GC statistics and tuning

**Files to Create**:
- `c-vm/src/gc.h` - GC API
- `c-vm/src/gc.c` - GC implementation

#### Phase 3: Object System (2-3 weeks)

**Week 11-12: Object Types**
- [ ] Implement String objects
- [ ] Implement Array objects
- [ ] Implement Block/closure objects
- [ ] Implement Class objects (basic)

**Files to Create**:
- `c-vm/src/string.c` - String implementation
- `c-vm/src/array.c` - Array implementation
- `c-vm/src/block.c` - Block/closure implementation
- `c-vm/src/class.c` - Class representation

#### Phase 4: Testing and Optimization (2-3 weeks)

**Week 13: Testing**
- [ ] Unit tests for VM operations
- [ ] Integration tests with existing bytecode
- [ ] Memory leak detection
- [ ] Stress testing

**Week 14-15: Optimization**
- [ ] Profile and optimize hot paths
- [ ] Optimize memory usage
- [ ] Add inline caching for message dispatch
- [ ] Benchmark against Go VM

**Total Estimated Time**: 11-16 weeks (3-4 months)

### Resource Requirements

**Development Requirements**:
- C compiler (GCC, Clang)
- Cross-compilation toolchain for target MCUs
- Debugging tools (GDB, valgrind)
- Memory profiling tools
- Unit testing framework (Check, Unity, or custom)

**Target Hardware Testing**:
- Development boards for target MCUs
- JTAG debugger
- Logic analyzer (optional)
- Power measurement tools (optional)

### Memory Footprint Estimation

**Code Size** (Flash):
- VM core: ~20-30 KB
- Garbage collector: ~5-10 KB
- Object system: ~10-15 KB
- Standard library primitives: ~15-25 KB
- **Total: ~50-80 KB**

**Runtime Memory** (RAM):
- VM state: ~15 KB (stack + locals + metadata)
- Object heap: Configurable (16-128 KB typical)
- **Minimum: ~32 KB total**
- **Recommended: 64 KB or more**

**Suitable Platforms**:
- ✅ ARM Cortex-M3/M4 (64+ KB RAM)
- ✅ ESP32-C3 (400 KB RAM)
- ✅ RP2040 (264 KB RAM)
- ⚠️ ARM Cortex-M0+ (32 KB RAM) - Tight fit, limited heap
- ❌ ARM Cortex-M0 (16 KB RAM) - Insufficient

### Performance Characteristics

**Expected Performance**:
- Instruction throughput: 1-10 million instructions/sec (100 MHz MCU)
- GC pause: 10-100 ms (depends on heap size)
- Startup time: <10 ms (bytecode loading)
- Power consumption: 10-50 mA active (platform dependent)

**Compared to Go VM**:
- Memory: 5-10x reduction
- Speed: Similar (within 2x)
- Binary size: 10-20x reduction

## Option 2: TinyGo Compilation

### Overview

TinyGo is a Go compiler for embedded systems and WebAssembly, designed to produce small binaries with minimal runtime overhead. It supports a subset of Go and uses LLVM for optimization.

### TinyGo Capabilities

**Supported Features**:
- Most Go language features (interfaces, goroutines, channels)
- Garbage collection (conservative, non-moving)
- Reflection (limited)
- Standard library (subset)

**Limitations**:
- No full reflect package
- Limited runtime/debug support
- Some standard library packages unavailable
- No CGo support on all platforms
- Smaller stack sizes

**Garbage Collection**:
- Conservative, non-moving mark-and-sweep
- Optional: Extalloc (external allocator) or no GC
- GC pause: ~1-10 ms typical
- Minimal overhead

### Compatibility Analysis

#### Smog Codebase Review

**Current Dependencies** (from go.mod):
```
module github.com/kristofer/smog
go 1.24.11
```

**Standard Library Usage**:
- `fmt` - Formatting and printing ✅ Supported
- `errors` - Error handling ✅ Supported
- `os` - File operations ✅ Supported (limited)
- `io` - I/O interfaces ✅ Supported
- `strings` - String manipulation ✅ Supported
- `bytes` - Byte operations ✅ Supported

**Assessment**: ✅ All current dependencies are TinyGo compatible

#### Feature Compatibility

**VM Core**:
- ✅ Stack-based execution
- ✅ Value types (int64, float64, string, bool)
- ✅ Hash maps (globals)
- ✅ Slice operations (stack, locals, bytecode)
- ✅ Interfaces (value polymorphism)

**Bytecode System**:
- ✅ Binary file I/O
- ✅ Constant pool
- ✅ Instruction decoding

**Object System**:
- ✅ Interfaces for object types
- ✅ Method dispatch
- ⚠️ Reflection usage (need to check)

**Potential Issues**:
1. File I/O may need adaptation for embedded filesystems
2. Error messages may need size optimization
3. Some debug/introspection features may not work

### Implementation Plan

#### Phase 1: TinyGo Compatibility (2-3 weeks)

**Week 1: Dependency Audit**
- [ ] Audit all Go dependencies for TinyGo compatibility
- [ ] Identify unsupported features
- [ ] Create compatibility shims where needed
- [ ] Test compilation with TinyGo

**Week 2: Build System**
- [ ] Create TinyGo build targets
- [ ] Configure linker settings for target platforms
- [ ] Set up cross-compilation
- [ ] Create build scripts

**Week 3: Testing**
- [ ] Run existing test suite with TinyGo
- [ ] Fix compatibility issues
- [ ] Benchmark binary size and performance
- [ ] Document limitations

#### Phase 2: Embedded Optimization (2-3 weeks)

**Week 4: Size Optimization**
- [ ] Remove unused code paths
- [ ] Optimize string literals
- [ ] Minimize error messages
- [ ] Enable aggressive optimization flags

**Week 5: Memory Optimization**
- [ ] Tune GC parameters
- [ ] Reduce stack sizes where safe
- [ ] Pool allocations
- [ ] Pre-allocate common objects

**Week 6: Platform Support**
- [ ] Test on target development boards
- [ ] Create platform-specific builds
- [ ] Handle platform differences (filesystem, I/O)
- [ ] Document platform requirements

#### Phase 3: Runtime Adaptation (1-2 weeks)

**Week 7-8: Embedded Runtime**
- [ ] Replace `fmt.Println` with embedded-friendly output
- [ ] Adapt file I/O for embedded filesystems
- [ ] Add UART/serial I/O primitives
- [ ] Test on bare metal RTOS

**Total Estimated Time**: 5-8 weeks (1.5-2 months)

### Resource Requirements

**Development Requirements**:
- TinyGo compiler (latest version)
- LLVM toolchain
- Target platform SDK
- Flash programming tools
- Serial/UART interface

**Testing Hardware**:
- Development boards for target platforms
- USB-to-serial adapter
- Logic analyzer (optional)

### Memory Footprint Estimation

**Code Size** (Flash):
- TinyGo runtime: ~20-40 KB
- VM code: ~30-50 KB (compiled from Go)
- Standard library: ~10-20 KB
- **Total: ~60-110 KB**

**Runtime Memory** (RAM):
- VM state: ~15 KB
- TinyGo runtime: ~5-15 KB
- Object heap: Configurable (16-128 KB)
- **Minimum: ~40 KB total**
- **Recommended: 64 KB or more**

**Suitable Platforms**:
- ✅ ARM Cortex-M3/M4 (64+ KB RAM)
- ✅ ESP32-C3 (400 KB RAM)
- ✅ RP2040 (264 KB RAM)
- ⚠️ ARM Cortex-M0+ (32 KB RAM) - Very tight
- ❌ ARM Cortex-M0 (16 KB RAM) - Insufficient

### Performance Characteristics

**Expected Performance**:
- Similar to C VM (TinyGo compiles to native code)
- GC pause: 1-10 ms (TinyGo conservative GC)
- Startup time: <50 ms
- Instruction throughput: 1-10 million/sec

**Compared to Standard Go**:
- Binary size: 10-50x smaller
- Memory usage: 5-10x lower
- Startup time: 10-100x faster
- Runtime overhead: Minimal

## Comparison Matrix

| Aspect | C VM Implementation | TinyGo Compilation |
|--------|--------------------|--------------------|
| **Development Time** | 3-4 months | 1.5-2 months |
| **Implementation Complexity** | High | Low |
| **Code Size (Flash)** | 50-80 KB | 60-110 KB |
| **RAM Usage** | 32+ KB | 40+ KB |
| **Performance** | Native, optimized | Native, optimized |
| **GC Pause Time** | 10-100 ms | 1-10 ms |
| **Portability** | Maximum (pure C) | Limited to TinyGo targets |
| **Maintainability** | Two codebases | Single codebase |
| **Debugging** | Standard C tools | TinyGo/LLVM tools |
| **Learning Curve** | C + GC expertise needed | Minimal (already Go) |
| **Risk** | Medium-High | Low-Medium |
| **Platform Support** | Universal (any C compiler) | TinyGo supported platforms |
| **Future Optimization** | Full control | Limited by TinyGo |

## Detailed Comparison

### Development Effort

**C VM**:
- **Pros**: Full control over implementation, potential for maximum optimization
- **Cons**: Significant engineering effort, need to reimplement everything, potential for bugs
- **Effort**: 3-4 months initial implementation
- **Ongoing**: Maintain two implementations in parallel

**TinyGo**:
- **Pros**: Leverage existing codebase, minimal changes needed, proven toolchain
- **Cons**: Limited to TinyGo capabilities, less control over low-level details
- **Effort**: 1.5-2 months adaptation
- **Ongoing**: Single codebase maintenance

### Memory and Performance

**C VM**:
- **Code Size**: Smaller (50-80 KB) - hand-optimized C
- **RAM Usage**: Lower (32+ KB minimum) - explicit management
- **GC Control**: Full control over GC implementation and tuning
- **Performance**: Native code, optimized for specific use cases

**TinyGo**:
- **Code Size**: Slightly larger (60-110 KB) - includes runtime
- **RAM Usage**: Moderate (40+ KB minimum) - TinyGo runtime overhead
- **GC Control**: Limited to TinyGo's conservative GC
- **Performance**: Native code, LLVM optimized

### Portability

**C VM**:
- **Platforms**: Any platform with a C compiler
- **Customization**: Can adapt to any embedded environment
- **Bare Metal**: Easy to run without OS
- **Exotic Targets**: Can port to unusual architectures

**TinyGo**:
- **Platforms**: Limited to TinyGo-supported targets
- **Current Support**: ARM, AVR, RISC-V, WebAssembly, x86
- **Expanding**: TinyGo team actively adding platforms
- **Limitations**: May not support all desired targets

### Long-term Considerations

**C VM**:
- **Pros**: 
  - Complete control for future optimization
  - Can implement advanced GC strategies
  - No dependency on external project
  - Maximum portability
- **Cons**: 
  - Permanent maintenance burden of two implementations
  - Features must be implemented twice
  - Higher risk of divergence
  - More testing required

**TinyGo**:
- **Pros**: 
  - Single codebase = easier maintenance
  - Automatic benefit from Go VM improvements
  - TinyGo improvements benefit Smog
  - Lower maintenance burden
- **Cons**: 
  - Dependency on TinyGo project
  - Limited by TinyGo's capabilities
  - Platform support depends on TinyGo
  - Less control over low-level optimizations

## Recommendations

### Primary Recommendation: TinyGo (Phase 1)

**Start with TinyGo** for the following reasons:

1. **Faster Time to Market**: 1.5-2 months vs 3-4 months
2. **Lower Risk**: Leverage existing tested code
3. **Single Codebase**: Easier long-term maintenance
4. **Proven Technology**: TinyGo is mature and well-supported
5. **Good Enough**: Meets requirements for most target platforms
6. **Reversible**: Can still do C VM later if needed

**Approach**:
1. Implement TinyGo support in 1.5-2 months
2. Test on target platforms (ESP32, RP2040, ARM Cortex-M4)
3. Evaluate performance and limitations
4. Use in production if meets requirements
5. Gather real-world experience

### Contingency: C VM (Phase 2)

**Consider C VM if**:
- TinyGo doesn't support critical target platform
- GC pauses are unacceptable for application
- Memory footprint is too large
- Need maximum performance optimization
- Want to eliminate dependency on TinyGo

**Approach**:
1. Start C VM development only after TinyGo evaluation
2. Use TinyGo implementation as reference
3. Focus on platforms where TinyGo doesn't work
4. Consider as long-term strategic option

### Hybrid Approach

**Long-term Vision**:
1. **Year 1**: TinyGo implementation for most platforms
2. **Year 2**: Evaluate need for C VM based on real usage
3. **Year 3**: Implement C VM if strategic benefit identified

**Rationale**: Don't over-engineer upfront. Let actual usage and requirements drive the C VM decision.

## Implementation Roadmap

### Phase 1: TinyGo Support (Months 1-2)

**Month 1**:
- Week 1: Compatibility audit and testing
- Week 2: Build system and cross-compilation
- Week 3: Size optimization
- Week 4: Memory optimization

**Month 2**:
- Week 5-6: Platform-specific testing
- Week 7-8: Runtime adaptation and documentation

**Deliverables**:
- TinyGo-compatible Smog VM
- Build scripts for target platforms
- Platform-specific documentation
- Performance benchmarks
- Decision document on C VM need

### Phase 2: Evaluation and Decision (Month 3)

**Activities**:
- Deploy on target hardware
- Measure real-world performance
- Gather user feedback
- Evaluate limitations

**Decision Point**: Proceed with C VM or optimize TinyGo implementation?

### Phase 3: C VM (Conditional, Months 4-7)

Only if Phase 2 evaluation determines C VM is necessary:
- Months 4-5: Core VM implementation
- Month 6: Memory management and GC
- Month 7: Object system and testing

## Success Criteria

### TinyGo Implementation Success

- ✅ Runs on target platforms (ESP32, RP2040, ARM Cortex-M4)
- ✅ Binary size < 120 KB
- ✅ RAM usage < 64 KB minimum
- ✅ GC pause < 20 ms
- ✅ Passes all VM tests
- ✅ Performance within 2x of desktop Go VM

### C VM Implementation Success (if pursued)

- ✅ Runs on all C-compatible platforms
- ✅ Binary size < 100 KB
- ✅ RAM usage < 40 KB minimum
- ✅ GC pause configurable and predictable
- ✅ Performance equivalent to Go VM
- ✅ Passes all VM tests
- ✅ No memory leaks

## Risk Analysis

### TinyGo Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Platform not supported | Low | High | Check TinyGo platform list early |
| Performance inadequate | Low | Medium | Profile and optimize, fallback to C |
| GC pauses too long | Medium | Medium | Tune GC, reduce heap size |
| Binary too large | Low | Low | Optimize, remove unused code |
| TinyGo bug/limitation | Medium | Medium | Report upstream, workaround, or C VM |

### C VM Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Implementation bugs | High | High | Extensive testing, fuzzing |
| GC bugs/leaks | Medium | High | Valgrind, sanitizers, testing |
| Schedule overrun | Medium | Medium | Phased approach, MVP first |
| Maintenance burden | High | Medium | Good documentation, tests |
| Divergence from Go VM | High | Medium | Shared bytecode format, test suite |

## Conclusion

**Recommended Path**: Start with **TinyGo implementation** (Phase 1)

**Rationale**:
- Fastest path to embedded deployment
- Lowest risk and development effort
- Single codebase maintenance
- Adequate for most target platforms
- Can evaluate C VM need with real data

**Next Steps**:
1. Approve this plan
2. Begin Phase 1: TinyGo compatibility work
3. Set up target hardware for testing
4. Schedule Phase 2 evaluation in 3 months
5. Make C VM decision based on real-world data

The C VM remains a valuable option for the future, but should be pursued only after TinyGo has been tried and evaluated. This approach minimizes risk, accelerates time to market, and ensures any C VM investment is driven by actual need rather than speculation.

---

**Document Version**: 1.0.0  
**Date**: 2026-01-16  
**Status**: Draft - Awaiting Approval  
**Author**: Smog Development Team
