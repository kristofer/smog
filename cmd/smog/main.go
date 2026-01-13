package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kristofer/smog/pkg/bytecode"
	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

const version = "0.4.0"

func main() {
	if len(os.Args) < 2 {
		// No arguments - start REPL
		runREPL()
		return
	}

	switch os.Args[1] {
	case "version", "-v", "--version":
		fmt.Printf("smog version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	case "repl":
		runREPL()
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: no file specified")
			printUsage()
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "compile":
		// Compile a .smog file to .sg bytecode
		if len(os.Args) < 3 {
			fmt.Println("Error: no file specified")
			fmt.Println("\nUsage: smog compile <input.smog> [output.sg]")
			os.Exit(1)
		}
		inputFile := os.Args[2]
		outputFile := ""
		if len(os.Args) >= 4 {
			outputFile = os.Args[3]
		}
		compileFile(inputFile, outputFile)
	case "disassemble", "disasm":
		// Disassemble a .sg file to human-readable format
		if len(os.Args) < 3 {
			fmt.Println("Error: no file specified")
			fmt.Println("\nUsage: smog disassemble <file.sg>")
			os.Exit(1)
		}
		disassembleFile(os.Args[2])
	default:
		// Assume it's a file to run
		runFile(os.Args[1])
	}
}

func printUsage() {
	fmt.Println("smog - A simple object-oriented language")
	fmt.Println("\nUsage:")
	fmt.Println("  smog                       Start interactive REPL")
	fmt.Println("  smog [file]                Run a .smog or .sg file")
	fmt.Println("  smog run [file]            Run a .smog or .sg file")
	fmt.Println("  smog compile <in> [out]    Compile .smog to .sg bytecode")
	fmt.Println("  smog disassemble <file>    Disassemble .sg bytecode file")
	fmt.Println("  smog repl                  Start interactive REPL")
	fmt.Println("  smog version               Show version")
	fmt.Println("  smog help                  Show this help")
	fmt.Println("\nFile Extensions:")
	fmt.Println("  .smog   Source code files (text)")
	fmt.Println("  .sg     Compiled bytecode files (binary)")
}

// runFile runs a .smog source file or .sg bytecode file.
//
// This function automatically detects the file type based on extension:
//   - .sg files are loaded directly as bytecode (fast)
//   - .smog files are parsed and compiled first (slower)
//
// This allows users to pre-compile frequently-used programs to .sg format
// for faster startup time.
func runFile(filename string) {
	ext := filepath.Ext(filename)
	
	// Check if it's a compiled bytecode file
	if ext == ".sg" {
		runBytecodeFile(filename)
		return
	}
	
	// Otherwise, treat it as source code
	runSourceFile(filename)
}

// runSourceFile reads, parses, compiles, and executes a .smog source file.
//
// This is the traditional path: source → AST → bytecode → execution.
// It's slower than runBytecodeFile because it includes parsing and compilation.
func runSourceFile(filename string) {
	// Read the source file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the source code into an AST
	p := parser.New(string(data))
	program, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Compile the AST to bytecode
	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compile error: %v\n", err)
		os.Exit(1)
	}

	// Run the bytecode on the VM
	v := vm.New()
	err = v.Run(bc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}

// runBytecodeFile loads and executes a pre-compiled .sg bytecode file.
//
// This is the fast path: bytecode → execution (no parsing or compilation).
// It's significantly faster than runSourceFile for large programs.
//
// Performance benefit:
//   - No text parsing
//   - No AST construction
//   - No bytecode compilation
//   - Direct deserialization from binary format
func runBytecodeFile(filename string) {
	// Open the bytecode file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode the bytecode from the file
	bc, err := bytecode.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bytecode: %v\n", err)
		os.Exit(1)
	}

	// Run the bytecode on the VM
	v := vm.New()
	err = v.Run(bc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}

// compileFile compiles a .smog source file to a .sg bytecode file.
//
// This command allows users to pre-compile their programs for faster loading.
// The .sg file can be distributed and executed without the source code.
//
// Usage:
//   smog compile program.smog           -> creates program.sg
//   smog compile program.smog out.sg    -> creates out.sg
//
// Benefits of compilation:
//   - Faster program startup (no parsing/compilation at runtime)
//   - Smaller file size in some cases (binary format)
//   - Code distribution without exposing source
//   - Enables building multi-file programs with pre-compiled modules
func compileFile(inputFile, outputFile string) {
	// Default output filename: replace .smog extension with .sg
	if outputFile == "" {
		if filepath.Ext(inputFile) == ".smog" {
			outputFile = inputFile[:len(inputFile)-5] + ".sg"
		} else {
			outputFile = inputFile + ".sg"
		}
	}

	// Read the source file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the source code into an AST
	p := parser.New(string(data))
	program, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Compile the AST to bytecode
	c := compiler.New()
	bc, err := c.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compile error: %v\n", err)
		os.Exit(1)
	}

	// Write the bytecode to the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := bytecode.Encode(bc, outFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing bytecode: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compiled %s -> %s\n", inputFile, outputFile)
}

// disassembleFile prints a human-readable representation of a .sg bytecode file.
//
// This is a debugging tool that shows:
//   - The constant pool contents
//   - The instruction sequence with opcodes and operands
//   - Class and method definitions (if present)
//
// It's useful for:
//   - Understanding how source code compiles
//   - Debugging compiler issues
//   - Learning the bytecode format
//   - Verifying .sg file contents
//
// Example output:
//   Constants Pool:
//     [0] int64: 42
//     [1] string: "println"
//   Instructions:
//     0: PUSH 0
//     1: SEND (1<<8)|0
//     2: RETURN 0
func disassembleFile(filename string) {
	// Open the bytecode file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode the bytecode
	bc, err := bytecode.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bytecode: %v\n", err)
		os.Exit(1)
	}

	// Print disassembly
	fmt.Printf("=== Bytecode Disassembly: %s ===\n\n", filename)
	
	// Print constant pool
	fmt.Println("Constants Pool:")
	if len(bc.Constants) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, c := range bc.Constants {
			fmt.Printf("  [%d] %s\n", i, formatConstant(c, "  "))
		}
	}
	
	fmt.Println("\nInstructions:")
	if len(bc.Instructions) == 0 {
		fmt.Println("  (empty)")
	} else {
		for i, instr := range bc.Instructions {
			fmt.Printf("  %4d: %s", i, instr.Op)
			
			// Format operand based on opcode
			switch instr.Op {
			case bytecode.OpSend, bytecode.OpSuperSend:
				// Decode message send operand
				selectorIdx := instr.Operand >> bytecode.SelectorIndexShift
				argCount := instr.Operand & bytecode.ArgCountMask
				fmt.Printf(" selector=%d args=%d", selectorIdx, argCount)
			case bytecode.OpMakeClosure:
				// Decode closure operand
				codeIdx := instr.Operand >> bytecode.SelectorIndexShift
				paramCount := instr.Operand & bytecode.ArgCountMask
				fmt.Printf(" code=%d params=%d", codeIdx, paramCount)
			default:
				// Simple operand
				if instr.Operand != 0 {
					fmt.Printf(" %d", instr.Operand)
				}
			}
			fmt.Println()
		}
	}
}

// formatConstant returns a human-readable string representation of a constant.
//
// This helper function is used by disassembleFile to pretty-print constants.
// It handles all constant types including nested structures.
func formatConstant(c interface{}, indent string) string {
	switch v := c.(type) {
	case int64:
		return fmt.Sprintf("int64: %d", v)
	case float64:
		return fmt.Sprintf("float64: %f", v)
	case string:
		return fmt.Sprintf("string: %q", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	case nil:
		return "nil"
	case *bytecode.ClassDefinition:
		return fmt.Sprintf("class: %s (extends %s, %d fields, %d methods)",
			v.Name, v.SuperClass, len(v.Fields), len(v.Methods))
	case *bytecode.MethodDefinition:
		return fmt.Sprintf("method: %s (%d params, %d instructions)",
			v.Selector, len(v.Parameters), len(v.Code.Instructions))
	case *bytecode.Bytecode:
		return fmt.Sprintf("bytecode: %d instructions, %d constants",
			len(v.Instructions), len(v.Constants))
	default:
		return fmt.Sprintf("unknown: %T", c)
	}
}

// runREPL starts an interactive Read-Eval-Print Loop.
//
// The REPL allows users to enter smog expressions and see the results immediately.
// It's useful for experimentation, learning, and quick testing.
//
// Features:
//   - Multi-line input support (statements ending with period)
//   - Persistent VM state (variables and values carry over between inputs)
//   - Persistent compiler state (local variables persist across inputs)
//   - Error recovery (errors don't crash the REPL)
//   - Special commands: :quit, :exit, :help
//
// Example session:
//   smog> | x |
//   smog> x := 42.
//   smog> x + 8.
//   => 50
func runREPL() {
	fmt.Printf("smog REPL v%s\n", version)
	fmt.Println("Type ':help' for help, ':quit' or ':exit' to exit")
	fmt.Println()

	// Create a persistent VM for the REPL session
	v := vm.New()
	// Create a persistent compiler for the REPL session
	// This maintains the symbol table across evaluations so that
	// local variables declared in one input remain available in subsequent inputs
	c := compiler.New()
	scanner := bufio.NewScanner(os.Stdin)
	
	// Buffer for multi-line input
	var inputBuffer strings.Builder
	
	for {
		// Show prompt
		if inputBuffer.Len() == 0 {
			fmt.Print("smog> ")
		} else {
			fmt.Print("....> ")
		}
		
		// Read input
		if !scanner.Scan() {
			break
		}
		
		line := scanner.Text()
		
		// Handle special commands
		if inputBuffer.Len() == 0 {
			switch strings.TrimSpace(line) {
			case ":quit", ":exit":
				fmt.Println("Goodbye!")
				return
			case ":help":
				printREPLHelp()
				continue
			case "":
				continue
			}
		}
		
		// Add line to buffer
		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")
		
		// Check if we have a complete statement (ends with period)
		// or if the line is empty (just execute what we have)
		//
		// Note: This is a simple heuristic that checks for a trailing period.
		// It doesn't handle periods within string literals or comments.
		// A more robust implementation would integrate with the parser.
		// For typical REPL usage, this simple approach works well.
		input := strings.TrimSpace(inputBuffer.String())
		if !strings.HasSuffix(input, ".") && line != "" {
			// Not complete yet, continue reading
			continue
		}
		
		// We have complete input, try to execute it
		if input != "" {
			evalREPL(v, c, input)
		}
		
		// Clear buffer for next input
		inputBuffer.Reset()
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

// evalREPL evaluates a single REPL input.
//
// This function parses, compiles, and runs the input using the persistent VM
// and persistent compiler. The compiler maintains the symbol table so that
// local variables declared in previous inputs remain available.
// Errors are printed but don't stop the REPL.
func evalREPL(v *vm.VM, c *compiler.Compiler, input string) {
	// Parse the input
	p := parser.New(input)
	program, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		if len(p.Errors()) > 0 {
			for _, e := range p.Errors() {
				fmt.Fprintf(os.Stderr, "  %s\n", e)
			}
		}
		return
	}
	
	// Compile the AST to bytecode using incremental compilation
	// This preserves the symbol table across REPL inputs so that
	// local variables remain accessible
	bc, err := c.CompileIncremental(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compile error: %v\n", err)
		return
	}
	
	// Run the bytecode
	err = v.Run(bc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return
	}
	
	// Success - no output for now (could show result of last expression)
}

// printREPLHelp prints help information for the REPL.
func printREPLHelp() {
	fmt.Println("smog REPL Help")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  :help     Show this help message")
	fmt.Println("  :quit     Exit the REPL")
	fmt.Println("  :exit     Exit the REPL")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  - Enter smog expressions and press Enter")
	fmt.Println("  - Statements should end with a period (.)")
	fmt.Println("  - Use | vars | to declare variables")
	fmt.Println("  - Variables persist across statements")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  smog> | x |")
	fmt.Println("  smog> x := 42.")
	fmt.Println("  smog> x + 8.")
	fmt.Println()
	fmt.Println("  smog> 'Hello, World!' println.")
	fmt.Println()
}
