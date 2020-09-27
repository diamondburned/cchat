package repository

import "strings"

type Interface struct {
	Comment
	Name    string
	Embeds  []EmbeddedInterface
	Methods []Method // actual methods
}

type EmbeddedInterface struct {
	Comment
	InterfaceName string
}

// IsContainer returns true if the interface is a frontend container interface,
// that is when its name has "Container" at the end.
func (i Interface) IsContainer() bool {
	return strings.HasSuffix(i.Name, "Container")
}
