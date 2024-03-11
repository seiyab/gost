package defref

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/ssa"
)

type DefRef struct {
	funcs map[ast.Node]Func
}

func (d *DefRef) trackFunc(f *ssa.Function) {
	fn := Func{
		defs: make(map[string]Def),
	}
	fmt.Println(f.Name(), "----------------")
	for _, b := range f.Blocks {
		fn.trackBlock(b)
	}
	d.funcs[f.Syntax()] = fn
}

type Func struct {
	defs map[string]Def
	// blocks []Block
}

func (f *Func) trackBlock(b *ssa.BasicBlock) {
	for i, isr := range b.Instrs {
		if isr == nil {
			fmt.Println("  ", i, "nil")
			continue
		}
		fmt.Println("  ", i, isr)
		for _, op := range isr.Operands(nil) {
			if op == nil {
				fmt.Println("    ", "nil")
				continue
			}
			// o := *op
			// fmt.Println(
			//	"    ",
			// o.Name(),
			// o.Referrers(),
			// o.Type(),
			// "(", o.String(), ")",
			// )
		}
	}
}

// type Block struct {
// 	defs map[string]Def
// }

type Def struct {
}
