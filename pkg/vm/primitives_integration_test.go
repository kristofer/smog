package vm

import (
	"testing"

	"github.com/kristofer/smog/pkg/bytecode"
)

// TestPrimitivesViaSend tests that primitives work through the send mechanism
func TestPrimitivesViaSend(t *testing.T) {
	vm := &VM{
		globals: make(map[string]interface{}),
	}

	// Test crypto primitives via send
	t.Run("SHA256", func(t *testing.T) {
		result, err := vm.send(nil, "sha256:", []interface{}{"test"})
		if err != nil {
			t.Fatalf("sha256: failed: %v", err)
		}
		hash, ok := result.(string)
		if !ok || len(hash) != 64 {
			t.Errorf("sha256: returned invalid hash: %v", result)
		}
	})

	t.Run("Base64Encode", func(t *testing.T) {
		result, err := vm.send(nil, "base64Encode:", []interface{}{"Hello"})
		if err != nil {
			t.Fatalf("base64Encode: failed: %v", err)
		}
		encoded, ok := result.(string)
		if !ok || len(encoded) == 0 {
			t.Errorf("base64Encode: returned invalid result: %v", result)
		}

		// Decode it back
		decoded, err := vm.send(nil, "base64Decode:", []interface{}{encoded})
		if err != nil {
			t.Fatalf("base64Decode: failed: %v", err)
		}
		if decoded != "Hello" {
			t.Errorf("base64 round-trip failed: got %v, want Hello", decoded)
		}
	})

	t.Run("AESEncryptDecrypt", func(t *testing.T) {
		key := "12345678901234567890123456789012" // 32 bytes
		plaintext := "secret"

		encrypted, err := vm.send(nil, "aesEncrypt:key:", []interface{}{plaintext, key})
		if err != nil {
			t.Fatalf("aesEncrypt:key: failed: %v", err)
		}

		decrypted, err := vm.send(nil, "aesDecrypt:key:", []interface{}{encrypted, key})
		if err != nil {
			t.Fatalf("aesDecrypt:key: failed: %v", err)
		}

		if decrypted != plaintext {
			t.Errorf("AES round-trip failed: got %v, want %v", decrypted, plaintext)
		}
	})

	t.Run("GzipCompress", func(t *testing.T) {
		data := "test data for compression"
		compressed, err := vm.send(nil, "gzipCompress:", []interface{}{data})
		if err != nil {
			t.Fatalf("gzipCompress: failed: %v", err)
		}

		decompressed, err := vm.send(nil, "gzipDecompress:", []interface{}{compressed})
		if err != nil {
			t.Fatalf("gzipDecompress: failed: %v", err)
		}

		if decompressed != data {
			t.Errorf("gzip round-trip failed: got %v, want %v", decompressed, data)
		}
	})

	t.Run("FileIO", func(t *testing.T) {
		path := "/tmp/vm_test_file.txt"
		content := "test file content"

		// Write
		_, err := vm.send(nil, "fileWrite:content:", []interface{}{path, content})
		if err != nil {
			t.Fatalf("fileWrite:content: failed: %v", err)
		}

		// Check exists
		exists, err := vm.send(nil, "fileExists:", []interface{}{path})
		if err != nil {
			t.Fatalf("fileExists: failed: %v", err)
		}
		if exists != true {
			t.Error("fileExists: should return true")
		}

		// Read
		readContent, err := vm.send(nil, "fileRead:", []interface{}{path})
		if err != nil {
			t.Fatalf("fileRead: failed: %v", err)
		}
		if readContent != content {
			t.Errorf("file content mismatch: got %v, want %v", readContent, content)
		}

		// Delete
		_, err = vm.send(nil, "fileDelete:", []interface{}{path})
		if err != nil {
			t.Fatalf("fileDelete: failed: %v", err)
		}
	})

	t.Run("JSON", func(t *testing.T) {
		jsonStr := `{"name":"test","value":42}`
		parsed, err := vm.send(nil, "jsonParse:", []interface{}{jsonStr})
		if err != nil {
			t.Fatalf("jsonParse: failed: %v", err)
		}
		if parsed == nil {
			t.Error("jsonParse: returned nil")
		}

		generated, err := vm.send(nil, "jsonGenerate:", []interface{}{parsed})
		if err != nil {
			t.Fatalf("jsonGenerate: failed: %v", err)
		}
		if _, ok := generated.(string); !ok {
			t.Errorf("jsonGenerate: should return string, got %T", generated)
		}
	})

	t.Run("Regex", func(t *testing.T) {
		matched, err := vm.send(nil, "regexMatch:text:", []interface{}{`\d+`, "abc123"})
		if err != nil {
			t.Fatalf("regexMatch:text: failed: %v", err)
		}
		if matched != true {
			t.Error("regexMatch:text: should return true")
		}

		matches, err := vm.send(nil, "regexFindAll:text:", []interface{}{`\d+`, "1 2 3"})
		if err != nil {
			t.Fatalf("regexFindAll:text: failed: %v", err)
		}
		array, ok := matches.(*Array)
		if !ok {
			t.Errorf("regexFindAll:text: should return Array, got %T", matches)
		} else if len(array.Elements) != 3 {
			t.Errorf("regexFindAll:text: should find 3 matches, got %d", len(array.Elements))
		}

		replaced, err := vm.send(nil, "regexReplace:text:with:", []interface{}{`\d+`, "Price: 100", "XXX"})
		if err != nil {
			t.Fatalf("regexReplace:text:with: failed: %v", err)
		}
		if replaced != "Price: XXX" {
			t.Errorf("regexReplace:text:with: got %v, want 'Price: XXX'", replaced)
		}
	})

	t.Run("Random", func(t *testing.T) {
		num, err := vm.send(nil, "randomInt:max:", []interface{}{int64(1), int64(10)})
		if err != nil {
			t.Fatalf("randomInt:max: failed: %v", err)
		}
		n, ok := num.(int64)
		if !ok || n < 1 || n > 10 {
			t.Errorf("randomInt:max: should return int64 between 1 and 10, got %v", num)
		}

		f, err := vm.send(nil, "randomFloat", []interface{}{})
		if err != nil {
			t.Fatalf("randomFloat failed: %v", err)
		}
		flt, ok := f.(float64)
		if !ok || flt < 0 || flt > 1 {
			t.Errorf("randomFloat should return float64 between 0 and 1, got %v", f)
		}

		bytes, err := vm.send(nil, "randomBytes:", []interface{}{int64(16)})
		if err != nil {
			t.Fatalf("randomBytes: failed: %v", err)
		}
		if _, ok := bytes.(string); !ok {
			t.Errorf("randomBytes: should return string, got %T", bytes)
		}
	})

	t.Run("DateTime", func(t *testing.T) {
		now, err := vm.send(nil, "dateNow", []interface{}{})
		if err != nil {
			t.Fatalf("dateNow failed: %v", err)
		}
		timestamp, ok := now.(int64)
		if !ok || timestamp <= 0 {
			t.Errorf("dateNow should return positive int64, got %v", now)
		}

		formatted, err := vm.send(nil, "dateFormat:format:", []interface{}{timestamp, "iso8601"})
		if err != nil {
			t.Fatalf("dateFormat:format: failed: %v", err)
		}
		if _, ok := formatted.(string); !ok {
			t.Errorf("dateFormat:format: should return string, got %T", formatted)
		}

		year, err := vm.send(nil, "timeYear:", []interface{}{timestamp})
		if err != nil {
			t.Fatalf("timeYear: failed: %v", err)
		}
		y, ok := year.(int64)
		if !ok || y < 2000 || y > 2100 {
			t.Errorf("timeYear: should return reasonable year, got %v", year)
		}
	})
}

// TestPrimitivesInBytecode tests that primitives work in compiled bytecode
func TestPrimitivesInBytecode(t *testing.T) {
	// Test that we can execute a simple hash operation
	bc := &bytecode.Bytecode{
		Instructions: []bytecode.Instruction{
			{Op: bytecode.OpPushNil, Operand: 0}, // Push nil as receiver
			{Op: bytecode.OpPush, Operand: 0}, // Push constant 0 (data string)
			{Op: bytecode.OpSend, Operand: (1 << 8) | 1}, // Send sha256: with 1 arg
			{Op: bytecode.OpReturn, Operand: 0},
		},
		Constants: []interface{}{
			"test data",
			"sha256:",
		},
	}

	vm := &VM{
		stack:   make([]interface{}, 1024),
		sp:      0,
		locals:  make([]interface{}, 256),
		globals: make(map[string]interface{}),
		classes: make(map[string]*bytecode.ClassDefinition),
	}
	
	err := vm.Run(bc)
	if err != nil {
		t.Fatalf("Failed to run bytecode with primitive: %v", err)
	}

	result := vm.StackTop()
	hash, ok := result.(string)
	if !ok || len(hash) != 64 {
		t.Errorf("Bytecode primitive execution returned invalid hash: %v", result)
	}
}
