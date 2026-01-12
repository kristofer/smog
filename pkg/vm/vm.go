// Package vm implements the bytecode virtual machine for smog.
package vm

import (
	"fmt"

	"github.com/kristofer/smog/pkg/bytecode"
)

// VM represents the virtual machine
type VM struct {
	stack     []interface{}
	sp        int // stack pointer
	locals    []interface{}
	globals   map[string]interface{}
	constants []interface{}
}

// New creates a new virtual machine
func New() *VM {
	return &VM{
		stack:   make([]interface{}, 1024),
		sp:      0,
		locals:  make([]interface{}, 256),
		globals: make(map[string]interface{}),
	}
}

// Run executes bytecode on the virtual machine
func (vm *VM) Run(bc *bytecode.Bytecode) error {
	// Reset state for clean execution
	vm.sp = 0
	for i := range vm.locals {
		vm.locals[i] = nil
	}
	
	vm.constants = bc.Constants

	for ip := 0; ip < len(bc.Instructions); ip++ {
		inst := bc.Instructions[ip]

		switch inst.Op {
		case bytecode.OpPush:
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			if err := vm.push(vm.constants[inst.Operand]); err != nil {
				return err
			}

		case bytecode.OpPop:
			if _, err := vm.pop(); err != nil {
				return err
			}

		case bytecode.OpPushTrue:
			if err := vm.push(true); err != nil {
				return err
			}

		case bytecode.OpPushFalse:
			if err := vm.push(false); err != nil {
				return err
			}

		case bytecode.OpPushNil:
			if err := vm.push(nil); err != nil {
				return err
			}

		case bytecode.OpLoadLocal:
			if inst.Operand < 0 || inst.Operand >= len(vm.locals) {
				return fmt.Errorf("local variable index out of bounds: %d", inst.Operand)
			}
			if err := vm.push(vm.locals[inst.Operand]); err != nil {
				return err
			}

		case bytecode.OpStoreLocal:
			if inst.Operand < 0 || inst.Operand >= len(vm.locals) {
				return fmt.Errorf("local variable index out of bounds: %d", inst.Operand)
			}
			val, err := vm.pop()
			if err != nil {
				return err
			}
			vm.locals[inst.Operand] = val
			// Push the value back (assignment returns the value)
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpLoadGlobal:
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			name, ok := vm.constants[inst.Operand].(string)
			if !ok {
				return fmt.Errorf("expected string constant for global name")
			}
			val, ok := vm.globals[name]
			if !ok {
				return fmt.Errorf("undefined global variable: %s", name)
			}
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpStoreGlobal:
			if inst.Operand < 0 || inst.Operand >= len(vm.constants) {
				return fmt.Errorf("constant index out of bounds: %d", inst.Operand)
			}
			name, ok := vm.constants[inst.Operand].(string)
			if !ok {
				return fmt.Errorf("expected string constant for global name")
			}
			val, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[name] = val
			// Push the value back
			if err := vm.push(val); err != nil {
				return err
			}

		case bytecode.OpSend:
			// Decode selector index and arg count from operand using shared constants
			selectorIdx := inst.Operand >> bytecode.SelectorIndexShift
			argCount := inst.Operand & bytecode.ArgCountMask

			if selectorIdx < 0 || selectorIdx >= len(vm.constants) {
				return fmt.Errorf("selector index out of bounds: %d", selectorIdx)
			}
			selector, ok := vm.constants[selectorIdx].(string)
			if !ok {
				return fmt.Errorf("expected string constant for selector")
			}

			// Pop arguments (in reverse order)
			args := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				arg, err := vm.pop()
				if err != nil {
					return err
				}
				args[i] = arg
			}

			// Pop receiver
			receiver, err := vm.pop()
			if err != nil {
				return err
			}

			// Execute the message send
			result, err := vm.send(receiver, selector, args)
			if err != nil {
				return err
			}

			// Push result
			if err := vm.push(result); err != nil {
				return err
			}

		case bytecode.OpReturn:
			return nil

		default:
			return fmt.Errorf("unknown opcode: %v", inst.Op)
		}
	}

	return nil
}

// send executes a message send
func (vm *VM) send(receiver interface{}, selector string, args []interface{}) (interface{}, error) {
	// Handle primitive operations
	switch selector {
	case "+":
		return vm.add(receiver, args[0])
	case "-":
		return vm.subtract(receiver, args[0])
	case "*":
		return vm.multiply(receiver, args[0])
	case "/":
		return vm.divide(receiver, args[0])
	case "<":
		return vm.lessThan(receiver, args[0])
	case ">":
		return vm.greaterThan(receiver, args[0])
	case "<=":
		return vm.lessOrEqual(receiver, args[0])
	case ">=":
		return vm.greaterOrEqual(receiver, args[0])
	case "=":
		return vm.equal(receiver, args[0])
	case "~=":
		return vm.notEqual(receiver, args[0])
	case "println":
		fmt.Println(receiver)
		return receiver, nil
	case "print":
		fmt.Print(receiver)
		return receiver, nil
	default:
		return nil, fmt.Errorf("unknown message: %s", selector)
	}
}

// Primitive operations
func (vm *VM) add(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal + bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal + bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot add %T and %T", a, b)
}

func (vm *VM) subtract(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal - bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal - bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot subtract %T and %T", a, b)
}

func (vm *VM) multiply(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal * bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal * bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot multiply %T and %T", a, b)
}

func (vm *VM) divide(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			if bVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return aVal / bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			if bVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return aVal / bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot divide %T and %T", a, b)
}

func (vm *VM) lessThan(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal < bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal < bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (vm *VM) greaterThan(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal > bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal > bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (vm *VM) lessOrEqual(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal <= bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal <= bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (vm *VM) greaterOrEqual(a, b interface{}) (interface{}, error) {
	switch aVal := a.(type) {
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal >= bVal, nil
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal >= bVal, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (vm *VM) equal(a, b interface{}) (interface{}, error) {
	return a == b, nil
}

func (vm *VM) notEqual(a, b interface{}) (interface{}, error) {
	return a != b, nil
}

func (vm *VM) push(obj interface{}) error {
	if vm.sp >= len(vm.stack) {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

func (vm *VM) pop() (interface{}, error) {
	if vm.sp <= 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	vm.sp--
	return vm.stack[vm.sp], nil
}

// StackTop returns the top value on the stack without popping it
func (vm *VM) StackTop() interface{} {
	if vm.sp <= 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}
