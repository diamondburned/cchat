package cchat

import (
	"io"
	"time"

	"github.com/diamondburned/cchat/text"
)

// Service contains the bare minimum set of interface that a backend has to
// implement. Core can also implement Authenticator.
type Service interface {
	Server
	ServerList
}

// Configurator is what the backend can implement for an arbitrary configuration
// API.
type Configurator interface {
	Configuration() (map[string]string, error)
	SetConfiguration(map[string]string) error
}

// ErrInvalidConfigAtField is the structure for an error at a specific
// configuration field. Frontends can use this and highlight fields if the
// backends support it.
type ErrInvalidConfigAtField struct {
	Key string
	Err error
}

func (err *ErrInvalidConfigAtField) Error() string {
	return "Error at " + err.Key + ": " + err.Err.Error()
}

func (err *ErrInvalidConfigAtField) Unwrap() error {
	return err.Err
}

// Authenticator is what the backend can implement for authentication. A typical
// authentication frontend implementation would look like this:
//
//    for {
//        outputs := renderAuthForm(svc.AuthenticateForm())
//        if err := svc.Authenticate(outputs); err != nil {
//            log.Println("Error while authenticating:", err)
//            continue // retry
//        }
//        break // success
//    }
type Authenticator interface {
	// AuthenticateForm should return a list of authentication entries for
	// the frontend to render.
	AuthenticateForm() []AuthenticateEntry
	// Authenticate will be called with a list of values with indices
	// correspond to the returned slice of AuthenticateEntry.
	Authenticate([]string) error
}

// AuthenticateEntry represents a single authentication entry, usually an email
// or password prompt. Passwords or similar entries should have Secrets set to
// true, which should imply to frontends that the fields be masked.
type AuthenticateEntry struct {
	Name   string
	Secret bool
}

// Commander is an optional interface that a backend could implement for command
// support. This is different from just intercepting the SendMessage() API, as
// this extends the entire service.
type Commander interface {
	// RunCommand executes the given command, with the slice being already split
	// arguments, similar to os.Args. The function could return an output
	// stream, in which the frontend must display it live and close it on EOF.
	RunCommand([]string) (io.ReadCloser, error)
}

// CommandCompleter is an optional interface that a backend could implement for
// completion support. This also depends on whether or not the frontend supports
// it.
type CommandCompleter interface {
	// CompleteCommand is called with the line and current word, which the
	// backend should return with a list of new words.
	CompleteCommand(words []string, wordIndex int) []string
}

// Server is a single server-like entity that could translate to a guild, a
// channel, a chat-room, and such. A server must implement at least ServerList
// or ServerMessage, else the frontend must treat it as a no-op.
type Server interface {
	// Name returns the server's name or the service's name.
	Name() (string, error)
	// Implement ServerList and/or ServerMessage.
}

// ServerIcon is an extra interface that Server could implement for an icon.
type ServerIcon interface {
	IconURL() (string, error)
}

// ServerList is for servers that contain children servers. This is similar to
// guilds containing channels in Discord, or IRC servers containing channels.
//
// There isn't a similar LeaveServers() API like ServerMessage because all
// servers are expected to be listed. However, they could be hidden, such as
// collapsing a tree.
type ServerList interface {
	// Servers should call SetServers() on the given ServersContainer to render
	// all servers.
	Servers(ServersContainer) error
}

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

// ServerMessage is for servers that contain messages. This is similar to
// Discord or IRC channels.
type ServerMessage interface {
	// JoinServer should be called if Servers() returns nil, in which the
	// backend should connect to the server and start calling methods in the
	// container.
	JoinServer(MessagesContainer) error
	// LeaveServer indicates the backend to stop calling the controller over.
	// This should be called before any other JoinServer() calls are made.
	LeaveServer() error
	// SendMessage is called by the frontend to send a message to this channel.
	SendMessage(SendableMessage) error
}

// SendableMessage is the bare minimum interface of a sendable message, that is,
// a message that can be sent with SendMessage().
type SendableMessage interface {
	Content() string
}

// Worth pointing out that frontend container interfaces will not have an error
// handling API, as frontends can do that themselves.

// MessagesContainer is a frontend implementation for a message view, with
// synchronous callbacks to render those events.
type MessagesContainer interface {
	CreateMessage(MessageCreate)
	UpdateMessage(MessageUpdate)
	DeleteMessage(MessageDelete)
}

// MessageHeader implements the interface for any message event.
type MessageHeader interface {
	ID() string
	Time() time.Time
}

// MessageCreate is the interface for an incoming message.
type MessageCreate interface {
	MessageHeader
	Author() text.Rich
	Content() text.Rich
}

// MessageUpdate is the interface for a message update (or edit) event. If the
// returned text.Rich returns true for Empty(), then the element shouldn't be
// changed.
type MessageUpdate interface {
	MessageHeader
	Author() text.Rich  // optional
	Content() text.Rich // optional
}

// MessageDelete is the interface for a message delete event.
type MessageDelete interface {
	MessageHeader
}

// MessageNonce extends SendableMessage and MessageCreate to add nonce support.
// This is known in IRC as labeled responses. Clients could use this for
// various purposes, including knowing when a message has been sent
// successfully.
//
// Both the backend and frontend must implement this for the feature to work
// properly. The backend must check if SendableMessage implements MessageNonce,
// and the frontend must check if MessageCreate implements MessageNonce.
type MessageNonce interface {
	Nonce() string
}

// MessageAuthorAvatar is an optional interface that messages could implement. A
// frontend may optionally support this.
type MessageAuthorAvatar interface {
	MessageHeader
	AuthorAvatar() (url string)
}
