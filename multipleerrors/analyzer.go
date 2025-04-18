package multipleerrors

import (
	"go/ast"
	"go/token"

	"github.com/seiyab/gost/astpath"
	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "multipleErrors",
	Doc:  "detects senseless / uncommon use of error concatenations",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, astpath.WithPath(func(n ast.Node, path *astpath.Path) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isJoinLike(pass, call) {
				return true
			}
			if len(call.Args) < 2 && call.Ellipsis == token.NoPos {
				pass.Reportf(call.Pos(), "too few arguments")
				return true
			}
			if _, ok := path.Parent.Current.(*ast.ExprStmt); ok {
				pass.Reportf(call.Pos(), "result is not used")
				return true
			}
			if originalErrorIsDiscarded(call, path) {
				pass.Reportf(call.Pos(), "original error is discarded")
				return true
			}

			return true
		}))
	}
	return nil, nil
}

var joinLikes = []utils.PkgFuncCallMatcher{
	utils.NewPkgFuncCallMatcher("errors", "Join"),
	utils.NewPkgFuncCallMatcher("github.com/hashicorp/go-multierror", "Append"),
	utils.NewPkgFuncCallMatcher("github.com/uber-go/multierr ", "Append"),
}

func isJoinLike(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, matcher := range joinLikes {
		if matcher.Matches(pass, call) {
			return true
		}
	}
	return false
}

func originalErrorIsDiscarded(call *ast.CallExpr, path *astpath.Path) bool {
	asmt, ok := path.Parent.Current.(*ast.AssignStmt)
	if !ok || asmt.Tok != token.ASSIGN {
		return false
	}
	if len(asmt.Lhs) != 1 { // can't be
		return false
	}
	left, ok := asmt.Lhs[0].(*ast.Ident)
	if !ok { // not implemented yet
		return false
	}

	for _, arg := range call.Args {
		id, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}
		if id.Name == left.Name {
			return false
		}
	}

	return true
}
