package sliceinitiallength_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/sliceinitiallength"
)

func TestSliceInitialLength(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, sliceinitiallength.Analyzer)
}
