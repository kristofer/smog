#!/bin/bash

# Script to run all Smog example files

# Build the smog binary first
echo "Building smog..."
if ! go build -o bin/smog ./cmd/smog; then
    echo "ERROR: Failed to build smog"
    exit 1
fi
echo ""

# Track success/failure counts
success_count=0
failure_count=0
failed_examples=()

# Find all .smog files in examples/ directory, excluding syntax-only, and sort them
find examples -name "*.smog" -type f -not -path "*/syntax-only/*" | sort | while read -r example; do
    echo "========================================="
    echo "Running: $example"
    echo "========================================="

    if ./bin/smog "$example"; then
        echo "✓ SUCCESS"
        ((success_count++))
    else
        echo "✗ FAILED with exit code $?"
        ((failure_count++))
        failed_examples+=("$example")
    fi

    echo ""
    echo ""
done

# Summary
echo "========================================="
echo "Summary"
echo "========================================="
echo "Successful: $success_count"
echo "Failed: $failure_count"

if [ $failure_count -gt 0 ]; then
    echo ""
    echo "Failed examples:"
    printf '%s\n' "${failed_examples[@]}"
fi

echo ""
echo "All examples completed!"
