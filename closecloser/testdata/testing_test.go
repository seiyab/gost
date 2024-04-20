package testdata_test

import (
	"os"
	"testing"
)

// NOTE:
// As of now, the analyzer does not report test functions.
// It is because:
//     - resources can be implicitly closed by t.Cleanup()
//     - people tend to write tests that do not close resources

func TestMyFunc(t *testing.T) {
	f, _ := os.Open("file")
	f.Stat()

	myUtil(t)
}

func myUtil(t *testing.T) {
	f, _ := os.Open("file")
	f.Stat()
}
