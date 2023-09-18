# :ghost: gost
gost is a toy project where I
- experimentally implement static code checkers for Golang
- try the checkers into popular projects to find possible false-positives, insights and bugs

## Usage
NOTE: It's just note for myself. This is an experimental project so I don't recommend to utilize it.
```sh
# install
go install github.com/seiyab/gost/gost@latest

# run
go vet -vettool="$(which gost)" ./...
```
