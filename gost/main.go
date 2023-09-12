package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "temp",
	Doc:  "just for test",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		return nil, nil
	},
}

func main() {
	singlechecker.Main(Analyzer)
}
