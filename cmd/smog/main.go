package main

import (
	"fmt"
	"os"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
	"github.com/kristofer/smog/pkg/vm"
)

const version = "0.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version", "-v", "--version":
		fmt.Printf("smog version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
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
	fmt.Println("  smog [file]           Run a smog file")
	fmt.Println("  smog run [file]       Run a smog file")
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
