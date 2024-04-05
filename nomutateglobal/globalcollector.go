package nomutateglobal

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type key types.Object

type globalCollector struct {
	pass            *analysis.Pass
	identCollection map[key]struct{}
	fieldCollection fieldSet
}

func newGlobalCollector(pass *analysis.Pass) globalCollector {
	return globalCollector{
		pass:            pass,
		identCollection: make(map[key]struct{}),
		fieldCollection: newFieldSet(pass),
	}
}

func (c *globalCollector) visitAssignment(asmt *ast.AssignStmt) {
	for i, rhs := range asmt.Rhs {
		if !isGlobalVar(rhs, c.pass) {
			continue
		}
		lhs := asmt.Lhs[i]
		switch lhs := lhs.(type) {
		case *ast.Ident:
			o := c.pass.TypesInfo.ObjectOf(lhs)
			if o == nil {
				continue
			}
			c.identCollection[key(o)] = struct{}{}
		case *ast.SelectorExpr:
			c.fieldCollection.add(lhs)
		}
	}
}

func (c *globalCollector) contains(expr ast.Expr) bool {
	switch expr := expr.(type) {
	case *ast.Ident:
		o := c.pass.TypesInfo.ObjectOf(expr)
		if o == nil {
			return false
		}
		_, ok := c.identCollection[key(o)]
		return ok
	case *ast.SelectorExpr:
		return c.fieldCollection.has(expr)
	default:
		return false
	}
}
