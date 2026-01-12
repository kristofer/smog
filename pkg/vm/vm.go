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
	globals   map[string]interface{}
	constants []interface{}
}

// New creates a new virtual machine
func New() *VM {
	return &VM{
		stack:   make([]interface{}, 1024),
		sp:      0,
		globals: make(map[string]interface{}),
	}
}

// Run executes bytecode on the virtual machine
func (vm *VM) Run(bc *bytecode.Bytecode) error {
	vm.constants = bc.Constants
	
	for ip := 0; ip < len(bc.Instructions); ip++ {
		inst := bc.Instructions[ip]
		
		switch inst.Op {
		case bytecode.OpPush:
			vm.push(vm.constants[inst.Operand])
		case bytecode.OpPop:
			vm.pop()
		case bytecode.OpReturn:
			return nil
		// TODO: Implement all opcodes
		default:
			return fmt.Errorf("unknown opcode: %v", inst.Op)
		}
	}
	
	return nil
}

func (vm *VM) push(obj interface{}) {
	vm.stack[vm.sp] = obj
	vm.sp++
}

func (vm *VM) pop() interface{} {
	vm.sp--
	return vm.stack[vm.sp]
}
