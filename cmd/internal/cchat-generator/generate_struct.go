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
		stmt.Block(jen.Return(generateTmplString(errStruct.ErrorString, recv)))

		stmt.Line()
		stmt.Line()

		if wrap := errStruct.Wrapped(); wrap != "" {
			stmt.Func()
			stmt.Params(jen.Id(recv).Id(errStruct.Name))
			stmt.Id("Unwrap").Params().Error()
			stmt.Block(jen.Return(jen.Id(recv).Dot(wrap)))
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

	if !s.Stringer.IsEmpty() {
		stmt.Line()
		stmt.Line()

		var recv = genutils.RecvName(s.Name)

		if !s.Stringer.Comment.IsEmpty() {
			stmt.Comment(s.Stringer.Comment.GoString(1))
			stmt.Line()
		}

		stmt.Func()
		stmt.Params(jen.Id(recv).Id(s.Name))
		stmt.Id("String").Params().String()
		stmt.BlockFunc(func(g *jen.Group) {
			if s.Stringer.Format == "%s" {
				g.Return(jen.Id(recv).Dot(s.Stringer.Fields[0]))
				return
			}

			g.Return(generateTmplString(s.Stringer.TmplString, recv))
		})
	}

	return stmt
}

func generateTmplString(tmpl repository.TmplString, recv string) jen.Code {
	return jen.Qual("fmt", "Sprintf").CallFunc(func(g *jen.Group) {
		g.Lit(tmpl.Format)

		for _, field := range tmpl.Fields {
			g.Add(jen.Id(recv).Dot(field))
		}
	})
}
