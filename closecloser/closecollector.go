package closecloser

import (
	"go/ast"
	"go/types"

	"github.com/seiyab/gost/astpath"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type closeCollector struct {
	pass  *analysis.Pass
	calls map[types.Object]struct{}
}

func newCloseCollector(pass *analysis.Pass) *closeCollector {
	c := &closeCollector{
		pass:  pass,
		calls: make(map[types.Object]struct{}),
	}
	c.trace()
	return c
}

func (c *closeCollector) trace() {
	ipr := inspector.New(c.pass.Files)
	ipr.WithStack(
		[]ast.Node{
			&ast.CallExpr{},
			&ast.CompositeLit{},
			&ast.AssignStmt{},
			&ast.ReturnStmt{},
		},
		func(n ast.Node, push bool, stack []ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				c.traceCloseCall(n)
				c.traceCloserArg(n)
			case *ast.CompositeLit:
				c.traceCompositeLiteral(n)
			case *ast.AssignStmt:
				c.traceAssignment(n)
			case *ast.ReturnStmt:
				c.traceReturn(n, stack)
			}

			return true
		},
	)
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
	if ftt == nil {
		return
	}
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

func (c *closeCollector) traceCompositeLiteral(lit *ast.CompositeLit) {
	lt := c.pass.TypesInfo.TypeOf(lit)
	if lt != nil {
		lt = lt.Underlying()
	}
	for i, e := range lit.Elts {
		switch e := e.(type) {
		case *ast.KeyValueExpr:
			id, ok := e.Value.(*ast.Ident)
			if !ok {
				continue
			}
			kt := c.pass.TypesInfo.TypeOf(e.Key)
			if kt == nil || !implementsCloser(kt) {
				continue
			}
			c.storeIdent(id)
		case *ast.Ident:
			switch lt := lt.(type) {
			case *types.Struct:
				if i >= lt.NumFields() {
					continue
				}
				if !implementsCloser(lt.Field(i).Type()) {
					continue
				}
				c.storeIdent(e)
			}
		default:
			continue
		}
	}
}

func (c *closeCollector) traceAssignment(asm *ast.AssignStmt) {
	for i, rh := range asm.Rhs {
		id, ok := rh.(*ast.Ident)
		if !ok {
			continue
		}
		if i >= len(asm.Lhs) {
			continue
		}
		lh := asm.Lhs[i]
		lt := c.pass.TypesInfo.TypeOf(lh)
		if !implementsCloser(lt) {
			continue
		}
		c.storeIdent(id)
	}
}

func (c *closeCollector) traceReturn(asm *ast.ReturnStmt, stack []ast.Node) {
	var rts []types.Type
	if fn := astpath.Nearest[*ast.FuncDecl](stack); fn != nil {
		rts = c.returnTypesOf(fn)
	} else if fn := astpath.Nearest[*ast.FuncLit](stack); fn != nil {
		rts = c.returnTypesOf(fn)
	}
	for i, r := range asm.Results {
		id, ok := r.(*ast.Ident)
		if !ok {
			continue
		}
		if i >= len(rts) {
			continue
		}
		if !implementsCloser(rts[i]) {
			continue
		}
		c.storeIdent(id)
	}
}

func (c *closeCollector) returnTypesOf(fn ast.Node) []types.Type {
	var ft *ast.FuncType
	switch fn := fn.(type) {
	case *ast.FuncDecl:
		ft = fn.Type
	case *ast.FuncLit:
		ft = fn.Type
	}
	if ft.Results == nil {
		return nil
	}
	var tys []types.Type
	for _, f := range ft.Results.List {
		tys = append(tys, c.pass.TypesInfo.TypeOf(f.Type))
	}
	return tys
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
