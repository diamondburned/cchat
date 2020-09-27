package repository

import "encoding/gob"

func init() {
	gob.Register(ContainerMethod{})
	gob.Register(AsserterMethod{})
	gob.Register(GetterMethod{})
	gob.Register(SetterMethod{})
	gob.Register(IOMethod{})
}

type Method interface {
	UnderlyingName() string
	internalMethod()
}

type RegularMethod struct {
	Comment
	Name string
}

func (m RegularMethod) UnderlyingName() string { return m.Name }
func (m RegularMethod) internalMethod()        {}

// GetterMethod is a method that returns a regular value. These methods must not
// do any IO. An example of one would be ID() returning ID.
type GetterMethod struct {
	RegularMethod

	// Parameters is the list of parameters in the function.
	Parameters []NamedType
	// Returns is the list of named types returned from the function.
	Returns []NamedType
	// ReturnError is true if the function returns an error at the end of
	// returns.
	ReturnError bool
}

// SetterMethod is a method that sets values. These methods must not do IO, and
// they have to be non-blocking.
type SetterMethod struct {
	RegularMethod

	// Parameters is the list of parameters in the function. These parameters
	// should be the parameters to set.
	Parameters []NamedType
}

// IOMethod is a regular method that can do IO and thus is blocking. These
// methods usually always return errors.
type IOMethod struct {
	RegularMethod

	// Parameters is the list of parameters in the function.
	Parameters []NamedType
	// ReturnValue is the return value in the function.
	ReturnValue NamedType
	// ReturnError is true if the function returns an error at the end of
	// returns.
	ReturnError bool
}

// ContainerMethod is a method that uses a Container. These methods can do IO.
type ContainerMethod struct {
	RegularMethod

	// HasContext is true if the method accepts a context as its first argument.
	HasContext bool
	// ContainerType is the name of the container interface. The name will
	// almost always have "Container" as its suffix.
	ContainerType string
	// HasStopFn is true if the function returns a callback of type func() as
	// its first return. The function will return an error in addition. If this
	// is false, then only the error is returned.
	HasStopFn bool
}

// AsserterMethod is a method that allows the parent interface to extend itself
// into children interfaces. These methods must not do IO.
type AsserterMethod struct {
	// ChildType is the children type that is returned.
	ChildType string
}

func (m AsserterMethod) internalMethod() {}

// UnderlyingName returns the name of the method.
func (m AsserterMethod) UnderlyingName() string {
	return "As" + m.ChildType
}
