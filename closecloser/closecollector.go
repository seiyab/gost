package closecloser

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type closeCollector struct {
	pass      *analysis.Pass
	calls     map[types.Object]struct{}
	inspector *inspector.Inspector
}

func newCloseCollector(pass *analysis.Pass, ipr *inspector.Inspector) *closeCollector {
	c := &closeCollector{
		pass:      pass,
		calls:     make(map[types.Object]struct{}),
		inspector: ipr,
	}
	c.trace()
	return c
}

func (c *closeCollector) trace() {
	c.inspector.WithStack(
		[]ast.Node{
			&ast.Ident{},
		},
		func(n ast.Node, push bool, stack []ast.Node) bool {
			if !push {
				return true
			}
			switch n := n.(type) {
			case *ast.Ident:
				if mightClose(n, stack) {
					c.storeIdent(n)
				}
			}
			return true
		},
	)
}

func mightClose(n *ast.Ident, stack []ast.Node) bool {
	if len(stack) < 2 {
		return false
	}
	switch parent := stack[len(stack)-2].(type) {
	case *ast.SelectorExpr:
		return parent.Sel.Name == "Close"
	case *ast.ValueSpec:
		return false
	case *ast.AssignStmt:
		for _, ex := range parent.Lhs {
			if id, ok := ex.(*ast.Ident); ok && id == n {
				return false
			}
		}
		return true
	}
	return true
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
