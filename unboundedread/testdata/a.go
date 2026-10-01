package testdata

import (
	"bytes"
	"io"
	"net/http"
)

// NOTE: function named `_` won't get SSA representation

func ReadAll(req http.Request) {
	io.ReadAll(req.Body) // want ".+"

	var s Struct1
	io.ReadAll(s.req.Body) // want ".+"

	var r io.Reader
	io.ReadAll(r)
}

func CopyToBuffer(req http.Request) {
	var b bytes.Buffer
	io.Copy(&b, req.Body) // want ".+"

	var s Struct1
	io.Copy(&b, s.req.Body) // want ".+"

	var r io.Reader
	io.Copy(&b, r)
}

func CalleeReadsAll(req http.Request) {
	drain(req.Body) // want ".+"

	read100(req.Body)
}

func ComposedReader(req http.Request) {
	var w io.Writer
	t := io.TeeReader(req.Body, w)
	io.ReadAll(t) // want ".+"

	var r io.Reader
	m := io.MultiReader(r, req.Body)
	io.ReadAll(m) // want ".+"

	l := io.LimitReader(req.Body, 1_000_000)
	io.ReadAll(l)
}

type Struct1 struct {
	req http.Request
}

func drain(r io.Reader) {
	io.ReadAll(r)
}

func read100(r io.Reader) {
	buf := make([]byte, 100)
	io.ReadFull(r, buf)
}
