package testdata

import (
	"bytes"
	"io"
	"net/http"
)

// Test-only reads should not produce diagnostics, including through helpers.
func testOnlyReads(req *http.Request) {
	io.ReadAll(req.Body)
	var b bytes.Buffer
	io.Copy(&b, req.Body)
	drain(req.Body)
	func() {
		io.ReadAll(req.Body)
	}()
}
