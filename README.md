# :ghost: gost

gost is a static checker for Golang. It contains aggressive rules that aren't afraid false-positive as long as diagnostics are informative.

# Installation

I recommend to use gost via [reviewdog](https://github.com/reviewdog/reviewdog).
Complete example configuration:

- [.reviewdog.yml](./.reviewdog.yml)
- [.github/workflows/reviewdog.yml](./.github/workflows/reviewdog.yml)

To run it locally, run following:

```sh
# install
go install github.com/seiyab/gost@latest

# run
go vet -vettool="$(which gost)" ./...
```

# Analyzers

<!-- prettier-ignore -->
| name | description | practical discovery | inspired by |
| :----------------- | :--------------------------------------------------------------------------- | :---------------------------------------------- | :---------------------------------------------- |
| [closeCloser](./docs/rules/closeCloser.md) | report that closer isn't closed | https://github.com/reviewdog/reviewdog/pull/1692 | |
| [multipleErrors](./docs/rules/multipleErrors.md) | report suspicious error concatenation | https://github.com/opentofu/opentofu/issues/539 | |
| noDiscardError | report that error is discarded | https://github.com/cli/cli/issues/8026 | |
| noMutateGlobal | reports indirect mutation of global variable | | https://pkg.go.dev/vuln/GO-2024-2618 |
| openFileFlag | report suspicious combination of flags in `os.OpenFile()` | https://github.com/anchore/go-logger/pull/13 | |
| preferFilepath | report misuse of `"path"` package where `"path/filepath"` should be suitable | https://github.com/anchore/grype/pull/1767 | |
| sliceInitialLength | reports confusion between slice length and capacity | https://github.com/beego/beego/pull/5631 | https://github.com/dominikh/go-tools/issues/112 |
| [urlString](./docs/rules/urlString.md) | urlString reports unsafe construction of URL-like string. | | https://github.com/dominikh/go-tools/issues/730 |
| wrapError | report senseless error wrapping | | |
