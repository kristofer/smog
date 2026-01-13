package vm

import (
"testing"

"github.com/kristofer/smog/pkg/compiler"
"github.com/kristofer/smog/pkg/parser"
)

func TestVMIntegerLiteral(t *testing.T) {
input := "42"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(42) {
t.Errorf("Expected 42, got %v", result)
}
}

func TestVMStringLiteral(t *testing.T) {
input := "'hello'"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != "hello" {
t.Errorf("Expected 'hello', got %v", result)
}
}

func TestVMBooleanLiterals(t *testing.T) {
tests := []struct {
input    string
expected bool
}{
{"true", true},
{"false", false},
}

for _, tt := range tests {
p := parser.New(tt.input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error for %s: %v", tt.input, err)
}

result := vm.StackTop()
if result != tt.expected {
t.Errorf("Expected %v, got %v", tt.expected, result)
}
}
}

func TestVMNilLiteral(t *testing.T) {
input := "nil"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != nil {
t.Errorf("Expected nil, got %v", result)
}
}

func TestVMArithmetic(t *testing.T) {
tests := []struct {
input    string
expected int64
}{
{"3 + 4", 7},
{"10 - 5", 5},
{"6 * 7", 42},
{"20 / 4", 5},
}

for _, tt := range tests {
p := parser.New(tt.input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error for %s: %v", tt.input, err)
}

result := vm.StackTop()
if result != tt.expected {
t.Errorf("For %s, expected %v, got %v", tt.input, tt.expected, result)
}
}
}

func TestVMComparison(t *testing.T) {
tests := []struct {
input    string
expected bool
}{
{"3 < 5", true},
{"5 < 3", false},
{"3 > 5", false},
{"5 > 3", true},
{"3 = 3", true},
{"3 = 5", false},
}

for _, tt := range tests {
p := parser.New(tt.input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error for %s: %v", tt.input, err)
}

result := vm.StackTop()
if result != tt.expected {
t.Errorf("For %s, expected %v, got %v", tt.input, tt.expected, result)
}
}
}

func TestVMVariableDeclarationAndAssignment(t *testing.T) {
input := "| x | x := 42. x"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(42) {
t.Errorf("Expected 42, got %v", result)
}
}

func TestVMMultipleStatements(t *testing.T) {
input := "42. 'hello'. true"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != true {
t.Errorf("Expected true (last value), got %v", result)
}
}

func TestVMSimpleBlock(t *testing.T) {
input := "[ 42 ] value"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(42) {
t.Errorf("Expected 42, got %v", result)
}
}

func TestVMBlockWithOneParameter(t *testing.T) {
input := "[ :x | x * 2 ] value: 5"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(10) {
t.Errorf("Expected 10, got %v", result)
}
}

func TestVMBlockWithTwoParameters(t *testing.T) {
input := "[ :x :y | x + y ] value: 3 value: 7"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(10) {
t.Errorf("Expected 10, got %v", result)
}
}

func TestVMArrayLiteral(t *testing.T) {
input := "#(1 2 3) size"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(3) {
t.Errorf("Expected 3, got %v", result)
}
}

func TestVMArrayAt(t *testing.T) {
input := "#(10 20 30) at: 2"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(20) {
t.Errorf("Expected 20, got %v", result)
}
}

func TestVMIfTrue(t *testing.T) {
input := "true ifTrue: [ 42 ]"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(42) {
t.Errorf("Expected 42, got %v", result)
}
}

func TestVMIfFalse(t *testing.T) {
input := "false ifFalse: [ 99 ]"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

result := vm.StackTop()
if result != int64(99) {
t.Errorf("Expected 99, got %v", result)
}
}


func TestVMTimesRepeat(t *testing.T) {
input := "5 timesRepeat: [ 1 ]"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

// timesRepeat returns nil
result := vm.StackTop()
if result != nil {
t.Errorf("Expected nil, got %v", result)
}
}

func TestVMArrayDo(t *testing.T) {
input := "#(1 2 3) do: [ :x | x ]"

p := parser.New(input)
program, _ := p.Parse()
c := compiler.New()
bc, _ := c.Compile(program)

vm := New()
err := vm.Run(bc)

if err != nil {
t.Fatalf("VM error: %v", err)
}

// do: returns the array
result := vm.StackTop()
array, ok := result.(*Array)
if !ok {
t.Fatalf("Expected array, got %T", result)
}
if len(array.Elements) != 3 {
t.Errorf("Expected array with 3 elements, got %d", len(array.Elements))
}
}
