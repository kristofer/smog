// Package test provides integration tests for smog standard library.
package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runSmogFile executes a smog file and returns stdout
func runSmogFile(t *testing.T, relpath string) string {
	// Get the project root (parent of test directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	root := filepath.Dir(wd)
	
	// Build path to smog file and cmd
	smogFile := filepath.Join(root, relpath)
	cmdPath := filepath.Join(root, "cmd", "smog")
	
	cmd := exec.Command("go", "run", cmdPath, smogFile)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run %s: %v\nOutput: %s", relpath, err, string(output))
	}
	return string(output)
}

// TestStandardLibrary_Set tests the Set collection class
func TestStandardLibrary_Set(t *testing.T) {
	output := runSmogFile(t, "examples/stdlib/set_example.smog")
	
	// Check that example ran successfully
	if !strings.Contains(output, "Creating Sets...") {
		t.Error("Set example did not run successfully")
	}
	
	// Check set size is correct (3 unique fruits)
	if !strings.Contains(output, "Fruits size: 3") {
		t.Error("Set did not correctly track unique elements")
	}
	
	// Check includes method works
	if !strings.Contains(output, "fruits includes apple: true") {
		t.Error("Set includes method failed for existing element")
	}
	
	if !strings.Contains(output, "fruits includes carrot: false") {
		t.Error("Set includes method failed for non-existing element")
	}
	
	// Check union operation
	if !strings.Contains(output, "Union size: 6") {
		t.Error("Set union operation did not produce correct result")
	}
	
	// Check intersection operation
	if !strings.Contains(output, "Intersection size: 1") {
		t.Error("Set intersection operation did not produce correct result")
	}
}

// TestStandardLibrary_Math tests the Math utility class
func TestStandardLibrary_Math(t *testing.T) {
	output := runSmogFile(t, "examples/stdlib/math_example.smog")
	
	// Check constants
	if !strings.Contains(output, "Pi: 3.14159265359") {
		t.Error("Math.pi constant incorrect")
	}
	
	if !strings.Contains(output, "e (Euler): 2.71828182846") {
		t.Error("Math.e constant incorrect")
	}
	
	// Check abs
	if !strings.Contains(output, "abs(-42): 42") {
		t.Error("Math.abs failed")
	}
	
	// Check max/min
	if !strings.Contains(output, "max(15, 23): 23") {
		t.Error("Math.max failed")
	}
	
	if !strings.Contains(output, "min(15, 23): 15") {
		t.Error("Math.min failed")
	}
	
	// Check power
	if !strings.Contains(output, "2^8: 256") {
		t.Error("Math.power failed for 2^8")
	}
	
	if !strings.Contains(output, "3^4: 81") {
		t.Error("Math.power failed for 3^4")
	}
	
	// Check sqrt
	if !strings.Contains(output, "sqrt(16): 4") {
		t.Error("Math.sqrt failed for 16")
	}
	
	if !strings.Contains(output, "sqrt(25): 5") {
		t.Error("Math.sqrt failed for 25")
	}
	
	if !strings.Contains(output, "sqrt(100): 10") {
		t.Error("Math.sqrt failed for 100")
	}
	
	// Check factorial
	if !strings.Contains(output, "5!: 120") {
		t.Error("Math.factorial failed for 5")
	}
	
	if !strings.Contains(output, "7!: 5040") {
		t.Error("Math.factorial failed for 7")
	}
	
	// Check fibonacci
	if !strings.Contains(output, "fib(10) = 55") {
		t.Error("Math.fibonacci failed for 10")
	}
	
	// Check gcd
	if !strings.Contains(output, "gcd(48, 18): 6") {
		t.Error("Math.gcd failed for (48, 18)")
	}
}

// TestStandardLibrary_OrderedCollection tests the OrderedCollection class
func TestStandardLibrary_OrderedCollection(t *testing.T) {
	output := runSmogFile(t, "examples/stdlib/ordered_collection_example.smog")
	
	// Check basic operations
	if !strings.Contains(output, "Size: 10") {
		t.Error("OrderedCollection size incorrect")
	}
	
	if !strings.Contains(output, "First: 1") {
		t.Error("OrderedCollection.first failed")
	}
	
	if !strings.Contains(output, "Last: 10") {
		t.Error("OrderedCollection.last failed")
	}
	
	// Check collect operation (doubling)
	if !strings.Contains(output, "=== Collect: Double each number ===") {
		t.Error("OrderedCollection.collect header missing")
	}
	
	// Check select operations
	if !strings.Contains(output, "=== Select: Even numbers ===") {
		t.Error("OrderedCollection.select for evens header missing")
	}
	
	if !strings.Contains(output, "=== Select: Odd numbers ===") {
		t.Error("OrderedCollection.select for odds header missing")
	}
	
	// Check detect operation
	if !strings.Contains(output, "Found: 6") {
		t.Error("OrderedCollection.detect failed")
	}
	
	// Check anySatisfy
	if !strings.Contains(output, "Any number > 8? true") {
		t.Error("OrderedCollection.anySatisfy failed")
	}
	
	// Check allSatisfy
	if !strings.Contains(output, "All numbers > 0? true") {
		t.Error("OrderedCollection.allSatisfy with true condition failed")
	}
	
	if !strings.Contains(output, "All numbers > 8? false") {
		t.Error("OrderedCollection.allSatisfy with false condition failed")
	}
}

// TestStandardLibrary_Comprehensive tests multiple stdlib features together
func TestStandardLibrary_Comprehensive(t *testing.T) {
	output := runSmogFile(t, "examples/stdlib/comprehensive_example.smog")
	
	// Check statistics
	if !strings.Contains(output, "Sum: 39") {
		t.Error("Comprehensive example: sum calculation failed")
	}
	
	if !strings.Contains(output, "Mean: 4") {
		t.Error("Comprehensive example: mean calculation failed")
	}
	
	if !strings.Contains(output, "Maximum: 12") {
		t.Error("Comprehensive example: max calculation failed")
	}
	
	if !strings.Contains(output, "Minimum: -3") {
		t.Error("Comprehensive example: min calculation failed")
	}
	
	// Check filtering
	if !strings.Contains(output, "Even numbers:") {
		t.Error("Comprehensive example: even numbers filter missing")
	}
	
	if !strings.Contains(output, "Odd numbers:") {
		t.Error("Comprehensive example: odd numbers filter missing")
	}
	
	if !strings.Contains(output, "Positive numbers:") {
		t.Error("Comprehensive example: positive numbers filter missing")
	}
	
	// Check transformations
	if !strings.Contains(output, "Squares:") {
		t.Error("Comprehensive example: squares transformation missing")
	}
	
	// Check unique values (Set)
	if !strings.Contains(output, "Count of unique values: 5") {
		t.Error("Comprehensive example: unique values count incorrect")
	}
}
