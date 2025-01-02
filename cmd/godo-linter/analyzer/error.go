package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// ErrorAnalyzer checks for proper error handling patterns
var ErrorAnalyzer = &analysis.Analyzer{
	Name: "errorcheck",
	Doc:  "checks for proper error handling patterns and practices",
	Run:  runErrorCheck,
}

func runErrorCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				checkErrorHandlingInFunction(pass, node)
			case *ast.TypeSpec:
				if isCustomErrorType(node) {
					checkCustomErrorType(pass, node)
				}
			case *ast.ReturnStmt:
				checkErrorReturnWrapping(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func checkErrorHandlingInFunction(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			if isErrorCondition(node.Cond) {
				checkErrorHandlingInBlock(pass, node.Body)
			}
		case *ast.AssignStmt:
			checkErrorAssignment(pass, node)
		}
		return true
	})
}

func isErrorCondition(expr ast.Expr) bool {
	if binExpr, ok := expr.(*ast.BinaryExpr); ok {
		if ident, ok := binExpr.X.(*ast.Ident); ok {
			return ident.Name == "err" || strings.HasSuffix(ident.Name, "Error")
		}
	}
	return false
}

func checkErrorHandlingInBlock(pass *analysis.Pass, block *ast.BlockStmt) {
	var hasErrorWrapping bool
	var hasLogging bool
	var hasReturn bool

	for _, stmt := range block.List {
		switch node := stmt.(type) {
		case *ast.ReturnStmt:
			hasReturn = true
		case *ast.ExprStmt:
			if call, ok := node.X.(*ast.CallExpr); ok {
				if isWrappingCall(call) {
					hasErrorWrapping = true
				}
				if isLogCall(call) {
					hasLogging = true
				}
			}
		}
	}

	if hasReturn && !hasErrorWrapping {
		pass.Reportf(block.Pos(), "errors should be wrapped with context before returning")
	}
	if hasReturn && !hasLogging {
		pass.Reportf(block.Pos(), "errors should be logged before returning")
	}
}

func checkErrorAssignment(pass *analysis.Pass, assign *ast.AssignStmt) {
	for _, rhs := range assign.Rhs {
		if call, ok := rhs.(*ast.CallExpr); ok {
			if returnsError(call) {
				checkErrorHandlingAfterAssignment(pass, assign)
			}
		}
	}
}

func returnsError(call *ast.CallExpr) bool {
	// This is a simplified check. In a real implementation,
	// we would need to check the function's return type.
	return true
}

func checkErrorHandlingAfterAssignment(pass *analysis.Pass, assign *ast.AssignStmt) {
	parent := findParentBlock(assign)
	if parent == nil {
		return
	}

	var hasErrorCheck bool
	for _, stmt := range parent.List {
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			if isErrorCondition(ifStmt.Cond) {
				hasErrorCheck = true
				break
			}
		}
	}

	if !hasErrorCheck {
		pass.Reportf(assign.Pos(), "error return value must be checked")
	}
}

func findParentBlock(node ast.Node) *ast.BlockStmt {
	// This would need to traverse up the AST to find the containing block
	// Implementation omitted for brevity
	return nil
}

func isCustomErrorType(spec *ast.TypeSpec) bool {
	if _, ok := spec.Type.(*ast.InterfaceType); ok {
		return spec.Name.Name == "error" || strings.HasSuffix(spec.Name.Name, "Error")
	}
	if _, ok := spec.Type.(*ast.StructType); ok {
		return strings.HasSuffix(spec.Name.Name, "Error")
	}
	return false
}

func checkCustomErrorType(pass *analysis.Pass, spec *ast.TypeSpec) {
	// Check if custom error type implements Error() string method
	if st, ok := spec.Type.(*ast.StructType); ok {
		hasErrorMethod := false
		for _, field := range st.Fields.List {
			if len(field.Names) > 0 && field.Names[0].Name == "Error" {
				hasErrorMethod = true
				break
			}
		}
		if !hasErrorMethod {
			pass.Reportf(spec.Pos(), "custom error type should implement Error() string method")
		}
	}
}

func checkErrorReturnWrapping(pass *analysis.Pass, ret *ast.ReturnStmt) {
	for _, expr := range ret.Results {
		if isErrorExpr(expr) && !isWrappedErrorExpr(expr) {
			pass.Reportf(expr.Pos(), "returned errors should be wrapped with context")
		}
	}
}

func isErrorExpr(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "err" || strings.HasSuffix(ident.Name, "Error")
	}
	return false
}

func isWrappedErrorExpr(expr ast.Expr) bool {
	if call, ok := expr.(*ast.CallExpr); ok {
		return isWrappingCall(call)
	}
	return false
}

func isWrappingCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Wrap" ||
			name == "Wrapf" ||
			name == "WithStack" ||
			name == "WithMessage" ||
			strings.HasPrefix(name, "Wrap")
	}
	return false
}

func isLogCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		name := sel.Sel.Name
		return name == "Error" ||
			name == "Errorf" ||
			strings.HasPrefix(name, "Log") ||
			strings.HasPrefix(name, "Fatal")
	}
	return false
}
