package testdata

import (
	"io"
	"io/fs"
	"os"
	"os/exec"
)

func _() {
	var cmd exec.Cmd

	p1, _ := cmd.StdinPipe()
	p2, _ := cmd.StdoutPipe()
	p3, _ := cmd.StderrPipe()

	p1.Write(nil)
	p2.Read(nil)
	p3.Read(nil)
}

func _() {
	var cmd *exec.Cmd

	p1, _ := cmd.StdinPipe()
	p2, _ := cmd.StdoutPipe()
	p3, _ := cmd.StderrPipe()

	p1.Write(nil)
	p2.Read(nil)
	p3.Read(nil)
}

type o struct {
	wc io.WriteCloser
}

func (o o) alias() {
	// just define shorthands. no need to close
	wc := o.wc
	wc.Write(nil)

	f, _ := o.wc.(*os.File)
	f.Stat()
}

func _() {
	var o o
	o.alias()
}

func _(f fs.File) {
	ff, _ := f.(fs.ReadDirFile)
	ff.Read(nil)
}
