package urlstring

import (
	"go/ast"
	"go/token"
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
		if f.Name != nil && strings.HasSuffix(f.Name.Name, "_test") {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				checkCallExpr(pass, n)
			case *ast.BinaryExpr:
				checkBinaryExpr(pass, n)
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
	if !isURLLikeNode(n.X) {
		return
	}
	pass.Reportf(n.Pos(), "use *url.URL for URL creation to avoid injection vulnerabilities")
}
