package testdata

import "fmt"

func _(host, userUUID, query string) {
	_ = fmt.Sprintf("http://%s/api/v1/users/%s/comments?q=%s", host, userUUID, query) // want ".+"

	_ = "https://" + host + "/api" // want ".+"

	_ = "https://example.com"
}
