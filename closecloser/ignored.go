package closecloser

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func shouldAllowUnclosed(expr ast.Expr, pass *analysis.Pass) bool {
	switch expr := expr.(type) {
	case *ast.SelectorExpr, *ast.Ident:
		return true
	case *ast.CallExpr:
		return isAllowedCall(expr, pass)
	case *ast.TypeAssertExpr:
		return shouldAllowUnclosed(expr.X, pass)
	}
	return false
}

func isAllowedCall(call *ast.CallExpr, pass *analysis.Pass) bool {
	if isMethodCall(call, pass) {
		// closer returned by method call is allowed.
		// its lifecycle might be managed by the receiver.
		return true
	}
	if isAllowedFuncCall(call, pass) {
		return true
	}
	return false
}

func isMethodCall(call *ast.CallExpr, pass *analysis.Pass) bool {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := fun.X.(*ast.Ident)
	if !ok {
		return true
	}
	ty := pass.TypesInfo.ObjectOf(ident)
	if ty == nil {
		return false
	}

	_, ok = ty.(*types.PkgName)
	return !ok
}

func isAllowedFuncCall(call *ast.CallExpr, pass *analysis.Pass) bool {
	fn, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := fn.X.(*ast.Ident)
	if !ok {
		return false
	}
	obj, ok := pass.TypesInfo.ObjectOf(x).(*types.PkgName)
	if !ok {
		return false
	}
	f := fun{
		pkg:  obj.Imported().Path(),
		name: fn.Sel.Name,
	}
	_, ok = allowedFuns[f]
	return ok
}

var allowedFuns = map[fun]struct{}{}

type fun struct {
	pkg  string
	name string
}
