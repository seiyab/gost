# urlString

urlString reports unsafe construction of URL-like string.

Example:

```go
s = fmt.Sprintf("http://%s/api", host) // ⚠️ Causes vulnerability if `host` is not validated enough
s = "http://" + host + "/api" // ⚠️
```
