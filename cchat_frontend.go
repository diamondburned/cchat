package cchat

import "github.com/diamondburned/cchat/text"

// ServersContainer is a frontend implementation for a server view, with
// synchronous callbacks to render those events. The frontend is typically
// expected to reset the entire list, but it can do so with or without deleting
// everything and starting all over again.
type ServersContainer interface {
	// SetServer is called by the backend service to request a reset of the
	// server list. The frontend can choose to call Servers() on each of the
	// given servers, or it can call that later. The backend should handle both
	// cases.
	SetServers([]Server)
}

// MessagesContainer is a frontend implementation for a message view, with
// synchronous callbacks to render those events.
type MessagesContainer interface {
	CreateMessage(MessageCreate)
	UpdateMessage(MessageUpdate)
	DeleteMessage(MessageDelete)
}

// LabelContainer is a generic interface for any container that can hold texts.
// It's typically used for rich text labelling for usernames and server names.
//
// Methods that takes in a LabelContainer typically holds it in the state and
// may call SetLabel any time it wants. Thus, the frontend should synchronize
// calls with the main thread if needed.
type LabelContainer interface {
	SetLabel(text.Rich)
}

// IconContainer is a generic interface for any container that can hold an
// image. It's typically used for icons that can update itself.
//
// Methods may call SetIcon at any time in its main thread, so the frontend must
// do any I/O (including downloading the image) in another goroutine to avoid
// blocking the backend.
type IconContainer interface {
	SetIcon(url string)
}

// SendableMessage is the bare minimum interface of a sendable message, that is,
// a message that can be sent with SendMessage().
type SendableMessage interface {
	Content() string
}
