package multipleerrors

import (
	"go/ast"
	"go/token"

	"github.com/seiyab/gost/gost/astpath"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "multipleerrors",
	Doc:  "detects helpless / uncommon use of error concatenations",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, astpath.WithPath(func(n ast.Node, path *astpath.Path) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if !isJoinLike(sel) {
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

func isJoinLike(sel *ast.SelectorExpr) bool {
	var patterns = [][]string{
		{"errors", "Join"},
		{"multierror", "Append"},
		{"multierr", "Append"},
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
