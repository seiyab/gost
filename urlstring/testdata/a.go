package testdata

import (
	"fmt"
	url "net/url"
)

func _(host, userUUID, query string) {
	_ = fmt.Sprintf("http://%s/api/v1/users/%s/comments?q=%s", host, userUUID, query) // want ".+"

	_ = "https://" + host + "/api" // want ".+"

	_ = "https://example.com"
	_ = "https://example.com/search?q=" + url.QueryEscape(query)
	_ = "https://example.com/search?" + url.QueryEscape(query) + "=" + url.QueryEscape(userUUID)
	_ = "https://example.com/search?q=" + (url.QueryEscape(query)) + "&user=" + url.QueryEscape(userUUID)
	values := url.Values{"q": {query}}
	_ = "https://example.com/search?" + values.Encode()
	_ = "https://example.com/search?q=x&" + values.Encode()
	_ = "https://example.com/search?q=" + url.QueryEscape(query) + "&user=" + userUUID // want ".+"
	_ = "https://example.com/" + url.QueryEscape(query)                                // want ".+"
	_ = "https://" + url.QueryEscape(host)                                             // want ".+"
	_ = "https://example.com/#?q=" + url.QueryEscape(query)                            // want ".+"
	_ = "https://example.com/?q=" + values.Encode()                                    // want ".+"
}
