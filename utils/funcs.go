package utils

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

type CallMatcher interface {
	Matches(*analysis.Pass, *ast.CallExpr) bool
	MatchesSSA(*ssa.CallCommon) bool
}

type PkgFuncCallMatcher struct {
	PkgPath  string
	FuncName string
}

var _ CallMatcher = PkgFuncCallMatcher{}

func NewPkgFuncCallMatcher(pkgPath, funcName string) PkgFuncCallMatcher {
	return PkgFuncCallMatcher{PkgPath: pkgPath, FuncName: funcName}
}

func (m PkgFuncCallMatcher) Matches(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel.Name != m.FuncName {
		return false
	}
	objIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	o := pass.TypesInfo.ObjectOf(objIdent)
	if o == nil {
		return false
	}
	p, ok := o.(*types.PkgName)
	if !ok {
		return false
	}
	return p.Imported().Path() == m.PkgPath
}

func (m PkgFuncCallMatcher) MatchesSSA(common *ssa.CallCommon) bool {
	fn := common.StaticCallee()
	if fn == nil {
		return false
	}
	if fn.Pkg == nil {
		return false
	}
	return fn.Pkg.Pkg.Path() == m.PkgPath && fn.Name() == m.FuncName
}
