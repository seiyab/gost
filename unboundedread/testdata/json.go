package testdata

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func JSONDecode(req *http.Request) {
	var v any
	json.NewDecoder(req.Body).Decode(&v) // want "unbounded read of HTTP request body"
	decoder := json.NewDecoder(bufio.NewReader(req.Body))
	decoder.Decode(&v) // want "unbounded read of HTTP request body"
	var fixed struct{ Count int }
	json.NewDecoder(req.Body).Decode(&fixed) // want "unbounded read of HTTP request body"
}

func JSONBoundedDecode(w http.ResponseWriter, req *http.Request) {
	var v any
	json.NewDecoder(io.LimitReader(req.Body, 1024)).Decode(&v)
	json.NewDecoder(http.MaxBytesReader(w, req.Body, 1024)).Decode(&v)
	json.NewDecoder(&io.LimitedReader{R: req.Body, N: 1024}).Decode(&v)
	json.NewDecoder(strings.NewReader(`{"count":1}`)).Decode(&v)
}

func decodeJSON(r io.Reader) {
	var v any
	json.NewDecoder(r).Decode(&v)
}

func decodeWithDecoder(decoder *json.Decoder) {
	var v any
	decoder.Decode(&v)
}

func JSONHelpers(req *http.Request) {
	decodeJSON(req.Body) // want "unbounded read of HTTP request body"
	decodeJSON(io.LimitReader(req.Body, 1024))
	decodeWithDecoder(json.NewDecoder(req.Body)) // want "unbounded read of HTTP request body"
	decodeWithDecoder(json.NewDecoder(io.LimitReader(req.Body, 1024)))
}

func JSONCompressedDecode(req *http.Request) {
	r, err := gzip.NewReader(io.LimitReader(req.Body, 1024))
	if err != nil {
		return
	}
	defer r.Close()
	var v any
	json.NewDecoder(r).Decode(&v) // want "unbounded read of HTTP request body"
	json.NewDecoder(io.LimitReader(r, 1024)).Decode(&v)
}

func JSONConstructOnly(req *http.Request) {
	json.NewDecoder(req.Body)
}
