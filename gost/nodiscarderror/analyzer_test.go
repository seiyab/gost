package nodiscarderror_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/gost/nodiscarderror"
)

func TestOpenFileFlag(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nodiscarderror.Analyzer)
}
