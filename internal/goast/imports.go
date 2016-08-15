package goast

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// RemoveUnusedImports is a Transformer that removes unused imports from parsed
// Go files.
var RemoveUnusedImports Transformer = func(fset *token.FileSet, f *ast.File) error {
	unused, err := findUnusedImports(f)
	if err != nil {
		return err
	}

	for _, spec := range unused {
		name := ""
		ipath := strings.Trim(spec.Path.Value, `"`)
		if spec.Name != nil {
			name = spec.Name.Name
		}

		astutil.DeleteNamedImport(fset, f, name, ipath)
	}

	return nil
}

type unusedImports struct {
	unused map[string]*ast.ImportSpec // package name -> spec
	err    error
}

func findUnusedImports(f *ast.File) (map[string]*ast.ImportSpec, error) {
	finder := &unusedImports{unused: make(map[string]*ast.ImportSpec)}
	ast.Walk(finder, f)
	return finder.unused, finder.err
}

func (u *unusedImports) Visit(node ast.Node) ast.Visitor {
	// Failed. Stop running.
	if u.err != nil {
		return nil
	}
	if node == nil {
		return u
	}

	switch v := node.(type) {
	case *ast.ImportSpec:
		// import foo "..."
		if v.Name != nil {
			u.unused[v.Name.Name] = v
			break
		}

		// import "foo"
		ipath := strings.Trim(v.Path.Value, `"`)
		if ipath == "C" { // ignore cgo when removing unused imports
			break
		}

		name, err := determinePackageName(ipath)
		if err != nil {
			u.err = err
			return nil
		}

		u.unused[name] = v
	case *ast.SelectorExpr: // foo.Bar
		ident, ok := v.X.(*ast.Ident)
		if !ok { // foo is a complex expression
			break
		}

		if ident.Obj != nil {
			// foo is an object, not a package
			break
		}

		// TODO(abg): Do we need to check globals?
		delete(u.unused, ident.Name)
	}

	return u
}
