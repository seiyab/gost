package multipleerrors_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/gost/multipleerrors"
)

func TestOpenFileFlag(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, multipleerrors.Analyzer)
}
