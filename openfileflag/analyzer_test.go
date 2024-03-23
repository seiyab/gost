package openfileflag_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/openfileflag"
)

func TestOpenFileFlag(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, openfileflag.Analyzer)
}
