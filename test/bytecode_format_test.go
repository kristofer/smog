package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kristofer/smog/pkg/bytecode"
	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

// TestBytecodeFileRoundTrip tests the complete workflow:
// source → compile → save to .sg → load .sg → execute
func TestBytecodeFileRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "simple integer",
			source: "42.",
		},
		{
			name:   "string literal",
			source: "'Hello, World!'.",
		},
		{
			name:   "arithmetic",
			source: "3 + 4.",
		},
		{
			name:   "variable assignment",
			source: "| x | x := 10. x + 5.",
		},
		{
			name:   "message send",
			source: "'test' println.",
		},
		{
			name:   "block creation",
			source: "[ 42 ] value.",
		},
		{
			name:   "array literal",
			source: "#(1 2 3).",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse source
			p := parser.New(tt.source)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			// Compile to bytecode
			c := compiler.New()
			bc, err := c.Compile(program)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			// Save to temporary file
			tmpFile := filepath.Join(t.TempDir(), "test.sg")
			file, err := os.Create(tmpFile)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			if err := bytecode.Encode(bc, file); err != nil {
				file.Close()
				t.Fatalf("Encode failed: %v", err)
			}
			file.Close()

			// Load from file
			file, err = os.Open(tmpFile)
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			loadedBC, err := bytecode.Decode(file)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Verify instruction count matches
			if len(loadedBC.Instructions) != len(bc.Instructions) {
				t.Errorf("Instruction count mismatch: got %d, want %d",
					len(loadedBC.Instructions), len(bc.Instructions))
			}

			// Verify constant count matches
			if len(loadedBC.Constants) != len(bc.Constants) {
				t.Errorf("Constant count mismatch: got %d, want %d",
					len(loadedBC.Constants), len(bc.Constants))
			}

			// Both bytecodes should execute successfully
			// (We don't verify output here as that's tested elsewhere)
			v := vm.New()
			if err := v.Run(loadedBC); err != nil {
				t.Errorf("Failed to execute loaded bytecode: %v", err)
			}
		})
	}
}

// TestMultipleCompilations tests that multiple .smog files can be
// compiled to .sg files independently.
func TestMultipleCompilations(t *testing.T) {
	sources := map[string]string{
		"module1.smog": "| x | x := 10. x.",
		"module2.smog": "| y | y := 20. y.",
		"module3.smog": "'Hello'.",
	}

	tmpDir := t.TempDir()

	// Compile all sources
	for filename, source := range sources {
		// Parse
		p := parser.New(source)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse %s failed: %v", filename, err)
		}

		// Compile
		c := compiler.New()
		bc, err := c.Compile(program)
		if err != nil {
			t.Fatalf("Compile %s failed: %v", filename, err)
		}

		// Save
		sgFile := strings.Replace(filename, ".smog", ".sg", 1)
		path := filepath.Join(tmpDir, sgFile)
		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("Create %s failed: %v", path, err)
		}

		if err := bytecode.Encode(bc, file); err != nil {
			file.Close()
			t.Fatalf("Encode %s failed: %v", filename, err)
		}
		file.Close()
	}

	// Verify all .sg files exist and can be loaded
	for filename := range sources {
		sgFile := strings.Replace(filename, ".smog", ".sg", 1)
		path := filepath.Join(tmpDir, sgFile)

		file, err := os.Open(path)
		if err != nil {
			t.Fatalf("Open %s failed: %v", path, err)
		}

		_, err = bytecode.Decode(file)
		file.Close()
		if err != nil {
			t.Fatalf("Decode %s failed: %v", path, err)
		}
	}
}

// TestClassDefinitionSerialization tests that class definitions
// serialize and deserialize correctly.
func TestClassDefinitionSerialization(t *testing.T) {
	source := `
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
`

	// Parse and compile
	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Save to file
	tmpFile := filepath.Join(t.TempDir(), "class.sg")
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if err := bytecode.Encode(bc, file); err != nil {
		file.Close()
		t.Fatalf("Encode failed: %v", err)
	}
	file.Close()

	// Load from file
	file, err = os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	loadedBC, err := bytecode.Decode(file)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify we have a class definition in constants
	found := false
	for _, c := range loadedBC.Constants {
		if classDef, ok := c.(*bytecode.ClassDefinition); ok {
			found = true
			if classDef.Name != "Counter" {
				t.Errorf("Class name mismatch: got %s, want Counter", classDef.Name)
			}
			if classDef.SuperClass != "Object" {
				t.Errorf("Superclass mismatch: got %s, want Object", classDef.SuperClass)
			}
			if len(classDef.Fields) != 1 || classDef.Fields[0] != "count" {
				t.Errorf("Fields mismatch: got %v, want [count]", classDef.Fields)
			}
			if len(classDef.Methods) != 3 {
				t.Errorf("Method count mismatch: got %d, want 3", len(classDef.Methods))
			}
		}
	}

	if !found {
		t.Error("No ClassDefinition found in loaded bytecode")
	}
}

// TestNestedBlocksSerialization tests that blocks (closures)
// serialize correctly.
func TestNestedBlocksSerialization(t *testing.T) {
	source := `| block | block := [ :x | x + 1 ].`

	// Parse and compile
	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Save to file
	tmpFile := filepath.Join(t.TempDir(), "nested.sg")
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if err := bytecode.Encode(bc, file); err != nil {
		file.Close()
		t.Fatalf("Encode failed: %v", err)
	}
	file.Close()

	// Load from file
	file, err = os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	loadedBC, err := bytecode.Decode(file)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify nested bytecode exists
	found := false
	for _, c := range loadedBC.Constants {
		if _, ok := c.(*bytecode.Bytecode); ok {
			found = true
			// Found block bytecode - good!
			break
		}
	}

	if !found {
		t.Error("No block bytecode found in loaded bytecode")
	}
}

// TestLargeProgram tests that larger programs with many instructions
// and constants serialize correctly.
func TestLargeProgram(t *testing.T) {
	// Generate a program with many operations
	var sb strings.Builder
	sb.WriteString("| x y z |\n")
	for i := 0; i < 100; i++ {
		sb.WriteString("x := ")
		sb.WriteString(string(rune('0' + (i % 10))))
		sb.WriteString(".\n")
	}
	source := sb.String()

	// Parse and compile
	p := parser.New(source)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Save to file
	tmpFile := filepath.Join(t.TempDir(), "large.sg")
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if err := bytecode.Encode(bc, file); err != nil {
		file.Close()
		t.Fatalf("Encode failed: %v", err)
	}
	file.Close()

	// Load from file
	file, err = os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	loadedBC, err := bytecode.Decode(file)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify instruction count
	if len(loadedBC.Instructions) != len(bc.Instructions) {
		t.Errorf("Instruction count mismatch: got %d, want %d",
			len(loadedBC.Instructions), len(bc.Instructions))
	}

	// Execute to verify correctness
	v := vm.New()
	if err := v.Run(loadedBC); err != nil {
		t.Errorf("Failed to execute loaded bytecode: %v", err)
	}
}

// TestFileCorruption tests that corrupted .sg files are detected.
func TestFileCorruption(t *testing.T) {
	// Create a valid bytecode file
	bc := &bytecode.Bytecode{
		Instructions: []bytecode.Instruction{
			{Op: bytecode.OpPush, Operand: 0},
			{Op: bytecode.OpReturn, Operand: 0},
		},
		Constants: []interface{}{int64(42)},
	}

	tmpFile := filepath.Join(t.TempDir(), "corrupt.sg")
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if err := bytecode.Encode(bc, file); err != nil {
		file.Close()
		t.Fatalf("Encode failed: %v", err)
	}
	file.Close()

	// Corrupt the file by truncating it
	if err := os.Truncate(tmpFile, 10); err != nil {
		t.Fatalf("Failed to truncate file: %v", err)
	}

	// Try to load corrupted file
	file, err = os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	_, err = bytecode.Decode(file)
	if err == nil {
		t.Error("Expected error when loading corrupted file, got nil")
	}
}

// TestEmptyFile tests handling of empty .sg files.
func TestEmptyFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.sg")
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	file.Close()

	// Try to load empty file
	file, err = os.Open(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	_, err = bytecode.Decode(file)
	if err == nil {
		t.Error("Expected error when loading empty file, got nil")
	}
}
