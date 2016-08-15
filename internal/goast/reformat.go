package goast

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

var cfg = printer.Config{
	Mode:     printer.UseSpaces | printer.TabIndent,
	Tabwidth: 8,
}

// Reformat reformats the given Go source code, optionally transforming it before rewriting.
func Reformat(fset *token.FileSet, path string, contents []byte, transforms ...Transformer) ([]byte, error) {
	if fset == nil {
		fset = token.NewFileSet()
	}

	f, err := parser.ParseFile(fset, path, contents, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, t := range transforms {
		if err := t(fset, f); err != nil {
			return nil, err
		}
	}

	var buff bytes.Buffer
	if err := cfg.Fprint(&buff, fset, f); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// A Transformer transform the parsed source code.
type Transformer func(*token.FileSet, *ast.File) error
