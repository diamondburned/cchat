package repository

// MainNamespace is the name of the namespace that should be top level.
const MainNamespace = "cchat"

type Repositories map[string]Repository

type Repository struct {
	Enums       []Enumeration
	TypeAliases []TypeAlias
	Structs     []Struct
	ErrorTypes  []ErrorType
	Interfaces  []Interface
}

// Interface finds an interface. Nil is returned if none is found.
func (r Repository) Interface(name string) *Interface {
	for _, iface := range r.Interfaces {
		if iface.Name == name {
			return &iface
		}
	}
	return nil
}

type NamedType struct {
	Name string // optional
	Type string
}

// IsZero is true if t.Type is empty.
func (t NamedType) IsZero() bool {
	return t.Type == ""
}
