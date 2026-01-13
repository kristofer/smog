// Package bytecode provides serialization and deserialization for .sg bytecode files.
//
// File Format Specification:
//
// The .sg file format is a binary format for storing compiled Smog bytecode.
// It allows pre-compilation of .smog source files to bytecode for faster loading
// and execution. The format is designed to be:
//   - Compact: Efficient binary encoding
//   - Versioned: Support for format evolution
//   - Complete: Stores all information needed for execution
//
// Binary Format Layout:
//
//   [Header]
//     Magic Number (4 bytes): "SMOG" (0x534D4F47)
//     Version (4 bytes): Format version number (currently 1)
//     Flags (4 bytes): Reserved for future use
//
//   [Constants Section]
//     Count (4 bytes): Number of constants
//     For each constant:
//       Type (1 byte): Constant type identifier
//       Data (variable): Type-specific encoding
//
//   [Instructions Section]
//     Count (4 bytes): Number of instructions
//     For each instruction:
//       Opcode (1 byte): Operation code
//       Operand (4 bytes): Instruction operand
//
// Constant Types:
//   0x01 = Integer (int64, 8 bytes)
//   0x02 = Float (float64, 8 bytes)
//   0x03 = String (4-byte length + UTF-8 bytes)
//   0x04 = Boolean (1 byte: 0=false, 1=true)
//   0x05 = Nil (0 bytes)
//   0x06 = ClassDefinition (nested structure)
//   0x07 = MethodDefinition (nested structure)
//   0x08 = Bytecode (recursive structure for blocks/methods)
//
// Example:
//
//   Source: 'Hello' println. 42.
//
//   .sg file:
//     Header: SMOG 0x00000001 0x00000000
//     Constants: count=3
//       [0] String: "Hello"
//       [1] String: "println"
//       [2] Integer: 42
//     Instructions: count=5
//       PUSH 0
//       SEND (1<<8)|0
//       POP 0
//       PUSH 2
//       RETURN 0
//
// Design Rationale:
//
// Binary Format:
//   - Faster to parse than text formats
//   - Smaller file size
//   - Direct mapping to in-memory structures
//
// Magic Number:
//   - Identifies file type
//   - Prevents accidental execution of wrong files
//
// Version Number:
//   - Allows format evolution
//   - Future versions can add features while maintaining compatibility
//
// This format is inspired by:
//   - Java .class files
//   - Python .pyc files
//   - Smalltalk image formats
package bytecode

import (
	"encoding/binary"
	"fmt"
	"io"
)

// File format constants
const (
	// MagicNumber is the file signature for .sg files: "SMOG"
	MagicNumber uint32 = 0x534D4F47

	// FormatVersion is the current bytecode format version
	FormatVersion uint32 = 1

	// Reserved flags (currently unused, set to 0)
	formatFlags uint32 = 0
)

// Constant type identifiers for serialization
const (
	constTypeInteger   byte = 0x01
	constTypeFloat     byte = 0x02
	constTypeString    byte = 0x03
	constTypeBoolean   byte = 0x04
	constTypeNil       byte = 0x05
	constTypeClass     byte = 0x06
	constTypeMethod    byte = 0x07
	constTypeBytecode  byte = 0x08
)

// Encode serializes bytecode to binary format and writes it to w.
//
// This function takes compiled bytecode and writes it to an io.Writer
// (typically a file) in the .sg binary format. The output can be later
// loaded with Decode() and executed without re-parsing or re-compiling.
//
// Process:
//   1. Write header (magic number, version, flags)
//   2. Write constants section
//   3. Write instructions section
//
// Example usage:
//
//   // Compile source to bytecode
//   bc, _ := compiler.Compile(program)
//
//   // Save to .sg file
//   file, _ := os.Create("program.sg")
//   defer file.Close()
//   bytecode.Encode(bc, file)
//
// Returns an error if writing fails or if the bytecode contains
// unsupported types.
func Encode(bc *Bytecode, w io.Writer) error {
	// Write header
	if err := writeHeader(w); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write constants section
	if err := writeConstants(w, bc.Constants); err != nil {
		return fmt.Errorf("failed to write constants: %w", err)
	}

	// Write instructions section
	if err := writeInstructions(w, bc.Instructions); err != nil {
		return fmt.Errorf("failed to write instructions: %w", err)
	}

	return nil
}

// Decode deserializes bytecode from binary format.
//
// This function reads a .sg file and reconstructs the bytecode structure
// in memory, ready for execution by the VM. It's the inverse of Encode().
//
// Process:
//   1. Read and validate header
//   2. Read constants section
//   3. Read instructions section
//
// Example usage:
//
//   // Load .sg file
//   file, _ := os.Open("program.sg")
//   defer file.Close()
//   bc, _ := bytecode.Decode(file)
//
//   // Execute with VM
//   vm := vm.New()
//   vm.Run(bc)
//
// Returns an error if:
//   - Magic number is incorrect (not a .sg file)
//   - Version is unsupported
//   - File is corrupted
//   - Unexpected end of file
func Decode(r io.Reader) (*Bytecode, error) {
	// Read and validate header
	version, err := readHeader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Check version compatibility
	if version != FormatVersion {
		return nil, fmt.Errorf("unsupported bytecode version: %d (expected %d)", version, FormatVersion)
	}

	// Read constants section
	constants, err := readConstants(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read constants: %w", err)
	}

	// Read instructions section
	instructions, err := readInstructions(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read instructions: %w", err)
	}

	return &Bytecode{
		Instructions: instructions,
		Constants:    constants,
	}, nil
}

// writeHeader writes the file header to w.
//
// Header format:
//   - Magic number (4 bytes): File signature
//   - Version (4 bytes): Format version
//   - Flags (4 bytes): Reserved for future use
func writeHeader(w io.Writer) error {
	// Write magic number
	if err := binary.Write(w, binary.LittleEndian, MagicNumber); err != nil {
		return err
	}

	// Write version
	if err := binary.Write(w, binary.LittleEndian, FormatVersion); err != nil {
		return err
	}

	// Write flags (reserved, currently 0)
	if err := binary.Write(w, binary.LittleEndian, formatFlags); err != nil {
		return err
	}

	return nil
}

// readHeader reads and validates the file header from r.
//
// Returns the format version if successful, or an error if:
//   - Magic number doesn't match (wrong file type)
//   - Read fails (corrupted file or I/O error)
func readHeader(r io.Reader) (uint32, error) {
	// Read and verify magic number
	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return 0, err
	}

	if magic != MagicNumber {
		return 0, fmt.Errorf("invalid magic number: 0x%08X (expected 0x%08X)", magic, MagicNumber)
	}

	// Read version
	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return 0, err
	}

	// Read flags (currently ignored)
	var flags uint32
	if err := binary.Read(r, binary.LittleEndian, &flags); err != nil {
		return 0, err
	}

	return version, nil
}

// writeConstants writes the constants section to w.
//
// Format:
//   - Count (4 bytes): Number of constants
//   - For each constant: type byte + type-specific data
//
// Supported constant types:
//   - int64: 8 bytes (little-endian)
//   - float64: 8 bytes (IEEE 754)
//   - string: 4-byte length + UTF-8 bytes
//   - bool: 1 byte (0 or 1)
//   - nil: just the type byte
//   - ClassDefinition: nested structure
//   - MethodDefinition: nested structure
//   - *Bytecode: recursively encoded bytecode (for blocks/methods)
func writeConstants(w io.Writer, constants []interface{}) error {
	// Write count
	count := uint32(len(constants))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}

	// Write each constant
	for i, c := range constants {
		if err := writeConstant(w, c); err != nil {
			return fmt.Errorf("failed to write constant %d: %w", i, err)
		}
	}

	return nil
}

// writeConstant writes a single constant value to w.
//
// The format is: type byte followed by type-specific data.
// This function handles all the constant types that can appear
// in the constant pool.
func writeConstant(w io.Writer, c interface{}) error {
	switch v := c.(type) {
	case int64:
		// Integer: type byte + 8 bytes
		if err := binary.Write(w, binary.LittleEndian, constTypeInteger); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, v)

	case float64:
		// Float: type byte + 8 bytes (IEEE 754)
		if err := binary.Write(w, binary.LittleEndian, constTypeFloat); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, v)

	case string:
		// String: type byte + 4-byte length + UTF-8 bytes
		if err := binary.Write(w, binary.LittleEndian, constTypeString); err != nil {
			return err
		}
		length := uint32(len(v))
		if err := binary.Write(w, binary.LittleEndian, length); err != nil {
			return err
		}
		_, err := w.Write([]byte(v))
		return err

	case bool:
		// Boolean: type byte + 1 byte (0 or 1)
		if err := binary.Write(w, binary.LittleEndian, constTypeBoolean); err != nil {
			return err
		}
		var b byte
		if v {
			b = 1
		}
		return binary.Write(w, binary.LittleEndian, b)

	case nil:
		// Nil: just the type byte
		return binary.Write(w, binary.LittleEndian, constTypeNil)

	case *ClassDefinition:
		// ClassDefinition: complex nested structure
		if err := binary.Write(w, binary.LittleEndian, constTypeClass); err != nil {
			return err
		}
		return writeClassDefinition(w, v)

	case *MethodDefinition:
		// MethodDefinition: complex nested structure
		if err := binary.Write(w, binary.LittleEndian, constTypeMethod); err != nil {
			return err
		}
		return writeMethodDefinition(w, v)

	case *Bytecode:
		// Bytecode (for blocks/methods): recursively encode
		if err := binary.Write(w, binary.LittleEndian, constTypeBytecode); err != nil {
			return err
		}
		return Encode(v, w)

	default:
		return fmt.Errorf("unsupported constant type: %T", c)
	}
}

// readConstants reads the constants section from r.
//
// Returns a slice of constants that can contain:
//   - int64, float64, string, bool, nil values
//   - *ClassDefinition, *MethodDefinition
//   - *Bytecode (for blocks/methods)
func readConstants(r io.Reader) ([]interface{}, error) {
	// Read count
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, err
	}

	// Read each constant
	constants := make([]interface{}, count)
	for i := uint32(0); i < count; i++ {
		c, err := readConstant(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant %d: %w", i, err)
		}
		constants[i] = c
	}

	return constants, nil
}

// readConstant reads a single constant value from r.
//
// Reads the type byte first, then reads the appropriate data
// based on the type.
func readConstant(r io.Reader) (interface{}, error) {
	// Read type byte
	var constType byte
	if err := binary.Read(r, binary.LittleEndian, &constType); err != nil {
		return nil, err
	}

	// Read type-specific data
	switch constType {
	case constTypeInteger:
		var v int64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		return v, nil

	case constTypeFloat:
		var v float64
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		return v, nil

	case constTypeString:
		var length uint32
		if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
			return nil, err
		}
		buf := make([]byte, length)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return string(buf), nil

	case constTypeBoolean:
		var b byte
		if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
			return nil, err
		}
		return b != 0, nil

	case constTypeNil:
		return nil, nil

	case constTypeClass:
		return readClassDefinition(r)

	case constTypeMethod:
		return readMethodDefinition(r)

	case constTypeBytecode:
		return Decode(r)

	default:
		return nil, fmt.Errorf("unknown constant type: 0x%02X", constType)
	}
}

// writeInstructions writes the instructions section to w.
//
// Format:
//   - Count (4 bytes): Number of instructions
//   - For each instruction:
//       - Opcode (1 byte)
//       - Operand (4 bytes, signed)
func writeInstructions(w io.Writer, instructions []Instruction) error {
	// Write count
	count := uint32(len(instructions))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}

	// Write each instruction
	for i, instr := range instructions {
		// Write opcode (1 byte)
		if err := binary.Write(w, binary.LittleEndian, byte(instr.Op)); err != nil {
			return fmt.Errorf("failed to write instruction %d opcode: %w", i, err)
		}

		// Write operand (4 bytes, signed)
		if err := binary.Write(w, binary.LittleEndian, int32(instr.Operand)); err != nil {
			return fmt.Errorf("failed to write instruction %d operand: %w", i, err)
		}
	}

	return nil
}

// readInstructions reads the instructions section from r.
//
// Returns a slice of Instruction structs.
func readInstructions(r io.Reader) ([]Instruction, error) {
	// Read count
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, err
	}

	// Read each instruction
	instructions := make([]Instruction, count)
	for i := uint32(0); i < count; i++ {
		// Read opcode (1 byte)
		var op byte
		if err := binary.Read(r, binary.LittleEndian, &op); err != nil {
			return nil, fmt.Errorf("failed to read instruction %d opcode: %w", i, err)
		}

		// Read operand (4 bytes, signed)
		var operand int32
		if err := binary.Read(r, binary.LittleEndian, &operand); err != nil {
			return nil, fmt.Errorf("failed to read instruction %d operand: %w", i, err)
		}

		instructions[i] = Instruction{
			Op:      Opcode(op),
			Operand: int(operand),
		}
	}

	return instructions, nil
}

// writeClassDefinition writes a ClassDefinition to w.
//
// Format:
//   - Name (string: 4-byte length + UTF-8)
//   - SuperClass (string: 4-byte length + UTF-8)
//   - Field count (4 bytes) + field names (strings)
//   - ClassVar count (4 bytes) + classvar names (strings)
//   - Method count (4 bytes) + methods (MethodDefinitions)
//   - ClassMethod count (4 bytes) + class methods (MethodDefinitions)
func writeClassDefinition(w io.Writer, cd *ClassDefinition) error {
	// Write name
	if err := writeString(w, cd.Name); err != nil {
		return err
	}

	// Write superclass name
	if err := writeString(w, cd.SuperClass); err != nil {
		return err
	}

	// Write fields
	if err := writeStringSlice(w, cd.Fields); err != nil {
		return err
	}

	// Write class variables
	if err := writeStringSlice(w, cd.ClassVariables); err != nil {
		return err
	}

	// Write methods
	if err := writeMethodSlice(w, cd.Methods); err != nil {
		return err
	}

	// Write class methods
	if err := writeMethodSlice(w, cd.ClassMethods); err != nil {
		return err
	}

	return nil
}

// readClassDefinition reads a ClassDefinition from r.
func readClassDefinition(r io.Reader) (*ClassDefinition, error) {
	// Read name
	name, err := readString(r)
	if err != nil {
		return nil, err
	}

	// Read superclass name
	superClass, err := readString(r)
	if err != nil {
		return nil, err
	}

	// Read fields
	fields, err := readStringSlice(r)
	if err != nil {
		return nil, err
	}

	// Read class variables
	classVars, err := readStringSlice(r)
	if err != nil {
		return nil, err
	}

	// Read methods
	methods, err := readMethodSlice(r)
	if err != nil {
		return nil, err
	}

	// Read class methods
	classMethods, err := readMethodSlice(r)
	if err != nil {
		return nil, err
	}

	return &ClassDefinition{
		Name:           name,
		SuperClass:     superClass,
		Fields:         fields,
		ClassVariables: classVars,
		ClassVarValues: make(map[string]interface{}), // Initialize empty map
		Methods:        methods,
		ClassMethods:   classMethods,
	}, nil
}

// writeMethodDefinition writes a MethodDefinition to w.
//
// Format:
//   - Selector (string: 4-byte length + UTF-8)
//   - Parameter count (4 bytes) + parameter names (strings)
//   - Code (Bytecode, recursively encoded)
func writeMethodDefinition(w io.Writer, md *MethodDefinition) error {
	// Write selector
	if err := writeString(w, md.Selector); err != nil {
		return err
	}

	// Write parameters
	if err := writeStringSlice(w, md.Parameters); err != nil {
		return err
	}

	// Write code (bytecode)
	return Encode(md.Code, w)
}

// readMethodDefinition reads a MethodDefinition from r.
func readMethodDefinition(r io.Reader) (*MethodDefinition, error) {
	// Read selector
	selector, err := readString(r)
	if err != nil {
		return nil, err
	}

	// Read parameters
	params, err := readStringSlice(r)
	if err != nil {
		return nil, err
	}

	// Read code (bytecode)
	code, err := Decode(r)
	if err != nil {
		return nil, err
	}

	return &MethodDefinition{
		Selector:   selector,
		Parameters: params,
		Code:       code,
	}, nil
}

// Helper functions for reading/writing strings and slices

func writeString(w io.Writer, s string) error {
	length := uint32(len(s))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func readString(r io.Reader) (string, error) {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func writeStringSlice(w io.Writer, slice []string) error {
	count := uint32(len(slice))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}
	for _, s := range slice {
		if err := writeString(w, s); err != nil {
			return err
		}
	}
	return nil
}

func readStringSlice(r io.Reader) ([]string, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, err
	}
	slice := make([]string, count)
	for i := uint32(0); i < count; i++ {
		s, err := readString(r)
		if err != nil {
			return nil, err
		}
		slice[i] = s
	}
	return slice, nil
}

func writeMethodSlice(w io.Writer, slice []*MethodDefinition) error {
	count := uint32(len(slice))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}
	for _, md := range slice {
		if err := writeMethodDefinition(w, md); err != nil {
			return err
		}
	}
	return nil
}

func readMethodSlice(r io.Reader) ([]*MethodDefinition, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, err
	}
	slice := make([]*MethodDefinition, count)
	for i := uint32(0); i < count; i++ {
		md, err := readMethodDefinition(r)
		if err != nil {
			return nil, err
		}
		slice[i] = md
	}
	return slice, nil
}
