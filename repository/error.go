package repository

type ErrorType struct {
	Struct
	ErrorString TmplString // used for Error()
}

// Wrapped returns true if the error struct contains a field with the error
// type.
func (t ErrorType) Wrapped() bool {
	for _, ret := range t.Struct.Fields {
		if ret.Type == "error" {
			return true
		}
	}
	return false
}
