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
	unused := findUnusedImports(f)
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

type unusedImports map[string]*ast.ImportSpec // package name -> spec

func findUnusedImports(f *ast.File) map[string]*ast.ImportSpec {
	u := make(unusedImports)
	ast.Walk(u, f)
	return map[string]*ast.ImportSpec(u)
}

func (u unusedImports) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return u
	}

	switch v := node.(type) {
	case *ast.ImportSpec:
		// import foo "..."
		if v.Name != nil {
			name := v.Name.Name
			if name != "_" && name != "." {
				u[name] = v
				break
			}
		}

		// import "foo"
		ipath := strings.Trim(v.Path.Value, `"`)
		if ipath == "C" { // ignore cgo when removing unused imports
			break
		}

		name := DeterminePackageName(ipath)
		u[name] = v
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
		delete(u, ident.Name)
	}

	return u
}
