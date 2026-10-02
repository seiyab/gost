package unboundedread

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name:     "unboundedRead",
	Doc:      "report unbounded reads of HTTP request bodies into memory",
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

// Dependencies distinguish concrete request bodies from a helper's parameters.
// Summaries describe which parameters reach an unbounded memory sink.
type dependencies map[int]bool

const requestBody = -1

func merge(dst, src dependencies) bool {
	changed := false
	for key := range src {
		if !dst[key] {
			dst[key] = true
			changed = true
		}
	}
	return changed
}

func named(t types.Type, pkg, name string) bool {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	n, ok := t.(*types.Named)
	return ok && n.Obj().Pkg() != nil && n.Obj().Pkg().Path() == pkg && n.Obj().Name() == name
}

func buffer(v ssa.Value) bool {
	if wrapped, ok := v.(*ssa.MakeInterface); ok {
		v = wrapped.X
	}
	return named(v.Type(), "bytes", "Buffer")
}

func function(c *ssa.CallCommon, pkg, name string) bool {
	f := c.StaticCallee()
	return f != nil && f.Pkg != nil && f.Pkg.Pkg.Path() == pkg && f.Name() == name
}

func run(pass *analysis.Pass) (any, error) {
	functions := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	summaries := make(map[*ssa.Function]dependencies)
	for _, f := range functions {
		summaries[f] = dependencies{}
	}
	// Monotone summaries also handle mutually recursive local helpers.
	for changed := true; changed; {
		changed = false
		for _, f := range functions {
			_, sinks := analyze(f, summaries)
			for _, d := range sinks {
				changed = merge(summaries[f], d) || changed
			}
		}
	}
	for _, f := range functions {
		_, sinks := analyze(f, summaries)
		for call, d := range sinks {
			if strings.HasSuffix(pass.Fset.Position(call.Pos()).Filename, "_test.go") {
				continue
			}
			if d[requestBody] {
				pass.Reportf(call.Pos(), "unbounded read of HTTP request body into memory can exhaust memory; limit the reader before reading")
			}
		}
	}
	return nil, nil
}

func analyze(f *ssa.Function, summaries map[*ssa.Function]dependencies) (map[ssa.Value]dependencies, map[*ssa.Call]dependencies) {
	values := make(map[ssa.Value]dependencies)
	// Origins retain request provenance through limits: a compressed-byte limit
	// does not provide the same allocation bound on decompressed output.
	origins := make(map[ssa.Value]dependencies)
	origin := func(v ssa.Value) dependencies {
		if origins[v] == nil {
			origins[v] = dependencies{}
		}
		return origins[v]
	}
	get := func(v ssa.Value) dependencies {
		if values[v] == nil {
			values[v] = dependencies{}
		}
		return values[v]
	}
	for i, p := range f.Params {
		get(p)[i] = true
		origin(p)[i] = true
	}
	sinks := make(map[*ssa.Call]dependencies)
	for changed := true; changed; {
		changed = false
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				var inputs []ssa.Value
				var originInputs []ssa.Value
				switch v := instr.(type) {
				case *ssa.FieldAddr:
					st, ok := v.X.Type().Underlying().(*types.Pointer)
					if ok && named(st.Elem(), "net/http", "Request") {
						s := st.Elem().Underlying().(*types.Struct)
						if s.Field(v.Field).Name() == "Body" {
							d := dependencies{requestBody: true}
							changed = merge(get(v), d) || changed
							changed = merge(origin(v), d) || changed
							continue
						}
					}
					inputs = []ssa.Value{v.X}
				case *ssa.Field:
					if named(v.X.Type(), "net/http", "Request") && v.X.Type().Underlying().(*types.Struct).Field(v.Field).Name() == "Body" {
						changed = merge(get(v), dependencies{requestBody: true}) || changed
						changed = merge(origin(v), dependencies{requestBody: true}) || changed
						continue
					}
					inputs = []ssa.Value{v.X}
				case *ssa.Store:
					changed = merge(get(v.Addr), get(v.Val)) || changed
					changed = merge(origin(v.Addr), origin(v.Val)) || changed
					if index, ok := v.Addr.(*ssa.IndexAddr); ok {
						changed = merge(get(index.X), get(v.Val)) || changed
						changed = merge(origin(index.X), origin(v.Val)) || changed
					}
					continue
				case *ssa.Call:
					c := v.Common()
					d := dependencies{}
					if function(c, "io", "ReadAll") && len(c.Args) == 1 {
						merge(d, get(c.Args[0]))
					} else if function(c, "io", "Copy") && len(c.Args) == 2 && buffer(c.Args[0]) {
						merge(d, get(c.Args[1]))
					} else if summary := summaries[c.StaticCallee()]; summary != nil {
						for i := range summary {
							// A helper's intrinsic source is reported inside that helper.
							if i >= 0 && i < len(c.Args) {
								merge(d, get(c.Args[i]))
							}
						}
					}
					if len(d) > 0 {
						sinks[v] = d
					}
					if function(c, "io", "TeeReader") {
						inputs = c.Args[:1]
					}
					if function(c, "io", "MultiReader") {
						inputs = c.Args
					}
					if function(c, "bufio", "NewReader") || function(c, "bufio", "NewReaderSize") || function(c, "io", "NopCloser") {
						inputs = c.Args[:1]
					}
					if function(c, "io", "LimitReader") {
						originInputs = c.Args[:1]
					}
					if function(c, "net/http", "MaxBytesReader") {
						originInputs = c.Args[1:2]
					}
					if function(c, "compress/gzip", "NewReader") || function(c, "compress/zlib", "NewReader") {
						changed = merge(get(v), origin(c.Args[0])) || changed
						originInputs = c.Args[:1]
					}
				case *ssa.Extract:
					// Reader constructors return (reader, error); only the reader
					// carries input provenance.
					if v.Index == 0 {
						inputs = []ssa.Value{v.Tuple}
					}
				case *ssa.UnOp:
					inputs = []ssa.Value{v.X}
				case *ssa.MakeInterface:
					inputs = []ssa.Value{v.X}
				case *ssa.ChangeInterface:
					inputs = []ssa.Value{v.X}
				case *ssa.ChangeType:
					inputs = []ssa.Value{v.X}
				case *ssa.Convert:
					inputs = []ssa.Value{v.X}
				case *ssa.Phi:
					inputs = v.Edges
				case *ssa.IndexAddr:
					inputs = []ssa.Value{v.X}
				case *ssa.Index:
					inputs = []ssa.Value{v.X}
				case *ssa.Slice:
					inputs = []ssa.Value{v.X}
				}
				if value, ok := instr.(ssa.Value); ok {
					if originInputs == nil {
						originInputs = inputs
					}
					for _, input := range originInputs {
						changed = merge(origin(value), origin(input)) || changed
					}
					for _, input := range inputs {
						if wrapped, ok := value.(*ssa.MakeInterface); ok && named(wrapped.X.Type(), "io", "LimitedReader") {
							continue
						}
						changed = merge(get(value), get(input)) || changed
					}
				}
			}
		}
	}
	return values, sinks
}
