package preferfilepath

import (
	"fmt"
	"go/ast"

	"github.com/seiyab/gost/gost/utils"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "preferfilepath",
	Doc:  "warn when using path where path/filepath should be suitable",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			var checkArgs []int
			for matcher, args := range takesPath {
				if matcher.Matches(pass, call) {
					checkArgs = args
					fmt.Println("matches tekes path")
					break
				}
			}
			if len(checkArgs) == 0 {
				return true
			}
			for _, argPos := range checkArgs {
				arg := call.Args[argPos]
				call, ok := arg.(*ast.CallExpr)
				if !ok {
					continue
				}
				for _, matcher := range pathLikeMatchers {
					if matcher.Matches(pass, call) {
						pass.ReportRangef(arg, "should use filepath instead of path")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

var pfm func(string, string) utils.PkgFuncCallMatcher = utils.NewPkgFuncCallMatcher

var takesPath = map[utils.PkgFuncCallMatcher][]int{
	pfm("os", "Chdir"):     {0},
	pfm("os", "Chmod"):     {0},
	pfm("os", "Chown"):     {0},
	pfm("os", "Chtimes"):   {0},
	pfm("os", "DirFS"):     {0},
	pfm("os", "Lchown"):    {0},
	pfm("os", "Link"):      {0, 1},
	pfm("os", "Mkdir"):     {0},
	pfm("os", "MkdirAll"):  {0},
	pfm("os", "MkdirTemp"): {0},
	pfm("os", "ReadFile"):  {0},
	pfm("os", "ReadLink"):  {0},
	pfm("os", "Remove"):    {0},
	pfm("os", "RemoveAll"): {0},
	pfm("os", "Rename"):    {0, 1},
	pfm("os", "Symlink"):   {0, 1},
	pfm("os", "Truncate"):  {0},
	pfm("os", "WriteFile"): {0},

	pfm("os", "Readdir"): {0},

	pfm("os", "Create"):     {0},
	pfm("os", "CreateTemp"): {0},
	pfm("os", "Open"):       {0},
	pfm("os", "OpenFile"):   {0},
}

var pathLikeMatchers = []utils.PkgFuncCallMatcher{
	pfm("path", "Join"),
	pfm("path", "Dir"),
}
