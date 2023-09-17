package astpath

import "go/ast"

func WithPath(fn func(ast.Node, *Path) bool) func(n ast.Node) bool {
	var path Path
	return func(n ast.Node) bool {
		if n != nil {
			path.push(n)
		} else {
			defer path.pop()
		}
		return fn(n, &path)
	}
}
