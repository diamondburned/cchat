package repository

import (
	"fmt"
	"path"
	"strings"
)

// MakePath returns RootPath joined with relPath.
func MakePath(relPath string) string {
	return path.Join(RootPath, relPath)
}

// MakeQual returns a qualified identifier that is the full path and name of
// something.
func MakeQual(relPath, name string) string {
	return fmt.Sprintf("(%s).%s", MakePath(relPath), name)
}

// TrimRoot trims the root path and returns the path relative to root.
func TrimRoot(fullPath string) string {
	return strings.TrimPrefix(strings.TrimPrefix(fullPath, RootPath), "/")
}

// Packages maps Go module paths to packages.
type Packages map[string]Package

type Package struct {
	Comment Comment

	Enums        []Enumeration
	TypeAliases  []TypeAlias
	Structs      []Struct
	ErrorStructs []ErrorStruct
	Interfaces   []Interface
}

// Interface finds an interface. Nil is returned if none is found.
func (p Package) Interface(name string) *Interface {
	for i, iface := range p.Interfaces {
		if iface.Name == name {
			return &p.Interfaces[i]
		}
	}
	return nil
}

// Struct finds a struct or an error struct. Nil is returned if none is found.
func (p Package) Struct(name string) *Struct {
	for i, sstruct := range p.Structs {
		if sstruct.Name == name {
			return &p.Structs[i]
		}
	}
	for i, estruct := range p.ErrorStructs {
		if estruct.Name == name {
			return &p.ErrorStructs[i].Struct
		}
	}
	return nil
}

// Enum finds an enumeration. Nil is returned if none is found.
func (p Package) Enum(name string) *Enumeration {
	for i, enum := range p.Enums {
		if enum.Name == name {
			return &p.Enums[i]
		}
	}
	return nil
}

// TypeAlias finds a type alias. Nil is returned if none is found.
func (p Package) TypeAlias(name string) *TypeAlias {
	for i, alias := range p.TypeAliases {
		if alias.Name == name {
			return &p.TypeAliases[i]
		}
	}
	return nil
}

// FindType finds any type. Nil is returned if nothing is found; a pointer to
// any of the type is returned if name is found.
func (p Package) FindType(name string) interface{} {
	if iface := p.Interface(name); iface != nil {
		return iface
	}
	if sstr := p.Struct(name); sstr != nil {
		return sstr
	}
	if enum := p.Enum(name); enum != nil {
		return enum
	}
	if alias := p.TypeAlias(name); alias != nil {
		return alias
	}
	return nil
}

// NamedType is an optionally named value with a type.
type NamedType struct {
	Name string // optional
	Type string // import/path.Type OR (import/path).Type
}

// Qual splits the type name into the path and type name. Refer to TypeQual.
func (t NamedType) Qual() (path, typeName string) {
	return TypeQual(t.Type)
}

// IsZero is true if t.Type is empty.
func (t NamedType) IsZero() bool {
	return t.Type == ""
}

// TypeQual splits the type name into path and type name. It accepts inputs that
// are similar to the example below:
//
//     string
//     context.Context
//     github.com/diamondburned/cchat/text.Rich
//    (github.com/diamondburned/cchat/text).Rich
//
func TypeQual(typePath string) (path, typeName string) {
	parts := strings.Split(typePath, ".")
	if len(parts) > 1 {
		path = strings.Join(parts[:len(parts)-1], ".")
		path = strings.TrimPrefix(path, "(")
		path = strings.TrimSuffix(path, ")")
		typeName = parts[len(parts)-1]
		return
	}

	typeName = typePath
	return
}

// TmplString is a generation-time templated string. It is used for string
// concatenation.
//
// Given the following TmplString:
//
//    TmplString{Format: "Hello, %s", Fields: []string{"Foo()"}}
//
// The output should be the same as the output of
//
//    fmt.Sprintf("Hello, %s", v.Foo())
//
type TmplString struct {
	Format string   // printf format syntax
	Fields []string // list of struct fields
}

// IsEmpty returns true if TmplString is zero.
func (s TmplString) IsEmpty() bool {
	return s.Format == ""
}
