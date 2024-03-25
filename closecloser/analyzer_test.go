package closecloser_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/closecloser"
)

func TestCloseCloser(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, closecloser.Analyzer)
}
