package analyzer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// ErrorAnalyzer checks for proper error handling patterns
var ErrorAnalyzer = &analysis.Analyzer{
	Name: "errorcheck",
	Doc:  "checks for proper error handling patterns and provides auto-fixes",
	Run:  runErrorCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runErrorCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ReturnStmt:
				checkErrorReturnWithFix(pass, node)
			case *ast.IfStmt:
				if isErrorCheck(node.Cond) {
					checkErrorHandlingBlockWithFix(pass, node.Body)
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkErrorReturnWithFix(pass *analysis.Pass, ret *ast.ReturnStmt) {
	for _, expr := range ret.Results {
		if isErrorResult(expr) && !isWrappedError(expr) {
			suggestErrorWrapping(pass, expr)
		}
	}
}

func isErrorResult(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "err" || ident.Name == "error"
	}
	return false
}

func suggestErrorWrapping(pass *analysis.Pass, expr ast.Expr) {
	// Get function name from context
	funcName := getCurrentFunctionName(pass, expr)

	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: "errors should be wrapped with context before returning",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Wrap error with context",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(fmt.Sprintf("fmt.Errorf(\"%s: %%w\", %s)", funcName, expr)),
					},
				},
			},
		},
	})
}

func getCurrentFunctionName(pass *analysis.Pass, node ast.Node) string {
	var funcName string
	ast.Inspect(pass.Files[0], func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Pos() <= node.Pos() && node.Pos() <= fn.End() {
				funcName = fn.Name.Name
				return false
			}
		}
		return true
	})
	return funcName
}

func checkErrorHandlingBlockWithFix(pass *analysis.Pass, block *ast.BlockStmt) {
	hasLogging := false
	hasReturn := false
	var returnStmt *ast.ReturnStmt

	for _, stmt := range block.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				if isLoggingCall(callExpr) {
					hasLogging = true
				}
			}
		}
		if ret, ok := stmt.(*ast.ReturnStmt); ok {
			hasReturn = true
			returnStmt = ret
		}
	}

	if hasReturn && !hasLogging {
		suggestErrorLogging(pass, returnStmt)
	}
}

func suggestErrorLogging(pass *analysis.Pass, ret *ast.ReturnStmt) {
	// Find the error variable being returned
	var errVar string
	for _, expr := range ret.Results {
		if ident, ok := expr.(*ast.Ident); ok {
			if ident.Name == "err" || ident.Name == "error" {
				errVar = ident.Name
				break
			}
		}
	}

	if errVar != "" {
		pass.Report(analysis.Diagnostic{
			Pos:     ret.Pos(),
			Message: "errors should be logged before being returned",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Add error logging",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     ret.Pos(),
							End:     ret.Pos(),
							NewText: []byte(fmt.Sprintf("log.Error(\"error occurred\", \"error\", %s)\n", errVar)),
						},
					},
				},
			},
		})
	}
}
