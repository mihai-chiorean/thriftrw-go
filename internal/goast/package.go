package goast

import (
	"go/build"
	"path/filepath"
)

// determinePackageName determines the name of the package at the given import
// path.
func determinePackageName(importPath string) string {
	// TODO(abg): This can be a lot faster by using build.FindOnly and parsing one
	// of the .go files in the directory with parser.PackageClauseOnly set. See
	// how goimports does this:
	// https://github.com/golang/tools/blob/0e9f43fcb67267967af8c15d7dc54b373e341d20/imports/fix.go#L284

	pkg, err := build.Import(importPath, "", 0)
	if err != nil {
		return filepath.Base(importPath)
	}
	return pkg.Name
}
