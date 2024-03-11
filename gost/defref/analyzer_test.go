package defref_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/gost/defref"
)

func TestOpenFileFlag(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, defref.Analyzer)
}
