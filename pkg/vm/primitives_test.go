package vm

import (
	"strings"
	"testing"
)

// TestCryptoPrimitives tests the crypto primitives
func TestCryptoPrimitives(t *testing.T) {
	vm := &VM{}

	// Test SHA-256
	hash := vm.sha256Hash("test")
	if len(hash) != 64 { // SHA-256 produces 64 hex characters
		t.Errorf("SHA-256 hash length incorrect: got %d, want 64", len(hash))
	}

	// Test SHA-512
	hash512 := vm.sha512Hash("test")
	if len(hash512) != 128 { // SHA-512 produces 128 hex characters
		t.Errorf("SHA-512 hash length incorrect: got %d, want 128", len(hash512))
	}

	// Test MD5
	md5Hash := vm.md5Hash("test")
	if len(md5Hash) != 32 { // MD5 produces 32 hex characters
		t.Errorf("MD5 hash length incorrect: got %d, want 32", len(md5Hash))
	}

	// Test Base64 encode/decode
	original := "Hello, World!"
	encoded := vm.base64Encode(original)
	decoded, err := vm.base64Decode(encoded)
	if err != nil {
		t.Fatalf("Base64 decode failed: %v", err)
	}
	if decoded != original {
		t.Errorf("Base64 decode mismatch: got %q, want %q", decoded, original)
	}

	// Test AES encryption/decryption
	key := "12345678901234567890123456789012" // 32 bytes for AES-256
	plaintext := "Secret message"
	encrypted, err := vm.aesEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("AES encrypt failed: %v", err)
	}
	decrypted, err := vm.aesDecrypt(encrypted, key)
	if err != nil {
		t.Fatalf("AES decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("AES decrypt mismatch: got %q, want %q", decrypted, plaintext)
	}

	// Test AES key generation
	generatedKey, err := vm.aesGenerateKey()
	if err != nil {
		t.Fatalf("AES key generation failed: %v", err)
	}
	if len(generatedKey) == 0 {
		t.Error("Generated AES key is empty")
	}
}

// TestCompressionPrimitives tests the compression primitives
func TestCompressionPrimitives(t *testing.T) {
	vm := &VM{}

	// Test GZIP compression/decompression
	original := "This is a test string for compression. It should compress well when repeated. " + strings.Repeat("test ", 100)
	compressed, err := vm.gzipCompress(original)
	if err != nil {
		t.Fatalf("GZIP compress failed: %v", err)
	}
	if len(compressed) == 0 {
		t.Error("GZIP compressed data is empty")
	}

	decompressed, err := vm.gzipDecompress(compressed)
	if err != nil {
		t.Fatalf("GZIP decompress failed: %v", err)
	}
	if decompressed != original {
		t.Errorf("GZIP decompress mismatch: length got %d, want %d", len(decompressed), len(original))
	}

	// Test ZIP compression/decompression
	zipCompressed, err := vm.zipCompress(original)
	if err != nil {
		t.Fatalf("ZIP compress failed: %v", err)
	}
	if len(zipCompressed) == 0 {
		t.Error("ZIP compressed data is empty")
	}

	zipDecompressed, err := vm.zipDecompress(zipCompressed)
	if err != nil {
		t.Fatalf("ZIP decompress failed: %v", err)
	}
	if zipDecompressed != original {
		t.Errorf("ZIP decompress mismatch: length got %d, want %d", len(zipDecompressed), len(original))
	}
}

// TestJSONPrimitives tests the JSON primitives
func TestJSONPrimitives(t *testing.T) {
	vm := &VM{}

	// Test simple JSON parsing
	jsonStr := `{"name":"John","age":30,"active":true}`
	parsed, err := vm.jsonParse(jsonStr)
	if err != nil {
		t.Fatalf("JSON parse failed: %v", err)
	}
	if parsed == nil {
		t.Error("JSON parse returned nil")
	}

	// Test JSON generation from simple values
	value := map[string]interface{}{
		"name":   "Alice",
		"age":    int64(25),
		"active": true,
	}
	generated, err := vm.jsonGenerate(value)
	if err != nil {
		t.Fatalf("JSON generate failed: %v", err)
	}
	if !strings.Contains(generated, "Alice") {
		t.Errorf("JSON generate missing expected content: %s", generated)
	}
	if !strings.Contains(generated, "25") {
		t.Errorf("JSON generate missing expected content: %s", generated)
	}

	// Test round-trip
	parsed2, err := vm.jsonParse(generated)
	if err != nil {
		t.Fatalf("JSON parse round-trip failed: %v", err)
	}
	if parsed2 == nil {
		t.Error("JSON parse round-trip returned nil")
	}
}

// TestRegexPrimitives tests the regex primitives
func TestRegexPrimitives(t *testing.T) {
	vm := &VM{}

	// Test regex match
	matched, err := vm.regexMatch(`\d+`, "abc123def")
	if err != nil {
		t.Fatalf("Regex match failed: %v", err)
	}
	if !matched {
		t.Error("Regex match should have matched")
	}

	noMatch, err := vm.regexMatch(`\d+`, "abcdef")
	if err != nil {
		t.Fatalf("Regex match failed: %v", err)
	}
	if noMatch {
		t.Error("Regex match should not have matched")
	}

	// Test regex find all
	text := "The numbers are 123, 456, and 789"
	result, err := vm.regexFindAll(`\d+`, text)
	if err != nil {
		t.Fatalf("Regex find all failed: %v", err)
	}
	array, ok := result.(*Array)
	if !ok {
		t.Fatal("Regex find all should return an Array")
	}
	if len(array.Elements) != 3 {
		t.Errorf("Regex find all found %d matches, want 3", len(array.Elements))
	}

	// Test regex replace
	replaced, err := vm.regexReplace(`\d+`, "Price: 100 dollars", "XXX")
	if err != nil {
		t.Fatalf("Regex replace failed: %v", err)
	}
	if replaced != "Price: XXX dollars" {
		t.Errorf("Regex replace incorrect: got %q, want %q", replaced, "Price: XXX dollars")
	}
}

// TestRandomPrimitives tests the random number primitives
func TestRandomPrimitives(t *testing.T) {
	vm := &VM{}

	// Test random int
	for i := 0; i < 10; i++ {
		num, err := vm.randomInt(1, 10)
		if err != nil {
			t.Fatalf("Random int failed: %v", err)
		}
		if num < 1 || num > 10 {
			t.Errorf("Random int out of range: got %d, want 1-10", num)
		}
	}

	// Test random float
	for i := 0; i < 10; i++ {
		f, err := vm.randomFloat()
		if err != nil {
			t.Fatalf("Random float failed: %v", err)
		}
		if f < 0.0 || f > 1.0 {
			t.Errorf("Random float out of range: got %f, want 0.0-1.0", f)
		}
	}

	// Test random bytes
	bytes, err := vm.randomBytes(16)
	if err != nil {
		t.Fatalf("Random bytes failed: %v", err)
	}
	if len(bytes) == 0 {
		t.Error("Random bytes returned empty string")
	}
}

// TestDateTimePrimitives tests the date/time primitives
func TestDateTimePrimitives(t *testing.T) {
	vm := &VM{}

	// Test dateNow
	now := vm.dateNow()
	if now <= 0 {
		t.Errorf("dateNow returned invalid timestamp: %d", now)
	}

	// Test date formatting
	timestamp := int64(1640000000) // A known timestamp
	formatted := vm.dateFormat(timestamp, "iso8601")
	if len(formatted) == 0 {
		t.Error("dateFormat returned empty string")
	}

	// Test date parsing
	parsed, err := vm.dateParse("2021-12-20T12:00:00Z", "iso8601")
	if err != nil {
		t.Fatalf("dateParse failed: %v", err)
	}
	if parsed <= 0 {
		t.Errorf("dateParse returned invalid timestamp: %d", parsed)
	}

	// Test time component extraction
	year := vm.timeYear(timestamp)
	if year < 2000 || year > 2100 {
		t.Errorf("timeYear returned unexpected year: %d", year)
	}

	month := vm.timeMonth(timestamp)
	if month < 1 || month > 12 {
		t.Errorf("timeMonth returned invalid month: %d", month)
	}

	day := vm.timeDay(timestamp)
	if day < 1 || day > 31 {
		t.Errorf("timeDay returned invalid day: %d", day)
	}

	hour := vm.timeHour(timestamp)
	if hour < 0 || hour > 23 {
		t.Errorf("timeHour returned invalid hour: %d", hour)
	}

	minute := vm.timeMinute(timestamp)
	if minute < 0 || minute > 59 {
		t.Errorf("timeMinute returned invalid minute: %d", minute)
	}

	second := vm.timeSecond(timestamp)
	if second < 0 || second > 59 {
		t.Errorf("timeSecond returned invalid second: %d", second)
	}
}

// TestFileIOPrimitives tests the file I/O primitives
func TestFileIOPrimitives(t *testing.T) {
	vm := &VM{}

	// Create a temporary file for testing
	testPath := "/tmp/smog_test_file.txt"
	testContent := "This is test content"

	// Test file write
	err := vm.fileWrite(testPath, testContent)
	if err != nil {
		t.Fatalf("fileWrite failed: %v", err)
	}

	// Test file exists
	exists := vm.fileExists(testPath)
	if !exists {
		t.Error("fileExists returned false for existing file")
	}

	// Test file read
	readContent, err := vm.fileRead(testPath)
	if err != nil {
		t.Fatalf("fileRead failed: %v", err)
	}
	if readContent != testContent {
		t.Errorf("fileRead mismatch: got %q, want %q", readContent, testContent)
	}

	// Test file delete
	err = vm.fileDelete(testPath)
	if err != nil {
		t.Fatalf("fileDelete failed: %v", err)
	}

	// Verify file was deleted
	exists = vm.fileExists(testPath)
	if exists {
		t.Error("fileExists returned true for deleted file")
	}
}
