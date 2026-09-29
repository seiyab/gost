package testdata

import (
	"io"
	"net/http"
)

// NOTE: function named `_` won't get SSA representation

func A(req http.Request) {
	io.ReadAll(req.Body) // want ".+"
}

