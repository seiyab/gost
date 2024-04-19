package urlstring

import (
	"go/ast"
	"go/token"
	"strings"
)

func isURLLikeNode(node ast.Node) bool {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return false
	}
	if lit.Kind != token.STRING {
		return false
	}
	return isURLLike(strings.Trim(lit.Value, "\"`"))
}

func isURLLike(s string) bool {
	head := s[:min(10, len(s))]
	return strings.Contains(head, "://")
}
