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

func Delegated() (io.Closer, io.Reader, io.ReadCloser, *os.File) {
	f1, _ := os.Open("a")
	takeReadCloser(f1)

	f12, _ := os.Open("b")
	f13, _ := os.Open("b") // want ".+"
	f14, _ := os.Open("b")
	f15, _ := os.Open("b")

	o := st{
		c:  f12,
		r:  f13,
		rc: f14,
		f:  f15,
	}

	markAsUsed(o)

	f22, _ := os.Open("c")
	f23, _ := os.Open("c") // want ".+"
	f24, _ := os.Open("c")
	f25, _ := os.Open("c")

	o = st{f22, f23, f24, f25}

	markAsUsed(o)

	o = st{}
	f32, _ := os.Open("d")
	f33, _ := os.Open("d") // want ".+"
	f34, _ := os.Open("d")
	f35, _ := os.Open("d")

	o.c = f32
	o.r = f33
	o.rc = f34
	o.f = f35

	markAsUsed(o)

	f42, _ := os.Open("d")
	f43, _ := os.Open("d") // want ".+"
	f44, _ := os.Open("d")
	f45, _ := os.Open("d")

	return f42, f43, f44, f45
}

type st struct {
	c  io.Closer
	r  io.Reader
	rc io.ReadCloser
	f  *os.File
}

func markAsUsed(_ ...any) {}

func takeReadCloser(r io.ReadCloser) {}

func closer() io.Closer {
	return nil
}
