package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/cmd/internal/cchat-generator/genutils"
	"github.com/diamondburned/cchat/repository"
)

const OutputDir = "."

func init() {
	log.SetFlags(0)
}

var comment = repository.Comment{Raw: `
	Package empty provides no-op asserter method implementations of interfaces
	in cchat.
`}

func main() {
	pk, ok := repository.Main[repository.RootPath]
	if !ok {
		log.Fatalf("Failed to find main namespace %q\n", repository.RootPath)
	}

	gen := jen.NewFile("empty")
	gen.HeaderComment("DO NOT EDIT: THIS FILE IS GENERATED!")
	gen.PackageComment(comment.GoString(1))

	for _, iface := range pk.Interfaces {
		// Skip structs without asserter methods.
		if !hasAsserter(iface) {
			continue
		}

		gen.Commentf("%[1]s provides no-op asserters for cchat.%[1]s.", iface.Name)
		gen.Type().Id(iface.Name).Struct()
		gen.Line()

		for _, method := range iface.Methods {
			am, ok := method.(repository.AsserterMethod)
			if !ok {
				continue
			}

			name := fmt.Sprintf("As%s", am.ChildType)
			gen.Comment(fmt.Sprintf("%s returns nil.", name))

			stmt := jen.Func()
			stmt.Parens(jen.Id(iface.Name))
			stmt.Id(fmt.Sprintf("As%s", am.ChildType))
			stmt.Params()
			stmt.Add(genutils.GenerateExternType(am))
			stmt.Values(jen.Return(jen.Nil()))

			gen.Add(stmt)
		}

		gen.Line()
	}

	f, err := os.Create("empty.go")
	if err != nil {
		log.Fatalln("Failed to create output file:", err)
	}
	defer f.Close()

	if err := gen.Render(f); err != nil {
		log.Fatalln("Failed to render output:", err)
	}
}

func hasAsserter(iface repository.Interface) bool {
	for _, method := range iface.Methods {
		if _, isA := method.(repository.AsserterMethod); isA {
			return true
		}
	}

	return false
}
