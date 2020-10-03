package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/cmd/internal/cchat-generator/genutils"
	"github.com/diamondburned/cchat/repository"
)

func generateStructs(structs []repository.Struct) jen.Code {
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Name < structs[j].Name
	})

	var stmt = new(jen.Statement)

	for _, s := range structs {
		stmt.Add(generateStruct(s))
		stmt.Line()
		stmt.Line()
	}

	return stmt
}

func generateErrorStructs(errStructs []repository.ErrorStruct) jen.Code {
	sort.Slice(errStructs, func(i, j int) bool {
		return errStructs[i].Name < errStructs[j].Name
	})

	var stmt = new(jen.Statement)

	for _, errStruct := range errStructs {
		stmt.Add(generateStruct(errStruct.Struct))
		stmt.Line()
		stmt.Line()

		var recv = genutils.RecvName(errStruct.Name)

		stmt.Func()
		stmt.Params(jen.Id(recv).Id(errStruct.Name))
		stmt.Id("Error").Params().String()
		stmt.BlockFunc(func(g *jen.Group) {
			g.Return(jen.Qual("fmt", "Sprintf").CallFunc(func(g *jen.Group) {
				g.Lit(errStruct.ErrorString.Format)

				for _, field := range errStruct.ErrorString.Fields {
					g.Add(jen.Id(recv).Dot(field))
				}
			}))
		})

		stmt.Line()
		stmt.Line()

		if wrap := errStruct.Wrapped(); wrap != "" {
			stmt.Func()
			stmt.Params(jen.Id(recv).Id(errStruct.Name))
			stmt.Id("Unwrap").Params().Error()
			stmt.BlockFunc(func(g *jen.Group) {
				g.Return(jen.Id(recv).Dot(wrap))
			})
			stmt.Line()
			stmt.Line()
		}
	}

	return stmt
}

func generateStruct(s repository.Struct) jen.Code {
	var stmt = new(jen.Statement)
	if !s.Comment.IsEmpty() {
		stmt.Comment(s.Comment.GoString(1))
		stmt.Line()
	}

	stmt.Type().Id(s.Name).StructFunc(func(group *jen.Group) {
		for _, field := range s.Fields {
			var stmt = new(jen.Statement)
			if field.Name != "" {
				stmt.Id(field.Name)
			}
			stmt.Add(genutils.GenerateType(field))
			group.Add(stmt)
		}
	})

	return stmt
}
