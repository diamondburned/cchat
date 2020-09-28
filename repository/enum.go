package repository

type Enumeration struct {
	Comment Comment
	Name    string
	Values  []EnumValue
	Bitwise bool // uint32 if true, uint8 otherwise
}

// GoType returns uint8 for a normal enum and uint32 for a bitwise enum. It
// returns an empty string if the length of values is overbound.
//
// The maximum number of values in a normal enum is math.MaxUint8 or 255. The
// maximum number of values in a bitwise enum is 32 for 32 bits in a uint32.
func (e Enumeration) GoType() string {
	if !e.Bitwise {
		if len(e.Values) > 255 {
			return ""
		}
		return "uint8"
	}

	if len(e.Values) > 32 {
		return ""
	}
	return "uint32"
}

type EnumValue struct {
	Comment Comment
	Name    string // also return value from String()
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
