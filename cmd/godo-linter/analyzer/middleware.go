package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// MiddlewareAnalyzer checks for proper middleware patterns
var MiddlewareAnalyzer = &analysis.Analyzer{
	Name: "middlewarecheck",
	Doc:  "checks for proper middleware patterns and provides auto-fixes",
	Run:  runMiddlewareCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runMiddlewareCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if isMiddlewareFunc(fn) {
					checkMiddlewarePatterns(pass, fn)
				}
			}
			return true
		})
	}
	return nil, nil
}

func isMiddlewareFunc(fn *ast.FuncDecl) bool {
	// Check if function returns http.HandlerFunc
	if fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
		if sel, ok := fn.Type.Results.List[0].Type.(*ast.SelectorExpr); ok {
			return sel.Sel.Name == "HandlerFunc"
		}
	}
	return false
}

func checkMiddlewarePatterns(pass *analysis.Pass, fn *ast.FuncDecl) {
	var hasNextCall bool
	var preservesContext bool

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if checkNextHandlerCall(node) {
				hasNextCall = true
			}
		case *ast.SelectorExpr:
			if isContextAccess(node) {
				preservesContext = true
			}
		}
		return true
	})

	if !hasNextCall {
		suggestNextHandlerCall(pass, fn)
	}

	if !preservesContext {
		suggestContextPreservation(pass, fn)
	}
}

func checkNextHandlerCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "ServeHTTP"
	}
	return false
}

func isContextAccess(sel *ast.SelectorExpr) bool {
	if sel.Sel.Name == "Context" {
		return true
	}
	return false
}

func suggestNextHandlerCall(pass *analysis.Pass, fn *ast.FuncDecl) {
	// Find the return statement
	var returnPos token.Pos
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			returnPos = ret.Pos()
			return false
		}
		return true
	})

	if returnPos != token.NoPos {
		pass.Report(analysis.Diagnostic{
			Pos:     fn.Pos(),
			Message: "middleware should call next handler",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Add next handler call",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     returnPos,
							End:     returnPos,
							NewText: []byte("next.ServeHTTP(w, r)\n"),
						},
					},
				},
			},
		})
	}
}

func suggestContextPreservation(pass *analysis.Pass, fn *ast.FuncDecl) {
	// Find the function body opening brace
	bodyStart := fn.Body.Lbrace

	pass.Report(analysis.Diagnostic{
		Pos:     fn.Pos(),
		Message: "middleware should preserve request context",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Add context preservation",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     bodyStart + 1,
						End:     bodyStart + 1,
						NewText: []byte("\nctx := r.Context()\n// Add context values here\nr = r.WithContext(ctx)\n"),
					},
				},
			},
		},
	})
}
