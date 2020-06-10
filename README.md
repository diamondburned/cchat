# cchat

A set of stabilized interfaces for cchat implementations, joining the backend
and frontend together.



## Backend

Methods implemented by the backend that have frontend containers as arguments
can do IO. Frontends must NOT rely on individual backend caches and should always
assume that they will block.

Methods that do not return an error must NOT do any IO to prevent blocking the
main thread. Methods that do return an error may do IO, but they should be
documented per method. `ID()` and `Name()` must never do any IO.

Backend implementations have certain conditions that should be adhered to:

-   Storing MessagesContainer and ServersContainer are advised against; however,
	they should be done if need be.
-   Other containers such as LabelContainer and IconContainer should also not be
	stored; however, the same rule as above applies.
-   For the server list, icon updates and such that happen after their calls
    should use `SetServers()`.
-   For the nickname of the current server, the backend can store the state of
    the label container. It must, however, remove the container when the stop
	callback from `JoinServer()` is called.

**Note:** IO in most cases usually refer to networking, but they should files and
anything that could block, such as mutexes or semaphores.



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

#### Interfaces

-   Namer
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

#### Interfaces

-   Identifier
-   Namer
-   ServerList
-   Icon (optional)
-   Commander (optional)
-   SessionSaver (optional)



### Commander

The commander interface allows the backend to implement custom commands to
easily extend the API.

#### Interfaces

-   CommandCompleter (optional)



### Identifier

The identifier interface forces whatever interface that embeds this one to be
uniquely identifiable.



### Namer

The namer interface forces whatever interface that embeds it to have an ideally
human-friendly name. This is typically a username or a service name.



### Configurator

The configurator interface is a way for the frontend to display configuration
options that the backend has.



### Server

A server is any entity that is usually a channel or a guild.

#### Interfaces

-   Identifier
-   Namer
-   ServerList and/or ServerMessage
-   ServerNickname (optional)
-   Icon (optional)
-   ServerMessageEditor (optional)
-   ServerMessageActioner (optional)



### ServerMessage

A server message is an entity that contains messages to be displayed. An example
would be channels in Discord and IRC.

#### Interfaces

-   ServerMessageSender (optional): adds message sending capability.
-   ServerMessageSendCompleter (optional): adds message completion capability.



### Messages

#### Interfaces

-   MessageHeader: the minimum for a proper message.
-   MessageCreate or MessageUpdate or MessageDelete
-   MessageNonce (optional)
-   MessageMentionable (optional)



### MessageAuthor

MessageAuthor is the interface that a message author would implement. ID would
typically return the user ID, or the username if that's the unique identifier.

#### Interfaces

- MessageAuthorAvatar (optional)



## Frontend

Frontend contains all interfaces that a frontend can or must implement. The
backend may call these methods any time from any goroutine. Thus, they should
be thread-safe. They should also not block the call by doing so, as backends
may call these methods in its own main thread.

It is worth pointing out that frontend container interfaces will not have an
error handling API, as frontends can do that themselves. Errors returned by
backend methods will be errors from the backend itself and never the frontend
errors.



### LabelContainer 

A label container is a generic abstraction for any container that can hold
texts. It's typically used for labels that can update itself, such as usernames.



### IconContainer

The icon container is similar to the label container. Refer to above.



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

#### Interfaces (only known)

-   MessageNonce (optional)


