package analyzer

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// ArchitectureAnalyzer enforces Godo's architectural patterns
var ArchitectureAnalyzer = &analysis.Analyzer{
	Name: "archcheck",
	Doc:  "enforces Godo's architectural patterns and package structure",
	Run:  runArchitectureCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runArchitectureCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filePath := pass.Fset.Position(file.Pos()).Filename
		dir := filepath.Dir(filePath)

		// Check package location rules
		checkPackageLocation(pass, file, dir)

		// Check dependencies
		checkDependencyRules(pass, file, dir)

		// Check implementation patterns
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				checkTypePatterns(pass, node, dir)
			case *ast.FuncDecl:
				checkLayerPatterns(pass, node, dir)
			}
			return true
		})
	}
	return nil, nil
}

func checkPackageLocation(pass *analysis.Pass, file *ast.File, dir string) {
	// Enforce package structure rules
	if strings.Contains(dir, "internal/api") {
		// API layer rules
		if strings.Contains(dir, "internal/api/handler") {
			checkHandlerPackage(pass, file)
		} else if strings.Contains(dir, "internal/api/middleware") {
			checkMiddlewarePackage(pass, file)
		}
	} else if strings.Contains(dir, "internal/storage") {
		// Storage layer rules
		checkStoragePackage(pass, file)
	} else if strings.Contains(dir, "internal/service") {
		// Service layer rules
		checkServicePackage(pass, file)
	}
}

func checkHandlerPackage(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// Handlers should not import storage directly
			if hasStorageImport(file) {
				pass.Reportf(fn.Pos(), "handlers should not import storage directly, use service layer instead")
			}
			// Handlers should use service interfaces
			if !usesServiceInterface(fn) {
				pass.Reportf(fn.Pos(), "handlers should depend on service interfaces")
			}
		}
	}
}

func checkMiddlewarePackage(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// Middleware should be stateless
			if hasStateFields(fn) {
				pass.Reportf(fn.Pos(), "middleware should be stateless")
			}
			// Middleware should implement standard interface
			if !implementsMiddlewareInterface(fn) {
				pass.Reportf(fn.Pos(), "middleware must implement standard middleware interface")
			}
		}
	}
}

func checkStoragePackage(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			if node.Tok.String() == "type" {
				for _, spec := range node.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						// Storage types should implement Store interface
						if !implementsStoreInterface(ts) {
							pass.Reportf(ts.Pos(), "storage types must implement Store interface")
						}
						// Storage should use transactions
						if !usesTransactions(ts) {
							pass.Reportf(ts.Pos(), "storage operations should use transactions")
						}
					}
				}
			}
		}
	}
}

func checkServicePackage(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			if node.Tok.String() == "type" {
				for _, spec := range node.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						// Services should use storage interfaces
						if !usesStorageInterface(ts) {
							pass.Reportf(ts.Pos(), "services should depend on storage interfaces")
						}
						// Services should not expose storage types
						if exposesStorageTypes(ts) {
							pass.Reportf(ts.Pos(), "services should not expose storage implementation details")
						}
					}
				}
			}
		}
	}
}

func checkDependencyRules(pass *analysis.Pass, file *ast.File, dir string) {
	imports := make(map[string]bool)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		imports[path] = true
	}

	// Enforce dependency rules based on clean architecture
	if strings.Contains(dir, "internal/api") {
		// API layer can only depend on service interfaces
		for path := range imports {
			if strings.Contains(path, "internal/storage") {
				pass.Reportf(file.Pos(), "api layer cannot depend directly on storage layer")
			}
		}
	} else if strings.Contains(dir, "internal/service") {
		// Service layer can depend on storage interfaces
		for path := range imports {
			if strings.Contains(path, "internal/api") {
				pass.Reportf(file.Pos(), "service layer cannot depend on api layer")
			}
		}
	}
}

func checkTypePatterns(pass *analysis.Pass, typeSpec *ast.TypeSpec, dir string) {
	// Check if types follow our patterns based on their layer
	if strings.Contains(dir, "internal/api") {
		if st, ok := typeSpec.Type.(*ast.StructType); ok {
			// API types should not embed storage types
			for _, field := range st.Fields.List {
				if sel, ok := field.Type.(*ast.SelectorExpr); ok {
					if strings.Contains(sel.X.(*ast.Ident).Name, "storage") {
						pass.Reportf(field.Pos(), "api types should not embed storage types")
					}
				}
			}
		}
	}
}

func checkLayerPatterns(pass *analysis.Pass, fn *ast.FuncDecl, dir string) {
	// Check layer-specific patterns
	if strings.Contains(dir, "internal/api") {
		// API layer should use DTOs
		if !usesDTOs(fn) {
			pass.Reportf(fn.Pos(), "api layer should use DTOs for request/response")
		}
	} else if strings.Contains(dir, "internal/service") {
		// Service layer should handle business logic
		if !hasBusinessLogic(fn) {
			pass.Reportf(fn.Pos(), "service layer should contain business logic")
		}
	}
}

// Helper functions

func hasStorageImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if strings.Contains(imp.Path.Value, "storage") {
			return true
		}
	}
	return false
}

func usesServiceInterface(fn *ast.FuncDecl) bool {
	var uses bool
	ast.Inspect(fn, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if strings.HasSuffix(sel.Sel.Name, "Service") {
				uses = true
				return false
			}
		}
		return true
	})
	return uses
}

func hasStateFields(fn *ast.FuncDecl) bool {
	if fn.Recv != nil {
		if st, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
			if id, ok := st.X.(*ast.Ident); ok {
				return id.Obj != nil && id.Obj.Kind == ast.Var
			}
		}
	}
	return false
}

func implementsMiddlewareInterface(fn *ast.FuncDecl) bool {
	return fn.Type != nil && fn.Type.Results != nil &&
		len(fn.Type.Results.List) == 1 &&
		isHandlerType(fn.Type.Results.List[0].Type)
}

func implementsStoreInterface(ts *ast.TypeSpec) bool {
	if _, ok := ts.Type.(*ast.InterfaceType); ok {
		return strings.HasSuffix(ts.Name.Name, "Store")
	}
	return false
}

func usesTransactions(ts *ast.TypeSpec) bool {
	var uses bool
	ast.Inspect(ts, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "BeginTx" {
				uses = true
				return false
			}
		}
		return true
	})
	return uses
}

func usesStorageInterface(ts *ast.TypeSpec) bool {
	if st, ok := ts.Type.(*ast.StructType); ok {
		for _, field := range st.Fields.List {
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				if strings.HasSuffix(sel.Sel.Name, "Store") {
					return true
				}
			}
		}
	}
	return false
}

func exposesStorageTypes(ts *ast.TypeSpec) bool {
	if st, ok := ts.Type.(*ast.StructType); ok {
		for _, field := range st.Fields.List {
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				if x, ok := sel.X.(*ast.Ident); ok {
					if x.Name == "storage" {
						return true
					}
				}
			}
		}
	}
	return false
}

func usesDTOs(fn *ast.FuncDecl) bool {
	var uses bool
	ast.Inspect(fn, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if strings.HasSuffix(ts.Name.Name, "Request") ||
				strings.HasSuffix(ts.Name.Name, "Response") {
				uses = true
				return false
			}
		}
		return true
	})
	return uses
}

func hasBusinessLogic(fn *ast.FuncDecl) bool {
	// This is a simplified check. In reality, we'd look for more specific patterns
	var hasLogic bool
	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.SwitchStmt, *ast.ForStmt:
			hasLogic = true
			return false
		}
		return true
	})
	return hasLogic
}

func isHandlerType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "HandlerFunc" || sel.Sel.Name == "Handler"
	}
	return false
}
