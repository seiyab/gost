# unboundedRead

Reading an HTTP request body into memory without a size limit can exhaust memory
when a client sends a large body. This rule reports `io.ReadAll` and `io.Copy`
into a `*bytes.Buffer` when their reader originates from `http.Request.Body`.

For example:

```go
// Reported:
data, err := io.ReadAll(req.Body)

// Bounded:
data, err := io.ReadAll(io.LimitReader(req.Body, 1<<20))
```

Choose an appropriate limit for the application. `io.LimitReader` truncates at
the limit; rejecting oversized requests requires separate handling.

The analyzer follows local assignments, branches, `io.TeeReader`,
`io.MultiReader`, and calls to helpers with bodies in the analyzed package.
Generic readers are not assumed to be untrusted. Any request body is considered
potentially unbounded.

This is a deliberately narrow check, not general OOM detection. It does not
track arbitrary reader wrappers, readers returned by helpers, cross-package
helpers, interface dispatch, or other memory sinks. It does not prove that a
limit is small enough or that the body was not already limited upstream.
