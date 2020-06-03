# cchat

A set of stabilized interfaces for cchat implementations, joining the backend
and frontend together.


## Backend


### Service

A service is a complete service that's capable of multiple sessions. It has to
implement the `Authenticate()` method, which returns an implementation of
Authenticator.

A service can implement `SessionRestorer`, which would indicate the frontend
that it can restore past sessions. Sessions are saved using the `SessionSaver`
interface that `Session` can implement.

A service can also implement `Configurator` if it has additional configurations.
The current API is a flat key-value map, which can be parsed by the backend
itself into more meaningful data structures. All configurations must be
optional, as frontends may not implement a configurator UI.

```go
type ServiceFull interface {
	Service

	// Optional
	SessionRestorer
	Configurator
	Icon
}
```

#### Interfaces

-   SessionRestorer (optional)
-   Configurator (optional)
-   Icon (optional)


### Authenticator

The authenticator interface allows for a multistage initial authentication API
that the backend could use. Multistage is done by calling `AuthenticateForm`
then `Authenticate` again forever until no errors are returned.

#### Reference Implementation

```go
var s *cchat.Session
var err error

for {
	// Pseudo-function to render the form and return the results of those forms
	// when the user confirms it.
	outputs := renderAuthForm(svc.AuthenticateForm())

	s, err = svc.Authenticate(outputs)
	if err != nil {
		renderError(errors.Wrap(err, "Error while authenticating"))
		continue // retry
	}

	break // success
}
```


### Session

A session is returned after authentication on the service. Session implements
`Name()`, which should return the username most of the time. It also implements
`UserID()`, which might be used by frontends to check against
`MessageAuthor.ID()` and other things.

A session can implement `SessionSaver`, which would allow the frontend to save
the session into its keyring at any time. Whether the keyring is completely
secure or not is up to the frontend. For cchat-gtk, that would be using the
Gnome Keyring daemon.

```go
type SessionFull interface {
	// Session already inherits ServerList.
	Session

	Icon
	SessionSaver
	Commander
	CommandCompleter
}
```

#### Interfaces

-   ServerList
-   Icon (optional)
-   Commander (optional)
-   SessionSaver (optional)


### Commander

The commander interface allows the backend to implement custom commands to
easily extend the API.

```go
type CommanderFull interface {
	Commander

	// Optional
	CommandCompleter
}
```

#### Interfaces

-   CommandCompleter (optional)


### Identifier

The identifier interface forces whatever interface that embeds this one to be
uniquely identifiable.


### Configurator

The configurator interface is a way for the frontend to display configuration
options that the backend has.


### Server

A server is any entity that is usually a channel or a guild.

```go
type ServerFull interface {
	Server
	[ ServerList | (
		ServerMessage

		// Optional
		ServerMessageSender
		ServerMessageSendCompleter
	) ]

	// Optional
	Icon
}
```

#### Interfaces

-   ServerList and/or ServerMessage
-   Icon (optional)


### ServerMessage

A server message is an entity that contains messages to be displayed. An example
would be channels in Discord and IRC.

```go
type ServerMessageFull interface {
	ServerMessage

	// Optional
	ServerMessageSender
	ServerMessageSendCompleter
}
```

#### Interfaces

-   ServerMessageSender (optional): adds message sending capability.
-   ServerMessageSendCompleter (optional): adds message completion capability.


### Messages

```go
type MessageFull interface {
	// These 3 types already inherit MessageHeader.
	[ MessageCreate | MessageUpdate | MessageDelete ]

	// Optional
	MessageNonce
	MessageMentionable
}
```

#### Interfaces

-   MessageHeader: the minimum for a proper message.
-   MessageCreate or MessageUpdate or MessageDelete
-   MessageNonce (optional)
-   MessageMentionable (optional)


### MessageAuthor

MessageAuthor is the interface that a message author would implement. ID would
typically return the user ID, or the username if that's the unique identifier.

```go
type MessageAuthorFull interface {
	MessageAuthor

	// Optional
	MessageAuthorAvatar
}
```

#### Interfaces

- MessageAuthorAvatar (optional)


## Frontend


### ServersContainer

A servers container is any type of view that displays the list of servers. It
should implement a `SetServers([]Server)` that the backend could use to call
anytime the server list changes (at all).

Typically, most frontend should implement this interface onto a tree node
instead of a tree view, as servers can be infinitely nested.

This interface expects the frontend to handle its own errors.


### MessagesContainer

A messages container is a view implementation that displays a list of messages
live. This implements the 3 most common message events: `CreateMessage`,
`UpdateMessage` and `DeleteMessage`. The frontend must handle all 3.

Since this container interface extends a single Server, the frontend is allowed
to have multiple views. This is usually done with tabs or splits, but the
backend should update them all nonetheless.


### SendableMessage

The frontend can make its own send message data implementation to indicate extra
capabilities that the backend may want.

An example of this is `MessageNonce`, which is similar to IRCv3's [labeled
response extension](https://ircv3.net/specs/extensions/labeled-response).
The frontend could implement this interface and check if incoming
`MessageCreate` events implement the same interface.

```go
type SendableMessageFull interface {
	SendableMessage

	// Optional
	MessageNonce
}
```

#### Interfaces (only known)

-   MessageNonce (optional)
