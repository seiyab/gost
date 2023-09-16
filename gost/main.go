package main

import (
	"github.com/seiyab/gost/gost/nodiscarderror"
	"github.com/seiyab/gost/gost/openfileflag"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		openfileflag.Analyzer,
		nodiscarderror.Analyzer,
	)
}
