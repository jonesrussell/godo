package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// LoggingAnalyzer checks for proper logging patterns
var LoggingAnalyzer = &analysis.Analyzer{
	Name: "loggingcheck",
	Doc:  "checks for proper logging patterns and provides auto-fixes",
	Run:  runLoggingCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runLoggingCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				if isLoggingCall(node) {
					checkLoggingPattern(pass, node)
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkLoggingPattern(pass *analysis.Pass, call *ast.CallExpr) {
	sel := call.Fun.(*ast.SelectorExpr)

	// Check if it's using structured logging
	if !isStructuredLogging(call) {
		suggestStructuredLogging(pass, call, sel)
	}

	// Check if context is included
	if !hasContext(call) {
		suggestContextInclusion(pass, call)
	}
}

func isStructuredLogging(call *ast.CallExpr) bool {
	// Check if any argument is a key-value pair
	for _, arg := range call.Args {
		if _, ok := arg.(*ast.KeyValueExpr); ok {
			return true
		}
	}
	return false
}

func hasContext(call *ast.CallExpr) bool {
	// Check function parameters for context
	if len(call.Args) > 0 {
		if ident, ok := call.Args[0].(*ast.Ident); ok {
			return strings.Contains(strings.ToLower(ident.Name), "ctx")
		}
	}
	return false
}

func suggestStructuredLogging(pass *analysis.Pass, call *ast.CallExpr, sel *ast.SelectorExpr) {
	// Convert simple logging to structured logging
	var newArgs []ast.Expr
	if len(call.Args) > 0 {
		// First arg is usually the message
		newArgs = append(newArgs, call.Args[0])

		// Add structured fields for remaining args
		for i, arg := range call.Args[1:] {
			fieldName := fmt.Sprintf("field%d", i+1)
			newArgs = append(newArgs, &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: fieldName},
				Value: arg,
			})
		}
	}

	pass.Report(analysis.Diagnostic{
		Pos:     call.Pos(),
		Message: "use structured logging with key-value pairs",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Convert to structured logging",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     call.Pos(),
						End:     call.End(),
						NewText: []byte(formatStructuredLogging(sel, newArgs)),
					},
				},
			},
		},
	})
}

func suggestContextInclusion(pass *analysis.Pass, call *ast.CallExpr) {
	// Add context parameter if not present
	newArgs := []ast.Expr{&ast.Ident{Name: "ctx"}}
	newArgs = append(newArgs, call.Args...)

	pass.Report(analysis.Diagnostic{
		Pos:     call.Pos(),
		Message: "include relevant context in log messages",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Add context parameter",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     call.Lparen,
						End:     call.Lparen + 1,
						NewText: []byte("ctx, "),
					},
				},
			},
		},
	})
}

func formatStructuredLogging(sel *ast.SelectorExpr, args []ast.Expr) string {
	var parts []string
	parts = append(parts, sel.X.(*ast.Ident).Name+"."+sel.Sel.Name+"(")

	for i, arg := range args {
		if i > 0 {
			parts = append(parts, ", ")
		}

		switch v := arg.(type) {
		case *ast.KeyValueExpr:
			parts = append(parts, fmt.Sprintf("%s: %s", v.Key.(*ast.Ident).Name, formatExpr(v.Value)))
		default:
			parts = append(parts, formatExpr(arg))
		}
	}

	parts = append(parts, ")")
	return strings.Join(parts, "")
}

func formatExpr(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		return v.Value
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name)
	default:
		return "value"
	}
}
