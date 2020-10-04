package repository

type Struct struct {
	Comment  Comment
	Name     string
	Fields   []StructField
	Stringer Stringer // used for String()
}

type StructField struct {
	Comment
	NamedType
}

type Stringer struct {
	Comment
	TmplString
}

func (s Stringer) IsEmpty() bool { return s.TmplString.IsEmpty() }

// ErrorStruct are structs that implement the "error" interface and starts with
// "Err".
type ErrorStruct struct {
	Struct
	ErrorString TmplString // used for Error()
}

// Wrapped returns true if the error struct contains a field with the error
// type.
func (t ErrorStruct) Wrapped() (fieldName string) {
	for _, ret := range t.Struct.Fields {
		if ret.Type == "error" {
			return ret.Name
		}
	}
	return ""
}
