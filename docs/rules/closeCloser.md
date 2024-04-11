# closeCloser

closeCloser reports closer that isn't closed.

For example, variable `f` in following code is reported.

```go
func myFunc() error {
	f, err := os.Open("file.txt") // ⚠️
    if err != nil {
        return err
    }
    _, err = f.Write(([]byte)("Hello, world"))
    return err
}
```

Variables that is only used as method receiver is reported. For example, returned variable will be closed by calller so closeCloser doesn't report.

```go
func example1() (*os.File, error) {
    f, err := os.Open("file.txt") // It's OK.
    if err != nil {
        return nil, err
    }
    return f, nil
}

func example2() error {
    f, err := os.Open("file.txt") // It's not reported. someFunc might close f.
    if err != nil {
        return  err
    }
    if err := someFunc(f); err != nil {
        return err
    }
    return nil
}
```
