package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	default:
		// Assume it's a file to run
		runFile(os.Args[1])
	}
}

func printUsage() {
	fmt.Println("smog - A simple object-oriented language")
	fmt.Println("\nUsage:")
	fmt.Println("  smog                  Start interactive REPL")
	fmt.Println("  smog [file]           Run a smog file")
	fmt.Println("  smog run [file]       Run a smog file")
	fmt.Println("  smog repl             Start interactive REPL")
	fmt.Println("  smog version          Show version")
	fmt.Println("  smog help             Show this help")
}

func runFile(filename string) {
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

// runREPL starts an interactive Read-Eval-Print Loop.
//
// The REPL allows users to enter smog expressions and see the results immediately.
// It's useful for experimentation, learning, and quick testing.
//
// Features:
//   - Multi-line input support (statements ending with period)
//   - Persistent VM state (variables and values carry over between inputs)
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
			evalREPL(v, input)
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
// This function parses, compiles, and runs the input using the persistent VM.
// Errors are printed but don't stop the REPL.
func evalREPL(v *vm.VM, input string) {
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
	
	// Compile the AST to bytecode
	c := compiler.New()
	bc, err := c.Compile(program)
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
