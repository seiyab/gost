package wraperror

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "wraperror",
	Doc:  "detects senseless error wrapping",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if !isWrapLike(sel) {
				return true
			}
			if len(call.Args) < 1 { // can't be
				return true
			}
			argType, ok := pass.TypesInfo.Types[call.Args[0]]
			if !ok { // can't be
				return true
			}
			if argType.IsNil() {
				pass.Reportf(call.Pos(), "wrapping nil doesn't make sense")
				return true
			}
			return true
		})
	}
	return nil, nil
}

func isWrapLike(sel *ast.SelectorExpr) bool {
	patterns := [][]string{
		{"errors", "WithMessage"},
		{"errors", "WithMessagef"},
		{"errors", "WithStack"},
		{"errors", "Wrap"},
		{"errors", "Wrapf"},
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	for _, pattern := range patterns {
		if x.Name == pattern[0] && sel.Sel.Name == pattern[1] {
			return true
		}
	}
	return false
}
