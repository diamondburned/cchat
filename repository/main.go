package repository

// RootPath is the root Go module path. This path is prefixed in every package
// path.
const RootPath = "github.com/diamondburned/cchat"

var Main = Packages{
	MakePath("text"): {
		Comment: Comment{`
			Package text provides a rich text API for cchat interfaces to use.

			Asserting

			Although interfaces here contain asserter methods similarly to
			cchat, the backend should take care to not implement multiple
			interfaces that may seem conflicting. For example, if Avatarer is
			already implemented, then Imager shouldn't be.
		`},
		Enums: []Enumeration{{
			Comment: Comment{`
				Attribute is the type for basic rich text markup attributes.
			`},
			Name: "Attribute",
			Values: []EnumValue{{
				Comment: Comment{"Normal is a zero-value attribute."},
				Name:    "Normal",
			}, {
				Comment: Comment{"Bold represents bold text."},
				Name:    "Bold",
			}, {
				Comment: Comment{"Italics represents italicized text."},
				Name:    "Italics",
			}, {
				Comment: Comment{"Underline represents underlined text."},
				Name:    "Underline",
			}, {
				Comment: Comment{`
					Strikethrough represents struckthrough text.
				`},
				Name: "Strikethrough",
			}, {
				Comment: Comment{`
					Spoiler represents spoiler text, which usually looks blacked
					out until hovered or clicked on.
				`},
				Name: "Spoiler",
			}, {
				Comment: Comment{`
					Monospace represents monospaced text, typically for inline
					code.
				`},
				Name: "Monospace",
			}, {
				Comment: Comment{`
					Dimmed represents dimmed text, typically slightly less
					visible than other text.
				`},
				Name: "Dimmed",
			}},
			Bitwise: true,
		}},
		Structs: []Struct{{
			Comment: Comment{`
				Rich is a normal text wrapped with optional format segments.
			`},
			Name: "Rich",
			Fields: []StructField{
				{
					NamedType: NamedType{"Content", "string"},
				},
				{
					Comment: Comment{`
						Segments are optional rich-text segment markers.
					`},
					NamedType: NamedType{"Segments", "[]Segment"},
				},
			},
			Stringer: Stringer{
				Comment: Comment{`
					String returns the Content in plain text.
				`},
				TmplString: TmplString{
					Format: "%s",
					Fields: []string{"Content"},
				},
			},
		}},
		Interfaces: []Interface{{
			Comment: Comment{`
				Segment is the minimum requirement for a format segment.
				Frontends will use this to determine when the format starts
				and ends. They will also assert this interface to any other
				formatting interface, including Linker, Colorer and
				Attributor.

				Note that a segment may implement multiple interfaces. For
				example, a Mentioner may also implement Colorer.
			`},
			Name: "Segment",
			Methods: []Method{
				GetterMethod{
					method: method{Name: "Bounds"},
					Returns: []NamedType{
						{Name: "start", Type: "int"},
						{Name: "end", Type: "int"},
					},
				},
				AsserterMethod{ChildType: "Colorer"},
				AsserterMethod{ChildType: "Linker"},
				AsserterMethod{ChildType: "Imager"},
				AsserterMethod{ChildType: "Avatarer"},
				AsserterMethod{ChildType: "Mentioner"},
				AsserterMethod{ChildType: "Attributor"},
				AsserterMethod{ChildType: "Codeblocker"},
				AsserterMethod{ChildType: "Quoteblocker"},
				AsserterMethod{ChildType: "MessageReferencer"},
			},
		}, {
			Comment: Comment{`
				MessageReferencer is similar to Linker, except it references a
				message instead of an arbitrary URL. As such, its appearance may
				be formatted similarly to a link, but this is up to the frontend
				to decide. When clicked, the frontend should scroll to the
				message with the ID returned by MessageID() and highlight it,
				though this is also for appearance, so the frontend may decide
				in detail how to display it.
			`},
			Name: "MessageReferencer",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "MessageID"},
					Returns: []NamedType{{Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Linker is a hyperlink format that a segment could implement.
				This implies that the segment should be replaced with a
				hyperlink, similarly to the anchor tag with href being the URL
				and the inner text being the text string.
			`},
			Name: "Linker",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Link"},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Imager implies the segment should be replaced with a (possibly
				inlined) image.

				The Imager segment must return a bound of length zero, that is,
				the start and end bounds must be the same, unless the Imager
				segment covers something meaningful, as images must not
				substitute texts and only complement them.

				An example of the start and end bounds being the same would be
				any inline image, and an Imager that belongs to a Mentioner
				segment should have its bounds overlap. Normally,
				implementations with separated Mentioner and Imager
				implementations don't have to bother about this, since with
				Mentioner, the same Bounds will be shared, and with Imager, the
				Bounds method can easily return the same variable for start and
				end.

				For segments that also implement mentioner, the image should be
				treated as a square avatar.
			`},
			Name: "Imager",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Image returns the URL for the image.
						`},
						Name: "Image",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							ImageSize returns the requested dimension for the
							image. This function could return (0, 0), which the
							frontend should use the image's dimensions.
						`},
						Name: "ImageSize",
					},
					Returns: []NamedType{
						{Name: "w", Type: "int"},
						{Name: "h", Type: "int"},
					},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							ImageText returns the underlying text of the image.
							Frontends could use this for hovering or
							displaying the text instead of the image.
						`},
						Name: "ImageText",
					},
					Returns: []NamedType{{Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Avatarer implies the segment should be replaced with a
				rounded-corners image. This works similarly to Imager.

				For segments that also implement mentioner, the image should be
				treated as a round avatar.
			`},
			Name: "Avatarer",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Avatar returns the URL for the image.
						`},
						Name: "Avatar",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							AvatarSize returns the requested dimension for the
							image. This function could return (0, 0), which the
							frontend should use the avatar's dimensions.
						`},
						Name: "AvatarSize",
					},
					Returns: []NamedType{{Name: "size", Type: "int"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							AvatarText returns the underlying text of the image.
							Frontends could use this for hovering or
							displaying the text instead of the image.
						`},
						Name: "AvatarText",
					},
					Returns: []NamedType{{Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Colorer is a text color format that a segment could implement.
				This is to be applied directly onto the text.

				The Color method must return a valid 32-bit RGBA color. That
				is, if the text color is solid, then the alpha value must be
				0xFF. Frontends that support 32-bit colors must render alpha
				accordingly without any edge cases.
			`},
			Name: "Colorer",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Color returns a 32-bit RGBA color.
						`},
						Name: "Color",
					},
					Returns: []NamedType{{Type: "uint32"}},
				},
			},
		}, {
			Comment: Comment{`
				Mentioner implies that the segment can be clickable, and when
				clicked it should open up a dialog containing information from
				MentionInfo().

				It is worth mentioning that frontends should assume whatever
				segment that Mentioner highlighted to be the display name of
				that user. This would allow frontends to flexibly layout the
				labels.
			`},
			Name: "Mentioner",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							MentionInfo returns the popup information of the
							mentioned segment. This is typically user
							information or something similar to that context.
						`},
						Name: "MentionInfo",
					},
					Returns: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
			},
		}, {
			Comment: Comment{`
				Attributor is a rich text markup format that a segment could
				implement. This is to be applied directly onto the text.
			`},
			Name: "Attributor",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Attribute"},
					Returns: []NamedType{{Type: "Attribute"}},
				},
			},
		}, {
			Comment: Comment{`
				Codeblocker is a codeblock that supports optional syntax
				highlighting using the language given. Note that as this is a
				block, it will appear separately from the rest of the paragraph.

				This interface is equivalent to Markdown's codeblock syntax.
			`},
			Name: "Codeblocker",
			Methods: []Method{
				GetterMethod{
					method: method{Name: "CodeblockLanguage"},
					Returns: []NamedType{{
						Name: "language",
						Type: "string",
					}},
				},
			},
		}, {
			Comment: Comment{`
				Quoteblocker represents a quoteblock that behaves similarly to
				the blockquote HTML tag. The quoteblock may be represented
				typically by an actaul quoteblock or with green arrows prepended
				to each line.
			`},
			Name: "Quoteblocker",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							QuotePrefix returns the prefix that every line the
							segment covers have. This is typically the
							greater-than sign ">" in Markdown. Frontends could
							use this information to format the quote properly.
						`},
						Name: "QuotePrefix",
					},
					Returns: []NamedType{{Name: "prefix", Type: "string"}},
				},
			},
		}},
	},
	RootPath: {
		Comment: Comment{`
			Package cchat is a set of stabilized interfaces for cchat
			implementations, joining the backend and frontend together.

			Backend

			Almost anything in the backend comes with an ID. For example, a
			Server must have an ID, or a Session must have a user ID. The
			backend is required to guarantee that IDs are somehow unique. This
			should already be the case for most chat services; for example,
			Discord provides IDs for guilds, channels, members, and more. The
			only time that the backend should not guarantee ID uniqueness is
			across Sessions, because it doesn't make sense to do so. In this
			case, the frontend should guarantee uniqueness instead, either by
			discarding duplicated items, overriding them, or anything
			reasonable and explicit.

			Methods implemented by the backend that have frontend containers as
			arguments can do IO. Frontends must NOT rely on individual backend
			states and should always assume that they will block.

			Methods that do not return an error must NOT do any IO to prevent
			blocking the main thread. As such, ID() and Name() must never do any
			IO. Methods that do return an error may do IO, but they should be
			documented per method.

			Backend implementations have certain conditions that should be
			adhered to:

			   - Storing MessagesContainer and ServersContainer are advised
			   against; however, they should be done if need be.
			   - Other containers such as LabelContainer and IconContainer
			   should also not be stored; however, the same rule as above
			   applies.
			   - For the server list, icon updates and such that happen after
			   their calls should use SetServers().
			   - For the nickname of the current server, the backend can store
			   the state of the label container. It must, however, remove the
			   container when the stop callback from JoinServer() is called.
			   - Some methods that take in a container may take in a context as
			   well.  Although implementations don't have to use this context,
			   it should try to.

			Note: IO in most cases usually refer to networking, but they should
			files and anything that is blocking, such as mutexes or semaphores.

			Note: As mentioned above, contexts are optional for both the
			frontend and backend. The frontend may use it for cancellation, and
			the backend may ignore it.

			Some interfaces can be extended. Interfaces that are extendable will
			have methods starting with "As" and returns another interface type.
			The implementation may or may not return the same struct as the
			interface, but the caller should not have to type assert it to a
			struct. They can also return nil, which should indicate the
			backend that the feature is not implemented.

			To avoid confusing, when said "A implements B," it is mostly assumed
			that A has a method named "AsB." It does not mean that A can be
			type-asserted to B.

			For future references, these "As" methods will be called asserter
			methods.

			Note: Backends must not do IO in the "As" methods. Most of the time,
			it should only conditionally check the local state and return value
			or nil.

			Below is an example of checking for an extended interface.

			   if iconer := server.AsIconer(); iconer != nil {
			       println("Server implements Iconer.")
			   }

			Frontend

			Frontend contains all interfaces that a frontend can or must
			implement. The backend may call these methods any time from any
			goroutine. Thus, they should be thread-safe. They should also not
			block the call by doing so, as backends may call these methods in
			its own main thread.

			It is worth pointing out that frontend container interfaces will not
			have an error handling API, as frontends can do that themselves.
			Errors returned by backend methods will be errors from the
			backend itself and never the frontend errors.
		`},
		Enums: []Enumeration{{
			Comment: Comment{`
				Status represents a user's status. This might be used by the
				frontend to visually display the status.
			`},
			Name: "Status",
			Values: []EnumValue{
				{Comment{""}, "Unknown"},
				{Comment{""}, "Online"},
				{Comment{""}, "Idle"},
				{Comment{""}, "Busy"},
				{Comment{""}, "Away"},
				{Comment{""}, "Offline"},
				{Comment{"Invisible is reserved."}, "Invisible"},
			},
		}},
		TypeAliases: []TypeAlias{{
			Comment: Comment{`
				ID is the type alias for an ID string. This type is used for
				clarification and documentation purposes only. Implementations
				could either use this type or a string type.
			`},
			NamedType: NamedType{"ID", "string"},
		}},
		Structs: []Struct{{
			Comment: Comment{`
				AuthenticateEntry represents a single authentication entry,
				usually an email or password prompt. Passwords or similar
				entries should have Secrets set to true, which should imply to
				frontends that the fields be masked.
			`},
			Name: "AuthenticateEntry",
			Fields: []StructField{
				{NamedType: NamedType{"Name", "string"}},
				{NamedType: NamedType{"Placeholder", "string"}},
				{NamedType: NamedType{"Description", "string"}},
				{NamedType: NamedType{"Secret", "bool"}},
				{NamedType: NamedType{"Multiline", "bool"}},
			},
		}, {
			Comment: Comment{`
				CompletionEntry is a single completion entry returned by
				CompleteMessage. The icon URL field is optional.
			`},
			Name: "CompletionEntry",
			Fields: []StructField{{
				Comment: Comment{`
					Raw is the text to be replaced in the input box.
				`},
				NamedType: NamedType{"Raw", "string"},
			}, {
				Comment: Comment{`
					Text is the label to be displayed.
				`},
				NamedType: NamedType{
					Name: "Text",
					Type: MakeQual("text", "Rich"),
				},
			}, {
				Comment: Comment{`
					Secondary is the label to be displayed on the second line,
					on the right of Text, or not displayed at all. This should
					be optional. This text may be dimmed out as styling.
				`},
				NamedType: NamedType{
					Name: "Secondary",
					Type: MakeQual("text", "Rich"),
				},
			}, {
				Comment: Comment{`
					IconURL is the URL to the icon that will be displayed on the
					left of the text. This field is optional.
				`},
				NamedType: NamedType{"IconURL", "string"},
			}, {
				Comment: Comment{`
					Image returns whether or not the icon URL is actually an
					image, which indicates that the frontend should not do
					rounded corners.
				`},
				NamedType: NamedType{"Image", "bool"},
			}},
		}, {
			Comment: Comment{`
				MessageAttachment represents a single file attachment. If
				needed, the frontend will close the reader after the message is
				sent, that is when the SendMessage function returns. The backend
				must not use the reader after that.
			`},
			Name: "MessageAttachment",
			Fields: []StructField{
				{NamedType: NamedType{"", "io.Reader"}},
				{NamedType: NamedType{"Name", "string"}},
			},
		}, {
			Comment: Comment{`
				ReadIndication represents a read indication of a user/author in
				a messager server. It relates to a message ID within the server
				and is meant to imply that the user/author has read up to the
				given message ID.

				The frontend should override an existing author with the
				received ones. This could be treated as upsert operations.
			`},
			Name: "ReadIndication",
			Fields: []StructField{
				{NamedType: NamedType{"User", "User"}},
				{NamedType: NamedType{"MessageID", "ID"}},
			},
		}},
		ErrorStructs: []ErrorStruct{{
			Struct: Struct{
				Comment: Comment{`
					ErrInvalidConfigAtField is the structure for an error at a
					specific configuration field. Frontends can use this and
					highlight fields if the backends support it.
				`},
				Name: "ErrInvalidConfigAtField",
				Fields: []StructField{
					{NamedType: NamedType{"Key", "string"}},
					{NamedType: NamedType{"Err", "error"}},
				},
			},
			ErrorString: TmplString{
				Format: "Error at %s: %s",
				Fields: []string{"Key", "Err.Error()"},
			},
		}},
		Interfaces: []Interface{{
			Comment: Comment{`
				Identifier requires ID() to return a uniquely identifiable
				string for whatever this is embedded into. Typically, servers
				and messages have IDs. It is worth mentioning that IDs should be
				consistent throughout the lifespan of the program or maybe even
				forever.
			`},
			Name: "Identifier",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "ID"},
					Returns: []NamedType{{Type: "ID"}},
				},
			},
		}, {
			Comment: Comment{`
				Namer requires Name() to return the name of the object.
				Typically, this implies usernames for sessions or service
				names for services.

				Frontends can show the ID of the object when a name hasn't yet
				been set. The backend may immediately update the name
				afterwards, but assumptions should not be made.
			`},
			Name: "Namer",
			Methods: []Method{
				ContainerMethod{
					method: method{
						Comment: Comment{`
							Name sets the given container to contain the name of
							the parent context. The method has no stop method;
							stopping is implied to be dependent on the parent
							context. As such, it's only used for updating.
						`},
						Name: "Name",
					},
					HasContext:    true,
					ContainerType: "LabelContainer",
				},
			},
		}, {
			Comment: Comment{`
				Noncer adds nonce support. A nonce is defined in this context as
				a unique identifier from the frontend. This interface defines
				the common nonce getter.

				Nonces are useful for frontends to know if an incoming event is
				a reply from the server backend. As such, nonces should be
				roundtripped through the server. For example, IRC would use
				labeled responses.

				The Nonce method can return an empty string. This indicates that
				either the frontend or backend (or neither) supports nonces.

				Contrary to other interfaces that extend with an "Is" method,
				the Nonce method could return an empty string here.
			`},
			Name: "Noncer",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Nonce"},
					Returns: []NamedType{{Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				User is the interface for an identifiable author. The
				interface defines that an author always have an ID and a name.

				An example of where this interface is used would be in
				MessageCreate's User method or embedded in Typer. The returned
				ID may or may not be used by the frontend, but backends must
				guarantee that the User's ID is in fact a user ID.

				The frontend may use the ID to squash messages with the same
				author together.
			`},
			Name: "User",
			Embeds: []EmbeddedInterface{
				{InterfaceName: "Identifier"},
				{InterfaceName: "Namer"},
			},
		}, {
			Comment: Comment{`
				A service is a complete service that's capable of multiple
				sessions. It has to implement the Authenticate() method, which
				returns multiple implementations of Authenticator.

				A service can implement SessionRestorer, which would indicate
				the frontend that it can restore past sessions. Sessions are
				saved using the SessionSaver interface that Session can
				implement.

				A service can also implement Configurator if it has additional
				configurations. The current API is a flat key-value map, which
				can be parsed by the backend itself into more meaningful data
				structures. All configurations must be optional, as frontends
				may not implement a configurator UI.
			`},
			Name: "Service",
			Embeds: []EmbeddedInterface{{
				Comment: Comment{`
					Namer returns the name of the service.
				`},
				InterfaceName: "Namer",
			}},
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Authenticate"},
					Returns: []NamedType{{Type: "[]Authenticator"}},
				},
				AsserterMethod{
					ChildType: "Configurator",
				},
				AsserterMethod{
					ChildType: "SessionRestorer",
				},
			},
		}, {
			Comment: Comment{`
				AuthenticateError is the error returned when authenticating.
				This error interface extends the normal error to allow backends
				to implement multi-stage authentication if needed in a clean way
				without needing any loops.

				This interface satisfies the error interface.
			`},
			Name: "AuthenticateError",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Error returns the error as a string. This method
							makes AuthenticateError satisfy the built-in error
							interface.
						`},
						Name: "Error",
					},
					Returns: []NamedType{{Type: "string"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							NextStage optionally returns a slice of
							Authenticator interfaces if the authentication
							process requires another stage. It works similarly
							to Service's Authenticate method, both of which
							returns a slice of Authenticators.

							If the error returned is an actual error, and that
							the user should retry any of the authentication
							fields, then NextStage could return nil to signify
							the error. The frontend could reliably check nil on
							this field to determine whether or not it should
							recreate the authentication fields.
						`},
						Name: "NextStage",
					},
					Returns: []NamedType{{Type: "[]Authenticator"}},
				},
			},
		}, {
			Comment: Comment{`
				The authenticator interface allows for a multistage initial
				authentication API that the backend could use. Multistage is
				done by calling Authenticate and check for AuthenticateError's
				NextStage method.
			`},
			Name: "Authenticator",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Name returns a short and concise name of this
							Authenticator method. The name should not include
							the name of the Service.
						`},
						Name: "Name",
					},
					Returns: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							Description returns the description of this
							authenticator method.
						`},
						Name: "Description",
					},
					Returns: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							AuthenticateForm should return a list of
							authentication entries for the frontend to render.
						`},
						Name: "AuthenticateForm",
					},
					Returns: []NamedType{{Type: "[]AuthenticateEntry"}},
				},
				IOMethod{
					method: method{
						Comment: Comment{`
							Authenticate will be called with a list of values
							with indices correspond to the returned slice of
							AuthenticateEntry.
						`},
						Name: "Authenticate",
					},
					Parameters:  []NamedType{{Type: "[]string"}},
					ReturnValue: NamedType{Type: "Session"},
					ErrorType:   "AuthenticateError",
				},
			},
		}, {
			Comment: Comment{`
				SessionRestorer extends Service and is called by the frontend to
				restore a saved session. The frontend may call this at any time,
				but it's usually on startup.

				To save a session, refer to SessionSaver.
			`},
			Name: "SessionRestorer",
			Methods: []Method{
				IOMethod{
					method:      method{Name: "RestoreSession"},
					Parameters:  []NamedType{{Type: "map[string]string"}},
					ReturnValue: NamedType{Type: "Session"},
					ErrorType:   "error",
				},
			},
		}, {
			Comment: Comment{`
				Configurator is an interface which the backend can implement for a
				primitive configuration API. Since these methods do return an error,
				they are allowed to do IO. The frontend should handle this
				appropriately, including running them asynchronously.
			`},
			Name: "Configurator",
			Methods: []Method{
				IOMethod{
					method:      method{Name: "Configuration"},
					ReturnValue: NamedType{Type: "map[string]string"},
					ErrorType:   "error",
				},
				IOMethod{
					method:     method{Name: "SetConfiguration"},
					Parameters: []NamedType{{Type: "map[string]string"}},
					ErrorType:  "error",
				},
			},
		}, {
			Comment: Comment{`
				A session is returned after authentication on the service.
				Session implements Name(), which should return the username most
				of the time. It also implements ID(), which might be used by
				frontends to check against User.ID() and other things.

				A session can implement SessionSaver, which would allow the
				frontend to save the session into its keyring at any time.
				Whether the keyring is completely secure or not is up to the
				frontend. For a GTK client, that would be using the GNOME
				Keyring daemon.
			`},
			Name: "Session",
			Embeds: []EmbeddedInterface{{
				Comment: Comment{`
					Identifier should typically return the user ID.
				`},
				InterfaceName: "Identifier",
			}, {
				Comment: Comment{`
					Namer gives the name of the session, which is typically the
					username.
				`},
				InterfaceName: "Namer",
			}, {
				InterfaceName: "Lister",
			}},
			Methods: []Method{
				IOMethod{
					method: method{
						Comment: Comment{`
							Disconnect asks the service to disconnect. It does
							not necessarily mean removing the service.

							The frontend must cancel the active ServerMessage
							before disconnecting. The backend can rely on this
							behavior.

							The frontend will reuse the stored session data from
							SessionSaver to reconnect.

							When this function fails, the frontend may display
							the error upfront. However, it will treat the
							session as actually disconnected. If needed, the
							backend must implement reconnection by itself.
						`},
						Name: "Disconnect",
					},
					ErrorType: "error",
					Disposer:  true,
				},
				AsserterMethod{ChildType: "Commander"},
				AsserterMethod{ChildType: "SessionSaver"},
			},
		}, {
			Comment: Comment{`
				SessionSaver extends Session and is called by the frontend to
				save the current session. This is typically called right after
				authentication, but a frontend may call this any time, including
				when it's closing.

				The frontend can ask to restore a session using SessionRestorer,
				which extends Service.

				The SaveSession method must not do IO; if there are any reasons
				that cause SaveSession to fail, then a nil map should be
				returned.
			`},
			Name: "SessionSaver",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "SaveSession"},
					Returns: []NamedType{{Type: "map[string]string"}},
				},
			},
		}, {
			Comment: Comment{`
				Commander is an optional interface that a session could
				implement for command support. This is different from just
				intercepting the SendMessage() API, as this extends globally to
				the entire session.

				A very primitive use of this API would be to provide additional
				features that are not in cchat through a very basic terminal
				interface.
			`},
			Name: "Commander",
			Methods: []Method{
				IOMethod{
					method: method{
						Comment: Comment{`
							Run executes the given command, with the slice being
							already split arguments, similar to os.Args. The
							function can return both a []byte and an error
							value. The frontend should render the byte slice's
							value first, then display the error.

							This function can do IO.

							The client should make guarantees that an empty
							string (and thus a zero-length string slice) should
							be ignored. The backend should be able to assume
							that the argument slice is always length 1 or more.

							Words

							This interface and everything else inside this
							interface must abide by shell rules when splitting
							words. This is in contrary to the default behavior
							elsewhere, such as in Sender's Completer, where
							words are split by whitespace without care for
							quotes.

							For example, provided this string:

								echo "This is a string"

							This is the correct output:

								[]string{"echo", "This is a string"}

							This is the incorrect output:

								[]string{"echo", "\"This", "is", "a", "string\""}

							A helper function for this kind of behavior is
							available in package split, under the ArgsIndexed
							function. This implementation also provides the
							rough specifications.
						`},
						Name: "Run",
					},
					Parameters: []NamedType{
						{Name: "words", Type: "[]string"},
					},
					ReturnValue: NamedType{Type: "[]byte"},
					ErrorType:   "error",
				},
				AsserterMethod{ChildType: "Completer"},
			},
		}, {
			Comment: Comment{`
				Server is a single server-like entity that could translate to a
				guild, a channel, a chat-room, and such. A server must implement
				at least ServerList or ServerMessage, else the frontend must
				treat it as a no-op.
			`},
			Name: "Server",
			Embeds: []EmbeddedInterface{
				{InterfaceName: "Identifier"},
				{InterfaceName: "Namer"},
			},
			Methods: []Method{
				AsserterMethod{ChildType: "Lister"},
				AsserterMethod{ChildType: "Messenger"},
				AsserterMethod{ChildType: "Commander"},
				AsserterMethod{ChildType: "Columnator"},
				AsserterMethod{ChildType: "Configurator"},
			},
		}, {
			Comment: Comment{`
				Columnator is optionally used by servers to give different
				nested servers its own nesting values. Top-level servers must
				start at 1. The zero-value (0) indicates that the server that
				implements this interface is inherently the children of its
				parent server. This is also the behavior for servers that don't
				implement this interface.

				For example, in Discord, guilds can be placed in guild folders,
				but guilds and guild folders are put in the same column while
				guilds are actually children of the folders. To replicate this
				behavior, both guild and guild folders can implement
				ServerColumnator to both return 1.
			`},
			Name: "Columnator",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Column"},
					Returns: []NamedType{{Type: "int"}},
				},
			},
		}, {
			Comment: Comment{`
				Lister is for servers that contain children servers. This is
				similar to guilds containing channels in Discord, or IRC servers
				containing channels.

				There isn't a similar stop callback API unlike other interfaces
				because all servers are expected to be listed. However, they
				could be hidden, such as collapsing a tree.

				The backend should call both the container and other icon and
				label containers, if any.
			`},
			Name: "Lister",
			Methods: []Method{
				ContainerMethod{
					method: method{
						Comment: Comment{`
							Servers should call SetServers() on the given
							ServersContainer to render all servers. This
							function can do IO, and the frontend should run this
							in a goroutine.
						`},
						Name: "Servers",
					},
					ContainerType: "ServersContainer",
				},
			},
		}, {
			Comment: Comment{`
				Messenger is for servers that contain messages. This is similar
				to Discord or IRC channels.
			`},
			Name: "Messenger",
			Methods: []Method{
				ContainerMethod{
					method: method{
						Comment: Comment{`
							JoinServer joins a server that's capable of
							receiving messages. The server may not necessarily
							support sending messages.

							Frontends must never call JoinServer on the same
							server more than twice without calling the stop
							function first. This is the best of both worlds, as
							it greatly reduces complexity on both sides in most
							cases, therefore the backend can safely assume that
							there will only ever be one active JoinServer. If
							the frontend wishes to do this, it must keep its own
							shared message buffer.
						`},
						Name: "JoinServer",
					},
					HasContext:    true,
					ContainerType: "MessagesContainer",
					HasStopFn:     true,
				},
				AsserterMethod{ChildType: "Sender"},
				AsserterMethod{ChildType: "Editor"},
				AsserterMethod{ChildType: "Actioner"},
				AsserterMethod{ChildType: "Nicknamer"},
				AsserterMethod{ChildType: "Backlogger"},
				AsserterMethod{ChildType: "MemberLister"},
				AsserterMethod{ChildType: "UnreadIndicator"},
				AsserterMethod{ChildType: "TypingIndicator"},
			},
		}, {
			Comment: Comment{`
				Sender adds message sending to a messenger. Messengers that
				don't implement MessageSender will be considered read-only.
			`},
			Name: "Sender",
			Methods: []Method{
				IOMethod{
					method: method{
						Comment: Comment{`
							Send is called by the frontend to send a message to
							this channel.
						`},
						Name: "Send",
					},
					Parameters: []NamedType{
						{Type: "SendableMessage"},
					},
					ErrorType: "error",
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							CanAttach returns whether or not the client is
							allowed to upload files.
						`},
						Name: "CanAttach",
					},
					Returns: []NamedType{{Type: "bool"}},
				},
				AsserterMethod{ChildType: "Completer"},
			},
		}, {
			Comment: Comment{`
				Editor adds message editing to the messenger. Only EditMessage
				can do IO.
			`},
			Name: "Editor",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							IsEditable returns whether or not a message can be
							edited by the client. This method must not do IO.
						`},
						Name: "IsEditable",
					},
					Parameters: []NamedType{{Name: "id", Type: "ID"}},
					Returns:    []NamedType{{Type: "bool"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							RawContent gets the original message text for
							editing. This method must not do IO.
						`},
						Name: "RawContent",
					},
					Parameters: []NamedType{{Name: "id", Type: "ID"}},
					Returns:    []NamedType{{Type: "string"}},
					ErrorType:  "error",
				},
				IOMethod{
					method: method{
						Comment: Comment{`
							Edit edits the message with the given ID to the
							given content, which is the edited string from
							RawMessageContent. This method can do IO.
						`},
						Name: "Edit",
					},
					Parameters: []NamedType{
						{Name: "id", Type: "ID"},
						{Name: "content", Type: "string"},
					},
					ErrorType: "error",
				},
			},
		}, {
			Comment: Comment{`
				Actioner adds custom message actions into each message.
				Similarly to ServerMessageEditor, some of these methods may
				do IO.
			`},
			Name: "Actioner",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							MessageActions returns a list of possible actions to
							a message in pretty strings that the frontend will
							use to directly display. This method must not do IO.

							The string slice returned can be nil or empty.
						`},
						Name: "Actions",
					},
					Parameters: []NamedType{{Name: "id", Type: "ID"}},
					Returns:    []NamedType{{Type: "[]string"}},
				},
				IOMethod{
					method: method{
						Comment: Comment{`
							Do executes a message action on the given messageID,
							which would be taken from MessageHeader.ID(). This
							method is allowed to do IO; the frontend should take
							care of running it asynchronously.
						`},
						Name: "Do",
					},
					Parameters: []NamedType{
						{Name: "action", Type: "string"},
						{Name: "id", Type: "ID"},
					},
					ErrorType: "error",
				},
			},
		}, {
			Comment: Comment{`
				Nicknamer adds the current user's nickname.

				The frontend will not traverse up the server tree, meaning the
				backend must handle nickname inheritance. This also means that
				servers that don't implement ServerMessage also don't need to
				implement ServerNickname. By default, the session name should be
				used.
			`},
			Name: "Nicknamer",
			Methods: []Method{
				ContainerMethod{
					method:        method{Name: "Nickname"},
					HasContext:    true,
					ContainerType: "LabelContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				Backlogger adds message history capabilities into a message
				container. The backend should send old messages using the
				MessageCreate method of the MessagesContainer, and the frontend
				should automatically sort messages based on the timestamp.

				As there is no stop callback, if the backend needs to fetch
				messages asynchronously, it is expected to use the context to
				know when to cancel.

				The frontend should usually call this method when the user
				scrolls to the top. It is expected to guarantee not to call
				Backlogger more than once on the same ID. This can usually be
				done by deactivating the UI.

				Note that the optional usage of contexts also apply here. The
				frontend should deactivate the UI when the backend is working.
				However, the frontend can accomodate this by not deactivating
				until another event is triggered, then freeze the UI until the
				method is cancelled. This works even when the backend does not
				use the context.
			`},
			Name: "Backlogger",
			Methods: []Method{
				IOMethod{ // technically a ContainerMethod.
					method: method{
						Comment: Comment{`
							Backlog fetches messages before the given message ID
							into the MessagesContainer.

							This method is technically a ContainerMethod, but is
							listed as an IOMethod because of the additional
							message ID parameter.
						`},
						Name: "Backlog",
					},
					Parameters: []NamedType{
						{"ctx", "context.Context"},
						{"before", "ID"},
						{"msgc", "MessagesContainer"},
					},
					ErrorType: "error",
				},
			},
		}, {
			Comment: Comment{`
				MemberLister adds a member list into a message server.
			`},
			Name: "MemberLister",
			Methods: []Method{
				ContainerMethod{
					method: method{
						Comment: Comment{`
							ListMembers assigns the given container to the
							channel's member list.  The given context may be
							used to provide HTTP request cancellations, but
							frontends must not rely solely on this, as the
							general context rules applies.

							Further behavioral documentations may be in
							Messenger's JoinServer method.
						`},
						Name: "ListMembers",
					},
					HasContext:    true,
					ContainerType: "MemberListContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				ReadIndicator adds a read indicator API for frontends to show.
				An example of the read indicator is in Matrix, where each
				message can have a small avatar indicating that the user in the
				room has read the message.
			`},
			Name: "ReadIndicator",
			Methods: []Method{
				ContainerMethod{
					method: method{
						Comment: Comment{`
							ReadIndicate subscribes the given container for read
							activities. The backend must keep track of which
							read states to send over to not overwhelm the
							frontend, and the frontend must either keep track of
							them, or it should not display it at all.
						`},
						Name: "ReadIndicate",
					},
					ContainerType: "ReadContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				UnreadIndicator adds an unread state API for frontends to use.
				The unread state describes whether a channel has been read or
				not by the current user. It is not to be confused with
				ReadIndicator, which indicates the unread state of others.
			`},
			Name: "UnreadIndicator",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							MarkRead marks a message in the server messenger as
							read. Backends that implement the UnreadIndicator
							interface must give control of marking messages as
							read to the frontend if possible.

							This method is assumed to be a setter method that
							does not error out, because the frontend has no use
							in knowing the error. As such, marking messages as
							read is best-effort. The backend is in charge of
							synchronizing the read state with the server and
							coordinating it with reasonable rate limits, if
							needed.
						`},
						Name: "MarkRead",
					},
					Parameters: []NamedType{{"messageID", "ID"}},
				},
				ContainerMethod{
					method: method{
						Comment: Comment{`
							UnreadIndicate subscribes the given unread indicator
							for unread and mention events. Examples include when
							a new message is arrived and the backend needs to
							indicate that it's unread.

							This function must provide a way to remove
							callbacks, as clients must call this when the old
							server is destroyed, such as when Servers is called.
						`},
						Name: "UnreadIndicate",
					},
					ContainerType: "UnreadContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				TypingIndicator optionally extends ServerMessage to provide
				bidirectional typing indicating capabilities. This is similar to
				typing events on Discord and typing client tags on IRCv3.

				The client should remove a typer when a message is received with
				the same user ID, when RemoveTyper() is called by the backend or
				when the timeout returned from TypingTimeout() has been reached.
			`},
			Name: "TypingIndicator",
			Methods: []Method{
				IOMethod{
					method: method{
						Comment: Comment{`
							Typing is called by the client to indicate that the
							user is typing. This function can do IO calls, and
							the client must take care of calling it in a
							goroutine (or an asynchronous queue) as well as
							throttling it to TypingTimeout.
						`},
						Name: "Typing",
					},
					ErrorType: "error",
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							TypingTimeout returns the interval between typing
							events sent by the client as well as the timeout
							before the client should remove the typer.
							Typically, a constant should be returned.
						`},
						Name: "TypingTimeout",
					},
					Returns: []NamedType{{Type: "time.Duration"}},
				},
				ContainerMethod{
					method: method{
						Comment: Comment{`
							TypingSubscribe subscribes the given indicator to
							typing events sent by the backend. The added event
							handlers have to be removed by the backend when the
							stop() callback is called.

							This method does not take in a context, as it's
							supposed to only use event handlers and not do any
							IO calls.  Nonetheless, the client must treat it
							like it does and call it asynchronously.
						`},
						Name: "TypingSubscribe",
					},
					ContainerType: "TypingContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				Completer adds autocompletion into the message composer. IO is
				not allowed, and the backend should do that only in goroutines
				and update its state for future calls.

				Frontends could utilize the split package inside utils for
				splitting words and index. This is the de-facto standard
				implementation for splitting words, thus backends can rely on
				their behaviors.
			`},
			Name: "Completer",
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Complete returns the list of possible completion
							entries for the given word list and the current word
							index. It takes in a list of whitespace-split slice
							of string as well as the position of the cursor
							relative to the given string slice.
						`},
						Name: "Complete",
					},
					Parameters: []NamedType{
						{Name: "words", Type: "[]string"},
						{Name: "current", Type: "int64"},
					},
					Returns: []NamedType{
						{Type: "[]CompletionEntry"},
					},
				},
			},
		}, {
			Comment: Comment{`
				ServersContainer is any type of view that displays the list of
				servers. It should implement a SetServers([]Server) that the
				backend could use to call anytime the server list changes (at
				all).

				Typically, most frontends should implement this interface onto a
				tree node, as servers can be infinitely nested. Frontends should
				also reset the entire node and its children when SetServers is
				called again.
			`},
			Name: "ServersContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							SetServer is called by the backend service to
							request a reset of the server list. The frontend can
							choose to call Servers() on each of the given
							servers, or it can call that later. The backend
							should handle both cases.
						`},
						Name: "SetServers",
					},
					Parameters: []NamedType{{Type: "[]Server"}},
				},
				SetterMethod{
					method:     method{Name: "UpdateServer"},
					Parameters: []NamedType{{Type: "ServerUpdate"}},
				},
			},
		}, {
			Comment: Comment{`
				ServerUpdate represents a server update event.
			`},
			Name: "ServerUpdate",
			Embeds: []EmbeddedInterface{{
				Comment: Comment{`
					Server embeds a complete server. Unlike MessageUpdate, which
					only returns data on methods that are changed,
					ServerUpdate's methods must return the complete data even if
					they stay the same. As such, zero-value returns are treated
					as not updated, including the name.
				`},
				InterfaceName: "Server",
			}},
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							PreviousID returns the ID of the item, either to be
							replaced or to be inserted in front of.

							If replace is true, then the returned ID is the ID
							of the item to be replaced, and the frontend should
							only try to use the ID as-is to find the old server
							and replace.

							If replace is false, then the returned ID will be
							the ID of the item in front of the embedded server.
							If the ID is empty or the frontend cannot find the
							server from this ID, then it should assume and
							prepend the server to the start.
						`},
						Name: "PreviousID",
					},
					Returns: []NamedType{
						{Name: "serverID", Type: "ID"},
						{Name: "replace", Type: "bool"},
					},
				},
			},
		}, {
			Comment: Comment{`
				MessagesContainer is a view implementation that displays a list
				of messages live. This implements the 3 most common message
				events: CreateMessage, UpdateMessage and DeleteMessage. The
				frontend must handle all 3.

				Since this container interface extends a single Server, the
				frontend is allowed to have multiple views. This is usually done
				with tabs or splits, but the backend should update them all
				nonetheless.
			`},
			Name: "MessagesContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							CreateMessage inserts a message into the container.
							The frontend must guarantee that the messages are
							in order based on what's returned from Time().
						`},
						Name: "CreateMessage",
					},
					Parameters: []NamedType{{Type: "MessageCreate"}},
				},
				SetterMethod{
					method:     method{Name: "UpdateMessage"},
					Parameters: []NamedType{{Type: "MessageUpdate"}},
				},
				SetterMethod{
					method:     method{Name: "DeleteMessage"},
					Parameters: []NamedType{{Type: "MessageDelete"}},
				},
			},
		}, {
			Comment: Comment{`
				MessageHeader implements the minimum interface for any message
				event.
			`},
			Name:   "MessageHeader",
			Embeds: []EmbeddedInterface{{InterfaceName: "Identifier"}},
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Time"},
					Returns: []NamedType{{Type: "time.Time"}},
				},
			},
		}, {
			Comment: Comment{`
				MessageCreate is the interface for an incoming message.
			`},
			Name: "MessageCreate",
			Embeds: []EmbeddedInterface{
				{Comment{""}, "MessageHeader"},
				{Comment{"Noncer is optional."}, "Noncer"},
			},
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Author"},
					Returns: []NamedType{{Type: "User"}},
				},
				GetterMethod{
					method: method{Name: "Content"},
					Returns: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							Mentioned returns whether or not the message
							mentions the current user. If a backend does not
							implement mentioning, then false can be returned.
						`},
						Name: "Mentioned",
					},
					Returns: []NamedType{{Type: "bool"}},
				},
			},
		}, {
			Comment: Comment{`
				MessageUpdate is the interface for a message update (or edit)
				event. It behaves similarly to MessageCreate, except all fields
				are optional. The frontend is responsible for checking which
				field is not empty and check it.
			`},
			Name:   "MessageUpdate",
			Embeds: []EmbeddedInterface{{InterfaceName: "MessageCreate"}},
		}, {
			Comment: Comment{`
				MessageDelete is the interface for a message delete event.
			`},
			Name:   "MessageDelete",
			Embeds: []EmbeddedInterface{{InterfaceName: "MessageHeader"}},
		}, {
			Comment: Comment{`
				LabelContainer is a generic interface for any container that can
				hold texts. It's typically used for rich text labelling for
				usernames and server names.

				Methods that takes in a LabelContainer typically holds it in the
				state and may call SetLabel any time it wants. Thus, the
				frontend should synchronize calls with the main thread if
				needed.

				Labels given to the frontend may contain images or avatars, and
				the frontend has the choice to display them or not.
			`},
			Name: "LabelContainer",
			Methods: []Method{
				SetterMethod{
					method: method{Name: "SetLabel"},
					Parameters: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
			},
		}, {
			Comment: Comment{`
				ReadContainer is an interface that a frontend container can
				implement to show the read bubbles on messages. This container
				typically implies the message container, but that is up to the
				frontend's implementation.
			`},
			Name: "ReadContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							AddIndications adds a map of users/authors to the
							respective message ID of the server that implements
							ReadIndicator.
						`},
						Name: "AddIndications",
					},
					Parameters: []NamedType{{"", "[]ReadIndication"}},
				},
				SetterMethod{
					method: method{
						Comment: Comment{`
							DeleteIndications deletes a list of unused
							users/authors associated with their read indicators.
							The backend can use this to free up users/authors
							that are no longer in the server, for example when
							they are offline or have left the server.
						`},
						Name: "DeleteIndications",
					},
					Parameters: []NamedType{{"authorIDs", "[]ID"}},
				},
			},
		}, {
			Comment: Comment{`
				UnreadContainer is an interface that a single server container
				(such as a button or a tree node) can implement if it's capable
				of indicating the read and mentioned status for that channel.

				Server containers that implement this has to represent unread
				and mentioned differently. For example, a mentioned channel
				could have a red outline, while an unread channel could appear
				brighter.

				Server containers are expected to represent this information in
				their parent nodes as well. For example, if a server is unread,
				then its parent servers as well as the session node should
				indicate the same status. Highlighting the session and service
				nodes are, however, implementation details, meaning that this
				decision is up to the frontend to decide.
			`},
			Name: "UnreadContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							SetUnread sets the container's unread state to the
							given boolean. The frontend may choose how to
							represent this.
						`},
						Name: "SetUnread",
					},
					Parameters: []NamedType{
						{"unread", "bool"},
						{"mentioned", "bool"},
					},
				},
			},
		}, {
			Comment: Comment{`
				TypingContainer is a generic interface for any container that can display
				users typing in the current chatbox. The typing indicator must adhere to the
				TypingTimeout returned from ServerMessageTypingIndicator. The backend should
				assume that to be the case and send events appropriately.

				For more documentation, refer to TypingIndicator.
			`},
			Name: "TypingContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							AddTyper appends the typer (author) into the
							frontend's list of typers, or it pushes this typer
							on top of others. The frontend should assume current
							time every time AddTyper is called.
						`},
						Name: "AddTyper",
					},
					Parameters: []NamedType{{"", "User"}},
				},
				SetterMethod{
					method: method{
						Comment: Comment{`
							RemoveTyper explicitly removes the typer with the
							given user ID from the list of typers. This function
							is usually not needed, as the client will take care
							of removing them after TypingTimeout has been
							reached or other conditions listed in
							ServerMessageTypingIndicator are met.
						`},
						Name: "RemoveTyper",
					},
					Parameters: []NamedType{{Name: "authorID", Type: "ID"}},
				},
			},
		}, {
			Comment: Comment{`
				MemberListContainer is a generic interface for any container
				that can display a member list. This is similar to Discord's
				right-side member list or IRC's users list. Below is a visual
				representation of a typical member list container:

				   +-MemberList-----------\
				   | +-Section------------|
				   | |                    |
				   | | Header - Total     |
				   | |                    |
				   | | +-Member-----------|
				   | | | Name             |
				   | | |   Secondary      |
				   | | \__________________|
				   | |                    |
				   | | +-Member-----------|
				   | | | Name             |
				   | | |   Secondary      |
				   | | \__________________|
				   \_\____________________/
			`},
			Name: "MemberListContainer",
			Methods: []Method{
				SetterMethod{
					method: method{
						Comment: Comment{`
							SetSections (re)sets the list of sections to be the
							given slice. Members from the old section list
							should be transferred over to the new section entry
							if the section name's content is the same. Old
							sections that don't appear in the new slice should
							be removed.
						`},
						Name: "SetSections",
					},
					Parameters: []NamedType{
						{Name: "sections", Type: "[]MemberSection"},
					},
				},
				SetterMethod{
					method: method{
						Comment: Comment{`
							SetMember adds or updates (or upsert) a member into
							a section. This operation must not change the
							section's member count. As such, changes should be
							done separately in SetSection. If the section does
							not exist, then the client should ignore this
							member. As such, backends must call SetSections
							first before SetMember on a new section.
						`},
						Name: "SetMember",
					},
					Parameters: []NamedType{
						{"sectionID", "ID"},
						{"member", "ListMember"},
					},
				},
				SetterMethod{
					method: method{
						Comment: Comment{`
							RemoveMember removes a member from a section. If
							neither the member nor the section exists, then the
							client should ignore it.
						`},
						Name: "RemoveMember",
					},
					Parameters: []NamedType{
						{"sectionID", "ID"},
						{"memberID", "ID"},
					},
				},
			},
		}, {
			Comment: Comment{`
				ListMember represents a single member in the member list. This
				is a base interface that may implement more interfaces, such as
				Iconer for the user's avatar.

				Note that the frontend may give everyone an avatar regardless,
				or it may not show any avatars at all.
			`},
			Name: "ListMember",
			Embeds: []EmbeddedInterface{
				{InterfaceName: "User"},
			},
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Status returns the status of the member. The backend
							does not have to show offline members with the
							offline status if it doesn't want to show offline
							menbers at all.
						`},
						Name: "Status",
					},
					Returns: []NamedType{{Type: "Status"}},
				},
				GetterMethod{
					method: method{
						Comment: Comment{`
							Secondary returns the subtext of this member. This
							could be anything, such as a user's custom status or
							away reason.
						`},
						Name: "Secondary",
					},
					Returns: []NamedType{{
						Type: MakeQual("text", "Rich"),
					}},
				},
			},
		}, {
			Comment: Comment{`
				MemberSection represents a member list section. The section
				name's content must be unique among other sections from the same
				list regardless of the rich segments.
			`},
			Name: "MemberSection",
			Embeds: []EmbeddedInterface{
				{InterfaceName: "Identifier"},
				{InterfaceName: "Namer"},
			},
			Methods: []Method{
				GetterMethod{
					method: method{
						Comment: Comment{`
							Total returns the total member count.
						`},
						Name: "Total",
					},
					Returns: []NamedType{{Type: "int"}},
				},
				AsserterMethod{ChildType: "MemberDynamicSection"},
			},
		}, {
			Comment: Comment{`
				MemberDynamicSection represents a dynamically loaded member list
				section. The section behaves similarly to MemberSection, except
				the information displayed will be considered incomplete until
				LoadMore returns false.

				LoadLess can be called by the client to mark chunks as stale,
				which the server can then unsubscribe from.
			`},
			Name: "MemberDynamicSection",
			Methods: []Method{
				IOMethod{
					method: method{
						Comment: Comment{`
							LoadMore is a method which the client can call to
							ask for more members. This method can do IO.

							Clients may call this method on the last section in
							the section slice; however, calling this method on
							any section is allowed. Clients may not call this
							method if the number of members in this section is
							equal to Total.
						`},
						Name: "LoadMore",
					},
					ReturnValue: NamedType{Type: "bool"},
				},
				IOMethod{
					method: method{
						Comment: Comment{`
							LoadLess is a method which the client must call
							after it is done displaying entries that were added
							from calling LoadMore.

							The client can call this method exactly as many
							times as it has called LoadMore. However, false
							should be returned if the client should stop, and
							future calls without LoadMore should still return
							false.
						`},
						Name: "LoadLess",
					},
					ReturnValue: NamedType{Type: "bool"},
				},
			},
		}, {
			Comment: Comment{`
				SendableMessage is the bare minimum interface of a sendable
				message, that is, a message that can be sent with SendMessage().
				This allows the frontend to implement its own message data
				implementation.

				An example of extending this interface is MessageNonce, which is
				similar to IRCv3's labeled response extension or Discord's
				nonces. The frontend could implement this interface and check if
				incoming MessageCreate events implement the same interface.
			`},
			Name: "SendableMessage",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Content"},
					Returns: []NamedType{{Type: "string"}},
				},
				AsserterMethod{ChildType: "Noncer"},
				AsserterMethod{ChildType: "Replier"},
				AsserterMethod{ChildType: "Attacher"},
			},
		}, {
			Comment: Comment{`
				Replier indicates that the message being sent is a reply to
				something. Frontends that support replies can assume that all
				messages in a Sender can be replied to, and the backend can
				choose to do nothing to the replied ID.
			`},
			Name: "Replier",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "ReplyingTo"},
					Returns: []NamedType{{Type: "ID"}},
				},
			},
		}, {
			Comment: Comment{`
				Attacher adds attachments into the message being sent.
			`},
			Name: "Attacher",
			Methods: []Method{
				GetterMethod{
					method:  method{Name: "Attachments"},
					Returns: []NamedType{{Type: "[]MessageAttachment"}},
				},
			},
		}},
	},
}
