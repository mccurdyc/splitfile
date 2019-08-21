package splitfile_test

import (
	"testing"

	"github.com/mccurdyc/splitfile"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestFromFileSystem demonstrates how to test an analysis using input
// files stored in the file system.
//
// These tests have the advantages that test data can be edited
// directly, and that files named in error messages can be opened.
// However, they tend to spread a small number of lines of text across a
// rather deep directory hierarchy, and obscure similarities among
// related tests, especially when tests involve multiple packages, or
// multiple variants of a single scenario.
// func TestFromFileSystem(t *testing.T) {
// 	testdata := analysistest.TestData()
// 	analysistest.Run(t, testdata, splitfile.Analyzer, "abc")
// }

// TestFromStringLiterals demonstrates how to test an analysis using
// a table of string literals for each test case.
//
// Such tests are typically quite compact.
func TestFromStringLiterals(t *testing.T) {
	for _, test := range [...]struct {
		name    string
		pkgpath string
		files   map[string]string
	}{
		{
			name:    "single type",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int
`,
			},
		},
		{
			name:    "small file un-related types",
			pkgpath: "ab",
			files: map[string]string{"ab/ab.go": `package ab
type a int
type b int
`,
			},
		},
		{
			name:    "related types through struct fields",
			pkgpath: "ab",
			files: map[string]string{"ab/ab.go": `package ab
type a int

type b struct {
  a a
}
`,
			},
		},
		{
			name:    "related type and function param",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func fa(a a) {
}
`,
			},
		},
		{
			name:    "related type and function usage",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func fa() {
  _ = a(123)
}
`,
			},
		},
		{
			name:    "related type and multiple functions",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func fa(a a) {
}

func faa() {
  _ = a(123)
}
`,
			},
		},
		{
			name:    "related type and method receiver (non-pointer)",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func (a a) ma() {
}
`,
			},
		},
		{
			name:    "related type and method receiver (pointer)",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func (a *a) ma() {
}
`,
			},
		},
		{
			name:    "related type and method param",
			pkgpath: "ab",
			files: map[string]string{"ab/ab.go": `package ab
type a int
type b int

func (b b) mb(a a) {
}
`,
			},
		},
		{
			name:    "related type and method usage",
			pkgpath: "ab",
			files: map[string]string{"ab/ab.go": `package ab
type a int
type b int

func (b b) mb() {
  _ = a(123)
}
`,
			},
		},
		{
			name:    "related type and multiple methods",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
type a int

func ma(a a) {
}

func maa() {
  _ = a(123)
}
`,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir, cleanup, err := analysistest.WriteFiles(test.files)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()
			analysistest.Run(t, dir, splitfile.Analyzer, test.pkgpath)
		})
	}
}
