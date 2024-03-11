package testdata

import (
	"io"
	"os"
)

func _() {

}

func A() {
	f1, _ := os.Open("a") // want ".+"
	f2, _ := os.Open("b")
	defer f2.Close()

	c1 := closer() // want ".+"

	var a io.Closer // want ".+"

	markAsUsed(f1, c1, a)
}

func Shadow(b bool) {
	f, _ := os.Open("b") // want ".+"

	if f, err := os.Open("a"); err == nil {
		f.Close()
	}

	markAsUsed(f)
}

func Delegated() {
	f, _ := os.Open("a")
	takeReadCloser(f)
}

func markAsUsed(_ ...any) {}

func takeReadCloser(r io.ReadCloser) {}

func closer() io.Closer {
	return nil
}
