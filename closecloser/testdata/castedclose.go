package testdata

import (
	"io"
	"os"
)

func _() {
	w := myWriter()
	defer func() {
		if c, ok := w.(io.Closer); ok {
			c.Close()
		}
	}()
}

func myWriter() io.Writer {
	w, _ := os.OpenFile("a", os.O_RDWR, 0)
	return w
}
