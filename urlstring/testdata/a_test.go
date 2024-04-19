package testdata_test

import "fmt"

func _(host, userUUID, query string) {
	_ = fmt.Sprintf("http://%s/api/v1/users/%s/comments?q=%s", host, userUUID, query)

	_ = "https://" + host + "/api"

	_ = "https://example.com"
}
