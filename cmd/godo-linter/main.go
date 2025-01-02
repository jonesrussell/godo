// Package main provides the entry point for the godolinter tool.
// It initializes and runs the custom analyzer that enforces Godo project conventions.
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// maxInterfaceMethods defines the maximum number of methods allowed in an interface
const maxInterfaceMethods = 5

// Analyzer is the main entry point for the linter
var Analyzer = &analysis.Analyzer{
	Name: "godolinter",
	Doc:  "enforces Godo project-specific conventions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Check for proper error handling
			if ret, ok := n.(*ast.ReturnStmt); ok {
				checkErrorReturn(pass, ret)
			}

			// Check for proper interface usage
			if decl, ok := n.(*ast.GenDecl); ok {
				if decl.Tok.String() == "type" {
					checkInterfaceSegregation(pass, decl)
				}
			}

			// Check for proper test patterns
			if fn, ok := n.(*ast.FuncDecl); ok {
				if isTestFunction(fn) {
					checkTestPatterns(pass, fn)
				}
			}

			return true
		})
	}
	return nil, nil
}

func checkErrorReturn(pass *analysis.Pass, ret *ast.ReturnStmt) {
	for _, expr := range ret.Results {
		if ident, ok := expr.(*ast.Ident); ok {
			if ident.Name == "err" && !isWrappedError(expr) {
				pass.Reportf(expr.Pos(), "errors should be wrapped with context")
			}
		}
	}
}

func checkInterfaceSegregation(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				if len(iface.Methods.List) > maxInterfaceMethods {
					pass.Reportf(typeSpec.Pos(), "interface %s has too many methods, consider splitting it", typeSpec.Name.Name)
				}
			}
		}
	}
}

func checkTestPatterns(pass *analysis.Pass, fn *ast.FuncDecl) {
	hasAssert := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Assert" || sel.Sel.Name == "Require" {
					hasAssert = true
					return false
				}
			}
		}
		return true
	})

	if !hasAssert {
		pass.Reportf(fn.Pos(), "test function %s should use testify assertions", fn.Name.Name)
	}
}

func isTestFunction(fn *ast.FuncDecl) bool {
	return fn.Name.Name != "" && len(fn.Name.Name) > 4 && fn.Name.Name[:4] == "Test"
}

func isWrappedError(expr ast.Expr) bool {
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			return sel.Sel.Name == "Wrap" || sel.Sel.Name == "Wrapf"
		}
	}
	return false
}

func main() {
	singlechecker.Main(Analyzer)
}
