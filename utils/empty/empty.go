// DO NOT EDIT: THIS FILE IS GENERATED!

// Package empty provides no-op asserter method implementations of interfaces in
// cchat's root and text packages.
package empty

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat/text"
)

// TextSegment provides no-op asserters for cchat.TextSegment.
type TextSegment struct{}

// AsColorer returns nil.
func (TextSegment) AsColorer() text.Colorer { return nil }

// AsLinker returns nil.
func (TextSegment) AsLinker() text.Linker { return nil }

// AsImager returns nil.
func (TextSegment) AsImager() text.Imager { return nil }

// AsAvatarer returns nil.
func (TextSegment) AsAvatarer() text.Avatarer { return nil }

// AsMentioner returns nil.
func (TextSegment) AsMentioner() text.Mentioner { return nil }

// AsAttributor returns nil.
func (TextSegment) AsAttributor() text.Attributor { return nil }

// AsCodeblocker returns nil.
func (TextSegment) AsCodeblocker() text.Codeblocker { return nil }

// AsQuoteblocker returns nil.
func (TextSegment) AsQuoteblocker() text.Quoteblocker { return nil }

// Namer provides no-op asserters for cchat.Namer.
type Namer struct{}

// AsIconer returns nil.
func (Namer) AsIconer() cchat.Iconer { return nil }

// Service provides no-op asserters for cchat.Service.
type Service struct{}

// AsConfigurator returns nil.
func (Service) AsConfigurator() cchat.Configurator { return nil }

// AsSessionRestorer returns nil.
func (Service) AsSessionRestorer() cchat.SessionRestorer { return nil }

// Session provides no-op asserters for cchat.Session.
type Session struct{}

// AsCommander returns nil.
func (Session) AsCommander() cchat.Commander { return nil }

// AsSessionSaver returns nil.
func (Session) AsSessionSaver() cchat.SessionSaver { return nil }

// Commander provides no-op asserters for cchat.Commander.
type Commander struct{}

// AsCompleter returns nil.
func (Commander) AsCompleter() cchat.Completer { return nil }

// Server provides no-op asserters for cchat.Server.
type Server struct{}

// AsLister returns nil.
func (Server) AsLister() cchat.Lister { return nil }

// AsMessenger returns nil.
func (Server) AsMessenger() cchat.Messenger { return nil }

// AsCommander returns nil.
func (Server) AsCommander() cchat.Commander { return nil }

// AsConfigurator returns nil.
func (Server) AsConfigurator() cchat.Configurator { return nil }

// Messenger provides no-op asserters for cchat.Messenger.
type Messenger struct{}

// AsSender returns nil.
func (Messenger) AsSender() cchat.Sender { return nil }

// AsEditor returns nil.
func (Messenger) AsEditor() cchat.Editor { return nil }

// AsActioner returns nil.
func (Messenger) AsActioner() cchat.Actioner { return nil }

// AsNicknamer returns nil.
func (Messenger) AsNicknamer() cchat.Nicknamer { return nil }

// AsBacklogger returns nil.
func (Messenger) AsBacklogger() cchat.Backlogger { return nil }

// AsMemberLister returns nil.
func (Messenger) AsMemberLister() cchat.MemberLister { return nil }

// AsUnreadIndicator returns nil.
func (Messenger) AsUnreadIndicator() cchat.UnreadIndicator { return nil }

// AsTypingIndicator returns nil.
func (Messenger) AsTypingIndicator() cchat.TypingIndicator { return nil }

// Sender provides no-op asserters for cchat.Sender.
type Sender struct{}

// AsCompleter returns nil.
func (Sender) AsCompleter() cchat.Completer { return nil }

// MemberSection provides no-op asserters for cchat.MemberSection.
type MemberSection struct{}

// AsMemberDynamicSection returns nil.
func (MemberSection) AsMemberDynamicSection() cchat.MemberDynamicSection { return nil }

// SendableMessage provides no-op asserters for cchat.SendableMessage.
type SendableMessage struct{}

// AsNoncer returns nil.
func (SendableMessage) AsNoncer() cchat.Noncer { return nil }

// AsAttachments returns nil.
func (SendableMessage) AsAttachments() cchat.Attachments { return nil }
