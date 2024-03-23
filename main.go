package main

import (
	"github.com/seiyab/gost/closecloser"
	"github.com/seiyab/gost/multipleerrors"
	"github.com/seiyab/gost/nodiscarderror"
	"github.com/seiyab/gost/openfileflag"
	"github.com/seiyab/gost/preferfilepath"
	"github.com/seiyab/gost/wraperror"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		openfileflag.Analyzer,
		nodiscarderror.Analyzer,
		multipleerrors.Analyzer,
		wraperror.Analyzer,
		closecloser.Analyzer,
		preferfilepath.Analyzer,
	)
}
