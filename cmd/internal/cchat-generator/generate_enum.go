package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/cmd/internal/cchat-generator/genutils"
	"github.com/diamondburned/cchat/repository"
)

func generateEnums(enums []repository.Enumeration) jen.Code {
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name < enums[j].Name
	})

	var stmt = new(jen.Statement)

	for _, enum := range enums {
		if !enum.Comment.IsEmpty() {
			stmt.Comment(enum.Comment.GoString(1))
			stmt.Line()
		}

		stmt.Type().Id(enum.Name).Id(enum.GoType())
		stmt.Line()
		stmt.Line()

		stmt.Const().DefsFunc(func(group *jen.Group) {
			for i, value := range enum.Values {
				var c jen.Statement
				if !value.Comment.IsEmpty() {
					c.Comment(value.Comment.GoString(2))
					c.Line()
				}

				c.Id(enum.Name + value.Name)

				if i == 0 {
					c.Id(enum.Name).Op("=").Iota()
				}

				group.Add(&c)
			}
		})

		stmt.Line()
		stmt.Line()

		var recv = genutils.RecvName(enum.Name)

		if enum.Bitwise {
			fn := stmt.Func()
			fn.Params(jen.Id(recv).Id(enum.Name))
			fn.Id("Has")
			fn.Params(jen.Id("has").Id(enum.Name))
			fn.Bool()
			fn.BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(recv).Op("&").Id("has").Op("==").Id("has"))
			})
		} else {
			fn := stmt.Func()
			fn.Params(jen.Id(recv).Id(enum.Name))
			fn.Id("Is")
			fn.Params(jen.Id("is").Id(enum.Name))
			fn.Bool()
			fn.BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(recv).Id("==").Id("is"))
			})
		}

		stmt.Line()
		stmt.Line()
	}

	return stmt
}
