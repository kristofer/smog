#!/bin/bash

# Script to run all Smog example files

set -e

# Build the smog binary first
echo "Building smog..."
go build -o bin/smog ./cmd/smog
echo ""

# Find all .smog files in examples/ directory and sort them
find examples -name "*.smog" -type f | sort | while read -r example; do
    echo "========================================="
    echo "Running: $example"
    echo "========================================="
    ./bin/smog "$example"
    echo ""
    echo ""
done

echo "All examples completed!"
