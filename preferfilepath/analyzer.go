package preferfilepath

import (
	"github.com/pkg/errors"
	"github.com/seiyab/gost/utils"
	"github.com/seiyab/gost/utils/graph"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name:     "preferfilepath",
	Doc:      "warn when using path where path/filepath should be suitable",
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	s, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok {
		return nil, errors.Errorf("failed to get SSA")
	}
	for _, fn := range s.SrcFuncs {
		pls := pathLikes(fn)
		inspect(pass, fn, pls)
	}
	return nil, nil
}

func pathLikes(fn *ssa.Function) utils.Set[string] {
	explicitPathLikes := utils.NewSet[string]()
	g := graph.NewDirected[string]()
	utils.EachInstr(fn, func(instr ssa.Instruction) {
		switch instr := instr.(type) {
		case ssa.CallInstruction:
			cmn := instr.Common()
			for _, matcher := range pathLikeMatchers {
				if matcher.MatchesSSA(cmn) {
					explicitPathLikes.Add(instr.Value().Name())
					break
				}
			}
		case *ssa.Phi:
			for _, e := range instr.Edges {
				g.AddEdge(instr.Name(), e.Name())
			}
		}
	})
	return g.LookupBackward(explicitPathLikes.ToSlice()...)
}

func inspect(pass *analysis.Pass, fn *ssa.Function, pathLikes utils.Set[string]) {
	utils.EachInstr(fn, func(instr ssa.Instruction) {
		switch instr := instr.(type) {
		case ssa.CallInstruction:
			cmn := instr.Common()
			for matcher, idx := range takesPath {
				if matcher.MatchesSSA(cmn) {
					for _, i := range idx {
						if i >= len(cmn.Args) {
							continue
						}
						if pathLikes.Has(cmn.Args[i].Name()) {
							pass.Reportf(
								cmn.Pos(),
								`path built by package "path" isn't suitable for file path manipulation because it doesn't handle Windows paths correctly. Use "path/filepath" instead.`,
							)
							break
						}
					}
					break
				}
			}
		}
	})
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
