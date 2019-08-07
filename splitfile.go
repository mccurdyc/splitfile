package splitfile

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "splitfile",
	Doc:      "checks for clean splits of files in packages.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// removeImportDecls filters out the import statements from the declarations slice.
func removeImportDecls(decls []ast.Decl) []ast.Decl {
	for i, d := range decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok || decl.Tok != token.IMPORT {
			return decls[i:]
		}
	}

	return decls
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		decls := removeImportDecls(file.Decls)

		fmt.Println(decls)

		pass.Reportf(file.Pos(), "split found %q",
			render(pass.Fset, file))
	}
	return nil, nil
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
