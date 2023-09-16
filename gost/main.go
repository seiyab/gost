package main

import (
	"github.com/seiyab/gost/gost/openfileflag"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		openfileflag.Analyzer,
	)
}
