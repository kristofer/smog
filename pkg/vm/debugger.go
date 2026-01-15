// Package vm - debugger support
package vm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kristofer/smog/pkg/bytecode"
)

// Debugger provides interactive debugging capabilities for the VM.
type Debugger struct {
	vm          *VM                        // The VM being debugged
	breakpoints map[int]bool               // Instruction positions where execution should pause
	stepMode    bool                       // If true, pause after each instruction
	enabled     bool                       // If true, debugger is active
	bytecode    *bytecode.Bytecode         // Current bytecode being executed
}

// NewDebugger creates a new debugger instance.
func NewDebugger(vm *VM) *Debugger {
	return &Debugger{
		vm:          vm,
		breakpoints: make(map[int]bool),
		stepMode:    false,
		enabled:     false,
	}
}

// Enable activates the debugger.
func (d *Debugger) Enable() {
	d.enabled = true
}

// Disable deactivates the debugger.
func (d *Debugger) Disable() {
	d.enabled = false
}

// SetStepMode enables or disables step mode.
// In step mode, execution pauses after each instruction.
func (d *Debugger) SetStepMode(enabled bool) {
	d.stepMode = enabled
}

// AddBreakpoint adds a breakpoint at the specified instruction position.
func (d *Debugger) AddBreakpoint(ip int) {
	d.breakpoints[ip] = true
}

// RemoveBreakpoint removes a breakpoint at the specified instruction position.
func (d *Debugger) RemoveBreakpoint(ip int) {
	delete(d.breakpoints, ip)
}

// ClearBreakpoints removes all breakpoints.
func (d *Debugger) ClearBreakpoints() {
	d.breakpoints = make(map[int]bool)
}

// ShouldPause checks if execution should pause at the current instruction.
// Returns true if we're in step mode or at a breakpoint.
func (d *Debugger) ShouldPause() bool {
	if !d.enabled {
		return false
	}
	
	if d.stepMode {
		return true
	}
	
	return d.breakpoints[d.vm.ip]
}

// ShowCurrentInstruction displays the current instruction being executed.
func (d *Debugger) ShowCurrentInstruction() {
	if d.bytecode == nil || d.vm.ip >= len(d.bytecode.Instructions) {
		fmt.Println("No current instruction")
		return
	}
	
	inst := d.bytecode.Instructions[d.vm.ip]
	fmt.Printf("  %4d: %s", d.vm.ip, inst.Op)
	d.formatInstructionOperand(inst, d.bytecode.Constants)
	fmt.Println()
}

// formatInstructionOperand formats the operand of an instruction based on its opcode.
func (d *Debugger) formatInstructionOperand(inst bytecode.Instruction, constants []interface{}) {
	switch inst.Op {
	case bytecode.OpSend, bytecode.OpSuperSend:
		selectorIdx := inst.Operand >> bytecode.SelectorIndexShift
		argCount := inst.Operand & bytecode.ArgCountMask
		fmt.Printf(" selector=%d args=%d", selectorIdx, argCount)
		if selectorIdx < len(constants) {
			if sel, ok := constants[selectorIdx].(string); ok {
				fmt.Printf(" (%s)", sel)
			}
		}
	case bytecode.OpMakeClosure:
		codeIdx := inst.Operand >> bytecode.SelectorIndexShift
		paramCount := inst.Operand & bytecode.ArgCountMask
		fmt.Printf(" code=%d params=%d", codeIdx, paramCount)
	default:
		if inst.Operand != 0 {
			fmt.Printf(" %d", inst.Operand)
		}
	}
}

// ShowStack displays the current VM stack.
func (d *Debugger) ShowStack() {
	fmt.Println("Stack (top to bottom):")
	if d.vm.sp == 0 {
		fmt.Println("  (empty)")
		return
	}
	
	for i := d.vm.sp - 1; i >= 0; i-- {
		fmt.Printf("  [%d] %v (%T)\n", i, d.vm.stack[i], d.vm.stack[i])
	}
}

// ShowLocals displays the current local variables.
func (d *Debugger) ShowLocals() {
	fmt.Println("Local variables:")
	hasAny := false
	for i, val := range d.vm.locals {
		if val != nil {
			hasAny = true
			fmt.Printf("  [%d] %v (%T)\n", i, val, val)
		}
	}
	if !hasAny {
		fmt.Println("  (none set)")
	}
}

// ShowGlobals displays all global variables.
func (d *Debugger) ShowGlobals() {
	fmt.Println("Global variables:")
	if len(d.vm.globals) == 0 {
		fmt.Println("  (none)")
		return
	}
	
	for name, val := range d.vm.globals {
		fmt.Printf("  %s = %v (%T)\n", name, val, val)
	}
}

// ShowCallStack displays the current call stack.
func (d *Debugger) ShowCallStack() {
	fmt.Println("Call stack (top to bottom):")
	if len(d.vm.callStack) == 0 {
		fmt.Println("  (empty)")
		return
	}
	
	for i := len(d.vm.callStack) - 1; i >= 0; i-- {
		frame := d.vm.callStack[i]
		fmt.Printf("  %s", frame.Name)
		if frame.Selector != "" {
			fmt.Printf(" (selector: %s)", frame.Selector)
		}
		if frame.IP >= 0 {
			fmt.Printf(" [IP: %d]", frame.IP)
		}
		fmt.Println()
	}
}

// InteractivePrompt provides an interactive debugger prompt.
// This is called when execution pauses at a breakpoint or in step mode.
func (d *Debugger) InteractivePrompt(bc *bytecode.Bytecode) (continueExecution bool) {
	d.bytecode = bc
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("\n=== Debugger Paused ===")
	d.ShowCurrentInstruction()
	
	for {
		fmt.Print("debug> ")
		if !scanner.Scan() {
			return false
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		command := parts[0]
		
		switch command {
		case "help", "h", "?":
			d.printHelp()
			
		case "continue", "c":
			d.SetStepMode(false)
			return true
			
		case "step", "s":
			d.SetStepMode(true)
			return true
			
		case "next", "n":
			// Step one instruction
			return true
			
		case "stack", "st":
			d.ShowStack()
			
		case "locals", "l":
			d.ShowLocals()
			
		case "globals", "g":
			d.ShowGlobals()
			
		case "callstack", "cs":
			d.ShowCallStack()
			
		case "instruction", "i":
			d.ShowCurrentInstruction()
			
		case "breakpoint", "b":
			if len(parts) < 2 {
				fmt.Println("Usage: breakpoint <instruction_number>")
				continue
			}
			ip, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid instruction number")
				continue
			}
			d.AddBreakpoint(ip)
			fmt.Printf("Breakpoint added at instruction %d\n", ip)
			
		case "delete", "d":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <instruction_number>")
				continue
			}
			ip, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid instruction number")
				continue
			}
			d.RemoveBreakpoint(ip)
			fmt.Printf("Breakpoint removed at instruction %d\n", ip)
			
		case "list", "ls":
			d.listInstructions(bc)
			
		case "quit", "q":
			return false
			
		default:
			fmt.Printf("Unknown command: %s (type 'help' for commands)\n", command)
		}
	}
}

// printHelp displays available debugger commands.
func (d *Debugger) printHelp() {
	fmt.Println("Debugger Commands:")
	fmt.Println("  help, h, ?           Show this help")
	fmt.Println("  continue, c          Continue execution")
	fmt.Println("  step, s              Enable step mode (pause after each instruction)")
	fmt.Println("  next, n              Execute next instruction")
	fmt.Println("  stack, st            Show VM stack")
	fmt.Println("  locals, l            Show local variables")
	fmt.Println("  globals, g           Show global variables")
	fmt.Println("  callstack, cs        Show call stack")
	fmt.Println("  instruction, i       Show current instruction")
	fmt.Println("  breakpoint <n>, b    Add breakpoint at instruction n")
	fmt.Println("  delete <n>, d        Remove breakpoint at instruction n")
	fmt.Println("  list, ls             List all instructions")
	fmt.Println("  quit, q              Quit debugging (abort execution)")
}

// listInstructions displays all instructions in the bytecode.
func (d *Debugger) listInstructions(bc *bytecode.Bytecode) {
	fmt.Println("Instructions:")
	for i, inst := range bc.Instructions {
		marker := "  "
		if i == d.vm.ip {
			marker = "->"
		} else if d.breakpoints[i] {
			marker = "*"
		}
		
		fmt.Printf("%s %4d: %s", marker, i, inst.Op)
		d.formatInstructionOperand(inst, bc.Constants)
		fmt.Println()
	}
}
