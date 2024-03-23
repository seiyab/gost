package testdata

import (
	"os"
	"path"
	"path/filepath"
)

// NOTE: function named `_` won't get SSA representation

func A() {
	os.RemoveAll(path.Join("a", "b")) // want ".+"
	os.RemoveAll(filepath.Join("a", "b"))
}

func B(b bool) {
	var p string
	var q string
	if b {
		p = path.Join("a", "b")
		q = filepath.Join("a", "b")
	}
	os.RemoveAll(p) // want ".+"
	os.RemoveAll(q)
}

func C(b bool) {
	var p string
	var q string
	if b {
		p = ""
		q = ""
	} else {
		p = path.Join("c", "d")
		q = filepath.Join("c", "d")
	}
	os.RemoveAll(p) // want ".+"
	os.RemoveAll(q)
}

func D(s []string) {
	var p string
	var q string
	for _, v := range s {
		p = path.Join(p, v)
		q = filepath.Join(q, v)
	}
	os.RemoveAll(p) // want ".+"
	os.RemoveAll(q)
}

func E() {
	p := path.Join("a", "b")
	os.RemoveAll(p) // want ".+"
	p = filepath.Join("a", "b")
	os.RemoveAll(p)
}
