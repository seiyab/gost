# urlString

urlString reports unsafe construction of URL-like string.

Example:

```go
s = fmt.Sprintf("http://%s/api", host) // ⚠️ Causes vulnerability if `host` is not validated enough
s = "http://" + host + "/api" // ⚠️
```

Concatenation with directly escaped query components is allowed:

```go
s = "https://example.com/search?q=" + url.QueryEscape(query)
s = "https://example.com/search?q=" + url.QueryEscape(query) + "&page=" + url.QueryEscape(page)
s = "https://example.com/search?" + values.Encode() // values is a url.Values
```

`QueryEscape` must appear in the query, and `Values.Encode` must follow `?`
or `&`. Concatenations containing unescaped dynamic components still report a
diagnostic. Escaped values stored in variables and `fmt.Sprintf` remain outside
this exception.
