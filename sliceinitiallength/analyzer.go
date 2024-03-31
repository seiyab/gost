package sliceinitiallength

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "sliceInitialLength",
	Doc:  "prevents mistakenly specifying initial length of slice",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	muts := collectSliceMutations(pass)
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			a, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			if a.Tok != token.DEFINE {
				return true
			}
			for i, lhs := range a.Lhs {
				if i >= len(a.Rhs) {
					break
				}
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.ObjectOf(ident)
				if obj == nil {
					continue
				}
				mut, ok := muts[obj]
				if !ok || mut != appended {
					continue
				}
				r := a.Rhs[i]
				call, ok := r.(*ast.CallExpr)
				if !ok {
					continue
				}
				fn, ok := call.Fun.(*ast.Ident)
				if !ok {
					continue
				}
				if fn.Name != "make" {
					continue
				}
				if len(call.Args) != 2 {
					continue
				}
				l := call.Args[1]
				if lit, ok := l.(*ast.BasicLit); ok && lit.Value == "0" {
					continue
				}

				pass.Reportf(
					call.Pos(),
					fmt.Sprintf(
						"Slice `%s` is appended to, so typically it should be created with an initial length of 0. Did you mean to use `make([]T, 0, cap)?`",
						ident.Name,
					),
				)
			}
			return true
		})
	}
	return nil, nil
}
