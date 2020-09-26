package cchat

import (
	"fmt"
	"io"
	"time"

	"github.com/diamondburned/cchat/text"
)

// ServersContainer is any type of view that displays the list of servers. It
// should implement a SetServers([]Server) that the backend could use to call
// anytime the server list changes (at all).
//
// Typically, most frontends should implement this interface onto a tree node,
// as servers can be infinitely nested. Frontends should also reset the entire
// node and its children when SetServers is called again.
type ServersContainer interface {
	// SetServer is called by the backend service to request a reset of the
	// server list. The frontend can choose to call Servers() on each of the
	// given servers, or it can call that later. The backend should handle both
	// cases.
	SetServers([]Server)

	// Update
	UpdateServer(ServerUpdate)
}

type ServerUpdate interface {
	// Server embeds a complete server. Unlike MessageUpdate, which only
	// returns data on methods that are changed, ServerUpdate's methods must
	// return the complete data even if they stay the same. As such, zero-value
	// returns are treated as not updated, including the name.
	Server
	// PreviousID returns the ID of the item before this server.
	PreviousID() ID
}

// MessagesContainer is a view implementation that displays a list of messages
// live. This implements the 3 most common message events: CreateMessage,
// UpdateMessage and DeleteMessage. The frontend must handle all 3.
//
// Since this container interface extends a single Server, the frontend is
// allowed to have multiple views. This is usually done with tabs or splits, but
// the backend should update them all nonetheless.
type MessagesContainer interface {
	CreateMessage(MessageCreate)
	UpdateMessage(MessageUpdate)
	DeleteMessage(MessageDelete)
}

// MessagePrepender extends MessagesContainer for backlog implementations. The
// backend is expected to call this interface's method from latest to earliest.
type MessagePrepender interface {
	// PrependMessage prepends the given MessageCreate message into the top of
	// the chat buffer.
	PrependMessage(MessageCreate)
}

// MessageHeader implements the minimum interface for any message event.
type MessageHeader interface {
	Identifier
	Time() time.Time
}

// MessageCreate is the interface for an incoming message.
type MessageCreate interface {
	MessageHeader
	Author() Author
	Content() text.Rich

	// Optional interfaces and methods that can return zero-values.

	Noncer
	// Mentioned returns whether or not the message mentions the current user.
	Mentioned() bool
}

// MessageUpdate is the interface for a message update (or edit) event. If the
// returned text.Rich returns true for Empty(), then the element shouldn't be
// changed.
type MessageUpdate interface {
	MessageHeader
	Author() Author     // optional (nilable)
	Content() text.Rich // optional (rich.Content == "")
}

// MessageDelete is the interface for a message delete event.
type MessageDelete interface {
	MessageHeader
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

// UnreadIndicator is an interface that a single server container (such as a
// button or a tree node) can implement if it's capable of indicating the read
// and mentioned status for that channel.
//
// Server containers that implement this has to implement both SetUnread and
// SetMentioned, and they should also represent those statuses differently. For
// example, a mentioned channel could have a red outline, while an unread
// channel could appear brighter.
//
// Server containers are expected to represent this information in their parent
// nodes as well. For example, if a server is unread, then its parent servers as
// well as the session node should indicate the same status. Highlighting the
// session and service nodes are, however, implementation details, meaning that
// this decision is up to the frontend to decide.
type UnreadContainer interface {
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
type TypingContainer interface {
	// AddTyper appends the typer into the frontend's list of typers, or it
	// pushes this typer on top of others.
	AddTyper(Typer)
	// RemoveTyper explicitly removes the typer with the given user ID from the
	// list of typers. This function is usually not needed, as the client will
	// take care of removing them after TypingTimeout has been reached or other
	// conditions listed in ServerMessageTypingIndicator are met.
	RemoveTyper(ID)
}

// Typer is an individual user that's typing. This interface is used
// interchangably in TypingIndicator and thus ServerMessageTypingIndicator as
// well.
type Typer interface {
	Author
	Time() time.Time
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
//
type MemberListContainer interface {
	// SetSections (re)sets the list of sections to be the given slice. Members
	// from the old section list should be transferred over to the new section
	// entry if the section name's content is the same. Old sections that don't
	// appear in the new slice should be removed.
	SetSections(sections []MemberSection)
	// SetMember adds or updates (or upsert) a member into a section. This
	// operation must not change the section's member count. As such, changes
	// should be done separately in SetSection. If the section does not exist,
	// then the client should ignore this member. As such, backends must call
	// SetSections first before SetMember on a new section.
	SetMember(sectionID ID, member ListMember)
	// RemoveMember removes a member from a section. If neither the member nor
	// the section exists, then the client should ignore it.
	RemoveMember(sectionID, memberID ID)
}

// ListMember represents a single member in the member list. This is a base
// interface that may implement more interfaces, such as Iconer for the user's
// avatar.
//
// Note that the frontend may give everyone an avatar regardless, or it may not
// show any avatars at all.
type ListMember interface {
	// Identifier identifies the individual member. This works similarly to
	// MessageAuthor.
	Identifier
	// Namer returns the name of the member. This works similarly to a
	// MessageAuthor.
	Namer
	// Status returns the status of the member. The backend does not have to
	// show offline members with the offline status if it doesn't want to show
	// offline menbers at all.
	Status() UserStatus
	// Secondary returns the subtext of this member. This could be anything,
	// such as a user's custom status or away reason.
	Secondary() text.Rich
}

// UserStatus represents a user's status. This might be used by the frontend to
// visually display the status.
type UserStatus uint8

const (
	UnknownStatus UserStatus = iota
	OnlineStatus
	IdleStatus
	BusyStatus // a.k.a. Do Not Disturb
	AwayStatus
	OfflineStatus
	InvisibleStatus // reserved; currently unused
)

// String formats a user status as a title string, such as "Online" or
// "Unknown". It treats unknown constants as UnknownStatus.
func (s UserStatus) String() string {
	switch s {
	case OnlineStatus:
		return "Online"
	case IdleStatus:
		return "Idle"
	case BusyStatus:
		return "Busy"
	case AwayStatus:
		return "Away"
	case OfflineStatus:
		return "Offline"
	case InvisibleStatus:
		return "Invisible"
	case UnknownStatus:
		return "Unknown"
	default:
		return fmt.Sprintf("UserStatus(%d)", s)
	}
}

// MemberListSection represents a member list section. The section name's
// content must be unique among other sections from the same list regardless of
// the rich segments.
type MemberSection interface {
	// Identifier identifies the current section.
	Identifier
	// Namer returns the section name.
	Namer
	// Total returns the total member count.
	Total() int

	// Optionals.

	AsMemberDynamicSection() MemberDynamicSection
}

// MemberListDynamicSection represents a dynamically loaded member list section.
// The section behaves similarly to MemberListSection, except the information
// displayed will be considered incomplete until LoadMore returns false.
//
// LoadLess can be called by the client to mark chunks as stale, which the
// server can then unsubscribe from.
type MemberDynamicSection interface {
	// LoadMore is a method which the client can call to ask for more members.
	// This method can do IO.
	//
	// Clients may call this method on the last section in the section slice;
	// however, calling this method on any section is allowed. Clients may not
	// call this method if the number of members in this section is equal to
	// Total.
	LoadMore() bool
	// LoadLess is a method which the client must call after it is done
	// displaying entries that were added from calling LoadMore.
	//
	// The client can call this method exactly as many times as it has called
	// LoadMore. However, false should be returned if the client should stop,
	// and future calls without LoadMore should still return false.
	LoadLess() bool
}

// SendableMessage is the bare minimum interface of a sendable message, that is,
// a message that can be sent with SendMessage(). This allows the frontend to
// implement its own message data implementation.
//
// An example of extending this interface is MessageNonce, which is similar to
// IRCv3's labeled response extension or Discord's nonces. The frontend could
// implement this interface and check if incoming MessageCreate events implement
// the same interface.
//
// SendableMessage can implement the following interfaces:
//
//    - MessageNonce (optional)
//    - SendableMessageAttachments (optional): refer to ServerMessageAttachmentSender
type SendableMessage interface {
	Content() string

	// Optionals.

	AsNoncer() Noncer
	AsAttachments() Attachments
}

// SendableMessageAttachments extends SendableMessage which adds attachments
// into the message. Backends that can use this interface should implement
// ServerMessageAttachmentSender.
type Attachments interface {
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
