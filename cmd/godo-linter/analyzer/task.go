package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// TaskAnalyzer checks for proper Task struct field usage
var TaskAnalyzer = &analysis.Analyzer{
	Name: "taskcheck",
	Doc:  "checks for proper Task struct field usage and provides auto-fixes",
	Run:  runTaskCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runTaskCheck(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CompositeLit:
				checkTaskLiteral(pass, node)
			case *ast.SelectorExpr:
				checkTaskFieldAccess(pass, node)
			case *ast.CallExpr:
				checkTimeConversion(pass, node)
			}
			return true
		})
	}
	return nil, nil
}

func checkTaskLiteral(pass *analysis.Pass, lit *ast.CompositeLit) {
	if !isTaskType(pass, lit.Type) {
		return
	}

	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				switch ident.Name {
				case "Content":
					suggestRename(pass, ident, "Title")
				case "Done":
					suggestRename(pass, ident, "Completed")
				case "CreatedAt", "UpdatedAt":
					checkTimeFieldAssignment(pass, kv.Value)
				}
			}
		}
	}
}

func checkTaskFieldAccess(pass *analysis.Pass, sel *ast.SelectorExpr) {
	if !isTaskType(pass, sel.X) {
		return
	}

	switch sel.Sel.Name {
	case "Content":
		suggestRename(pass, sel.Sel, "Title")
	case "Done":
		suggestRename(pass, sel.Sel, "Completed")
	case "CreatedAt", "UpdatedAt":
		checkTimeFieldUsage(pass, sel)
	}
}

func checkTimeFieldAssignment(pass *analysis.Pass, expr ast.Expr) {
	t := pass.TypesInfo.Types[expr].Type
	if t == nil {
		return
	}

	if basic, ok := t.(*types.Basic); !ok || basic.Kind() != types.Int64 {
		if call, ok := expr.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == "time" && sel.Sel.Name == "Now" {
						suggestTimeConversion(pass, call)
						return
					}
				}
			}
		}
		pass.Reportf(expr.Pos(), "time fields must be int64 (Unix timestamp), use time.Now().Unix() instead of time.Now()")
	}
}

func suggestRename(pass *analysis.Pass, node *ast.Ident, newName string) {
	pass.Report(analysis.Diagnostic{
		Pos:     node.Pos(),
		Message: fmt.Sprintf("use %s instead of %s in Task struct", newName, node.Name),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("Rename to %s", newName),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     node.Pos(),
						End:     node.End(),
						NewText: []byte(newName),
					},
				},
			},
		},
	})
}

func suggestTimeConversion(pass *analysis.Pass, call *ast.CallExpr) {
	pass.Report(analysis.Diagnostic{
		Pos:     call.Pos(),
		Message: "use time.Now().Unix() for Task time fields",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Convert to Unix timestamp",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     call.End(),
						End:     call.End(),
						NewText: []byte(".Unix()"),
					},
				},
			},
		},
	})
}

func checkTimeFieldUsage(pass *analysis.Pass, sel *ast.SelectorExpr) {
	parent := findParentCallExpr(pass, sel)
	if parent == nil {
		return
	}

	if isAssertWithinDuration(parent) {
		suggestTimeUnixConversion(pass, sel)
	}
}

func suggestTimeUnixConversion(pass *analysis.Pass, sel *ast.SelectorExpr) {
	pass.Report(analysis.Diagnostic{
		Pos:     sel.Pos(),
		Message: "convert Unix timestamp to time.Time using time.Unix() before using assert.WithinDuration",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Convert to time.Time",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     sel.Pos(),
						End:     sel.End(),
						NewText: []byte(fmt.Sprintf("time.Unix(%s, 0)", pass.Fset.Position(sel.Pos()).String())),
					},
				},
			},
		},
	})
}

func checkTimeConversion(pass *analysis.Pass, call *ast.CallExpr) {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if ident.Name == "time" && sel.Sel.Name == "Now" {
				parent := findParentKeyValueExpr(pass, call)
				if parent != nil && isTaskTimeField(parent) {
					pass.Reportf(call.Pos(), "use time.Now().Unix() for Task time fields")
				}
			}
		}
	}
}

func isTaskType(pass *analysis.Pass, expr ast.Expr) bool {
	if t := pass.TypesInfo.Types[expr].Type; t != nil {
		return strings.HasSuffix(t.String(), "storage.Task")
	}
	return false
}

func isTaskTimeField(kv *ast.KeyValueExpr) bool {
	if ident, ok := kv.Key.(*ast.Ident); ok {
		return ident.Name == "CreatedAt" || ident.Name == "UpdatedAt"
	}
	return false
}

func isAssertWithinDuration(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "assert" && sel.Sel.Name == "WithinDuration"
		}
	}
	return false
}

func findParentCallExpr(pass *analysis.Pass, node ast.Node) *ast.CallExpr {
	var parent *ast.CallExpr
	ast.Inspect(pass.Files[0], func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			for _, arg := range call.Args {
				if sameNode(arg, node) {
					parent = call
					return false
				}
			}
		}
		return true
	})
	return parent
}

func findParentKeyValueExpr(pass *analysis.Pass, node ast.Node) *ast.KeyValueExpr {
	var parent *ast.KeyValueExpr
	ast.Inspect(pass.Files[0], func(n ast.Node) bool {
		if kv, ok := n.(*ast.KeyValueExpr); ok {
			if sameNode(kv.Value, node) {
				parent = kv
				return false
			}
		}
		return true
	})
	return parent
}

func sameNode(n1, n2 ast.Node) bool {
	return n1.Pos() == n2.Pos() && n1.End() == n2.End()
}
