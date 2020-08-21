package godok

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// Comment adds comment groups to an ast.File
func Comment(f *ast.File) {
	pre := func(cur *astutil.Cursor) bool {
		switch c := cur.Node().(type) {
		case *ast.GenDecl:
			// got doc, don't traverse children
			if c.Doc != nil {
				return false
			}
			// apply inline comments to parenthesized declaration
			if c.Lparen.IsValid() {
				gd := commentGenDeclParen(cur)
				cur.Replace(gd)
				return true
			}
			gd := commentGenDecl(cur)
			cur.Replace(gd)
			return true
		case *ast.FuncDecl:
			// got doc, don't traverse children
			if c.Doc != nil {
				return false
			}
			fd := commentFuncDecl(cur)
			cur.Replace(fd)
			return true
		}
		return true
	}
	astutil.Apply(f, pre, nil)
}

func commentGenDecl(cur *astutil.Cursor) *ast.GenDecl {
	var com *ast.Comment
	pre := func(cur *astutil.Cursor) bool {
		switch n := cur.Node().(type) {
		case *ast.ValueSpec:
			nam := n.Names
			var names []string
			for _, n := range nam {
				names = append(names, n.Name)
			}
			com = &ast.Comment{
				Text: fmt.Sprintf("\n// %s ...\n", strings.Join(names, ", ")),
			}
		case *ast.TypeSpec:
			com = &ast.Comment{
				Text: fmt.Sprintf("\n// %s ...\n", n.Name),
			}
		}
		return true
	}
	post := func(cur *astutil.Cursor) bool {
		cg := &ast.CommentGroup{
			List: []*ast.Comment{},
		}
		if com == nil {
			return true
		}
		p, ok := cur.Parent().(*ast.GenDecl)
		if ok && p.Doc == nil {
			cg.List = append(cg.List, com)
			cur.Replace(cg)
			return true
		}
		return true
	}
	astutil.Apply(cur.Node(), pre, nil)
	return astutil.Apply(cur.Node(), post, nil).(*ast.GenDecl)
}

func commentGenDeclParen(cur *astutil.Cursor) *ast.GenDecl {
	com := make(map[string]*ast.Comment)
	pre := func(cur *astutil.Cursor) bool {
		switch n := cur.Node().(type) {
		case *ast.ValueSpec:
			if n.Doc != nil || n.Comment != nil {
				return false
			}
			n.Comment = &ast.CommentGroup{}
			nam := n.Names
			var names []string
			for _, n := range nam {
				names = append(names, n.Name)
			}
			com[strings.Join(names, ", ")] = &ast.Comment{
				Text: fmt.Sprintf("\t// %s ...\n", strings.Join(names, ", ")),
			}
		}
		return true
	}
	post := func(cur *astutil.Cursor) bool {
		cg := &ast.CommentGroup{
			List: []*ast.Comment{},
		}
		if com == nil {
			return true
		}
		vs, ok := cur.Parent().(*ast.ValueSpec)
		_, cok := cur.Node().(*ast.CommentGroup)
		if ok && cok {
			nam := vs.Names
			var names []string
			for _, n := range nam {
				names = append(names, n.Name)
			}
			cg.List = append(cg.List, com[strings.Join(names, ", ")])
			cur.Replace(cg)
			return true
		}
		return true
	}
	astutil.Apply(cur.Node(), pre, nil)
	return astutil.Apply(cur.Node(), post, nil).(*ast.GenDecl)
}

func commentFuncDecl(cur *astutil.Cursor) *ast.FuncDecl {
	com := make(map[int]*ast.Comment)
	pre := func(cur *astutil.Cursor) bool {
		switch n := cur.Node().(type) {
		case *ast.Ident:
			if _, ok := cur.Parent().(*ast.FuncDecl); !ok {
				return true
			}
			nam := n.Name
			does := ToComment(n.Name)
			com[0] = &ast.Comment{
				Text: fmt.Sprintf("// %s %s\n", nam, does),
			}
		case *ast.FuncType:
			var par []string
			for _, p := range n.Params.List {
				switch typ := p.Type.(type) {
				case *ast.Ident:
					par = append(par, typ.Name)
				case *ast.StarExpr:
					switch x := typ.X.(type) {
					case *ast.Ident:
						par = append(par, fmt.Sprintf("pointer to %s", x.Name))
					}
				}
			}
			if len(par) > 0 {
				com[1] = &ast.Comment{
					Text: fmt.Sprintf("// It takes %s as parameters\n", strings.Join(par, ",")),
				}
			}
			if n.Results != nil {
				var res []string
				for _, r := range n.Results.List {
					switch typ := r.Type.(type) {
					case *ast.Ident:
						res = append(res, typ.Name)
					case *ast.StarExpr:
						switch x := typ.X.(type) {
						case *ast.Ident:
							res = append(res, fmt.Sprintf("pointer to %s", x.Name))
						}
					}
				}
				com[3] = &ast.Comment{
					Text: fmt.Sprintf("// It returns %s\n", strings.Join(res, ",")),
				}
			}
		case *ast.BlockStmt:
			if n == nil {
				com[2] = nil
			}
		}
		return true
	}
	post := func(cur *astutil.Cursor) bool {
		cg := &ast.CommentGroup{
			List: []*ast.Comment{},
		}
		p, ok := cur.Parent().(*ast.FuncDecl)
		if ok && p.Doc == nil {
			for i := 0; i < 4; i++ {
				_, ok := com[i]
				if !ok {
					continue
				}
				cg.List = append(cg.List, com[i])
			}
			cur.Replace(cg)
			return true
		}
		return true
	}
	astutil.Apply(cur.Node(), pre, nil)
	return astutil.Apply(cur.Node(), post, nil).(*ast.FuncDecl)
}
