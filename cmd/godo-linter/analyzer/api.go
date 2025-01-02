package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// APIAnalyzer checks for proper API implementation patterns
var APIAnalyzer = &analysis.Analyzer{
	Name: "apicheck",
	Doc:  "checks for proper API implementation patterns",
	Run:  runAPICheck,
}

func runAPICheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				if isHTTPHandler(node) {
					checkHandlerPattern(pass, node)
				}
				if isMiddleware(node) {
					checkMiddlewarePattern(pass, node)
				}
			case *ast.CallExpr:
				checkResponseWriting(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func isHTTPHandler(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 0 {
		return false
	}
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 2 {
		return false
	}

	// Check if parameters match http.ResponseWriter and *http.Request
	params := fn.Type.Params.List
	return isResponseWriter(params[0].Type) && isHTTPRequest(params[1].Type)
}

func isMiddleware(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	// Check if return type is http.HandlerFunc
	return isHandlerFunc(fn.Type.Results.List[0].Type)
}

func checkHandlerPattern(pass *analysis.Pass, fn *ast.FuncDecl) {
	var hasRequestValidation bool
	var hasErrorHandling bool
	var hasContextUsage bool

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if isValidationCall(node) {
				hasRequestValidation = true
			}
			if isErrorHandlingCall(node) {
				hasErrorHandling = true
			}
			if isContextUsageCall(node) {
				hasContextUsage = true
			}
		}
		return true
	})

	if !hasRequestValidation {
		pass.Reportf(fn.Pos(), "handler should validate request input")
	}
	if !hasErrorHandling {
		pass.Reportf(fn.Pos(), "handler should implement error handling")
	}
	if !hasContextUsage {
		pass.Reportf(fn.Pos(), "handler should use request context")
	}
}

func checkMiddlewarePattern(pass *analysis.Pass, fn *ast.FuncDecl) {
	var preservesContext bool
	var hasNextCall bool

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if isContextPreservingCall(node) {
				preservesContext = true
			}
			if isNextHandlerCall(node) {
				hasNextCall = true
			}
		}
		return true
	})

	if !preservesContext {
		pass.Reportf(fn.Pos(), "middleware should preserve request context")
	}
	if !hasNextCall {
		pass.Reportf(fn.Pos(), "middleware should call next handler")
	}
}

func checkResponseWriting(pass *analysis.Pass, call *ast.CallExpr) {
	if isResponseWriteCall(call) {
		checkResponseHeaders(pass, call)
		checkStatusCode(pass, call)
	}
}

func isResponseWriter(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "ResponseWriter"
	}
	return false
}

func isHTTPRequest(expr ast.Expr) bool {
	if star, ok := expr.(*ast.StarExpr); ok {
		if sel, ok := star.X.(*ast.SelectorExpr); ok {
			return sel.Sel.Name == "Request"
		}
	}
	return false
}

func isHandlerFunc(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "HandlerFunc"
	}
	return false
}

func isValidationCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return strings.Contains(sel.Sel.Name, "Validate") ||
			strings.Contains(sel.Sel.Name, "Decode")
	}
	return false
}

func isErrorHandlingCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return strings.Contains(sel.Sel.Name, "Error") ||
			strings.Contains(sel.Sel.Name, "WriteError")
	}
	return false
}

func isContextUsageCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return strings.Contains(sel.Sel.Name, "Context") ||
			strings.Contains(sel.Sel.Name, "WithValue")
	}
	return false
}

func isContextPreservingCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return strings.Contains(sel.Sel.Name, "WithContext") ||
			strings.Contains(sel.Sel.Name, "Context")
	}
	return false
}

func isNextHandlerCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "ServeHTTP" ||
			strings.Contains(sel.Sel.Name, "Handle")
	}
	return false
}

func isResponseWriteCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Write" ||
			sel.Sel.Name == "WriteHeader"
	}
	return false
}

func checkResponseHeaders(pass *analysis.Pass, call *ast.CallExpr) {
	// Check if Content-Type is set
	// Check if security headers are set
	// TODO: Implement header checks
}

func checkStatusCode(pass *analysis.Pass, call *ast.CallExpr) {
	// Check if status code is appropriate
	// TODO: Implement status code validation
}
