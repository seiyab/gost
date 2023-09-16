package nodiscarderror

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nodiscarderror",
	Doc:      "prevents `if err != nil { return nil }`",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	ti := pass.TypesInfo
	for _, file := range pass.Files {
		var stack Stack[ast.Node]
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				stack.Pop()
				return true
			} else {
				stack.Push(n)
			}

			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
			if !ok {
				return true
			}
			if binExpr.Op != token.NEQ {
				return true
			}
			if ti.TypeOf(binExpr.X).String() != "error" {
				return true
			}
			y, ok := ti.Types[binExpr.Y]
			if !ok || !y.IsNil() {
				return true
			}
			list := ifStmt.Body.List
			if len(list) != 1 {
				return true
			}
			ret, ok := list[0].(*ast.ReturnStmt)
			if !ok {
				return true
			}
			fn := enclosingFunction(stack)
			if fn == nil {
				return true
			}
			results := fn.Type.Results
			if results == nil {
				return true
			}
			retTypes := results.List
			if len(retTypes) != len(ret.Results) {
				return true
			}
			for i, result := range ret.Results {
				t := ti.TypeOf(retTypes[i].Type)
				if t.String() != "error" {
					continue
				}
				res, ok := ti.Types[result]
				if !ok {
					continue
				}
				if res.IsNil() {
					pass.Reportf(ret.Pos(), "error is discarded")
					return true
				}
			}
			return true
		})
	}
	return nil, nil
}

func enclosingFunction(stack Stack[ast.Node]) *ast.FuncDecl {
	if fn, ok := stack.Top.(*ast.FuncDecl); ok {
		return fn
	}
	if stack.Next == nil {
		return nil
	}
	return enclosingFunction(*stack.Next)
}

type Stack[T any] struct {
	Top  T
	Next *Stack[T]
}

func (s *Stack[T]) Push(v T) {
	n := *s
	s.Top, s.Next = v, &n
}

func (s *Stack[T]) Pop() (T, bool) {
	if s == nil {
		var v T
		return v, false
	}
	ret := s.Top
	s.Top, s.Next = s.Next.Top, s.Next.Next
	return ret, true
}
