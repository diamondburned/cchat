package repository

type Struct struct {
	Comment Comment
	Name    string
	Fields  []StructField
}

type StructField struct {
	Comment
	NamedType
}

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
