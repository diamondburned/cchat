package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/repository"
)

func generateTypeAlises(aliases []repository.TypeAlias) jen.Code {
	sort.Slice(aliases, func(i, j int) bool {
		return aliases[i].Name < aliases[j].Name
	})

	var stmt = new(jen.Statement)

	for _, alias := range aliases {
		if !alias.Comment.IsEmpty() {
			stmt.Comment(alias.Comment.GoString(1))
			stmt.Line()
		}

		stmt.Type().Id(alias.Name).Op("=").Add(generateType(alias))
		stmt.Line()
		stmt.Line()
	}

	return stmt
}

type qualer interface {
	Qual() (path, name string)
}

func generateType(t qualer) jen.Code {
	path, name := t.Qual()
	if path == "" {
		return jen.Id(name)
	}
	return jen.Qual(path, name)
}
