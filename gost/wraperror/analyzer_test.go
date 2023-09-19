package wraperror_test

import (
	"os/exec"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/seiyab/gost/gost/wraperror"
)

func TestWrapError(t *testing.T) {
	make := exec.Command("make")
	make.Dir = "./testdata"
	if err := make.Run(); err != nil {
		t.Fatal(err)
	}
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, wraperror.Analyzer)
}
