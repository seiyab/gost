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
	switch node := node.(type) {
	case *ast.CallExpr:
		c.traceCloseCall(node)
		c.traceCloserArg(node)
	}
	return c
}

func (c *closeCollector) traceCloseCall(call *ast.CallExpr) {
	if len(call.Args) != 0 {
		return
	}
	ti, ok := c.pass.TypesInfo.Types[call]
	if !ok || ti.Type.String() != "error" {
		return
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Close" {
		return
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}
	c.storeIdent(x)
}

func (c *closeCollector) traceCloserArg(call *ast.CallExpr) {
	ftt := c.pass.TypesInfo.TypeOf(call.Fun)
	sig, ok := ftt.Underlying().(*types.Signature)
	if !ok {
		return
	}
	ps := sig.Params()
	for i, a := range call.Args {
		if i >= ps.Len() {
			break
		}
		p := ps.At(i)
		if !implementsCloser(p.Type()) {
			continue
		}

		id, ok := a.(*ast.Ident)
		if !ok {
			continue
		}
		c.storeIdent(id)

	}
}

func (c *closeCollector) storeIdent(x *ast.Ident) {
	o := c.pass.TypesInfo.ObjectOf(x)
	if o == nil {
		return
	}
	c.calls[o] = struct{}{}
}

func (c *closeCollector) isClosed(o types.Object) bool {
	_, ok := c.calls[o]
	return ok
}
