package closecloser

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "closeCloser",
	Doc:  "report unclosed variables that implement io.Closer",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	ti := pass.TypesInfo
	ipr := inspector.New(pass.Files)
	c := newCloseCollector(pass, ipr)

	ipr.Nodes([]ast.Node{
		&ast.FuncDecl{},
		&ast.AssignStmt{},
		&ast.ValueSpec{},
	},
		func(n ast.Node, push bool) bool {
			if !push {
				return true
			}
			var ids []*ast.Ident
			switch n := n.(type) {
			case *ast.FuncDecl:
				// NOTE:
				// As of now, the analyzer does not report test functions.
				// It is because:
				//     - resources can be implicitly closed by t.Cleanup()
				//     - people tend to write tests that do not close resources
				// I think it is not desirable but this is practically helpful anyway.
				if isTestOrTestHelper(n, pass) {
					return false
				}
			case *ast.AssignStmt:
				if n.Tok != token.DEFINE {
					return true
				}
				if len(n.Rhs) == 1 && shouldAllowUnclosed(n.Rhs[0], pass) {
					return true
				}
				for _, ex := range n.Lhs {
					if id, ok := ex.(*ast.Ident); ok {
						ids = append(ids, id)
					}
				}
			case *ast.ValueSpec:
				ids = n.Names
			default:
				return true
			}
			for _, id := range ids {
				dt := ti.TypeOf(id)
				if id.Name == "_" {
					continue
				}
				if !implementsCloser(dt) {
					continue
				}

				o := ti.ObjectOf(id)
				if c.isClosed(o) {
					continue
				}
				pass.Reportf(id.Pos(), "variable %s is not closed", id.Name)
			}

			return true
		},
	)
	return nil, nil
}

func implementsCloser(t types.Type) bool {
	if t == nil {
		return false
	}
	switch v := t.(type) {
	case *types.Interface:
		return methoderImplementsCloser(v)
	case *types.Named:
		ul := v.Underlying()
		switch ul := ul.(type) {
		case *types.Interface:
			return methoderImplementsCloser(ul)
		default:
			return methoderImplementsCloser(v)
		}
	case *types.Pointer:
		return implementsCloser(v.Elem())
	}
	ul := t.Underlying()
	if t == ul {
		return false
	}
	return implementsCloser(ul)
}

func methoderImplementsCloser(ifc methoder) bool {
	var close *types.Func
	for i := 0; i < ifc.NumMethods(); i++ {
		m := ifc.Method(i)
		if m.Name() == "Close" {
			close = m
			break
		}
	}
	if close == nil {
		return false
	}
	sgn, ok := close.Type().(*types.Signature)
	if !ok {
		return false
	}
	if sgn.Params().Len() != 0 {
		return false
	}
	if sgn.Results().Len() != 1 {
		return false
	}
	ret := sgn.Results().At(0)
	return utils.IsError(ret.Type())
}

type methoder interface {
	NumMethods() int
	Method(int) *types.Func
}
