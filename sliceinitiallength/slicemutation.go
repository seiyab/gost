package sliceinitiallength

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type mutation int

const (
	appended mutation = iota + 1
	assigned
	mixed
)

func collectSliceMutations(pass *analysis.Pass) map[types.Object]mutation {
	var muts = make(map[types.Object]mutation)
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.AssignStmt:
				am := visitAssign(n, pass)
				for k, v := range am {
					m, ok := muts[k]
					if !ok {
						muts[k] = v
					} else {
						if m != v {
							muts[k] = mixed
						}
					}
				}
				return true
			default:
				return true
			}
		})
	}
	return muts
}

func visitAssign(n *ast.AssignStmt, pass *analysis.Pass) map[types.Object]mutation {
	var muts = make(map[types.Object]mutation)
	for i, lhs := range n.Lhs {
		if i >= len(n.Rhs) {
			break
		}
		if ident, ok := lhs.(*ast.Ident); ok {
			t := pass.TypesInfo.TypeOf(ident)
			if t == nil {
				continue
			}
			o := pass.TypesInfo.ObjectOf(ident)
			if o == nil {
				continue
			}
			if _, ok := t.Underlying().(*types.Slice); !ok {
				continue
			}
			call, ok := n.Rhs[i].(*ast.CallExpr)
			if !ok {
				continue
			}
			fn, ok := call.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			if fn.Name != "append" {
				continue
			}
			muts[o] = appended
			continue
		}

		if sel, ok := lhs.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				t := pass.TypesInfo.TypeOf(ident)
				if t == nil {
					continue
				}
				o := pass.TypesInfo.ObjectOf(ident)
				if o == nil {
					continue
				}
				if _, ok := t.Underlying().(*types.Slice); !ok {
					continue
				}
				muts[o] = assigned
				continue
			}
		}
	}
	return muts
}
