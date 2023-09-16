package openfileflag

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "openfileflag",
	Doc:  "prevents forgetting to specify O_TRUNC / O_APPEND / O_EXCL flags",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			slct, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			obj, ok := slct.X.(*ast.Ident)
			if !ok {
				return true
			}
			if slct.Sel.Name != "OpenFile" || obj.Name != "os" {
				return true
			}
			if len(call.Args) != 3 {
				return true
			}
			flags := FileFlagsOf(call.Args[1])
			if !flags.IsPure {
				return true
			}
			if !flags.Has("O_WRONLY") {
				return true
			}
			if !flags.Has("O_TRUNC") && !flags.Has("O_APPEND") && !flags.Has("O_EXCL") {
				pass.Reportf(call.Pos(), "O_TRUNC / O_APPEND / O_EXCL flags are not specified")
			}
			return true
		})
	}
	return nil, nil
}

type FileFlags struct {
	Flags  []*ast.SelectorExpr
	IsPure bool
}

func FileFlagsOf(expr ast.Expr) FileFlags {
	switch expr := expr.(type) {
	case *ast.SelectorExpr:
		return FileFlags{
			Flags:  []*ast.SelectorExpr{expr},
			IsPure: true,
		}
	case *ast.BinaryExpr:
		if expr.Op != token.OR {
			return FileFlags{
				Flags:  nil,
				IsPure: false,
			}
		}
		left := FileFlagsOf(expr.X)
		right := FileFlagsOf(expr.Y)
		return FileFlags{
			Flags:  append(left.Flags, right.Flags...),
			IsPure: left.IsPure && right.IsPure,
		}
	default:
		return FileFlags{
			Flags:  nil,
			IsPure: false,
		}
	}
}

func (f FileFlags) Has(flag string) bool {
	for _, fl := range f.Flags {
		if fl.Sel.Name == flag {
			return true
		}
	}
	return false
}
