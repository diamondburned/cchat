// Package cchat is a set of stabilized interfaces for cchat implementations,
// joining the backend and frontend together.
//
// For detailed explanations, refer to the README.
package cchat

import (
	"io"
	"time"

	"github.com/diamondburned/cchat/text"
)

// Service contains the bare minimum set of interface that a backend has to
// implement. Core can also implement Authenticator.
type Service interface {
	// Namer returns the name of the service.
	Namer
	// Authenticate begins the authentication process. It's put into a method so
	// backends can easily restart the entire process.
	Authenticate() Authenticator
}

// SessionRestorer extends Service and is called by the frontend to restore a
// saved session. The frontend may call this at any time, but it's usually on
// startup.
//
// To save a session, refer to SessionSaver which extends Session.
type SessionRestorer interface {
	RestoreSession(map[string]string) (Session, error)
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
	Authenticate([]string) (Session, error)
}

// AuthenticateEntry represents a single authentication entry, usually an email
// or password prompt. Passwords or similar entries should have Secrets set to
// true, which should imply to frontends that the fields be masked.
type AuthenticateEntry struct {
	Name      string
	Secret    bool
	Multiline bool
}

// Identifier requires ID() to return a uniquely identifiable string for
// whatever this is embedded into. Typically, servers and messages have IDs.
type Identifier interface {
	ID() string
}

// Namer requires Name() to return the name of the object. Typically, this
// implies usernames for sessions or service names for services.
type Namer interface {
	Name() text.Rich
}

// Service contains the bare minimum set of interface that a backend has to
// implement. Core can also implement Authenticator.
type Session interface {
	// Identifier should typically return the user ID.
	Identifier
	// Namer gives the name of the session, which is typically the username.
	Namer

	// Disconnect asks the service to disconnect. It does not necessarily mean
	// removing the service. The frontend must cancel the active ServerMessage
	// before disconnecting. The backend can rely on this behavior.
	//
	// The frontend will reuse the stored session data from SessionSaver to
	// reconnect.
	Disconnect() error

	ServerList
}

// SessionSaver extends Session and is called by the frontend to save the
// current session. This is typically called right after authentication, but a
// frontend may call this any time, including when it's closing.
//
// The frontend can ask to restore a session using SessionRestorer, which
// extends Service.
type SessionSaver interface {
	Save() (map[string]string, error)
}

// Commander is an optional interface that a session could implement for command
// support. This is different from just intercepting the SendMessage() API, as
// this extends the entire session.
type Commander interface {
	// RunCommand executes the given command, with the slice being already split
	// arguments, similar to os.Args. The function could return an output
	// stream, in which the frontend must display it live and close it on EOF.
	RunCommand([]string) (io.ReadCloser, error)
}

// CommandCompleter is an optional interface that a session could implement for
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
	Identifier
	Namer

	// Implement ServerList and/or ServerMessage.
}

// ServerNickname extends Server to add a specific user nickname into a server.
// The frontend should not traverse up the server tree, and thus the backend
// must handle nickname inheritance. This also means that servers that don't
// implement ServerMessage also don't need to implement ServerNickname. By
// default, the session name should be used.
type ServerNickname interface {
	Nickname(LabelContainer) error
}

// Icon is an extra interface that an interface could implement for an icon.
// Typically, Service would return the service logo, Session would return the
// user's avatar, and Server would return the server icon.
//
// For session, the avatar should be the same as the one returned by messages
// sent by the current user.
type Icon interface {
	Icon(IconContainer) error
}

// ServerList is for servers that contain children servers. This is similar to
// guilds containing channels in Discord, or IRC servers containing channels.
//
// There isn't a similar stop callback API unlike other interfaces because all
// servers are expected to be listed. However, they could be hidden, such as
// collapsing a tree.
//
// The backend should call both the container and other icon and label
// containers, if any.
type ServerList interface {
	// Servers should call SetServers() on the given ServersContainer to render
	// all servers. This function can do IO, and the frontend should run this in
	// a goroutine.
	Servers(ServersContainer) error
}

// ServerMessage is for servers that contain messages. This is similar to
// Discord or IRC channels.
type ServerMessage interface {
	// JoinServer should be called if Servers() returns nil, in which the
	// backend should connect to the server and start calling methods in the
	// container.
	JoinServer(MessagesContainer) (stop func(), err error)
}

// ServerMessageSender optionally extends ServerMessage to add message sending
// capability to the server.
type ServerMessageSender interface {
	// SendMessage is called by the frontend to send a message to this channel.
	SendMessage(SendableMessage) error
}

// ServerMessageEditor optionally extends ServerMessage to add message editing
// capability to the server. These functions can have IO, and the frontend
// should take care of running them in goroutines.
type ServerMessageEditor interface {
	// RawMessageContent gets the original message text for editing.
	RawMessageContent(id string) (string, error)
	// EditMessage edits the message with the given ID to the given content,
	// which is the edited string from RawMessageContent.
	EditMessage(id, content string) error
}

// ServerMessageActioner optionally extends ServerMessage to add custom message
// action capabilities to the server. Similarly to ServerMessageEditor, these
// functions can have IO.
type ServerMessageActioner interface {
	// MessageActions returns a list of possible actions in pretty strings that
	// the frontend will use to directly display.
	MessageActions() []string
	// DoMessageAction executes a message action on the given messageID, which
	// would be taken from MessageHeader.ID(). This function should not do any
	// IO. If it must, then it has to do it inside a goroutine, or
	// asynchronously.
	//
	// This method must not store MessagesContainer. It should ideally only use
	// the container interface exactly once.
	DoMessageAction(c MessagesContainer, action, messageID string) error
}

// ServerMessageSendCompleter optionally extends ServerMessageSender to add
// autocompletion into the message composer. IO is not allowed and the backend
// should do that only in goroutines and update its state for future calls.
//
// Frontends could utilize the split package inside utils for splitting words
// and index.
type ServerMessageSendCompleter interface {
	// CompleteMessage returns the list of possible completion entries for the
	// given word list and the current word index. It takes in a list of
	// whitespace-split slice of string as well as the position of the cursor
	// relative to the given string slice.
	CompleteMessage(words []string, current int) []CompletionEntry
}

// CompletionEntry is a single completion entry returned by CompleteMessage. The
// icon URL field is optional.
type CompletionEntry struct {
	// Raw is the text to be replaced in the input box.
	Raw string
	// Text is the label to be displayed.
	Text text.Rich
	// IconURL is the URL to the icon that will be displayed on the left of the
	// text. This field is optional.
	IconURL string
}

// MessageHeader implements the interface for any message event.
type MessageHeader interface {
	Identifier
	Time() time.Time
}

// MessageCreate is the interface for an incoming message.
type MessageCreate interface {
	MessageHeader
	Author() MessageAuthor
	Content() text.Rich
}

// MessageUpdate is the interface for a message update (or edit) event. If the
// returned text.Rich returns true for Empty(), then the element shouldn't be
// changed.
type MessageUpdate interface {
	MessageHeader
	Author() MessageAuthor // optional (nilable)
	Content() text.Rich    // optional (rich.Content == "")
}

// MessageAuthor is the interface for an identifiable message author. The
// returned ID may or may not be used by the frontend, but clients must
// guarantee uniqueness for intended behaviors.
//
// The frontend may also use this to squash messages with the same author
// together.
type MessageAuthor interface {
	Identifier
	Name() text.Rich
}

// MessageAuthorAvatar is an optional interface that MessageAuthor could
// implement. A frontend may optionally support this. A backend may return an
// empty string, in which the frontend must handle, perhaps by using a
// placeholder.
type MessageAuthorAvatar interface {
	Avatar() (url string)
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

// MessageMentioned extends MessageCreate to add mentioning support. The
// frontend may or may not implement this. If it does, the frontend will
// typically format the message into a notification and play a sound.
type MessageMentioned interface {
	// Mentioned returns whether or not the message mentions the current user.
	Mentioned() bool
}
