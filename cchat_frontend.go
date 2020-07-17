package cchat

import (
	"io"

	"github.com/diamondburned/cchat/text"
)

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
// thread-safe callbacks to render those events.
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
// image. It's typically used for icons that can update itself. Frontends should
// round these icons. For images that shouldn't be rounded, use ImageContainer.
//
// Methods may call SetIcon at any time in its main thread, so the frontend must
// do any I/O (including downloading the image) in another goroutine to avoid
// blocking the backend.
type IconContainer interface {
	SetIcon(url string)
}

// ImageContainer does nothing; it's reserved for future API usages. Typically,
// images don't have round corners while icons do.
type ImageContainer interface {
	SetImage(url string)
}

// UnreadIndicator is a generic interface for any container that can have
// different styles to indicate an unread and/or mentioned server.
//
// Servers that have this highlighted must traverse up the tree and highlight
// their parent servers too, if needed.
//
// Methods that have this interface as its arguments can do IO.
type UnreadIndicator interface {
	// Unread sets the container's unread state to the given boolean. The
	// frontend may choose how to represent this.
	SetUnread(unread, mentioned bool)
}

// TypingIndicator is a generic interface for any container that can display
// users typing in the current chatbox. The typing indicator must adhere to the
// TypingTimeout returned from ServerMessageTypingIndicator. The backend should
// assume that to be the case and send events appropriately.
//
// For more documentation, refer to ServerMessageTypingIndicator.
type TypingIndicator interface {
	// AddTyper appends the typer into the frontend's list of typers, or it
	// pushes this typer on top of others.
	AddTyper(Typer)
	// RemoveTyper explicitly removes the typer with the given user ID from the
	// list of typers. This function is usually not needed, as the client will
	// take care of removing them after TypingTimeout has been reached or other
	// conditions listed in ServerMessageTypingIndicator are met.
	RemoveTyper(id string)
}

// SendableMessage is the bare minimum interface of a sendable message, that is,
// a message that can be sent with SendMessage().
type SendableMessage interface {
	Content() string
}

// SendableMessageAttachments extends SendableMessage which adds attachments
// into the message. Backends that can use this interface should implement
// ServerMessageAttachmentSender.
type SendableMessageAttachments interface {
	Attachments() []MessageAttachment
}

// MessageAttachment represents a single file attachment.
//
// If needed, the frontend will close the reader after the message is sent, that
// is when the SendMessage function returns. The backend must not use the reader
// after that.
type MessageAttachment struct {
	io.Reader
	Name string
}
