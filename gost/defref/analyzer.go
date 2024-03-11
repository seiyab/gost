package defref

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var Analyzer = &analysis.Analyzer{
	Name:       "defref",
	Doc:        "provides definition / reference information",
	Run:        run,
	Requires:   []*analysis.Analyzer{buildssa.Analyzer},
	ResultType: reflect.TypeOf(new(DefRef)),
}

func run(pass *analysis.Pass) (any, error) {
	return typedRun(pass)
}

func typedRun(pass *analysis.Pass) (*DefRef, error) {
	ti := pass.TypesInfo
	var defRef *DefRef = &DefRef{
		funcs: make(map[ast.Node]Func),
		// blocks: make([]Block, 0),
	}

	r, ok := pass.ResultOf[buildssa.Analyzer]
	if !ok {
		return nil, errors.New("no result of buildssa")
	}
	ssaPkg, ok := r.(*buildssa.SSA)
	if !ok {
		return nil, errors.New("result of buildssa is not *buildssa.SSA")
	}

	for _, f := range ssaPkg.SrcFuncs {
		defRef.trackFunc(f)
	}

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.Ident:
				fmt.Println(
					n.Name,
					ti.Defs[n],
					ti.Uses[n],
				)
			}

			return true
		})
	}

	return defRef, nil
}
