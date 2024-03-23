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
	if isAllowedMethodCall(call, pass) {
		return true
	}
	if isAllowedFuncCall(call, pass) {
		return true
	}
	return false
}

func isAllowedMethodCall(call *ast.CallExpr, pass *analysis.Pass) bool {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ty := pass.TypesInfo.TypeOf(fun.X)
	if ty == nil {
		return false
	}
	pt, ok := ty.(*types.Pointer)
	if ok {
		ty = pt.Elem()
	}

	nt, ok := ty.(*types.Named)
	if !ok {
		return false
	}
	obj := nt.Obj()
	if obj.Pkg() == nil {
		return false
	}
	m := method{
		pkg:  obj.Pkg().Path(),
		typ:  nt.Obj().Name(),
		name: fun.Sel.Name,
	}
	_, ok = allowedMethods[m]
	return ok
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

var allowedMethods = map[method]struct{}{
	{"os/exec", "Cmd", "StdinPipe"}:  {},
	{"os/exec", "Cmd", "StdoutPipe"}: {},
	{"os/exec", "Cmd", "StderrPipe"}: {},
}

type method struct {
	pkg  string
	typ  string
	name string
}

var allowedFuns = map[fun]struct{}{}

type fun struct {
	pkg  string
	name string
}
