package nodiscarderror

import (
	"go/ast"
	"go/token"

	"github.com/seiyab/gost/astpath"
	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noDiscardError",
	Doc:  "prevents `if err != nil { return nil }`",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	ti := pass.TypesInfo
	for _, file := range pass.Files {
		ast.Inspect(file, astpath.WithPath(func(n ast.Node, path *astpath.Path) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
			if !ok {
				return true
			}
			if binExpr.Op != token.NEQ {
				return true
			}
			if !utils.IsError(ti.TypeOf(binExpr.X)) {
				return true
			}
			y, ok := ti.Types[binExpr.Y]
			if !ok || !y.IsNil() {
				return true
			}
			list := ifStmt.Body.List
			if len(list) != 1 {
				return true
			}
			ret, ok := list[0].(*ast.ReturnStmt)
			if !ok {
				return true
			}
			fn, _ := astpath.FindNearest[*ast.FuncDecl](path)
			if fn == nil {
				return true
			}
			results := fn.Type.Results
			if results == nil {
				return true
			}
			retTypes := results.List
			if len(retTypes) != len(ret.Results) {
				return true
			}
			for i, result := range ret.Results {
				if !utils.IsError(ti.TypeOf(retTypes[i].Type)) {
					continue
				}
				res, ok := ti.Types[result]
				if !ok {
					continue
				}
				if res.IsNil() {
					pass.Reportf(ret.Pos(), "error is discarded")
					return true
				}
			}
			return true
		}))
	}
	return nil, nil
}
