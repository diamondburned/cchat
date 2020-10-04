package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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
	in cchat's root and text packages.
`}

func main() {
	gen := jen.NewFile("empty")
	gen.HeaderComment("DO NOT EDIT: THIS FILE IS GENERATED!")
	gen.PackageComment(comment.GoString(1))

	for pkgpath, pk := range repository.Main {
		gen.ImportName(pkgpath, path.Base(pkgpath))

		for _, iface := range pk.Interfaces {
			// Skip structs without asserter methods.
			if !hasAsserter(iface) {
				continue
			}

			var ifaceName = newIfaceName(pkgpath, iface)

			gen.Commentf("%[1]s provides no-op asserters for cchat.%[1]s.", ifaceName)
			gen.Type().Id(ifaceName).Struct()
			gen.Line()

			for _, embed := range iface.Embeds {
				if iface := pk.Interface(embed.InterfaceName); iface != nil {
					genIfaceMethods(gen, *iface, ifaceName, pkgpath)
				}
			}

			genIfaceMethods(gen, iface, ifaceName, pkgpath)

			gen.Line()
		}
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

func newIfaceName(pkgpath string, iface repository.Interface) string {
	if pkgpath == repository.RootPath {
		return iface.Name
	} else {
		return strings.Title(repository.TrimRoot(pkgpath)) + iface.Name
	}
}

func genIfaceMethods(gen *jen.File, iface repository.Interface, ifaceName, pkgpath string) {
	for _, method := range iface.Methods {
		am, ok := method.(repository.AsserterMethod)
		if !ok {
			continue
		}

		name := fmt.Sprintf("As%s", am.ChildType)
		gen.Comment(fmt.Sprintf("%s returns nil.", name))

		stmt := jen.Func()
		stmt.Parens(jen.Id(ifaceName))
		stmt.Id(name)
		stmt.Params()
		stmt.Add(genutils.GenerateExternType(pkgpath, am))
		stmt.Values(jen.Return(jen.Nil()))

		gen.Add(stmt)
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
