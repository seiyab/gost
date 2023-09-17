package nodiscarderror

import (
	"go/ast"
	"go/token"

	"github.com/seiyab/gost/gost/astpath"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nodiscarderror",
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
			if ti.TypeOf(binExpr.X).String() != "error" {
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
				t := ti.TypeOf(retTypes[i].Type)
				if t.String() != "error" {
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
