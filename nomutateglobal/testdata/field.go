package testdata

import "net/http"

type MyClient struct {
	Client *http.Client
}

// https://pkg.go.dev/vuln/GO-2024-2618
func New(client *http.Client) *MyClient {
	c := &MyClient{}
	c.Client.Transport = nil
	if client != nil {
		c.Client = http.DefaultClient
	}
	c.Client.Timeout = 0 // want ".+"
	return c
}
