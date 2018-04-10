package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

type Lookup struct {
	Target  string
	Types   map[string]types.Type
	comment *ast.CommentGroup
}

type MyEnum string
const (
	White MyEnum = "white"
	Blue MyEnum = "blue"
)

const hello = `
package main

import "fmt"

type MyEnum string
const (
   White MyEnum = "white"
   Blue MyEnum = "blue"
)

// append
func main() {
        // fmt
        fmt.Println("Hello, world")
        // main
        main, x := 1, 2
        // main
        print(main, x)
        // x
}
// x
`

var helloFile = hello

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", helloFile, parser.ParseComments)
	if err != nil {
		log.Fatal(err) // parse error
	}

	conf := types.Config{Importer: importer.Default()}
	_, err = conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err) // type error
	}

	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.CONST {
				fmt.Printf("At %s,\t%q\n", fset.Position(gd.Pos()), gd.Tok)

				for _, s := range gd.Specs {
					if vs, ok := s.(*ast.ValueSpec); ok {
						for _, name := range vs.Names {
							fmt.Printf("At %s,\t%q\t%v\n", fset.Position(name.Pos()), name.Name, vs.Type.(*ast.Ident).Name)
						}
					}
				}
			}
		}
	}
}
