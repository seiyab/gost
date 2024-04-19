package utils

import (
	"go/ast"
	"strings"
)

func IsTestRelatedFile(file *ast.File) bool {
	if file.Name != nil && strings.HasSuffix(file.Name.Name, "_test") {
		return true
	}
	for _, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, "\"") == "testing" {
			return true
		}
	}
	return false
}
