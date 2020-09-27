package repository

// TmplString is a generation-time templated string. It is used for string
// concatenation.
//
// Given the following TmplString:
//
//    TmplString{Receiver: "v", Template: "Hello, {v.Foo()}"}
//
// The output of String() should be the same as the output of
//
//    "Hello, " + v.Foo()
//
type TmplString struct {
	Receiver string
	Template string
}

func (s TmplString) String() string {
	panic("TODO")
}
