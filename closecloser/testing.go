package closecloser

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func isTestOrTestHelper(n *ast.FuncDecl, pass *analysis.Pass) bool {
	ty := n.Type
	if ty == nil || ty.Params == nil || len(ty.Params.List) == 0 {
		return false
	}
	t := ty.Params.List[0].Type
	return isTestingTB(t, pass)
}

func isTestingTB(t ast.Expr, pass *analysis.Pass) bool {
	s, ok := t.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := s.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	pkgObj := pass.TypesInfo.ObjectOf(pkg)
	if pkgObj == nil {
		return false
	}
	pn, ok := pkgObj.(*types.PkgName)
	if !ok {
		return false
	}
	if pn.Imported().Path() != "testing" {
		return false
	}

	name := sel.Sel.Name
	return name == "T" || name == "TB" || name == "B"
}
