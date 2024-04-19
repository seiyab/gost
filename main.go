package main

import (
	"github.com/seiyab/gost/closecloser"
	"github.com/seiyab/gost/multipleerrors"
	"github.com/seiyab/gost/nodiscarderror"
	"github.com/seiyab/gost/nomutateglobal"
	"github.com/seiyab/gost/openfileflag"
	"github.com/seiyab/gost/preferfilepath"
	"github.com/seiyab/gost/sliceinitiallength"
	"github.com/seiyab/gost/urlstring"
	"github.com/seiyab/gost/wraperror"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		openfileflag.Analyzer,
		nodiscarderror.Analyzer,
		nomutateglobal.Analyzer,
		multipleerrors.Analyzer,
		wraperror.Analyzer,
		closecloser.Analyzer,
		preferfilepath.Analyzer,
		sliceinitiallength.Analyzer,
		urlstring.Analyzer,
	)
}
