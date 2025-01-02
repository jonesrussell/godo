package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// StorageAnalyzer checks for proper storage implementation patterns
var StorageAnalyzer = &analysis.Analyzer{
	Name: "storagecheck",
	Doc:  "checks for proper storage implementation patterns",
	Run:  runStorageCheck,
}

func runStorageCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				if node.Type != nil {
					if iface, ok := node.Type.(*ast.InterfaceType); ok {
						checkStorageInterface(pass, iface, node)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkStorageInterface(pass *analysis.Pass, iface *ast.InterfaceType, typeSpec *ast.TypeSpec) {
	if isStoreInterface(typeSpec.Name.Name) {
		checkStoreMethods(pass, iface, typeSpec.Pos())
	}
}

func isStoreInterface(name string) bool {
	return strings.HasSuffix(name, "Store")
}

func checkStoreMethods(pass *analysis.Pass, iface *ast.InterfaceType, pos token.Pos) {
	requiredMethods := []string{"BeginTx", "Close"}
	foundMethods := make(map[string]bool)

	for _, method := range iface.Methods.List {
		if len(method.Names) > 0 {
			foundMethods[method.Names[0].Name] = true
		}
	}

	for _, required := range requiredMethods {
		if !foundMethods[required] {
			pass.Reportf(pos, "storage interface must implement %s method", required)
		}
	}
}
