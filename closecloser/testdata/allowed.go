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

	markAsUsed(p1, p2, p3)
}

func _() {
	var cmd *exec.Cmd

	p1, _ := cmd.StdinPipe()
	p2, _ := cmd.StdoutPipe()
	p3, _ := cmd.StderrPipe()

	markAsUsed(p1, p2, p3)
}

type o struct {
	wc io.WriteCloser
}

func (o o) alias() {
	// just define shorthands. no need to close
	wc := o.wc
	f, _ := o.wc.(*os.File)
	markAsUsed(wc, f)
}

func _() {
	var o o
	markAsUsed(o, o.alias)
}

func _(f fs.File) {
	ff, _ := f.(fs.ReadDirFile)
	markAsUsed(ff)
}
