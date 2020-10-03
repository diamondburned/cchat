package genutils

import (
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/diamondburned/cchat/repository"
)

type Qualer interface {
	Qual() (path, name string)
}

// GenerateType generates a jen.Qual from the given Qualer.
func GenerateType(t Qualer) jen.Code {
	path, name := t.Qual()
	if path == "" {
		return jen.Id(name)
	}
	return jen.Qual(path, name)
}

// GenerateExternType generates a jen.Qual from the given Qualer, except if
// the path is empty, RootPath is used instead.
func GenerateExternType(t Qualer) jen.Code {
	path, name := t.Qual()
	if path == "" {
		return jen.Qual(repository.RootPath, name)
	}
	return jen.Qual(path, name)
}

// RecvName is used to get the receiver variable name. It returns the first
// letter lower-cased. It does NOT do length checking. It only works with ASCII.
func RecvName(name string) string {
	return string(unicode.ToLower(rune(name[0])))
}
