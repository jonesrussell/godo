package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// ArchitectureAnalyzer checks for architectural pattern compliance
var ArchitectureAnalyzer = &analysis.Analyzer{
	Name: "architecture",
	Doc:  "checks for compliance with architectural patterns",
	Run:  runArchitectureCheck,
}

func runArchitectureCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				checkDomainTypes(pass, node)
				checkInterfaces(pass, node)
				checkStorageImplementation(pass, node)
			case *ast.FuncDecl:
				checkErrorHandling(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func checkDomainTypes(pass *analysis.Pass, ts *ast.TypeSpec) {
	if isDomainType(ts) {
		// Domain types should be in domain package
		if !strings.Contains(pass.Pkg.Path(), "domain") {
			pass.Reportf(ts.Pos(), "domain type %s should be in domain package", ts.Name.Name)
		}

		// Domain types should have validation methods
		if !hasValidationMethod(pass, ts.Name.Name) {
			pass.Reportf(ts.Pos(), "domain type %s should implement Validate() method", ts.Name.Name)
		}
	}
}

func checkInterfaces(pass *analysis.Pass, ts *ast.TypeSpec) {
	if iface, ok := ts.Type.(*ast.InterfaceType); ok {
		// Check interface size
		if iface.Methods != nil && len(iface.Methods.List) > 5 {
			pass.Reportf(ts.Pos(), "interface %s has too many methods (max 5)", ts.Name.Name)
		}

		// Check interface naming
		if !strings.HasSuffix(ts.Name.Name, "er") && !strings.HasSuffix(ts.Name.Name, "Service") {
			pass.Reportf(ts.Pos(), "interface %s should end with 'er' or 'Service'", ts.Name.Name)
		}
	}
}

func checkStorageImplementation(pass *analysis.Pass, ts *ast.TypeSpec) {
	if isStoreImplementation(ts) {
		// Check error wrapping
		if !usesCustomErrors(pass, ts) {
			pass.Reportf(ts.Pos(), "store implementation %s should use domain-specific errors", ts.Name.Name)
		}

		// Check transaction support
		if !implementsTransactions(pass, ts) {
			pass.Reportf(ts.Pos(), "store implementation %s should support transactions", ts.Name.Name)
		}
	}
}

func checkErrorHandling(pass *analysis.Pass, fn *ast.FuncDecl) {
	ast.Inspect(fn, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			for _, expr := range ret.Results {
				if isErrorReturn(expr) && !usesCustomErrorWrapping(pass, expr) {
					pass.Reportf(expr.Pos(), "errors should be wrapped with domain-specific error types")
				}
			}
		}
		return true
	})
}

func isDomainType(ts *ast.TypeSpec) bool {
	if st, ok := ts.Type.(*ast.StructType); ok {
		// Check if type has typical domain type fields
		hasID := false
		hasTimestamps := false
		if st.Fields != nil {
			for _, field := range st.Fields.List {
				if len(field.Names) > 0 {
					name := field.Names[0].Name
					if name == "ID" {
						hasID = true
					}
					if name == "CreatedAt" || name == "UpdatedAt" {
						hasTimestamps = true
					}
				}
			}
		}
		return hasID && hasTimestamps
	}
	return false
}

func hasValidationMethod(pass *analysis.Pass, typeName string) bool {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name.Name == "Validate" && fn.Recv != nil && len(fn.Recv.List) == 1 {
					if id, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
						if ident, ok := id.X.(*ast.Ident); ok {
							if ident.Name == typeName {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func isStoreImplementation(ts *ast.TypeSpec) bool {
	return strings.HasSuffix(ts.Name.Name, "Store") || strings.HasSuffix(ts.Name.Name, "Repository")
}

func usesCustomErrors(pass *analysis.Pass, ts *ast.TypeSpec) bool {
	var uses bool
	ast.Inspect(ts, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "Error" {
				uses = true
				return false
			}
		}
		return true
	})
	return uses
}

func implementsTransactions(pass *analysis.Pass, ts *ast.TypeSpec) bool {
	var implements bool
	ast.Inspect(ts, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == "BeginTx" {
				implements = true
				return false
			}
		}
		return true
	})
	return implements
}

func isErrorReturn(expr ast.Expr) bool {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name == "error" || strings.HasSuffix(id.Name, "Error")
	}
	return false
}

func usesCustomErrorWrapping(pass *analysis.Pass, expr ast.Expr) bool {
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			return sel.Sel.Name == "Error"
		}
	}
	return false
}
