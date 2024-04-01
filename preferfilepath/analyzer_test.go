package preferfilepath_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/preferfilepath"
)

func TestPreferFilepath(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, preferfilepath.Analyzer)
}
