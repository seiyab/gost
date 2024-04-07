package testdata

import (
	"flag"
	"net"
	"net/http"
)

var MyGlobal *struct {
	Field string
}

func _() {
	client := http.DefaultClient
	client.Transport = nil // want ".+"

	resolver := net.DefaultResolver
	resolver.PreferGo = true // want ".+"

	x := MyGlobal
	x.Field = "foo"
}

func _(b bool) {
	var f1, f2 *flag.FlagSet
	f1.Usage = nil
	f2.Usage = nil
	if b {
		f1 = flag.CommandLine
	} else {
		f2 = flag.CommandLine
	}
	f1.Usage = nil // want ".+"
	f2.Usage = nil // want ".+"
}
