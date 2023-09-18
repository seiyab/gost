package main

import (
	"github.com/seiyab/gost/gost/multipleerrors"
	"github.com/seiyab/gost/gost/nodiscarderror"
	"github.com/seiyab/gost/gost/openfileflag"
	"github.com/seiyab/gost/gost/wraperror"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		openfileflag.Analyzer,
		nodiscarderror.Analyzer,
		multipleerrors.Analyzer,
		wraperror.Analyzer,
	)
}
