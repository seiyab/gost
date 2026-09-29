package unboundedread_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/unboundedread"
)

func TestPreferFilepath(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, unboundedread.Analyzer)
}

