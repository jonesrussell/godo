package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	maxStructFields = 10
)

// StandardAnalyzer enforces consistent code patterns and clean code practices
var StandardAnalyzer = &analysis.Analyzer{
	Name: "standardcheck",
	Doc:  "enforces consistent code patterns and clean code practices",
	Run:  runStandardCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runStandardCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				checkFunctionStandards(pass, node)
			case *ast.GenDecl:
				checkDeclarationStandards(pass, node)
			case *ast.File:
				checkFileStandards(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func checkFunctionStandards(pass *analysis.Pass, fn *ast.FuncDecl) {
	// Check function length (max 50 lines)
	if fn.Body != nil && len(fn.Body.List) > 50 {
		pass.Reportf(fn.Pos(), "function %s is too long (> 50 lines), consider breaking it down", fn.Name.Name)
	}

	// Check parameter count (max 5)
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 5 {
		pass.Reportf(fn.Pos(), "function %s has too many parameters (> 5), consider using a config struct", fn.Name.Name)
	}

	// Check return values (max 3)
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 3 {
		pass.Reportf(fn.Pos(), "function %s has too many return values (> 3), consider using a result struct", fn.Name.Name)
	}

	// Check function naming conventions
	checkFunctionNaming(pass, fn)
}

func checkFunctionNaming(pass *analysis.Pass, fn *ast.FuncDecl) {
	name := fn.Name.Name

	// Method naming conventions
	if fn.Recv != nil {
		if strings.HasPrefix(name, "Get") && !returnsValue(fn) {
			pass.Reportf(fn.Pos(), "getter method %s should return a value", name)
		}
		if strings.HasPrefix(name, "Set") && len(fn.Type.Params.List) == 0 {
			pass.Reportf(fn.Pos(), "setter method %s should take a parameter", name)
		}
	}

	// Handler naming conventions
	if strings.HasSuffix(name, "Handler") && !hasResponseWriterParam(fn) {
		pass.Reportf(fn.Pos(), "handler function %s should take http.ResponseWriter and *http.Request parameters", name)
	}
}

func checkDeclarationStandards(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			checkTypeStandards(pass, s)
		case *ast.ValueSpec:
			checkConstantStandards(pass, s, decl.Tok.String() == "const")
		}
	}
}

func checkTypeStandards(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	// Check struct standards
	if st, ok := typeSpec.Type.(*ast.StructType); ok {
		checkStructStandards(pass, typeSpec.Name.Name, st)
	}

	// Check interface standards
	if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		checkInterfaceStandards(pass, typeSpec.Name.Name, iface)
	}
}

func checkStructStandards(pass *analysis.Pass, name string, st *ast.StructType) {
	// Check struct field count
	if len(st.Fields.List) > maxStructFields {
		pass.Reportf(st.Pos(), "struct %s has too many fields (> %d), consider breaking it down", name, maxStructFields)
	}

	// Check struct field ordering (exported fields first)
	var lastExported bool
	for i, field := range st.Fields.List {
		if len(field.Names) > 0 {
			isExported := field.Names[0].IsExported()
			if i > 0 && !lastExported && isExported {
				pass.Reportf(field.Pos(), "exported fields should be declared before unexported fields")
			}
			lastExported = isExported
		}
	}
}

func checkInterfaceStandards(pass *analysis.Pass, name string, iface *ast.InterfaceType) {
	// Check interface method count
	if len(iface.Methods.List) > 5 {
		pass.Reportf(iface.Pos(), "interface %s has too many methods (> 5), consider splitting it", name)
	}

	// Check interface naming
	if !strings.HasSuffix(name, "er") && !strings.HasSuffix(name, "Service") {
		pass.Reportf(iface.Pos(), "interface %s should end with 'er' or 'Service'", name)
	}
}

func checkConstantStandards(pass *analysis.Pass, val *ast.ValueSpec, isConst bool) {
	for _, name := range val.Names {
		// Constants should be ALL_CAPS
		if isConst && !isAllCaps(name.Name) {
			pass.Reportf(name.Pos(), "constant %s should be ALL_CAPS", name.Name)
		}

		// Variables should be camelCase
		if !isConst && !isCamelCase(name.Name) {
			pass.Reportf(name.Pos(), "variable %s should be camelCase", name.Name)
		}
	}
}

func checkFileStandards(pass *analysis.Pass, file *ast.File) {
	// Check package naming (should be single lowercase word)
	if !isValidPackageName(file.Name.Name) {
		pass.Reportf(file.Name.Pos(), "package name %s should be a single lowercase word", file.Name.Name)
	}

	// Check file imports
	checkImportStandards(pass, file)
}

func checkImportStandards(pass *analysis.Pass, file *ast.File) {
	var stdlibCount, thirdPartyCount, internalCount int

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		switch {
		case !strings.Contains(path, "."):
			stdlibCount++
		case strings.Contains(path, "github.com/jonesrussell/godo"):
			internalCount++
		default:
			thirdPartyCount++
		}
	}

	// Check import grouping
	if stdlibCount > 0 && thirdPartyCount > 0 && internalCount > 0 {
		var lastGroup string
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, "\"")
			var currentGroup string
			switch {
			case !strings.Contains(path, "."):
				currentGroup = "stdlib"
			case strings.Contains(path, "github.com/jonesrussell/godo"):
				currentGroup = "internal"
			default:
				currentGroup = "third-party"
			}

			if lastGroup != "" && currentGroup != lastGroup {
				if (lastGroup == "stdlib" && currentGroup != "third-party") ||
					(lastGroup == "third-party" && currentGroup != "internal") {
					pass.Reportf(imp.Pos(), "imports should be grouped: stdlib > third-party > internal")
				}
			}
			lastGroup = currentGroup
		}
	}
}

// Helper functions

func returnsValue(fn *ast.FuncDecl) bool {
	return fn.Type.Results != nil && len(fn.Type.Results.List) > 0
}

func hasResponseWriterParam(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil {
		return false
	}
	for _, param := range fn.Type.Params.List {
		if sel, ok := param.Type.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "ResponseWriter" {
				return true
			}
		}
	}
	return false
}

func isAllCaps(name string) bool {
	return strings.ToUpper(name) == name && !strings.Contains(name, " ")
}

func isCamelCase(name string) bool {
	return !strings.Contains(name, "_") && !strings.Contains(name, " ")
}

func isValidPackageName(name string) bool {
	return !strings.Contains(name, "_") && strings.ToLower(name) == name
}
