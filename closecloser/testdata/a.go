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

	f1.Stat()

	c1 := closer() // want ".+"
	c1.Read(nil)

	var a io.WriteCloser // want ".+"
	a.Write(([]byte)("foobar"))
}

func Shadow(b bool) {
	f, _ := os.Open("b") // want ".+"

	if f, err := os.Open("a"); err == nil {
		f.Close()
	}

	f.Stat()
}

func closer() io.ReadCloser {
	return nil
}
