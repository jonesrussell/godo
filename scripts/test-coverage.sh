#!/bin/bash

# Run tests with coverage
go test -race -coverprofile=coverage.out -coverpkg=./... ./...

# Check coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
THRESHOLD=80

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "Test coverage is below threshold. Current: $COVERAGE%. Required: $THRESHOLD%"
    exit 1
else
    echo "Test coverage is good. Current: $COVERAGE%. Required: $THRESHOLD%"
fi
