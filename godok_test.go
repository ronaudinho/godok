package godok_test

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ronaudinho/godok"
)

type test struct {
	name string
	in   *ast.File
	out  bytes.Buffer
}

type want struct {
	line int
	in   string
	out  string
}

func TestComment(t *testing.T) {
	infset := token.NewFileSet()
	tf, err := parseDir(infset, "./tests/", parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tf {
		t.Run(tt.name, func(t *testing.T) {
			godok.Comment(tt.in)
			var buf bytes.Buffer
			format.Node(&buf, infset, tt.in)
			if buf.String() != tt.out.String() {
				// t.Errorf("got:\n%s\nwanted:\n%s", buf.String(), tt.out.String())
				wants := errLine(buf.String(), tt.out.String())
				for _, w := range wants {
					t.Errorf("\nline %d:\ngot: %s\nwanted: %s", w.line, w.in, w.out)
				}
			}
		})
	}
}

func parseDir(infset *token.FileSet, path string, mod parser.Mode) ([]test, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var tf []test
	for _, fil := range dir {
		if !strings.HasSuffix(fil.Name(), ".go") {
			continue
		}
		inpath := filepath.Join(path, fil.Name())
		in, err := parser.ParseFile(infset, inpath, nil, mod)
		if err != nil {
			return tf, err
		}
		outpath := filepath.Join(path, fil.Name()+"dok")
		outfset := token.NewFileSet()
		out, err := parser.ParseFile(outfset, outpath, nil, mod)
		if err != nil {
			return tf, err
		}
		var buf bytes.Buffer
		format.Node(&buf, outfset, out)
		tf = append(tf, test{
			name: fil.Name(),
			in:   in,
			out:  buf,
		})
	}
	return tf, nil
}

func errLine(in, out string) []want {
	var w []want
	ins := strings.Split(in, "\n")
	outs := strings.Split(out, "\n")
	for i := range ins {
		if ins[i] != outs[i] {
			w = append(w, want{
				line: i,
				in:   ins[i],
				out:  outs[i],
			})
		}
	}
	return w
}
