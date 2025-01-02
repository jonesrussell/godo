package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// LoggingAnalyzer checks for proper logging patterns
var LoggingAnalyzer = &analysis.Analyzer{
	Name: "loggingcheck",
	Doc:  "checks for proper logging patterns and practices",
	Run:  runLoggingCheck,
}

func runLoggingCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				if isLoggerCall(node) {
					checkLoggingPattern(pass, node)
				}
			case *ast.FuncDecl:
				checkFunctionLogging(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func isLoggerCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return strings.HasPrefix(name, "Log") ||
			name == "Error" ||
			name == "Info" ||
			name == "Debug" ||
			name == "Warn" ||
			name == "Fatal"
	}
	return false
}

func checkLoggingPattern(pass *analysis.Pass, call *ast.CallExpr) {
	// Check for structured logging
	if !hasStructuredFields(call) {
		pass.Reportf(call.Pos(), "use structured logging with key-value pairs")
	}

	// Check for appropriate log level
	if !hasAppropriateLogLevel(call) {
		pass.Reportf(call.Pos(), "use appropriate log level (Error, Info, Debug, Warn)")
	}

	// Check for context inclusion
	if !hasContextField(call) {
		pass.Reportf(call.Pos(), "include relevant context in log messages")
	}

	// Check for sensitive data
	if hasSensitiveData(call) {
		pass.Reportf(call.Pos(), "avoid logging sensitive information")
	}
}

func checkFunctionLogging(pass *analysis.Pass, fn *ast.FuncDecl) {
	var hasErrorLogging bool
	var hasErrorReturn bool

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ReturnStmt:
			for _, result := range node.Results {
				if isErrorType(result) {
					hasErrorReturn = true
				}
			}
		case *ast.CallExpr:
			if isErrorLoggingCall(node) {
				hasErrorLogging = true
			}
		}
		return true
	})

	if hasErrorReturn && !hasErrorLogging {
		pass.Reportf(fn.Pos(), "log errors before returning them")
	}
}

func hasStructuredFields(call *ast.CallExpr) bool {
	if len(call.Args) < 2 {
		return false
	}

	// Check if arguments are key-value pairs
	for i := 1; i < len(call.Args); i += 2 {
		if !isStringLiteral(call.Args[i-1]) {
			return false
		}
	}
	return true
}

func hasAppropriateLogLevel(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		switch sel.Sel.Name {
		case "Error", "Info", "Debug", "Warn", "Fatal":
			return true
		}
	}
	return false
}

func hasContextField(call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		if isStringLiteral(arg) {
			str := getStringValue(arg)
			if strings.Contains(strings.ToLower(str), "ctx") ||
				strings.Contains(strings.ToLower(str), "context") ||
				strings.Contains(strings.ToLower(str), "correlation") ||
				strings.Contains(strings.ToLower(str), "request_id") {
				return true
			}
		}
	}
	return false
}

func hasSensitiveData(call *ast.CallExpr) bool {
	sensitivePatterns := []string{
		"password", "token", "secret", "key", "auth",
		"credential", "private", "cert", "ssh",
	}

	for _, arg := range call.Args {
		if isStringLiteral(arg) {
			str := strings.ToLower(getStringValue(arg))
			for _, pattern := range sensitivePatterns {
				if strings.Contains(str, pattern) {
					return true
				}
			}
		}
	}
	return false
}

func isErrorLoggingCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Error" ||
			(sel.Sel.Name == "Log" && hasErrorArg(call))
	}
	return false
}

func hasErrorArg(call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		if isErrorType(arg) {
			return true
		}
	}
	return false
}

func isErrorType(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "error" || ident.Name == "err"
	}
	return false
}

func isStringLiteral(expr ast.Expr) bool {
	_, ok := expr.(*ast.BasicLit)
	return ok
}

func getStringValue(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}
