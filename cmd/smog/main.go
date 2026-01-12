package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

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
	// TODO: Implement file execution
	// This will involve:
	// 1. Reading the source file
	// 2. Parsing it into an AST
	// 3. Compiling the AST to bytecode
	// 4. Running the bytecode on the VM
	
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("TODO: Execute smog code from %s\n", filename)
	fmt.Printf("Source length: %d bytes\n", len(data))
}
