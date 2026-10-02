package testdata

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
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

func HelperChain(req *http.Request) {
	drainAgain(req.Body) // want "unbounded read of HTTP request body"
	drainAgain(io.LimitReader(req.Body, 1024))
}

func drainAgain(r io.Reader) {
	drain(r)
}

func BoundedAndStreaming(req *http.Request) {
	var b bytes.Buffer
	io.Copy(&b, io.LimitReader(req.Body, 1024))
	io.Copy(io.Discard, req.Body)
	io.CopyN(&b, req.Body, 1024)
	io.ReadAll(bytes.NewReader([]byte("local data")))
}

func JoinedReaders(req *http.Request, condition bool) {
	var r io.Reader = bytes.NewReader(nil)
	if condition {
		r = req.Body
	}
	io.ReadAll(r) // want "unbounded read of HTTP request body"
}

func BufferedReaders(req *http.Request) {
	io.ReadAll(bufio.NewReader(req.Body))           // want "unbounded read of HTTP request body"
	io.ReadAll(bufio.NewReaderSize(req.Body, 4096)) // want "unbounded read of HTTP request body"
	io.ReadAll(io.NopCloser(req.Body))              // want "unbounded read of HTTP request body"
}

func ExplicitlyLimitedReaders(w http.ResponseWriter, req *http.Request) {
	io.ReadAll(http.MaxBytesReader(w, req.Body, 1024))
	io.ReadAll(&io.LimitedReader{R: req.Body, N: 1024})
	io.ReadAll(io.LimitReader(bufio.NewReader(req.Body), 1024))
	io.ReadAll(bufio.NewReader(io.LimitReader(req.Body, 1024)))
}

func MultiReaderLimits(req *http.Request) {
	limitedBody := io.LimitReader(req.Body, 1024)
	io.ReadAll(io.MultiReader(limitedBody, req.Body)) // want "unbounded read of HTTP request body"
	limitedA := io.LimitReader(req.Body, 1024)
	limitedB := io.LimitReader(req.Body, 2048)
	io.ReadAll(io.MultiReader(limitedA, limitedB))
}

func MemoryBackedReaders() {
	io.ReadAll(strings.NewReader("local data"))
	io.ReadAll(bytes.NewReader([]byte("local data")))
}

func GzipReader(req *http.Request) {
	r, err := gzip.NewReader(req.Body)
	if err != nil {
		return
	}
	defer r.Close()
	io.ReadAll(r) // want "unbounded read of HTTP request body"
}

func GzipLimitBeforeDecompression(req *http.Request) {
	r, err := gzip.NewReader(io.LimitReader(req.Body, 1024))
	if err != nil {
		return
	}
	defer r.Close()
	// Limit decompressed output to prevent excessive allocation from expansion.
	io.ReadAll(r) // want "unbounded read of HTTP request body"
}

func GzipLimitAfterDecompression(req *http.Request) {
	r, err := gzip.NewReader(req.Body)
	if err != nil {
		return
	}
	defer r.Close()
	io.ReadAll(io.LimitReader(r, 1024))
}

func ZlibReader(req *http.Request) {
	r, err := zlib.NewReader(req.Body)
	if err != nil {
		return
	}
	defer r.Close()
	io.ReadAll(r) // want "unbounded read of HTTP request body"
}

func ZlibLimitBeforeDecompression(req *http.Request) {
	r, err := zlib.NewReader(io.LimitReader(req.Body, 1024))
	if err != nil {
		return
	}
	defer r.Close()
	io.ReadAll(r) // want "unbounded read of HTTP request body"
}

func ZlibLimitAfterDecompression(req *http.Request) {
	r, err := zlib.NewReader(req.Body)
	if err != nil {
		return
	}
	defer r.Close()
	io.ReadAll(io.LimitReader(r, 1024))
}
