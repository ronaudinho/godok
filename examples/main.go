package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ronaudinho/godok"
)

func main() {
	fset := token.NewFileSet()
	pkg, err := parseDir(fset, "../tests/", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, p := range pkg {
		for nam, fil := range p.Files {
			fmt.Println("adding comments to", nam)
			godok.Comment(fil)
			var buf bytes.Buffer
			format.Node(&buf, fset, fil)
			fmt.Println(buf.String())
			// writeToFile(nam, buf)
		}
	}
}

func parseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mod parser.Mode) (map[string]*ast.Package, error) {
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var ferr error
	pkgs := make(map[string]*ast.Package)
	for _, d := range lst {
		if d.IsDir() {
			pkg, err := parseDir(fset, filepath.Join(path, d.Name()), filter, mod)
			if err == nil {
				for k, v := range pkg {
					pkgs[k] = v
				}
			} else if ferr == nil {
				ferr = err
			}
		}
		// omit test files
		if strings.HasSuffix(d.Name(), "_test.go") && (filter == nil || filter(d)) {
			continue
		}
		// omit godok files
		if strings.HasSuffix(d.Name(), ".godok") && (filter == nil || filter(d)) {
			continue
		}
		if strings.HasSuffix(d.Name(), ".go") && (filter == nil || filter(d)) {
			filename := filepath.Join(path, d.Name())
			src, err := parser.ParseFile(fset, filename, nil, mod)
			if err == nil {
				nam := src.Name.Name
				pkg, found := pkgs[nam]
				if !found {
					pkg = &ast.Package{
						Name:  nam,
						Files: make(map[string]*ast.File),
					}
					pkgs[nam] = pkg
				}
				pkg.Files[filename] = src
			} else if ferr == nil {
				ferr = err
			}
		}
	}
	return pkgs, ferr
}

func writeToFile(nam string, buf bytes.Buffer) {
	f, _ := os.Open(nam)
	s, _ := f.Stat()
	mod := s.Mode()
	f.Close()
	err := ioutil.WriteFile(nam, buf.Bytes(), mod)
	if err != nil {
		fmt.Printf("write to file: %s: %v", nam, err)
	}
}
