package sliceinitiallength

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
)

type mutation int

const (
	appended mutation = iota + 1
	assigned
	replaced
	copied
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
			case *ast.CallExpr:
				if copyMatcher.Matches(pass, n) {
					dst := n.Args[0]
					id, ok := dst.(*ast.Ident)
					if !ok {
						break
					}
					obj := pass.TypesInfo.ObjectOf(id)
					if obj == nil {
						break
					}
					m, ok := muts[obj]
					if !ok {
						muts[obj] = copied
					} else {
						if m != copied {
							muts[obj] = mixed
						}
					}
				}
			}
			return true
		})
	}
	return muts
}

func visitAssign(n *ast.AssignStmt, pass *analysis.Pass) map[types.Object]mutation {
	if n.Tok != token.ASSIGN {
		return nil
	}
	var muts = make(map[types.Object]mutation)
	for i, lhs := range n.Lhs {
		if i >= len(n.Rhs) {
			break
		}
		switch lhs := lhs.(type) {
		case *ast.Ident:
			t := pass.TypesInfo.TypeOf(lhs)
			if t == nil {
				continue
			}
			o := pass.TypesInfo.ObjectOf(lhs)
			if o == nil {
				continue
			}
			if _, ok := t.Underlying().(*types.Slice); !ok {
				continue
			}
			call, ok := n.Rhs[i].(*ast.CallExpr)
			if !ok {
				muts[o] = replaced
				continue
			}
			if !appendMatcher.Matches(pass, call) {
				muts[o] = replaced
				continue
			}
			orig, ok := call.Args[0].(*ast.Ident)
			if !ok {
				muts[o] = replaced
				continue
			}
			if pass.TypesInfo.ObjectOf(orig) != o {
				muts[o] = replaced
				continue
			}
			muts[o] = appended
			continue
		case *ast.IndexExpr:
			if ident, ok := lhs.X.(*ast.Ident); ok {
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

var (
	appendMatcher = utils.NewBuiltinFuncCallMatcher("append")
	copyMatcher   = utils.NewBuiltinFuncCallMatcher("copy")
)
