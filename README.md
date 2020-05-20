# cchat

A set of stabilized interfaces for cchat implementations, joining the backend
and frontend together.

## Backend

### Service

A service is a server with extra services implemented.

#### Interfaces

- Server
- ServerList
- ServerIcon (optional)
- Configurator (optional)
- Authenticator (optional)

### Authenticator

The authenticator interface allows for a multistage initial authentication API
that the backend could use. Multistage is done by calling `AuthenticateForm`
then `Authenticate` again forever until no errors are returned.

```go
for {
	// Pseudo-function to render the form and return the results of those forms
	// when the user confirms it.
	outputs := renderAuthForm(svc.AuthenticateForm())

	if err := svc.Authenticate(outputs); err != nil {
		renderError(errors.Wrap(err, "Error while authenticating"))
		continue // retry
	}

	break // success
}
```

### Commander

The commander interface allows the backend to implement custom commands to
easily extend the API.

#### Interfaces

- CommandCompleter (optional)

### Server

A server is any entity that is usually a channel or a guild.

#### Interfaces

- ServerList and/or ServerMessage
- ServerIcon (optional)

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
