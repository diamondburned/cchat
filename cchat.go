// Code generated by ./cmd/internal. DO NOT EDIT.

// Package cchat is a set of stabilized interfaces for cchat implementations,
// joining the backend and frontend together.
//
//
// Backend
//
// Almost anything in the backend comes with an ID. For example, a Server must
// have an ID, or a Session must have a user ID. The backend is required to
// guarantee that IDs are somehow unique. This should already be the case for
// most chat services; for example, Discord provides IDs for guilds, channels,
// members, and more. The only time that the backend should not guarantee ID
// uniqueness is across Sessions, because it doesn't make sense to do so. In
// this case, the frontend should guarantee uniqueness instead, either by
// discarding duplicated items, overriding them, or anything reasonable and
// explicit.
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
//    - Storing MessagesContainer and ServersContainer are advised
//    against; however, they should be done if need be.
//    - Other containers such as LabelContainer and IconContainer
//    should also not be stored; however, the same rule as above
//    applies.
//    - For the server list, icon updates and such that happen after
//    their calls should use SetServers().
//    - For the nickname of the current server, the backend can store
//    the state of the label container. It must, however, remove the
//    container when the stop callback from JoinServer() is called.
//    - Some methods that take in a container may take in a context as
//    well.  Although implementations don't have to use this context,
//    it should try to.
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
//    if iconer := server.AsIconer(); iconer != nil {
//        println("Server implements Iconer.")
//    }
//
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
package cchat

import (
	"context"
	"fmt"
	text "github.com/diamondburned/cchat/text"
	"io"
	"time"
)

// ID is the type alias for an ID string. This type is used for clarification
// and documentation purposes only. Implementations could either use this type
// or a string type.
type ID = string

// Status represents a user's status. This might be used by the frontend to
// visually display the status.
type Status uint8

const (
	StatusUnknown Status = iota
	StatusOnline
	StatusIdle
	StatusBusy
	StatusAway
	StatusOffline
	// Invisible is reserved.
	StatusInvisible
)

func (s Status) Is(is Status) bool {
	return s == is
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

// CompletionEntry is a single completion entry returned by CompleteMessage. The
// icon URL field is optional.
type CompletionEntry struct {
	Raw       string
	Text      text.Rich
	Secondary text.Rich
	IconURL   string
	Image     bool
}

// MessageAttachment represents a single file attachment. If needed, the
// frontend will close the reader after the message is sent, that is when the
// SendMessage function returns. The backend must not use the reader after that.
type MessageAttachment struct {
	io.Reader
	Name string
}

// ReadIndication represents a read indication of a user/author in a messager
// server. It relates to a message ID within the server and is meant to imply
// that the user/author has read up to the given message ID.
//
// The frontend should override an existing author with the received ones. This
// could be treated as upsert operations.
type ReadIndication struct {
	User      User
	MessageID ID
}

// ErrInvalidConfigAtField is the structure for an error at a specific
// configuration field. Frontends can use this and highlight fields if the
// backends support it.
type ErrInvalidConfigAtField struct {
	Key string
	Err error
}

func (e ErrInvalidConfigAtField) Error() string {
	return fmt.Sprintf("Error at %s: %s", e.Key, e.Err.Error())
}

func (e ErrInvalidConfigAtField) Unwrap() error {
	return e.Err
}

// Actioner adds custom message actions into each message. Similarly to
// ServerMessageEditor, some of these methods may do IO.
type Actioner interface {
	// Do executes a message action on the given messageID, which would be taken
	// from MessageHeader.ID(). This method is allowed to do IO; the frontend should
	// take care of running it asynchronously.
	Do(ctx context.Context, action string, id ID) error // Blocking
	// MessageActions returns a list of possible actions to a message in pretty
	// strings that the frontend will use to directly display. This method must not
	// do IO.
	//
	// The string slice returned can be nil or empty.
	Actions(id ID) []string
}

// Attacher adds attachments into the message being sent.
type Attacher interface {
	Attachments() []MessageAttachment
}

// AuthenticateError is the error returned when authenticating. This error
// interface extends the normal error to allow backends to implement multi-stage
// authentication if needed in a clean way without needing any loops.
//
// This interface satisfies the error interface.
type AuthenticateError interface {
	// NextStage optionally returns a slice of Authenticator interfaces if the
	// authentication process requires another stage. It works similarly to
	// Service's Authenticate method, both of which returns a slice of
	// Authenticators.
	//
	// If the error returned is an actual error, and that the user should retry any
	// of the authentication fields, then NextStage could return nil to signify the
	// error. The frontend could reliably check nil on this field to determine
	// whether or not it should recreate the authentication fields.
	NextStage() []Authenticator
	// Error returns the error as a string. This method makes AuthenticateError
	// satisfy the built-in error interface.
	Error() string
}

// The authenticator interface allows for a multistage initial authentication
// API that the backend could use. Multistage is done by calling Authenticate
// and check for AuthenticateError's NextStage method.
type Authenticator interface {
	// Authenticate will be called with a list of values with indices correspond to
	// the returned slice of AuthenticateEntry.
	Authenticate(context.Context, []string) (Session, AuthenticateError) // Blocking
	// AuthenticateForm should return a list of authentication entries for the
	// frontend to render.
	AuthenticateForm() []AuthenticateEntry
	// Description returns the description of this authenticator method.
	Description() text.Rich
	// Name returns a short and concise name of this Authenticator method. The name
	// should not include the name of the Service.
	Name() text.Rich
}

// Backlogger adds message history capabilities into a message container. The
// backend should send old messages using the MessageCreate method of the
// MessagesContainer, and the frontend should automatically sort messages based
// on the timestamp.
//
// As there is no stop callback, if the backend needs to fetch messages
// asynchronously, it is expected to use the context to know when to cancel.
//
// The frontend should usually call this method when the user scrolls to the
// top. It is expected to guarantee not to call Backlogger more than once on the
// same ID. This can usually be done by deactivating the UI.
//
// Note that the optional usage of contexts also apply here. The frontend should
// deactivate the UI when the backend is working. However, the frontend can
// accomodate this by not deactivating until another event is triggered, then
// freeze the UI until the method is cancelled. This works even when the backend
// does not use the context.
type Backlogger interface {
	// Backlog fetches messages before the given message ID into the
	// MessagesContainer.
	//
	// This method is technically a ContainerMethod, but is listed as an IOMethod
	// because of the additional message ID parameter.
	Backlog(ctx context.Context, before ID, msgc MessagesContainer) error // Blocking
}

// Commander is an optional interface that a session could implement for command
// support. This is different from just intercepting the SendMessage() API, as
// this extends globally to the entire session.
//
// A very primitive use of this API would be to provide additional features that
// are not in cchat through a very basic terminal interface.
type Commander interface {
	// Run executes the given command, with the slice being already split arguments,
	// similar to os.Args. The function can return both a []byte and an error value.
	// The frontend should render the byte slice's value first, then display the
	// error.
	//
	// This function can do IO.
	//
	// The client should make guarantees that an empty string (and thus a
	// zero-length string slice) should be ignored. The backend should be able to
	// assume that the argument slice is always length 1 or more.
	//
	//
	// Words
	//
	// This interface and everything else inside this interface must abide by shell
	// rules when splitting words. This is in contrary to the default behavior
	// elsewhere, such as in Sender's Completer, where words are split by whitespace
	// without care for quotes.
	//
	// For example, provided this string:
	//
	//    echo "This is a string"
	//
	// This is the correct output:
	//
	//    []string{"echo", "This is a string"}
	//
	// This is the incorrect output:
	//
	//    []string{"echo", "\"This", "is", "a", "string\""}
	//
	// A helper function for this kind of behavior is available in package split,
	// under the ArgsIndexed function. This implementation also provides the rough
	// specifications.
	Run(ctx context.Context, words []string) ([]byte, error) // Blocking

	// Asserters.

	AsCompleter() Completer // Optional
}

// Completer adds autocompletion into the message composer. IO is not allowed,
// and the backend should do that only in goroutines and update its state for
// future calls.
//
// Frontends could utilize the split package inside utils for splitting words
// and index. This is the de-facto standard implementation for splitting words,
// thus backends can rely on their behaviors.
type Completer interface {
	// Complete returns the list of possible completion entries for the given word
	// list and the current word index. It takes in a list of whitespace-split slice
	// of string as well as the position of the cursor relative to the given string
	// slice.
	Complete(words []string, current int64) []CompletionEntry
}

// Configurator is an interface which the backend can implement for a primitive
// configuration API.
type Configurator interface {
	SetConfiguration(map[string]string) error
	Configuration() map[string]string
}

// Editor adds message editing to the messenger. Only EditMessage can do IO.
type Editor interface {
	// Edit edits the message with the given ID to the given content, which is the
	// edited string from RawMessageContent. This method can do IO.
	Edit(ctx context.Context, id ID, content string) error // Blocking
	// RawContent gets the original message text for editing. This method must not
	// do IO.
	RawContent(id ID) (string, error)
	// IsEditable returns whether or not a message can be edited by the client. This
	// method must not do IO.
	IsEditable(id ID) bool
}

// Identifier requires ID() to return a uniquely identifiable string for
// whatever this is embedded into. Typically, servers and messages have IDs. It
// is worth mentioning that IDs should be consistent throughout the lifespan of
// the program or maybe even forever.
type Identifier interface {
	ID() ID
}

// LabelContainer is a generic interface for any container that can hold texts.
// It's typically used for rich text labelling for usernames and server names.
//
// Methods that takes in a LabelContainer typically holds it in the state and
// may call SetLabel any time it wants. Thus, the frontend should synchronize
// calls with the main thread if needed.
//
// Labels given to the frontend may contain images or avatars, and the frontend
// has the choice to display them or not.
type LabelContainer interface {
	SetLabel(context.Context, text.Rich)
}

// ListMember represents a single member in the member list. Note that this
// interface should be treated as a static container: updating a member will
// involve a completely new ListMember instance with the same ID.
//
// Note that the frontend may give everyone an avatar regardless, or it may not
// show any avatars at all.
type ListMember interface {
	Identifier

	// Secondary returns the subtext of this member. This could be anything, such as
	// a user's custom status or away reason.
	Secondary() text.Rich
	// Status returns the status of the member. The backend does not have to show
	// offline members with the offline status if it doesn't want to show offline
	// menbers at all.
	Status() Status
	// Name returns the username or the nickname of the member, whichever the
	// backend should prefer.
	Name() text.Rich
}

// Lister is for servers that contain children servers. This is similar to
// guilds containing channels in Discord, or IRC servers containing channels.
//
// There isn't a similar stop callback API unlike other interfaces because all
// servers are expected to be listed. However, they could be hidden, such as
// collapsing a tree.
//
// The backend should call both the container and other icon and label
// containers, if any.
type Lister interface {
	// Servers should call SetServers() on the given ServersContainer to render all
	// servers. This function can do IO, and the frontend should run this in a
	// goroutine.
	Servers(context.Context, ServersContainer) error
	// Columnate is optionally used by servers to tell the frontend whether or not
	// its children should be put onto a new column instead of underneath it within
	// the same tree. If the method returns false, then the frontend can treat its
	// children as normal and show it as children within the same tree.
	//
	// For example, in Discord, guilds can be placed in guild folders, but guilds
	// and guild folders are put in the same column while guilds are actually
	// children of the folders. To replicate this behavior, guild folders should
	// return false, and guilds should return true. Both channels and categories can
	// return false.
	Columnate() bool
}

// MemberDynamicSection represents a dynamically loaded member list section. The
// section behaves similarly to MemberSection, except the information displayed
// will be considered incomplete until LoadMore returns false.
//
// LoadLess can be called by the client to mark chunks as stale, which the
// server can then unsubscribe from.
type MemberDynamicSection interface {
	// LoadLess is a method which the client must call after it is done displaying
	// entries that were added from calling LoadMore.
	//
	// The client can call this method exactly as many times as it has called
	// LoadMore. However, false should be returned if the client should stop, and
	// future calls without LoadMore should still return false.
	LoadLess(context.Context) bool // Blocking
	// LoadMore is a method which the client can call to ask for more members. This
	// method can do IO.
	//
	// Clients may call this method on the last section in the section slice;
	// however, calling this method on any section is allowed. Clients may not call
	// this method if the number of members in this section is equal to Total.
	LoadMore(context.Context) bool // Blocking
}

// MemberListContainer is a generic interface for any container that can display
// a member list. This is similar to Discord's right-side member list or IRC's
// users list. Below is a visual representation of a typical member list
// container:
//
//    +-MemberList-----------\
//    | +-Section------------|
//    | |                    |
//    | | Header - Total     |
//    | |                    |
//    | | +-Member-----------|
//    | | | Name             |
//    | | |   Secondary      |
//    | | \__________________|
//    | |                    |
//    | | +-Member-----------|
//    | | | Name             |
//    | | |   Secondary      |
//    | | \__________________|
//    \_\____________________/
type MemberListContainer interface {
	// RemoveMember removes a member from a section. If neither the member nor the
	// section exists, then the client should ignore it.
	RemoveMember(ctx context.Context, sectionID ID, memberID ID)
	// SetMember adds or updates (or upsert) a member into a section. This operation
	// must not change the section's member count. As such, changes should be done
	// separately in SetSection. If the section does not exist, then the client
	// should ignore this member, so, backends must call SetSections first before
	// SetMember on a new section.
	//
	// Typically, the backend should try and avoid calling this method and instead
	// update the labeler in the name. This method should only be used for adding
	// members.
	SetMember(ctx context.Context, sectionID ID, member ListMember)
	// SetSections (re)sets the list of sections to be the given slice. Members from
	// the old section list should be transferred over to the new section entry if
	// the section name's content is the same. Old sections that don't appear in the
	// new slice should be removed.
	SetSections(ctx context.Context, sections []MemberSection)
}

// MemberLister adds a member list into a message server.
type MemberLister interface {
	// ListMembers assigns the given container to the channel's member list. The
	// given context may be used to provide HTTP request cancellations, but
	// frontends must not rely solely on this, as the general context rules applies.
	//
	// Further behavioral documentations may be in Messenger's JoinServer method.
	ListMembers(context.Context, MemberListContainer) error
}

// MemberSection represents a member list section. The section name's content
// must be unique among other sections from the same list regardless of the rich
// segments.
type MemberSection interface {
	Identifier
	Namer

	// Total returns the total member count.
	Total() int

	// Asserters.

	AsMemberDynamicSection() MemberDynamicSection // Optional
}

// MessageCreate is the interface for an incoming message.
type MessageCreate interface {
	MessageHeader
	// Noncer is optional.
	Noncer

	// Mentioned returns whether or not the message mentions the current user. If a
	// backend does not implement mentioning, then false can be returned.
	Mentioned() bool
	Content() text.Rich
	Author() User
}

// MessageDelete is the interface for a message delete event.
type MessageDelete interface {
	MessageHeader
}

// MessageHeader implements the minimum interface for any message event.
type MessageHeader interface {
	Identifier

	Time() time.Time
}

// MessageUpdate is the interface for a message update (or edit) event. It is
// only responsible for updating a message's content. The author's name should
// be updated using MessageCreate's Author.
type MessageUpdate interface {
	MessageHeader

	Content() text.Rich
}

// MessagesContainer is a view implementation that displays a list of messages
// live. This implements the 3 most common message events: CreateMessage,
// UpdateMessage and DeleteMessage. The frontend must handle all 3.
//
// Since this container interface extends a single Server, the frontend is
// allowed to have multiple views. This is usually done with tabs or splits, but
// the backend should update them all nonetheless.
type MessagesContainer interface {
	DeleteMessage(context.Context, MessageDelete)
	UpdateMessage(context.Context, MessageUpdate)
	// CreateMessage inserts a message into the container. The frontend must
	// guarantee that the messages are in order based on what's returned from
	// Time().
	CreateMessage(context.Context, MessageCreate)
}

// Messenger is for servers that contain messages. This is similar to Discord or
// IRC channels.
type Messenger interface {
	// JoinServer joins a server that's capable of receiving messages. The server
	// may not necessarily support sending messages.
	//
	// Frontends must never call JoinServer on the same server more than twice
	// without calling the stop function first. This is the best of both worlds, as
	// it greatly reduces complexity on both sides in most cases, therefore the
	// backend can safely assume that there will only ever be one active JoinServer.
	// If the frontend wishes to do this, it must keep its own shared message
	// buffer.
	JoinServer(context.Context, MessagesContainer) error

	// Asserters.

	AsSender() Sender                   // Optional
	AsEditor() Editor                   // Optional
	AsActioner() Actioner               // Optional
	AsNicknamer() Nicknamer             // Optional
	AsBacklogger() Backlogger           // Optional
	AsMemberLister() MemberLister       // Optional
	AsUnreadIndicator() UnreadIndicator // Optional
	AsTypingIndicator() TypingIndicator // Optional
}

// Namer requires Name() to return the name of the object. Typically, this
// implies usernames for sessions or service names for services.
//
// Frontends can show the ID of the object when a name hasn't yet been set. The
// backend may immediately update the name afterwards, but assumptions should
// not be made.
type Namer interface {
	// Name sets the given container to contain the name of the parent context. The
	// method has no stop method; stopping is implied to be dependent on the parent
	// context. As such, it's only used for updating.
	Name(context.Context, LabelContainer) error
}

// Nicknamer adds the current user's nickname.
//
// The frontend will not traverse up the server tree, meaning the backend must
// handle nickname inheritance. This also means that servers that don't
// implement ServerMessage also don't need to implement ServerNickname. By
// default, the session name should be used.
type Nicknamer interface {
	Namer
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

// ReadContainer is an interface that a frontend container can implement to show
// the read bubbles on messages. This container typically implies the message
// container, but that is up to the frontend's implementation.
type ReadContainer interface {
	// DeleteIndications deletes a list of unused users/authors associated with
	// their read indicators. The backend can use this to free up users/authors that
	// are no longer in the server, for example when they are offline or have left
	// the server.
	DeleteIndications(ctx context.Context, authorIDs []ID)
	// AddIndications adds a map of users/authors to the respective message ID of
	// the server that implements ReadIndicator.
	AddIndications(context.Context, []ReadIndication)
}

// ReadIndicator adds a read indicator API for frontends to show. An example of
// the read indicator is in Matrix, where each message can have a small avatar
// indicating that the user in the room has read the message.
type ReadIndicator interface {
	// ReadIndicate subscribes the given container for read activities. The backend
	// must keep track of which read states to send over to not overwhelm the
	// frontend, and the frontend must either keep track of them, or it should not
	// display it at all.
	ReadIndicate(context.Context, ReadContainer) error
}

// Replier indicates that the message being sent is a reply to something.
// Frontends that support replies can assume that all messages in a Sender can
// be replied to, and the backend can choose to do nothing to the replied ID.
type Replier interface {
	ReplyingTo() ID
}

// SendableMessage is the bare minimum interface of a sendable message, that is,
// a message that can be sent with SendMessage(). This allows the frontend to
// implement its own message data implementation.
//
// An example of extending this interface is MessageNonce, which is similar to
// IRCv3's labeled response extension or Discord's nonces. The frontend could
// implement this interface and check if incoming MessageCreate events implement
// the same interface.
type SendableMessage interface {
	Content() string

	// Asserters.

	AsNoncer() Noncer     // Optional
	AsReplier() Replier   // Optional
	AsAttacher() Attacher // Optional
}

// Sender adds message sending to a messenger. Messengers that don't implement
// MessageSender will be considered read-only.
type Sender interface {
	// CanAttach returns whether or not the client is allowed to upload files.
	CanAttach() bool
	// Send is called by the frontend to send a message to this channel.
	Send(context.Context, SendableMessage) error // Blocking

	// Asserters.

	AsCompleter() Completer // Optional
}

// Server is a single server-like entity that could translate to a guild, a
// channel, a chat-room, and such. A server must implement at least ServerList
// or ServerMessage, else the frontend must treat it as a no-op.
//
// Note that the Server is allowed to implement both Lister and Messenger. This
// is useful when the messenger contains sub-servers, such as threads.
type Server interface {
	Identifier
	Namer

	// Asserters.

	AsLister() Lister             // Optional
	AsMessenger() Messenger       // Optional
	AsCommander() Commander       // Optional
	AsConfigurator() Configurator // Optional
}

// ServerUpdate represents a server update event.
type ServerUpdate interface {
	// Server embeds a complete server. Unlike MessageUpdate, which only returns
	// data on methods that are changed, ServerUpdate's methods must return the
	// complete data even if they stay the same. As such, zero-value returns are
	// treated as not updated, including the name.
	Server

	// PreviousID returns the ID of the item, either to be replaced or to be
	// inserted in front of.
	//
	// If replace is true, then the returned ID is the ID of the item to be
	// replaced, and the frontend should only try to use the ID as-is to find the
	// old server and replace.
	//
	// If replace is false, then the returned ID will be the ID of the item in front
	// of the embedded server. If the ID is empty or the frontend cannot find the
	// server from this ID, then it should assume and prepend the server to the
	// start.
	PreviousID() (serverID ID, replace bool)
}

// ServersContainer is any type of view that displays the list of servers. It
// should implement a SetServers([]Server) that the backend could use to call
// anytime the server list changes (at all).
//
// Typically, most frontends should implement this interface onto a tree node,
// as servers can be infinitely nested. Frontends should also reset the entire
// node and its children when SetServers is called again.
type ServersContainer interface {
	UpdateServer(context.Context, ServerUpdate)
	// SetServer is called by the backend service to request a reset of the server
	// list. The frontend can choose to call Servers() on each of the given servers,
	// or it can call that later. The backend should handle both cases.
	//
	// If the backend sets a nil server slice, then the frontend should take that as
	// an unavailable server list rather than an empty server list. The server list
	// should only be considered empty if it's an empty non-nil slice. An
	// unavailable list, on the other hand, can be treated as backend issues, e.g. a
	// connection issue.
	SetServers(context.Context, []Server)
}

// Service is a complete service that's capable of multiple sessions. It has to
// implement the Authenticate() method, which returns multiple implementations
// of Authenticator.
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
	// Identifier returns the unique identifier for the service. There is no
	// enforced representation, but services are recommended to follow the Reverse
	// Domain Name Notation for consistency. An example of that would be:
	//
	//    com.github.diamondburned.cchat-discord
	//    com.github.username.service
	Identifier
	// Namer returns the name of the service.
	Namer

	Authenticate() []Authenticator

	// Asserters.

	AsConfigurator() Configurator       // Optional
	AsSessionRestorer() SessionRestorer // Optional
}

// Session is returned after authentication on the service. It implements
// Name(), which should return the username most of the time. It also implements
// ID(), which might be used by frontends to check against User.ID() and other
// things.
//
// A session can implement SessionSaver, which would allow the frontend to save
// the session into its keyring at any time. Whether the keyring is completely
// secure or not is up to the frontend. For a GTK client, that would be using
// the GNOME Keyring daemon.
type Session interface {
	// Identifier should typically return the user ID.
	Identifier
	// Namer gives the name of the session, which is typically the username.
	Namer
	Lister

	// Disconnect asks the service to disconnect. It does not necessarily mean
	// removing the service.
	//
	// The frontend must cancel the active ServerMessage before disconnecting. The
	// backend can rely on this behavior.
	//
	// The frontend will reuse the stored session data from SessionSaver to
	// reconnect.
	//
	// When this function fails, the frontend may display the error upfront.
	// However, it will treat the session as actually disconnected. If needed, the
	// backend must implement reconnection by itself.
	Disconnect(context.Context) error // Blocking, Disposer

	// Asserters.

	AsCommander() Commander       // Optional
	AsSessionSaver() SessionSaver // Optional
}

// SessionRestorer extends Service and is called by the frontend to restore a
// saved session. The frontend may call this at any time, but it's usually on
// startup.
//
// To save a session, refer to SessionSaver.
type SessionRestorer interface {
	RestoreSession(context.Context, map[string]string) (Session, error) // Blocking
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

// TypingContainer is a generic interface for any container that can display
// users typing in the current chatbox. The typing indicator must adhere to the
// TypingTimeout returned from ServerMessageTypingIndicator. The backend should
// assume that to be the case and send events appropriately.
//
// For more documentation, refer to TypingIndicator.
type TypingContainer interface {
	// RemoveTyper explicitly removes the typer with the given user ID from the list
	// of typers. This function is usually not needed, as the client will take care
	// of removing them after TypingTimeout has been reached or other conditions
	// listed in ServerMessageTypingIndicator are met.
	RemoveTyper(ctx context.Context, authorID ID)
	// AddTyper appends the typer (author) into the frontend's list of typers, or it
	// pushes this typer on top of others. The frontend should assume current time
	// every time AddTyper is called.
	AddTyper(context.Context, User)
}

// TypingIndicator optionally extends ServerMessage to provide bidirectional
// typing indicating capabilities. This is similar to typing events on Discord
// and typing client tags on IRCv3.
//
// The client should remove a typer when a message is received with the same
// user ID, when RemoveTyper() is called by the backend or when the timeout
// returned from TypingTimeout() has been reached.
type TypingIndicator interface {
	// TypingSubscribe subscribes the given indicator to typing events sent by the
	// backend. The added event handlers have to be removed by the backend when the
	// stop() callback is called.
	//
	// This method does not take in a context, as it's supposed to only use event
	// handlers and not do any IO calls. Nonetheless, the client must treat it like
	// it does and call it asynchronously.
	TypingSubscribe(context.Context, TypingContainer) error
	// TypingTimeout returns the interval between typing events sent by the client
	// as well as the timeout before the client should remove the typer. Typically,
	// a constant should be returned.
	TypingTimeout() time.Duration
	// Typing is called by the client to indicate that the user is typing. This
	// function can do IO calls, and the client must take care of calling it in a
	// goroutine (or an asynchronous queue) as well as throttling it to
	// TypingTimeout.
	Typing(context.Context) error // Blocking
}

// UnreadContainer is an interface that a single server container (such as a
// button or a tree node) can implement if it's capable of indicating the read
// and mentioned status for that channel.
//
// Server containers that implement this has to represent unread and mentioned
// differently. For example, a mentioned channel could have a red outline, while
// an unread channel could appear brighter.
//
// Server containers are expected to represent this information in their parent
// nodes as well. For example, if a server is unread, then its parent servers as
// well as the session node should indicate the same status. Highlighting the
// session and service nodes are, however, implementation details, meaning that
// this decision is up to the frontend to decide.
type UnreadContainer interface {
	// SetUnread sets the container's unread state to the given boolean. The
	// frontend may choose how to represent this.
	SetUnread(ctx context.Context, unread bool, mentioned bool)
}

// UnreadIndicator adds an unread state API for frontends to use. The unread
// state describes whether a channel has been read or not by the current user.
// It is not to be confused with ReadIndicator, which indicates the unread state
// of others.
type UnreadIndicator interface {
	// UnreadIndicate subscribes the given unread indicator for unread and mention
	// events. Examples include when a new message is arrived and the backend needs
	// to indicate that it's unread.
	//
	// This function must provide a way to remove callbacks, as clients must call
	// this when the old server is destroyed, such as when Servers is called.
	UnreadIndicate(context.Context, UnreadContainer) error
	// MarkRead marks a message in the server messenger as read. Backends that
	// implement the UnreadIndicator interface must give control of marking messages
	// as read to the frontend if possible.
	//
	// This method is assumed to be a setter method that does not error out, because
	// the frontend has no use in knowing the error. As such, marking messages as
	// read is best-effort. The backend is in charge of synchronizing the read state
	// with the server and coordinating it with reasonable rate limits, if needed.
	MarkRead(ctx context.Context, messageID ID)
}

// User is the interface for an identifiable author. The interface defines that
// an author always have an ID and a name.
//
// An example of where this interface is used would be in MessageCreate's User
// method or embedded in Typer. The returned ID may or may not be used by the
// frontend, but backends must guarantee that the User's ID is in fact a user
// ID.
//
// The frontend may use the ID to squash messages with the same author together.
type User interface {
	Identifier
	Namer
}
