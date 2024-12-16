#!/usr/bin/env bash

set -e  # Exit on any error

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Display coverage statistics
go tool cover -func=coverage.out

# Check minimum coverage threshold
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
THRESHOLD=80

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
    exit 1
fi
