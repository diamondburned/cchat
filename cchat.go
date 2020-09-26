// Package cchat is a set of stabilized interfaces for cchat implementations,
// joining the backend and frontend together.
//
// Backend
//
// Methods implemented by the backend that have frontend containers as arguments
// can do IO. Frontends must NOT rely on individual backend states and should
// always assume that they will block.
//
// Methods that do not return an error must NOT do any IO to prevent blocking
// the main thread. As such, ID() and Name() must never do any IO. Methods that
// do return an error may do IO, but they should be documented per method.
//
// Backend implementations have certain conditions that should be adhered to:
//
//    - Storing MessagesContainer and ServersContainer are advised against;
//    however, they should be done if need be.
//    - Other containers such as LabelContainer and IconContainer should also
//    not be stored; however, the same rule as above applies.
//    - For the server list, icon updates and such that happen after their calls
//    should use SetServers().
//    - For the nickname of the current server, the backend can store the state
//    of the label container. It must, however, remove the container when the
//    stop callback from JoinServer() is called.
//    - Some methods that take in a container may take in a context as well.
//    Although implementations don't have to use this context, it should try to.
//
// Note: IO in most cases usually refer to networking, but they should files and
// anything that is blocking, such as mutexes or semaphores.
//
// Note: As mentioned above, contexts are optional for both the frontend and
// backend. The frontend may use it for cancellation, and the backend may ignore
// it.
//
// Some interfaces can be extended. Interfaces that are extendable will have
// methods starting with "As" and returns another interface type. The
// implementation may or may not return the same struct as the interface, but
// the caller should not have to type assert it to a struct. They can also
// return nil, which should indicate the backend that the feature is not
// implemented.
//
// To avoid confusing, when said "A implements B," it is mostly assumed that A
// has a method named "AsB." It does not mean that A can be type-asserted to B.
//
// For future references, these "As" methods will be called asserter methods.
//
// Note: Backends must not do IO in the "As" methods. Most of the time, it
// should only conditionally check the local state and return value or nil.
//
// Below is an example of checking for an extended interface.
//
//    if backlogger := server.AsBacklogger(); backlogger != nil {
//        println("Server implements Backlogger.")
//    }
//
// Frontend
//
// Frontend contains all interfaces that a frontend can or must implement. The
// backend may call these methods any time from any goroutine. Thus, they should
// be thread-safe. They should also not block the call by doing so, as backends
// may call these methods in its own main thread.
//
// It is worth pointing out that frontend container interfaces will not have an
// error handling API, as frontends can do that themselves. Errors returned by
// backend methods will be errors from the backend itself and never the frontend
// errors.
//
package cchat

import (
	"context"
	"io"
	"time"

	"github.com/diamondburned/cchat/text"
)

// ID is the type alias for an ID string. This type is used for clarification
// and documentation purposes only; implementations could either use this type
// or a string type.
type ID = string

// Identifier requires ID() to return a uniquely identifiable string for
// whatever this is embedded into. Typically, servers and messages have IDs. It
// is worth mentioning that IDs should be consistent throughout the lifespan of
// the program or maybe even forever.
type Identifier interface {
	ID() ID
}

// Namer requires Name() to return the name of the object. Typically, this
// implies usernames for sessions or service names for services.
type Namer interface {
	Name() text.Rich

	AsIconer() Iconer
}

// Iconer adds icon support into Namer, which in turn is returned by other
// interfaces. Typically, Service would return the service logo, Session would
// return the user's avatar, and Server would return the server icon.
//
// For session, the avatar should be the same as the one returned by messages
// sent by the current user.
type Iconer interface {
	Icon(context.Context, IconContainer) (stop func(), err error)
}

// Noncer adds nonce support. A nonce is defined in this context as a unique
// identifier from the frontend. This interface defines the common nonce getter.
//
// Nonces are useful for frontends to know if an incoming event is a reply from
// the server backend. As such, nonces should be roundtripped through the
// server. For example, IRC would use labeled responses.
//
// The Nonce method can return an empty string. This indicates that either the
// frontend or backend (or neither) supports nonces.
//
// Contrary to other interfaces that extend with an "Is" method, the Nonce
// method could return an empty string here.
type Noncer interface {
	Nonce() string
}

// Author is the interface for an identifiable author. The interface defines
// that an author always have an ID and a name.
//
// An example of where this interface is used would be in MessageCreate's Author
// method or embedded in Typer. The returned ID may or may not be used by the
// frontend, but backends must guarantee that the Author's ID is in fact a user
// ID.
//
// The frontend may use the ID to squash messages with the same author together.
type Author interface {
	Identifier
	Namer
}

// A service is a complete service that's capable of multiple sessions. It has
// to implement the Authenticate() method, which returns an implementation of
// Authenticator.
//
// A service can implement SessionRestorer, which would indicate the frontend
// that it can restore past sessions. Sessions are saved using the SessionSaver
// interface that Session can implement.
//
// A service can also implement Configurator if it has additional
// configurations. The current API is a flat key-value map, which can be parsed
// by the backend itself into more meaningful data structures. All
// configurations must be optional, as frontends may not implement a
// configurator UI.
type Service interface {
	// Namer returns the name of the service.
	Namer
	// Authenticate begins the authentication process. It's put into a method so
	// backends can easily restart the entire process.
	Authenticate() Authenticator

	AsConfigurator() Configurator
	AsSessionRestorer() SessionRestorer
}

// The authenticator interface allows for a multistage initial authentication
// API that the backend could use. Multistage is done by calling
// AuthenticateForm then Authenticate again forever until no errors are
// returned.
//
//    var s *cchat.Session
//    var err error
//
//    for {
//        // Pseudo-function to render the form and return the results of those forms
//        // when the user confirms it.
//        outputs := renderAuthForm(svc.AuthenticateForm())
//
//        s, err = svc.Authenticate(outputs)
//        if err != nil {
//            renderError(errors.Wrap(err, "Error while authenticating"))
//            continue // retry
//        }
//
//        break // success
//    }
//
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
	Name        string
	Placeholder string
	Description string
	Secret      bool
	Multiline   bool
}

// SessionRestorer extends Service and is called by the frontend to restore a
// saved session. The frontend may call this at any time, but it's usually on
// startup.
//
// To save a session, refer to SessionSaver.
type SessionRestorer interface {
	RestoreSession(map[string]string) (Session, error)
}

// Configurator is an interface which the backend can implement for a primitive
// configuration API. Since these methods do return an error, they are allowed
// to do IO. The frontend should handle this appropriately, including running
// them asynchronously.
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

var _ error = (*ErrInvalidConfigAtField)(nil)

// Error formats the error; it satisfies the error interface.
func (err *ErrInvalidConfigAtField) Error() string {
	return "Error at " + err.Key + ": " + err.Err.Error()
}

// Unwrap returns the underlying error.
func (err *ErrInvalidConfigAtField) Unwrap() error {
	return err.Err
}

// A session is returned after authentication on the service. Session implements
// Name(), which should return the username most of the time. It also implements
// ID(), which might be used by frontends to check against MessageAuthor.ID()
// and other things.
//
// A session can implement SessionSaver, which would allow the frontend to save
// the session into its keyring at any time. Whether the keyring is completely
// secure or not is up to the frontend. For a Gtk client, that would be using
// the GNOME Keyring daemon.
type Session interface {
	// Identifier should typically return the user ID.
	Identifier
	// Namer gives the name of the session, which is typically the username.
	Namer

	// Disconnect asks the service to disconnect. It does not necessarily mean
	// removing the service.
	//
	// The frontend must cancel the active ServerMessage before disconnecting.
	// The backend can rely on this behavior.
	//
	// The frontend will reuse the stored session data from SessionSaver to
	// reconnect.
	//
	// When this function fails, the frontend may display the error upfront.
	// However, it will treat the session as actually disconnected. If needed,
	// the backend must implement reconnection by itself.
	Disconnect() error

	Lister

	AsCommander() Commander
	AsSessionSaver() SessionSaver
}

// SessionSaver extends Session and is called by the frontend to save the
// current session. This is typically called right after authentication, but a
// frontend may call this any time, including when it's closing.
//
// The frontend can ask to restore a session using SessionRestorer, which
// extends Service.
//
// The SaveSession method must not do IO; if there are any reasons that cause
// SaveSession to fail, then a nil map should be returned.
type SessionSaver interface {
	SaveSession() map[string]string
}

// Commander is an optional interface that a session could implement for command
// support. This is different from just intercepting the SendMessage() API, as
// this extends globally to the entire session.
//
// A very primitive use of this API would be to provide additional features that
// are not in cchat through a very basic terminal interface.
type Commander interface {
	// RunCommand executes the given command, with the slice being already split
	// arguments, similar to os.Args. The function could return an output
	// stream, in which the frontend must display it live and close it on EOF.
	//
	// The function must not do any IO; if it does, then they have to be in a
	// goroutine and stream their results to the ReadCloser.
	//
	// The client should make guarantees that an empty string (and thus a
	// zero-length string slice) should be ignored. The backend should be able
	// to assume that the argument slice is always length 1 or more.
	RunCommand([]string) (io.ReadCloser, error)

	AsCompleter() Completer
}

// Server is a single server-like entity that could translate to a guild, a
// channel, a chat-room, and such. A server must implement at least ServerList
// or ServerMessage, else the frontend must treat it as a no-op.
type Server interface {
	Identifier
	Namer

	// Implement either one of those only.
	AsLister() Lister
	AsMessenger() Messenger

	AsCommander() Commander
	AsConfigurator() Configurator
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
type Lister interface {
	// Servers should call SetServers() on the given ServersContainer to render
	// all servers. This function can do IO, and the frontend should run this in
	// a goroutine.
	Servers(ServersContainer) error
}

// ServerMessage is for servers that contain messages. This is similar to
// Discord or IRC channels.
type Messenger interface {
	// JoinServer joins a server that's capable of receiving messages. The
	// server may not necessarily support sending messages.
	JoinServer(context.Context, MessagesContainer) (stop func(), err error)

	AsSender() Sender
	AsEditor() Editor
	AsActioner() Actioner
	AsNicknamer() Nicknamer
	AsBacklogger() Backlogger
	AsMemberLister() MemberLister
	AsUnreadIndicator() UnreadIndicator
	AsTypingIndicator() TypingIndicator
}

// MessageSender adds message sending to a messenger. Messengers that don't
// implement MessageSender will be considered read-only.
type Sender interface {
	// Send is called by the frontend to send a message to this channel.
	Send(SendableMessage) error
	// CanAttach returns whether or not the client is allowed to upload files.
	CanAttach() bool

	AsCompleter() Completer
}

// Editor adds message editing to the messenger. Only EditMessage can do IO.
type Editor interface {
	// MessageEditable returns whether or not a message can be edited by the
	// client. This method must not do IO.
	MessageEditable(id ID) bool
	// RawMessageContent gets the original message text for editing. This method
	// must not do IO.
	RawMessageContent(id ID) (string, error)
	// EditMessage edits the message with the given ID to the given content,
	// which is the edited string from RawMessageContent. This method can do IO.
	EditMessage(id ID, content string) error
}

// Actioner adds custom message actions into each message. Similarly to
// ServerMessageEditor, some of these methods may do IO.
type Actioner interface {
	// MessageActions returns a list of possible actions in pretty strings that
	// the frontend will use to directly display. This method must not do IO.
	//
	// The string slice returned can be nil or empty.
	Actions(id string) []string
	// DoMessageAction executes a message action on the given messageID, which
	// would be taken from MessageHeader.ID(). This method is allowed to do IO;
	// the frontend should take care of running it asynchronously.
	DoAction(action, id string) error
}

// Nicknamer adds the current user's nickname.
//
// The frontend will not traverse up the server tree, meaning the backend
// must handle nickname inheritance. This also means that servers that don't
// implement ServerMessage also don't need to implement ServerNickname. By
// default, the session name should be used.
type Nicknamer interface {
	Nickname(context.Context, LabelContainer) (stop func(), err error)
}

// Backlogger adds message history capabilities into a message container. The
// frontend should typically call this method when the user scrolls to the top.
//
// As there is no stop callback, if the backend needs to fetch messages
// asynchronously, it is expected to use the context to know when to cancel.
//
// The frontend should usually call this method when the user scrolls to the
// top. It is expected to guarantee not to call MessagesBefore more than once on
// the same ID. This can usually be done by deactivating the UI.
//
// Note: Although backends might rely on this context, the frontend is still
// expected to invalidate the given container when the channel is changed.
type Backlogger interface {
	// MessagesBefore fetches messages before the given message ID into the
	// MessagesContainer.
	MessagesBefore(ctx context.Context, before ID, c MessagePrepender) error
}

// MemberLister adds a member list into a message server.
type MemberLister interface {
	// ListMembers assigns the given container to the channel's member list.
	// The given context may be used to provide HTTP request cancellations, but
	// frontends must not rely solely on this, as the general context rules
	// applies.
	//
	// Further behavioral documentations may be in Messenger's JoinServer
	// method.
	ListMembers(context.Context, MemberListContainer) (stop func(), err error)
}

// UnreadIndicator adds an unread state API for frontends to use.
type UnreadIndicator interface {
	// UnreadIndicate subscribes the given unread indicator for unread and
	// mention events. Examples include when a new message is arrived and the
	// backend needs to indicate that it's unread.
	//
	// This function must provide a way to remove callbacks, as clients must
	// call this when the old server is destroyed, such as when Servers is
	// called.
	UnreadIndicate(UnreadContainer) (stop func(), err error)
}

// ServerMessageTypingIndicator optionally extends ServerMessage to provide
// bidirectional typing indicating capabilities. This is similar to typing
// events on Discord and typing client tags on IRCv3.
//
// The client should remove a typer when a message is received with the same
// user ID, when RemoveTyper() is called by the backend or when the timeout
// returned from TypingTimeout() has been reached.
type TypingIndicator interface {
	// Typing is called by the client to indicate that the user is typing. This
	// function can do IO calls, and the client must take care of calling it in
	// a goroutine (or an asynchronous queue) as well as throttling it to
	// TypingTimeout.
	Typing() error
	// TypingTimeout returns the interval between typing events sent by the
	// client as well as the timeout before the client should remove the typer.
	// Typically, a constant should be returned.
	TypingTimeout() time.Duration
	// TypingSubscribe subscribes the given indicator to typing events sent by
	// the backend. The added event handlers have to be removed by the backend
	// when the stop() callback is called.
	//
	// This method does not take in a context, as it's supposed to only use
	// event handlers and not do any IO calls. Nonetheless, the client must
	// treat it like it does and call it asynchronously.
	TypingSubscribe(TypingContainer) (stop func(), err error)
}

// MessageCompleter adds autocompletion into the message composer. IO is not
// allowed, and the backend should do that only in goroutines and update its
// state for future calls.
//
// Frontends could utilize the split package inside utils for splitting words
// and index. This is the de-facto standard implementation for splitting words,
// thus backends can rely on their behaviors.
type Completer interface {
	// Complete returns the list of possible completion entries for the
	// given word list and the current word index. It takes in a list of
	// whitespace-split slice of string as well as the position of the cursor
	// relative to the given string slice.
	Complete(words []string, current int) []CompletionEntry
}

// CompletionEntry is a single completion entry returned by CompleteMessage. The
// icon URL field is optional.
type CompletionEntry struct {
	// Raw is the text to be replaced in the input box.
	Raw string
	// Text is the label to be displayed.
	Text text.Rich
	// Secondary is the label to be displayed on the second line, on the right
	// of Text, or not displayed at all. This should be optional. This text may
	// be dimmed out as styling.
	Secondary text.Rich
	// IconURL is the URL to the icon that will be displayed on the left of the
	// text. This field is optional.
	IconURL string
	// Image returns whether or not the icon URL is actually an image, which
	// indicates that the frontend should not do rounded corners.
	Image bool
}
