package nomutateglobal

import (
	"go/ast"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type fieldSet struct {
	pass *analysis.Pass
	set  map[types.Object]fieldSetField
}

type fieldSetField struct {
	// NOTE: use "" as self
	children map[string]fieldSetField
}

func newFieldSet(pass *analysis.Pass) fieldSet {
	return fieldSet{
		pass: pass,
		set:  make(map[types.Object]fieldSetField),
	}
}

func (s *fieldSet) add(sel *ast.SelectorExpr) bool {
	o, path, ok := s.retrievePath(sel)
	if !ok {
		return false
	}
	otr, ok := s.set[o]
	if !ok {
		otr = fieldSetField{children: make(map[string]fieldSetField)}
		s.set[o] = otr
	}
	for _, p := range path {
		if _, ok := otr.children[p]; !ok {
			otr.children[p] = fieldSetField{children: make(map[string]fieldSetField)}
		}
		otr = otr.children[p]
	}
	otr.children[""] = fieldSetField{}
	return true
}

func (s *fieldSet) has(sel *ast.SelectorExpr) bool {
	o, path, ok := s.retrievePath(sel)
	if !ok {
		return false
	}
	otr, ok := s.set[o]
	if !ok {
		return false
	}
	for _, p := range path {
		otr, ok = otr.children[p]
		if !ok {
			return false
		}
	}
	_, ok = otr.children[""]
	return ok
}

func (s *fieldSet) retrievePath(sel *ast.SelectorExpr) (types.Object, []string, bool) {
	currentSel := sel
	var revPath []string
	var id *ast.Ident
	for id == nil {
		revPath = append(revPath, currentSel.Sel.Name)
		switch x := currentSel.X.(type) {
		case *ast.SelectorExpr:
			currentSel = x
		case *ast.Ident:
			id = x
		default:
			return nil, nil, false
		}
	}
	o := s.pass.TypesInfo.ObjectOf(id)
	if o == nil {
		return nil, nil, false
	}
	slices.Reverse(revPath)
	path := revPath
	return o, path, true
}
