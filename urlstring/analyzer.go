package urlstring

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/seiyab/gost/utils"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "urlString",
	Doc:  "reports unsafe URL string creation",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		if utils.IsTestRelatedFile(f) {
			continue
		}
		var stack []ast.Node
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			var parent ast.Node
			for i := len(stack) - 1; i >= 0; i-- {
				if _, ok := stack[i].(*ast.ParenExpr); !ok {
					parent = stack[i]
					break
				}
			}
			stack = append(stack, n)
			switch n := n.(type) {
			case *ast.CallExpr:
				checkCallExpr(pass, n)
			case *ast.BinaryExpr:
				if p, ok := parent.(*ast.BinaryExpr); !ok || p.Op != token.ADD {
					checkBinaryExpr(pass, n)
				}
			}
			return true
		})
	}
	return nil, nil
}

var sprintfMatcher = utils.NewPkgFuncCallMatcher("fmt", "Sprintf")

func checkCallExpr(pass *analysis.Pass, n *ast.CallExpr) {
	if !sprintfMatcher.Matches(pass, n) {
		return
	}
	if len(n.Args) < 2 {
		return
	}
	if !isURLLikeNode(n.Args[0]) {
		return
	}
	pass.Reportf(n.Pos(), "use *url.URL for URL creation to avoid injection vulnerabilities")
}

func checkBinaryExpr(pass *analysis.Pass, n *ast.BinaryExpr) {
	if n.Op != token.ADD {
		return
	}
	var parts []ast.Expr
	var flatten func(ast.Expr)
	flatten = func(expr ast.Expr) {
		if paren, ok := expr.(*ast.ParenExpr); ok {
			flatten(paren.X)
		} else if binary, ok := expr.(*ast.BinaryExpr); ok && binary.Op == token.ADD {
			flatten(binary.X)
			flatten(binary.Y)
		} else {
			parts = append(parts, expr)
		}
	}
	flatten(n)
	if !isURLLikeNode(parts[0]) {
		return
	}
	var prefix strings.Builder
	for _, part := range parts {
		if lit, ok := part.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			s, err := strconv.Unquote(lit.Value)
			if err == nil {
				prefix.WriteString(s)
				continue
			}
		}
		p := prefix.String()
		query := strings.IndexByte(p, '?')
		call, ok := part.(*ast.CallExpr)
		if query < 0 || strings.Contains(p, "#") || !ok || !isEscapedQuery(pass, call, p) {
			pass.Reportf(n.Pos(), "use *url.URL for URL creation to avoid injection vulnerabilities")
			return
		}
		prefix.WriteString("escaped")
	}
}

var queryEscapeMatcher = utils.NewPkgFuncCallMatcher("net/url", "QueryEscape")

func isEscapedQuery(pass *analysis.Pass, call *ast.CallExpr, prefix string) bool {
	if queryEscapeMatcher.Matches(pass, call) {
		return true
	}
	if !strings.HasSuffix(prefix, "?") && !strings.HasSuffix(prefix, "&") {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	selection := pass.TypesInfo.Selections[sel]
	if selection == nil {
		return false
	}
	method := selection.Obj()
	sig, ok := method.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return false
	}
	receiver, ok := sig.Recv().Type().(*types.Named)
	return ok && receiver.Obj().Pkg() != nil && receiver.Obj().Pkg().Path() == "net/url" && receiver.Obj().Name() == "Values" && method.Name() == "Encode"
}
