package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// InterfaceLocationAnalyzer enforces that interfaces must be defined in specific files
var InterfaceLocationAnalyzer = &analysis.Analyzer{
	Name: "interfacelocation",
	Doc:  "Checks that interfaces are defined in the correct files",
	Run:  runInterfaceLocationCheck,
}

func runInterfaceLocationCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()

		// Skip interfaces.go files as they are the correct location
		if strings.HasSuffix(fileName, "interfaces.go") {
			continue
		}

		// Check each declaration in the file
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if it's an interface declaration
				if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
					pass.Reportf(typeSpec.Pos(), "interface %q must be defined in interfaces.go", typeSpec.Name.Name)
				}
			}
		}
	}

	return nil, nil
}
