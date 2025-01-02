// Package main provides the entry point for the godolinter tool.
// It initializes and runs the custom analyzer that enforces Godo project conventions.
package main

import (
	"github.com/jonesrussell/godo/cmd/godo-linter/analyzer"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		analyzer.APIAnalyzer,
		analyzer.StorageAnalyzer,
		analyzer.LoggingAnalyzer,
		analyzer.ErrorAnalyzer,
		analyzer.TaskAnalyzer,
		analyzer.MiddlewareAnalyzer,
		analyzer.StandardAnalyzer,
		analyzer.ArchitectureAnalyzer,
		analyzer.InterfaceLocationAnalyzer,
	)
}
