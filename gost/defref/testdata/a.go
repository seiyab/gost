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

	markAsUsed(f1, c1)
}

func C(b bool) {
	f2, _ := os.Open("b")
	defer f2.Close()

	if b {
		f2 = nil
	}

	c1 := closer() // want ".+"

	markAsUsed(c1, f2)
}

func markAsUsed(_ ...any) {}

func closer() io.Closer {
	return nil
}
