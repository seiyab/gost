package nomutateglobal_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/nomutateglobal"
)

func TestNoMutateGlobal(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nomutateglobal.Analyzer)
}
