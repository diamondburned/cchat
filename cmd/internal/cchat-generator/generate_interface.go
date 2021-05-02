package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/cmd/internal/cchat-generator/genutils"
	"github.com/diamondburned/cchat/repository"
)

func generateInterfaces(ifaces []repository.Interface) jen.Code {
	sort.Slice(ifaces, func(i, j int) bool {
		return ifaces[i].Name < ifaces[j].Name
	})

	var stmt = new(jen.Statement)

	for _, iface := range ifaces {
		if !iface.Comment.IsEmpty() {
			stmt.Comment(iface.Comment.GoString(1))
			stmt.Line()
		}

		stmt.Type().Id(iface.Name).InterfaceFunc(func(group *jen.Group) {
			if len(iface.Embeds) > 0 {
				for _, embed := range iface.Embeds {
					if !embed.Comment.IsEmpty() {
						group.Comment(embed.Comment.GoString(1))
					}
					group.Id(embed.InterfaceName)
				}

				group.Line()
			}

			// Put Asserter methods last.
			sort.SliceStable(iface.Methods, func(i, j int) bool {
				_, assert := iface.Methods[i].(repository.AsserterMethod)
				return !assert
			})

			// Boolean to only write the Asserter comment once.
			var writtenComment bool

			for _, method := range iface.Methods {
				if !writtenComment {
					if _, ok := method.(repository.AsserterMethod); ok {
						group.Line()
						group.Comment("// Asserters.")
						group.Line()
						writtenComment = true
					}
				}

				var stmt = new(jen.Statement)
				if comment := method.UnderlyingComment(); !comment.IsEmpty() {
					stmt.Comment(comment.GoString(1))
					stmt.Line()
				}

				stmt.Id(method.UnderlyingName())

				switch method := method.(type) {
				case repository.GetterMethod:
					stmt.Params(generateFuncParams(method.Parameters, "")...)
					stmt.Params(generateFuncParams(method.Returns, method.ErrorType)...)
				case repository.SetterMethod:
					stmt.Params(generateFuncParams(method.Parameters, "")...)
					stmt.Params(generateFuncParamsErr(repository.NamedType{}, method.ErrorType)...)
				case repository.IOMethod:
					stmt.Params(generateFuncParamsCtx(method.Parameters, "")...)
					stmt.Params(generateFuncParamsErr(method.ReturnValue, method.ErrorType)...)
					var comment = "Blocking"
					if method.Disposer {
						comment += ", Disposer"
					}
					stmt.Comment("// " + comment)
				case repository.ContainerMethod:
					stmt.Params(generateContainerFuncParams(method)...)
					stmt.Params(generateContainerFuncReturns(method)...)
				case repository.AsserterMethod:
					stmt.Params()
					stmt.Params(genutils.GenerateType(method))
					stmt.Comment("// Optional")
				default:
					continue
				}

				group.Add(stmt)
			}
		})

		stmt.Line()
		stmt.Line()
	}

	return stmt
}

func generateFuncParamsErr(param repository.NamedType, errorType string) []jen.Code {
	stmt := make([]jen.Code, 0, 2)

	if !param.IsZero() {
		stmt = append(stmt, generateFuncParam(param))
	}

	if errorType != "" {
		if param.Name == "" {
			stmt = append(stmt, jen.Id(errorType))
		} else {
			stmt = append(stmt, jen.Err().Id(errorType))
		}
	}

	return stmt
}

func generateFuncParam(param repository.NamedType) jen.Code {
	if param.Name == "" {
		return genutils.GenerateType(param)
	}
	return jen.Id(param.Name).Add(genutils.GenerateType(param))
}

func generateFuncParamsCtx(params []repository.NamedType, errorType string) []jen.Code {
	var name string
	if len(params) > 0 && params[0].Name != "" {
		name = "ctx"
	}

	p := []repository.NamedType{{Name: name, Type: "context.Context"}}
	p = append(p, params...)

	return generateFuncParams(p, errorType)
}

func generateFuncParams(params []repository.NamedType, errorType string) []jen.Code {
	if len(params) == 0 {
		return nil
	}

	var stmt jen.Statement
	for _, param := range params {
		stmt.Add(generateFuncParam(param))
	}

	if errorType != "" {
		if params[0].Name != "" {
			stmt.Add(jen.Err().Id(errorType))
		} else {
			stmt.Add(jen.Id(errorType))
		}
	}

	return stmt
}

func generateContainerFuncReturns(method repository.ContainerMethod) []jen.Code {
	return []jen.Code{jen.Error()}
}

func generateContainerFuncParams(method repository.ContainerMethod) []jen.Code {
	var stmt jen.Statement

	stmt.Qual("context", "Context")
	stmt.Add(genutils.GenerateType(method))

	return stmt
}
