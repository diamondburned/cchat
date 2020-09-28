package main

import (
	"sort"

	"github.com/dave/jennifer/jen"
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
					stmt.Params(generateFuncParams(method.Parameters, false)...)
					stmt.Params(generateFuncParams(method.Returns, method.ReturnError)...)
				case repository.SetterMethod:
					stmt.Params(generateFuncParams(method.Parameters, false)...)
				case repository.IOMethod:
					stmt.Params(generateFuncParams(method.Parameters, false)...)
					stmt.Params(generateFuncParamErr(method.ReturnValue, method.ReturnError)...)
					stmt.Comment("// Blocking")
				case repository.ContainerMethod:
					stmt.Params(generateContainerFuncParams(method)...)
					stmt.Params(generateContainerFuncReturns(method)...)
				case repository.AsserterMethod:
					stmt.Params()
					stmt.Params(generateType(method))
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

func generateFuncParamErr(param repository.NamedType, genErr bool) []jen.Code {
	stmt := make([]jen.Code, 0, 2)

	if !param.IsZero() {
		stmt = append(stmt, generateFuncParam(param))
	}

	if genErr {
		if param.Name == "" {
			stmt = append(stmt, jen.Error())
		} else {
			stmt = append(stmt, jen.Err().Error())
		}
	}

	return stmt
}

func generateFuncParam(param repository.NamedType) jen.Code {
	if param.Name == "" {
		return generateType(param)
	}
	return jen.Id(param.Name).Add(generateType(param))
}

func generateFuncParams(params []repository.NamedType, withError bool) []jen.Code {
	if len(params) == 0 {
		return nil
	}

	var stmt jen.Statement
	for _, param := range params {
		stmt.Add(generateFuncParam(param))
	}

	if withError {
		if params[0].Name != "" {
			stmt.Add(jen.Err().Error())
		} else {
			stmt.Add(jen.Error())
		}
	}

	return stmt
}

func generateContainerFuncReturns(method repository.ContainerMethod) []jen.Code {
	var stmt jen.Statement

	if method.HasStopFn {
		stmt.Add(jen.Id("stop").Func().Params())
	}
	stmt.Add(jen.Err().Error())

	return stmt
}

func generateContainerFuncParams(method repository.ContainerMethod) []jen.Code {
	var stmt jen.Statement

	if method.HasContext {
		stmt.Qual("context", "Context")
	}
	stmt.Add(generateType(method))

	return stmt
}
