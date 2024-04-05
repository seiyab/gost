package nomutateglobal

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type key types.Object

type globalCollector struct {
	collection map[key]struct{}
}

func newGlobalCollector() globalCollector {
	return globalCollector{
		collection: make(map[key]struct{}),
	}
}

func (c *globalCollector) visitAssignment(asmt *ast.AssignStmt, pass *analysis.Pass) {
	for i, rhs := range asmt.Rhs {
		if !isGlobalVar(rhs, pass) {
			continue
		}
		lhs := asmt.Lhs[i]
		idt, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		o := pass.TypesInfo.ObjectOf(idt)
		if o == nil {
			continue
		}
		c.add(o)
	}
}

func (c *globalCollector) add(o types.Object) {
	c.collection[key(o)] = struct{}{}
}

func (c *globalCollector) contains(o types.Object) bool {
	_, ok := c.collection[key(o)]
	return ok
}
