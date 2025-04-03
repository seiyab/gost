package utils

import "go/types"

func IsError(t types.Type) bool {
	e := types.Universe.Lookup("error")
	return types.Identical(t, e.Type())
}

type TypeMatcher struct {
	PkgPath  string
	TypeName string
}

func (m TypeMatcher) Matches(t types.Type) bool {
	if t == nil {
		return false
	}
	if p, ok := t.(*types.Pointer); ok {
		return m.Matches(p.Elem())
	}

	if n, ok := t.(*types.Named); ok {
		o := n.Obj()
		if o == nil {
			return false
		}
		pkg := o.Pkg()
		if pkg == nil {
			return false
		}
		return pkg.Path() == m.PkgPath && o.Name() == m.TypeName
	}
	return false
}
