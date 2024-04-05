package testdata

import "net/http"

var MyGlobal *struct {
	Field string
}

func _() {
	client := http.DefaultClient
	client.Transport = nil // want ".+"

	x := MyGlobal
	x.Field = "foo"
}
