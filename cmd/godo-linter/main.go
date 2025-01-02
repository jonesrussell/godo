// Package main provides the entry point for the godolinter tool.
// It initializes and runs the custom analyzer that enforces Godo project conventions.
package main

import (
	"go/ast"
	"strings"

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
					checkInterfaceNaming(pass, decl)
					checkInterfaceLocation(pass, decl, file)
				}
			}

			// Check for proper test patterns
			if fn, ok := n.(*ast.FuncDecl); ok {
				if isTestFunction(fn) {
					checkTestPatterns(pass, fn)
				}
				// Check error handling in functions
				checkFunctionErrorHandling(pass, fn)
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

func checkInterfaceNaming(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				name := typeSpec.Name.Name
				if !strings.HasSuffix(name, "er") && !strings.HasSuffix(name, "Service") {
					pass.Reportf(typeSpec.Pos(), "interface %s should follow naming convention (end with 'er' or 'Service')", name)
				}
			}
		}
	}
}

func checkInterfaceLocation(pass *analysis.Pass, decl *ast.GenDecl, file *ast.File) {
	// Check if interface is defined in consumer package
	if strings.Contains(pass.Fset.Position(file.Pos()).String(), "internal/interfaces") {
		pass.Reportf(decl.Pos(), "interfaces should be defined in consumer packages, not in a central interfaces package")
	}
}

func checkFunctionErrorHandling(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if ifStmt, ok := n.(*ast.IfStmt); ok {
			if isErrorCheck(ifStmt.Cond) {
				checkErrorHandlingBlock(pass, ifStmt.Body)
			}
		}
		return true
	})
}

func isErrorCheck(expr ast.Expr) bool {
	if binExpr, ok := expr.(*ast.BinaryExpr); ok {
		if ident, ok := binExpr.X.(*ast.Ident); ok {
			return ident.Name == "err"
		}
	}
	return false
}

func checkErrorHandlingBlock(pass *analysis.Pass, block *ast.BlockStmt) {
	// Check if error is being logged before being returned
	hasLogging := false
	hasReturn := false

	for _, stmt := range block.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				if isLoggingCall(callExpr) {
					hasLogging = true
				}
			}
		}
		if _, ok := stmt.(*ast.ReturnStmt); ok {
			hasReturn = true
		}
	}

	if hasReturn && !hasLogging {
		pass.Reportf(block.Pos(), "errors should be logged before being returned")
	}
}

func isLoggingCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return strings.HasPrefix(sel.Sel.Name, "Log") ||
			strings.HasPrefix(sel.Sel.Name, "Error") ||
			strings.HasPrefix(sel.Sel.Name, "Warn") ||
			strings.HasPrefix(sel.Sel.Name, "Info") ||
			strings.HasPrefix(sel.Sel.Name, "Debug")
	}
	return false
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
