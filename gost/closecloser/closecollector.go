package closecloser

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type closeCollector struct {
	pass  *analysis.Pass
	calls map[types.Object]struct{}
}

func newCloseCollector(pass *analysis.Pass) *closeCollector {
	return &closeCollector{
		pass:  pass,
		calls: make(map[types.Object]struct{}),
	}
}

func (c *closeCollector) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok || len(call.Args) != 0 {
		return c
	}
	ti, ok := c.pass.TypesInfo.Types[call]
	if !ok || ti.Type.String() != "error" {
		return c
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Close" {
		return c
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return c
	}
	o := c.pass.TypesInfo.ObjectOf(x)
	if o == nil {
		return c
	}
	c.calls[o] = struct{}{}
	return c
}

func (c *closeCollector) isClosed(o types.Object) bool {
	_, ok := c.calls[o]
	return ok
}
