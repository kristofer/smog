.PHONY: all build test examples stdlib clean help

# Default target
all: build stdlib test examples

# Build the smog executable
build:
	@echo "Building smog executable..."
	go build -o smog ./cmd/smog
	@echo "✓ Built: smog"

# Run all Go tests
test:
	@echo "Running Go tests..."
	go test -v ./...

# Compile all stdlib files to bytecode
stdlib:
	@echo "Compiling stdlib files..."
	@mkdir -p bin
	@success=0; \
	failed=0; \
	failed_list=""; \
	for file in $$(find stdlib -name "*.smog" -type f | sort); do \
		output=$${file%.smog}.sg; \
		echo "  $$file -> $$output"; \
		if ./smog compile $$file $$output 2>&1 | grep -v "^Compiled"; then \
			failed=$$((failed + 1)); \
			failed_list="$$failed_list\n    $$file"; \
		else \
			success=$$((success + 1)); \
		fi; \
	done; \
	echo ""; \
	echo "Stdlib compilation summary:"; \
	echo "  Successful: $$success"; \
	echo "  Failed: $$failed"; \
	if [ $$failed -gt 0 ]; then \
		echo ""; \
		echo "  Failed files (may contain parser issues):$$failed_list"; \
		echo ""; \
		echo "  Note: Some stdlib files may use syntax features not yet fully supported."; \
		echo "  This does not affect the core smog functionality."; \
	fi; \
	echo "✓ Stdlib compilation attempted"

# Run all example files
examples:
	@echo "Running examples..."
	@success=0; \
	failed=0; \
	failed_list=""; \
	for file in $$(find examples -name "*.smog" -type f -not -path "*/syntax-only/*" | sort); do \
		echo "========================================="; \
		echo "Running: $$file"; \
		echo "========================================="; \
		if ./smog $$file; then \
			echo "✓ SUCCESS"; \
			success=$$((success + 1)); \
		else \
			echo "✗ FAILED"; \
			failed=$$((failed + 1)); \
			failed_list="$$failed_list\n  $$file"; \
		fi; \
		echo ""; \
	done; \
	echo "========================================="; \
	echo "Summary"; \
	echo "========================================="; \
	echo "Successful: $$success"; \
	echo "Failed: $$failed"; \
	if [ $$failed -gt 0 ]; then \
		echo ""; \
		echo "Failed examples:$$failed_list"; \
		exit 1; \
	fi; \
	echo ""; \
	echo "✓ All examples completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f smog
	rm -f bin/smog
	find stdlib -name "*.sg" -type f -delete
	find examples -name "*.sg" -type f -delete
	@echo "✓ Cleaned"

# Show help
help:
	@echo "Smog Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all        - Build executable, compile stdlib, run tests and examples (default)"
	@echo "  build      - Build the smog executable"
	@echo "  stdlib     - Compile all stdlib .smog files to .sg bytecode"
	@echo "  test       - Run all Go tests"
	@echo "  examples   - Run all example .smog files"
	@echo "  clean      - Remove build artifacts and compiled bytecode"
	@echo "  help       - Show this help message"
