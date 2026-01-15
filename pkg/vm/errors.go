// Package vm - error handling with stack traces
package vm

import (
	"fmt"
	"strings"
)

// StackFrame represents a single frame in the call stack.
// It captures information about where execution is occurring.
type StackFrame struct {
	Name       string // Method/function name or description
	Selector   string // Message selector (for message sends)
	IP         int    // Instruction pointer at time of call
	SourceLine int    // Source line number (0 if unknown)
	SourceCol  int    // Source column number (0 if unknown)
}

// RuntimeError represents a runtime error with stack trace information.
// This provides detailed context about where an error occurred.
type RuntimeError struct {
	Message    string       // Error message
	StackTrace []StackFrame // Call stack at time of error
}

// Error implements the error interface.
// It formats the error message with a stack trace.
func (e *RuntimeError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	
	if len(e.StackTrace) > 0 {
		b.WriteString("\n\nStack trace:")
		for i := len(e.StackTrace) - 1; i >= 0; i-- {
			frame := e.StackTrace[i]
			b.WriteString(fmt.Sprintf("\n  at %s", frame.Name))
			if frame.Selector != "" {
				b.WriteString(fmt.Sprintf(" (selector: %s)", frame.Selector))
			}
			if frame.SourceLine > 0 {
				b.WriteString(fmt.Sprintf(" [line %d:%d]", frame.SourceLine, frame.SourceCol))
			}
			if frame.IP >= 0 {
				b.WriteString(fmt.Sprintf(" [IP: %d]", frame.IP))
			}
		}
	}
	
	return b.String()
}

// newRuntimeError creates a new RuntimeError with the given message.
func newRuntimeError(message string, stack []StackFrame) *RuntimeError {
	return &RuntimeError{
		Message:    message,
		StackTrace: stack,
	}
}
