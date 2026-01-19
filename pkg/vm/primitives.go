// Package vm implements stdlib primitives for the virtual machine.
//
// This file contains VM primitive implementations for standard library
// functionality including HTTP, crypto, compression, file I/O, JSON, regex,
// date/time, and random number generation.
package vm

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// HTTP Primitives

// httpGet performs an HTTP GET request
func (vm *VM) httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP GET failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

// httpPost performs an HTTP POST request
func (vm *VM) httpPost(url string, body string) (string, error) {
	resp, err := http.Post(url, "text/plain", strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("HTTP POST failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(respBody), nil
}

// Crypto Primitives

// aesEncrypt encrypts data using AES-256
func (vm *VM) aesEncrypt(data string, key string) (string, error) {
	// Ensure key is 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("AES key must be 32 bytes, got %d", len(keyBytes))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %v", err)
	}

	// Encrypt the data
	plaintext := []byte(data)
	// Pad to block size
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	paddedData := make([]byte, len(plaintext)+padding)
	copy(paddedData, plaintext)
	for i := len(plaintext); i < len(paddedData); i++ {
		paddedData[i] = byte(padding)
	}

	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// Prepend IV to ciphertext
	result := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// aesDecrypt decrypts AES-256 encrypted data
func (vm *VM) aesDecrypt(data string, key string) (string, error) {
	// Ensure key is 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("AES key must be 32 bytes, got %d", len(keyBytes))
	}

	// Decode base64
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Extract IV from the beginning
	iv := encrypted[:aes.BlockSize]
	ciphertext := encrypted[aes.BlockSize:]

	// Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	padding := int(plaintext[len(plaintext)-1])
	if padding > len(plaintext) || padding > aes.BlockSize {
		return "", fmt.Errorf("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return string(plaintext), nil
}

// aesGenerateKey generates a random 32-byte key for AES-256
func (vm *VM) aesGenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate key: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// sha256Hash computes SHA-256 hash
func (vm *VM) sha256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// sha512Hash computes SHA-512 hash
func (vm *VM) sha512Hash(data string) string {
	hash := sha512.Sum512([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// md5Hash computes MD5 hash (deprecated but included for compatibility)
func (vm *VM) md5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// base64Encode encodes data to base64
func (vm *VM) base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// base64Decode decodes base64 data
func (vm *VM) base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}
	return string(decoded), nil
}

// Compression Primitives

// zipCompress compresses data using ZIP
func (vm *VM) zipCompress(data string) (string, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	f, err := w.Create("data")
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %v", err)
	}

	if _, err := f.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("failed to write to zip: %v", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close zip: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// zipDecompress decompresses ZIP data
func (vm *VM) zipDecompress(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %v", err)
	}

	if len(r.File) == 0 {
		return "", fmt.Errorf("zip archive is empty")
	}

	f, err := r.File[0].Open()
	if err != nil {
		return "", fmt.Errorf("failed to open zip entry: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read zip entry: %v", err)
	}

	return string(content), nil
}

// gzipCompress compresses data using GZIP
func (vm *VM) gzipCompress(data string) (string, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	if _, err := w.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("failed to write to gzip: %v", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// gzipDecompress decompresses GZIP data
func (vm *VM) gzipDecompress(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	r, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return "", fmt.Errorf("failed to open gzip: %v", err)
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read gzip: %v", err)
	}

	return string(content), nil
}

// File I/O Primitives

// fileRead reads entire file contents
func (vm *VM) fileRead(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(content), nil
}

// fileWrite writes content to a file
func (vm *VM) fileWrite(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

// fileExists checks if a file exists
func (vm *VM) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// fileDelete deletes a file
func (vm *VM) fileDelete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// JSON Primitives

// jsonParse parses JSON string to a value
func (vm *VM) jsonParse(data string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	return vm.convertJSONValue(result), nil
}

// jsonGenerate generates JSON string from a value
func (vm *VM) jsonGenerate(value interface{}) (string, error) {
	data, err := json.Marshal(vm.convertToJSONValue(value))
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %v", err)
	}
	return string(data), nil
}

// convertJSONValue converts JSON value to VM types
func (vm *VM) convertJSONValue(value interface{}) interface{} {
	switch v := value.(type) {
	case float64:
		// JSON numbers are float64, convert to int64 if whole number
		if v == float64(int64(v)) {
			return int64(v)
		}
		return v
	case []interface{}:
		// Convert to Array
		elements := make([]interface{}, len(v))
		for i, elem := range v {
			elements[i] = vm.convertJSONValue(elem)
		}
		return &Array{Elements: elements}
	case map[string]interface{}:
		// Keep as map for now (Dictionary type not yet implemented)
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = vm.convertJSONValue(val)
		}
		return result
	default:
		return v
	}
}

// convertToJSONValue converts VM types to JSON-compatible values
func (vm *VM) convertToJSONValue(value interface{}) interface{} {
	switch v := value.(type) {
	case *Array:
		result := make([]interface{}, len(v.Elements))
		for i, elem := range v.Elements {
			result[i] = vm.convertToJSONValue(elem)
		}
		return result
	case map[string]interface{}:
		// Handle map (used when Dictionary type not yet implemented)
		result := make(map[string]interface{})
		for k, val := range v {
			result[k] = vm.convertToJSONValue(val)
		}
		return result
	default:
		return v
	}
}

// Regular Expression Primitives

// regexMatch checks if pattern matches string
func (vm *VM) regexMatch(pattern string, text string) (bool, error) {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %v", err)
	}
	return matched, nil
}

// regexFindAll finds all matches of pattern in text
func (vm *VM) regexFindAll(pattern string, text string) (interface{}, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	matches := re.FindAllString(text, -1)
	elements := make([]interface{}, len(matches))
	for i, m := range matches {
		elements[i] = m
	}
	return &Array{Elements: elements}, nil
}

// regexReplace replaces all matches of pattern in text
func (vm *VM) regexReplace(pattern string, text string, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %v", err)
	}
	return re.ReplaceAllString(text, replacement), nil
}

// Random Number Generation Primitives

// randomInt generates a random integer between min and max (inclusive)
func (vm *VM) randomInt(min int64, max int64) (int64, error) {
	if min > max {
		return 0, fmt.Errorf("min must be <= max")
	}
	diff := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %v", err)
	}
	return n.Int64() + min, nil
}

// randomFloat generates a cryptographically secure random float in [0, 1) using crypto/rand
func (vm *VM) randomFloat() (float64, error) {
	// Generate 8 random bytes
	bytes := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return 0, fmt.Errorf("failed to generate random float: %v", err)
	}
	
	// Convert to uint64 and then to float in range [0, 1)
	// Use the high 53 bits for the mantissa
	n := uint64(bytes[0])<<56 | uint64(bytes[1])<<48 | uint64(bytes[2])<<40 | uint64(bytes[3])<<32 |
		uint64(bytes[4])<<24 | uint64(bytes[5])<<16 | uint64(bytes[6])<<8 | uint64(bytes[7])
	
	// Mask to 53 bits and convert to float in [0, 1)
	return float64(n>>11) / float64(1<<53), nil
}

// randomBytes generates random bytes
func (vm *VM) randomBytes(length int64) (string, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Date and Time Primitives

// dateNow returns the current Unix timestamp
func (vm *VM) dateNow() int64 {
	return time.Now().Unix()
}

// dateFormat formats a Unix timestamp to a string
func (vm *VM) dateFormat(timestamp int64, format string) string {
	t := time.Unix(timestamp, 0)
	// Support common format patterns
	switch format {
	case "iso8601", "ISO8601":
		return t.Format(time.RFC3339)
	case "rfc3339", "RFC3339":
		return t.Format(time.RFC3339)
	case "date":
		return t.Format("2006-01-02")
	case "time":
		return t.Format("15:04:05")
	case "datetime":
		return t.Format("2006-01-02 15:04:05")
	default:
		// Use Go's time format layout
		return t.Format(format)
	}
}

// dateParse parses a date string to Unix timestamp
func (vm *VM) dateParse(dateStr string, format string) (int64, error) {
	var t time.Time
	var err error

	switch format {
	case "iso8601", "ISO8601", "rfc3339", "RFC3339":
		t, err = time.Parse(time.RFC3339, dateStr)
	case "date":
		t, err = time.Parse("2006-01-02", dateStr)
	case "time":
		t, err = time.Parse("15:04:05", dateStr)
	case "datetime":
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
	default:
		t, err = time.Parse(format, dateStr)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %v", err)
	}
	return t.Unix(), nil
}

// timeYear extracts year from Unix timestamp
func (vm *VM) timeYear(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Year())
}

// timeMonth extracts month from Unix timestamp
func (vm *VM) timeMonth(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Month())
}

// timeDay extracts day from Unix timestamp
func (vm *VM) timeDay(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Day())
}

// timeHour extracts hour from Unix timestamp
func (vm *VM) timeHour(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Hour())
}

// timeMinute extracts minute from Unix timestamp
func (vm *VM) timeMinute(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Minute())
}

// timeSecond extracts second from Unix timestamp
func (vm *VM) timeSecond(timestamp int64) int64 {
	return int64(time.Unix(timestamp, 0).Second())
}
