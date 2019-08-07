package splitfile_test

import (
	"testing"

	"github.com/mccurdyc/splitfile"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestFromStringLiterals demonstrates how to test an analysis using
// a table of string literals for each test case.
//
// Such tests are typically quite compact.
func TestFromStringLiterals(t *testing.T) {

	for _, test := range [...]struct {
		desc    string
		pkgpath string
		files   map[string]string
	}{
		{
			desc:    "SimpleTest",
			pkgpath: "main",
			files: map[string]string{"main/main.go": `package main
func main() {
	println("hello") // not split found
	print("goodbye") // not split found
}`,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			dir, cleanup, err := analysistest.WriteFiles(test.files)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()
			analysistest.Run(t, dir, splitfile.Analyzer, test.pkgpath)
		})
	}
}

// TestFromFileSystem demonstrates how to test an analysis using input
// files stored in the file system.
//
// These tests have the advantages that test data can be edited
// directly, and that files named in error messages can be opened.
// However, they tend to spread a small number of lines of text across a
// rather deep directory hierarchy, and obscure similarities among
// related tests, especially when tests involve multiple packages, or
// multiple variants of a single scenario.
func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, splitfile.Analyzer, "a") // loads testdata/src/a/a.go.
}
