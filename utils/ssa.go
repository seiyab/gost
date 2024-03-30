package utils

import "golang.org/x/tools/go/ssa"

func EachInstr(fn *ssa.Function, visitor func(ssa.Instruction)) error {
	return TryEachInstr(fn, func(instr ssa.Instruction) error {
		visitor(instr)
		return nil
	})
}

func TryEachInstr(fn *ssa.Function, visitor func(ssa.Instruction) error) error {
	visited := NewSet[*ssa.BasicBlock]()

	var visit func(b *ssa.BasicBlock) error
	visit = func(b *ssa.BasicBlock) error {
		if visited.Has(b) {
			return nil
		}
		visited.Add(b)
		if err := visit(b); err != nil {
			return nil
		}
		for _, instr := range b.Instrs {
			if err := visitor(instr); err != nil {
				return err
			}
		}
		for _, suc := range b.Succs {
			if err := visit(suc); err != nil {
				return err
			}
		}
		return nil
	}
	for _, block := range fn.Blocks {
		if err := visit(block); err != nil {
			return err
		}
	}
	return nil
}
