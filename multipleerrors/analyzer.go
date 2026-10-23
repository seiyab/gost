package multipleerrors

import (
	"go/ast"
	"go/token"
	"go/types"

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
			if originalErrorIsDiscarded(pass, call, path) {
				pass.Reportf(call.Pos(), "original error is discarded")
			}
			return true
		}))
	}
	return nil, nil
}

func originalErrorIsDiscarded(pass *analysis.Pass, call *ast.CallExpr, path *astpath.Path) bool {
	assignment, ok := path.Parent.Current.(*ast.AssignStmt)
	if !ok || assignment.Tok != token.ASSIGN || len(assignment.Lhs) != 1 {
		return false
	}
	left, ok := assignment.Lhs[0].(*ast.Ident)
	if !ok {
		return false
	}
	object := pass.TypesInfo.ObjectOf(left)
	isOriginal := func(id *ast.Ident) bool { return pass.TypesInfo.ObjectOf(id) == object }
	for _, arg := range call.Args {
		if id, ok := ast.Unparen(arg).(*ast.Ident); ok && isOriginal(id) {
			return false
		}
	}

	// A preceding guard may already have handled this error. Only inspect
	// siblings in the same block, and stop at writes to avoid treating a guard
	// for an older value as handling the current one.
	block, ok := path.Parent.Parent.Current.(*ast.BlockStmt)
	if !ok {
		return true
	}
	index := 0
	for index < len(block.List) && block.List[index] != assignment {
		index++
	}
	for index--; index >= 0; index-- {
		stmt := block.List[index]
		writes := false
		ast.Inspect(stmt, func(n ast.Node) bool {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			if a, ok := n.(*ast.AssignStmt); ok {
				for _, lhs := range a.Lhs {
					if id, ok := lhs.(*ast.Ident); ok && isOriginal(id) {
						writes = true
					}
				}
			}
			return true
		})
		if writes {
			break
		}
		guard, ok := stmt.(*ast.IfStmt)
		if !ok {
			continue
		}
		checked := false
		ast.Inspect(guard.Cond, func(n ast.Node) bool {
			comparison, ok := n.(*ast.BinaryExpr)
			if !ok || comparison.Op != token.NEQ {
				return true
			}
			for _, pair := range [][2]ast.Expr{{comparison.X, comparison.Y}, {comparison.Y, comparison.X}} {
				id, ok := ast.Unparen(pair[0]).(*ast.Ident)
				nilID, nilOK := ast.Unparen(pair[1]).(*ast.Ident)
				if ok && isOriginal(id) && nilOK && pass.TypesInfo.ObjectOf(nilID) == types.Universe.Lookup("nil") {
					checked = true
				}
			}
			return true
		})
		if !checked || len(guard.Body.List) == 0 {
			continue
		}
		if ret, ok := guard.Body.List[len(guard.Body.List)-1].(*ast.ReturnStmt); ok {
			for _, result := range ret.Results {
				if id, ok := ast.Unparen(result).(*ast.Ident); ok && isOriginal(id) {
					return false
				}
			}
		}
	}
	return true
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
