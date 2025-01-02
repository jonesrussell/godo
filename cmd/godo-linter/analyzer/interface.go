package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// InterfaceAnalyzer checks for proper interface implementations
var InterfaceAnalyzer = &analysis.Analyzer{
	Name: "interfacecheck",
	Doc:  "checks for proper interface implementations",
	Run:  runInterfaceCheck,
}

func runInterfaceCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				checkInterfaceImplementation(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func checkInterfaceImplementation(pass *analysis.Pass, spec *ast.TypeSpec) {
	// Get type information
	typeInfo := pass.TypesInfo
	if typeInfo == nil {
		return
	}

	// Check if type implements storage.Store
	if implementsInterface(pass, spec, "storage.Store") {
		checkRequiredMethods(pass, spec, []string{
			"BeginTx",
			"CreateTask",
			"GetTask",
			"ListTasks",
			"UpdateTask",
			"DeleteTask",
			"Close",
		})
	}

	// Check if type implements gui.MainWindowManager
	if implementsInterface(pass, spec, "gui.MainWindowManager") {
		checkRequiredMethods(pass, spec, []string{
			"GetWindow",
			"Show",
			"Hide",
			"Close",
		})
	}
}

func implementsInterface(pass *analysis.Pass, spec *ast.TypeSpec, interfaceName string) bool {
	typeObj := pass.TypesInfo.Defs[spec.Name]
	if typeObj == nil {
		return false
	}

	// Check if type is named and get its methods
	named, ok := typeObj.Type().(*types.Named)
	if !ok {
		return false
	}

	// Get method set
	methodSet := types.NewMethodSet(types.NewPointer(named))
	for i := 0; i < methodSet.Len(); i++ {
		method := methodSet.At(i)
		if method.Obj().Name() == "BeginTx" && interfaceName == "storage.Store" {
			return true
		}
		if method.Obj().Name() == "GetWindow" && interfaceName == "gui.MainWindowManager" {
			return true
		}
	}
	return false
}

func checkRequiredMethods(pass *analysis.Pass, spec *ast.TypeSpec, methods []string) {
	typeObj := pass.TypesInfo.Defs[spec.Name]
	if typeObj == nil {
		return
	}

	named, ok := typeObj.Type().(*types.Named)
	if !ok {
		return
	}

	methodSet := types.NewMethodSet(types.NewPointer(named))
	for _, required := range methods {
		found := false
		for i := 0; i < methodSet.Len(); i++ {
			if methodSet.At(i).Obj().Name() == required {
				found = true
				break
			}
		}
		if !found {
			pass.Reportf(spec.Pos(), "type %s must implement method %s", spec.Name.Name, required)
		}
	}
}
