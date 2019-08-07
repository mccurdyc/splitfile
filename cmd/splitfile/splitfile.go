package main

import (
	"github.com/mccurdyc/splitfile"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(splitfile.Analyzer)
}
