package urlstring_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/urlstring"
)

func TestURLString(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, urlstring.Analyzer)
}
