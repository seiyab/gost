package nomutateglobal

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noMutateGlobal",
	Doc:  "prevents mutation of global variables",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	c := newGlobalCollector(pass)
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			asmt, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			for _, lhs := range asmt.Lhs {
				sel, ok := lhs.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				if !c.contains(sel.X) {
					continue
				}
				pass.Reportf(
					asmt.Pos(),
					"This assignment might mutate a global variable. Lhs can be a pointer to a global variable.",
				)
			}
			c.visitAssignment(asmt)
			return true
		})
	}
	return nil, nil
}

func isGlobalVar(expr ast.Expr, pass *analysis.Pass) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	o := pass.TypesInfo.ObjectOf(ident)
	if o == nil {
		return false
	}
	_, ok = o.(*types.PkgName)
	return ok
}
