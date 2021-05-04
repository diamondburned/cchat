package repository

import (
	"encoding/gob"
	"strings"
)

func init() {
	gob.Register(ContainerMethod{})
	gob.Register(AsserterMethod{})
	gob.Register(GetterMethod{})
	gob.Register(SetterMethod{})
	gob.Register(ContainerUpdaterMethod{})
	gob.Register(IOMethod{})
}

type Interface struct {
	Comment Comment
	Name    string
	Embeds  []EmbeddedInterface
	Methods []Method // actual methods
}

type EmbeddedInterface struct {
	Comment       Comment
	InterfaceName string
}

// IsContainer returns true if the interface is a frontend container interface,
// that is when its name has "Container" at the end.
func (i Interface) IsContainer() bool {
	return strings.HasSuffix(i.Name, "Container")
}

type Method interface {
	UnderlyingName() string
	UnderlyingComment() Comment
	internalMethod()
}

type method struct {
	Comment Comment
	Name    string
}

func (m method) UnderlyingName() string     { return m.Name }
func (m method) UnderlyingComment() Comment { return m.Comment }
func (m method) internalMethod()            {}

// GetterMethod is a method that returns a regular value. These methods must not
// do any IO. An example of one would be ID() returning ID.
type GetterMethod struct {
	method

	// Parameters is the list of parameters in the function.
	Parameters []NamedType
	// Returns is the list of named types returned from the function.
	Returns []NamedType
	// ErrorType is non-empty if the function returns an error at the end of
	// returns. For the most part, this field should be "error" if that is the
	// case, but some methods may choose to extend the error base type.
	ErrorType string
}

// ReturnError returns true if the method can error out.
func (m GetterMethod) ReturnError() bool {
	return m.ErrorType != ""
}

// SetterMethod is a method that sets values. These methods must not do IO, and
// they have to be non-blocking.
type SetterMethod struct {
	method

	// Parameters is the list of parameters in the function. These parameters
	// should be the parameters to set.
	Parameters []NamedType
	// ErrorType is non-empty if the function returns an error at the end of
	// returns. An error may be returned from the backend if the input is
	// invalid, but it must not do IO. Frontend setters must never error.
	ErrorType string
}

// ContainerUpdaterMethod is a SetterMethod that passes to the container the
// current context to prevent race conditions when synchronizing.
// The rule of thumb is that any setter method done inside a method with a
// context is usually this type of method.
type ContainerUpdaterMethod struct {
	method

	// Parameters is the list of parameters in the function. These parameters
	// should be the parameters to set.
	Parameters []NamedType
	// ErrorType is non-empty if the function returns an error at the end of
	// returns. An error may be returned from the backend if the input is
	// invalid, but it must not do IO. Frontend setters must never error.
	ErrorType string
}

// IOMethod is a regular method that can do IO and thus is blocking. These
// methods usually always return errors. IOMethods must always have means of
// cancelling them in the API, but implementations don't have to use it; as
// such, the user should always have a timeout to gracefully wait.
type IOMethod struct {
	method

	// Parameters is the list of parameters in the function.
	Parameters []NamedType
	// ReturnValue is the return value in the function.
	ReturnValue NamedType
	// ErrorType is non-empty if the function returns an error at the end of
	// returns. For the most part, this field should be "error" if that is the
	// case, but some methods may choose to extend the error base type.
	ErrorType string
	// Disposer indicates that this method signals the disposal of the interface
	// that implements it. This is used similarly to stop functions, except all
	// disposer functions can be synchronous, and the frontend should handle
	// indicating such. The frontend can also ignore the result and run the
	// method in a dangling goroutine, but it must gracefully wait for it to be
	// done on exit.
	//
	// Similarly to the stop function, the instance that the disposer method belongs
	// to will also be considered invalid and should be freed once the function
	// returns regardless of the error.
	Disposer bool
}

// ReturnError returns true if the method can error out.
func (m IOMethod) ReturnError() bool {
	return m.ErrorType != ""
}

// ContainerMethod is a method that uses a Container. These methods can do IO
// and always return a stop callback and an error.
type ContainerMethod struct {
	method

	// HasContext is true if the method accepts a context as its first argument.
	HasContext bool
	// ContainerType is the name of the container interface. The name will
	// almost always have "Container" as its suffix.
	ContainerType string
}

// Qual returns what TypeQual returns with m.ContainerType.
func (m ContainerMethod) Qual() (path, name string) {
	return TypeQual(m.ContainerType)
}

// AsserterMethod is a method that allows the parent interface to extend itself
// into children interfaces. These methods must not do IO.
type AsserterMethod struct {
	// ChildType is the children type that is returned.
	ChildType string
}

func (m AsserterMethod) internalMethod()            {}
func (m AsserterMethod) UnderlyingComment() Comment { return Comment{} }

// UnderlyingName returns the name of the method.
func (m AsserterMethod) UnderlyingName() string {
	return "As" + m.ChildType
}

// Qual returns what TypeQual returns with m.ChildType.
func (m AsserterMethod) Qual() (path, name string) {
	return TypeQual(m.ChildType)
}
