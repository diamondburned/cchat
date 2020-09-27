package repository

type Enumeration struct {
	Comment
	Name    string
	Values  []EnumValue
	Bitwise bool
}

type EnumValue struct {
	Comment
	Name string // also return value from String()
}

// IsPlaceholder returns true if the enumeration value is meant to be a
// placeholder. In Go, it would look like this:
//
//    const (
//        _ EnumType = iota // IsPlaceholder() == true
//        V1
//    )
//
func (v EnumValue) IsPlaceholder() bool {
	return v.Name == ""
}
