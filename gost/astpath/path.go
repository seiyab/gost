package astpath

import "go/ast"

type Path struct {
	Current ast.Node
	Parent  *Path
}

func (s *Path) push(v ast.Node) {
	n := *s
	s.Current, s.Parent = v, &n
}

func (s *Path) pop() (ast.Node, bool) {
	if s == nil {
		return nil, false
	}
	ret := s.Current
	s.Current, s.Parent = s.Parent.Current, s.Parent.Parent
	return ret, true
}

func FindNearest[T ast.Node](p *Path) (T, *Path) {
	if p == nil {
		var zero T
		return zero, nil
	}
	x, ok := p.Current.(T)
	if ok {
		return x, p
	}
	return FindNearest[T](p.Parent)
}
