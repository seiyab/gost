package wraperror

import (
	"go/ast"

	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "wrapError",
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
			if !isWrapLike(pass, call) {
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

var wrapLikes = []utils.PkgFuncCallMatcher{
	utils.NewPkgFuncCallMatcher("github.com/pkg/errors", "WithMessage"),
	utils.NewPkgFuncCallMatcher("github.com/pkg/errors", "WithMessagef"),
	utils.NewPkgFuncCallMatcher("github.com/pkg/errors", "WithStack"),
	utils.NewPkgFuncCallMatcher("github.com/pkg/errors", "Wrap"),
	utils.NewPkgFuncCallMatcher("github.com/pkg/errors", "Wrapf"),
}

func isWrapLike(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, matcher := range wrapLikes {
		if matcher.Matches(pass, call) {
			return true
		}
	}
	return false
}
