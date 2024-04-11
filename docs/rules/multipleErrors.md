# multipleErrors

multipleErrors reports suspicious error concatenation.

[errors.Join](https://pkg.go.dev/errors#Join), [multierror.Append](https://pkg.go.dev/github.com/hashicorp/go-multierror#Append) and [multierr.Append](https://pkg.go.dev/go.uber.org/multierr#Append) are supported.

Example:

```go
errs = errors.Join(err) // ⚠️ Call with single error doesn't make sense

errs = errors.Join(e, err) // ⚠️ Original one goes away. Typically, it should be like following
errs = errors.Join(errs, err)
```
