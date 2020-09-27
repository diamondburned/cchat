package repository

var Main = Repositories{
	"text": {
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
			Fields: []StructField{{
				NamedType: NamedType{"Content", "string"},
			}, {
				Comment: Comment{`
					Segments are optional rich-text segment markers.
				`},
				NamedType: NamedType{"Segments", "[]Segment"},
			}},
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
					RegularMethod: RegularMethod{Name: "Bounds"},
					Returns: []NamedType{
						{Name: "start", Type: "int"},
						{Name: "end", Type: "int"},
					},
				},
			},
		}, {
			Comment: Comment{`
				Linker is a hyperlink format that a segment could implement.
				This implies that the segment should be replaced with a
				hyperlink, similarly to the anchor tag with href being the URL
				and the inner text being the text string.
			`},
			Name:   "Linker",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{Name: "Link"},
					Returns:       []NamedType{{Name: "url", Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Imager implies the segment should be replaced with a (possibly
				inlined) image. Only the starting bound matters, as images
				cannot substitute texts.
			`},
			Name:   "Imager",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Image returns the URL for the image.
						`},
						Name: "Image",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
				GetterMethod{
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
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
			`},
			Name:   "Avatarer",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Avatar returns the URL for the image.
						`},
						Name: "Avatar",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
				GetterMethod{
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
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
				Mentioner implies that the segment can be clickable, and when
				clicked it should open up a dialog containing information from
				MentionInfo().
				
				It is worth mentioning that frontends should assume whatever
				segment that Mentioner highlighted to be the display name of
				that user. This would allow frontends to flexibly layout the
				labels.
			`},
			Name:   "Mentioner",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							MentionInfo returns the popup information of the
							mentioned segment. This is typically user
							information or something similar to that context.
						`},
						Name: "MentionInfo",
					},
					Returns: []NamedType{{Type: "Rich"}},
				},
			},
		}, {
			Comment: Comment{`
				MentionerImage extends Mentioner to give the mentioned object an
				image. This interface allows the frontend to be more flexible
				in layouting. A Mentioner can only implement EITHER
				MentionedImage or MentionedAvatar.
			`},
			Name:   "MentionerImage",
			Embeds: []EmbeddedInterface{{InterfaceName: "Mentioner"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Image returns the mentioned object's image URL.
						`},
						Name: "Image",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				MentionerAvatar extends Mentioner to give the mentioned object
				an avatar.  This interface allows the frontend to be more
				flexible in layouting. A Mentioner can only implement EITHER
				MentionedImage or MentionedAvatar.
			`},
			Name:   "MentionerAvatar",
			Embeds: []EmbeddedInterface{{InterfaceName: "Mentioner"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Avatar returns the mentioned object's avatar URL.
						`},
						Name: "Avatar",
					},
					Returns: []NamedType{{Name: "url", Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Colorer is a text color format that a segment could implement.
				This is to be applied directly onto the text.
			`},
			Name:   "Colorer",
			Embeds: []EmbeddedInterface{{InterfaceName: "Mentioner"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Color returns a 24-bit RGB or 32-bit RGBA color.
						`},
						Name: "Color",
					},
					Returns: []NamedType{{Type: "uint32"}},
				},
			},
		}, {
			Comment: Comment{`
				Attributor is a rich text markup format that a segment could
				implement. This is to be applied directly onto the text.
			`},
			Name:   "Attributor",
			Embeds: []EmbeddedInterface{{InterfaceName: "Mentioner"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{Name: "Attribute"},
					Returns:       []NamedType{{Type: "Attribute"}},
				},
			},
		}, {
			Comment: Comment{`
				Codeblocker is a codeblock that supports optional syntax
				highlighting using the language given. Note that as this is a
				block, it will appear separately from the rest of the paragraph.
				
				This interface is equivalent to Markdown's codeblock syntax.
			`},
			Name:   "Codeblocker",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{Name: "CodeblockLanguage"},
				},
			},
		}, {
			Comment: Comment{`
				Quoteblocker represents a quoteblock that behaves similarly to
				the blockquote HTML tag. The quoteblock may be represented
				typically by an actaul quoteblock or with green arrows prepended
				to each line.
			`},
			Name:   "Quoteblocker",
			Embeds: []EmbeddedInterface{{InterfaceName: "Segment"}},
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
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
	"cchat": {
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
			Name: "ID",
			Type: "string",
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
				NamedType: NamedType{"Text", "text.Rich"},
			}, {
				Comment: Comment{`
					Secondary is the label to be displayed on the second line,
					on the right of Text, or not displayed at all. This should
					be optional. This text may be dimmed out as styling.
				`},
				NamedType: NamedType{"Secondary", "text.Rich"},
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
		}},
		ErrorTypes: []ErrorType{{
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
				Receiver: "err",
				Template: "Error at {err.Key}: {err.Err.Error()}",
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
					RegularMethod: RegularMethod{Name: "ID"},
					Returns:       []NamedType{{Type: "ID"}},
				},
			},
		}, {
			Comment: Comment{`
				Namer requires Name() to return the name of the object.
				Typically, this implies usernames for sessions or service
				names for services.
			`},
			Name: "Namer",
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{Name: "Name"},
					Returns:       []NamedType{{Type: "text.Rich"}},
				},
				AsserterMethod{
					ChildType: "Iconer",
				},
			},
		}, {
			Comment: Comment{`
				Iconer adds icon support into Namer, which in turn is returned
				by other interfaces. Typically, Service would return the service
				logo, Session would return the user's avatar, and Server would
				return the server icon.
				
				For session, the avatar should be the same as the one returned
				by messages sent by the current user.
			`},
			Name: "Iconer",
			Methods: []Method{
				ContainerMethod{
					RegularMethod: RegularMethod{Name: "Icon"},
					HasContext:    true,
					ContainerType: "IconContainer",
					HasStopFn:     true,
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
					RegularMethod: RegularMethod{Name: "Nonce"},
					Returns:       []NamedType{{Type: "string"}},
				},
			},
		}, {
			Comment: Comment{`
				Author is the interface for an identifiable author. The
				interface defines that an author always have an ID and a name.
				
				An example of where this interface is used would be in
				MessageCreate's Author method or embedded in Typer. The returned
				ID may or may not be used by the frontend, but backends must
				guarantee that the Author's ID is in fact a user ID.
				
				The frontend may use the ID to squash messages with the same
				author together.
			`},
			Name: "Author",
			Embeds: []EmbeddedInterface{
				{InterfaceName: "Identifier"},
				{InterfaceName: "Namer"},
			},
		}, {
			Comment: Comment{`
				A service is a complete service that's capable of multiple
				sessions. It has to implement the Authenticate() method, which
				returns an implementation of Authenticator.
				
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
					RegularMethod: RegularMethod{Name: "Authenticate"},
					Returns:       []NamedType{{Type: "Authenticator"}},
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
				The authenticator interface allows for a multistage initial
				authentication API that the backend could use. Multistage is
				done by calling AuthenticateForm then Authenticate again forever
				until no errors are returned.

					var s *cchat.Session
					var err error

					for {
						// Pseudo-function to render the form and return the results of those
						// forms when the user confirms it.
						outputs := renderAuthForm(svc.AuthenticateForm())

						s, err = svc.Authenticate(outputs)
						if err != nil {
							renderError(errors.Wrap(err, "Error while authenticating"))
							continue // retry
						}

						break // success
					}
			`},
			Name: "Authenticator",
			Methods: []Method{
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							AuthenticateForm should return a list of
							authentication entries for the frontend to render.
						`},
						Name: "AuthenticateForm",
					},
					Returns: []NamedType{{Type: "[]AuthenticateEntry"}},
				},
				IOMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Authenticate will be called with a list of values
							with indices correspond to the returned slice of
							AuthenticateEntry.
						`},
						Name: "Authenticate",
					},
					Parameters:  []NamedType{{Type: "[]string"}},
					ReturnValue: NamedType{Type: "Session"},
					ReturnError: true,
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
					RegularMethod: RegularMethod{Name: "RestoreSession"},
					Parameters:    []NamedType{{Type: "map[string]string"}},
					ReturnValue:   NamedType{Type: "Session"},
					ReturnError:   true,
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
					RegularMethod: RegularMethod{Name: "Configuration"},
					ReturnValue:   NamedType{Type: "map[string]string"},
					ReturnError:   true,
				},
				IOMethod{
					RegularMethod: RegularMethod{Name: "SetConfiguration"},
					Parameters:    []NamedType{{Type: "map[string]string"}},
					ReturnError:   true,
				},
			},
		}, {
			Comment: Comment{`
				A session is returned after authentication on the service.
				Session implements Name(), which should return the username
				most of the time. It also implements ID(), which might be
				used by frontends to check against MessageAuthor.ID() and
				other things.
				
				A session can implement SessionSaver, which would allow the
				frontend to save the session into its keyring at any time.
				Whether the keyring is completely secure or not is up to the
				frontend. For a Gtk client, that would be using the GNOME
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
					RegularMethod: RegularMethod{
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
					ReturnError: true,
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
					RegularMethod: RegularMethod{Name: "SaveSession"},
					Returns:       []NamedType{{Type: "map[string]string"}},
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							RunCommand executes the given command, with the
							slice being already split arguments, similar to
							os.Args. The function could return an output stream,
							in which the frontend must display it live and close
							it on EOF.
							
							The function can do IO, and outputs should be
							written to the given io.Writer.
							
							The client should make guarantees that an empty
							string (and thus a zero-length string slice) should
							be ignored.  The backend should be able to assume
							that the argument slice is always length 1 or more.
						`},
						Name: "RunCommand",
					},
					Parameters: []NamedType{
						{Type: "[]string"},
						{Type: "io.Writer"},
					},
					ReturnError: true,
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
				AsserterMethod{ChildType: "Configurator"},
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
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							JoinServer joins a server that's capable of
							receiving messages. The server may not necessarily
							support sending messages.
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Send is called by the frontend to send a message to
							this channel.
						`},
						Name: "Send",
					},
					Parameters: []NamedType{
						{Type: "SendableMessage"},
					},
					ReturnError: true,
				},
				GetterMethod{
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							MessageEditable returns whether or not a message can
							be edited by the client. This method must not do IO.
						`},
						Name: "MessageEditable",
					},
					Parameters: []NamedType{{Name: "id", Type: "ID"}},
					Returns:    []NamedType{{Type: "bool"}},
				},
				GetterMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							RawMessageContent gets the original message text for
							editing. This method must not do IO.
						`},
						Name: "RawMessageContent",
					},
					Parameters:  []NamedType{{Name: "id", Type: "ID"}},
					Returns:     []NamedType{{Type: "string"}},
					ReturnError: true,
				},
				IOMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							EditMessage edits the message with the given ID to
							the given content, which is the edited string from
							RawMessageContent. This method can do IO.
						`},
						Name: "EditMessage",
					},
					Parameters: []NamedType{
						{Name: "id", Type: "ID"},
						{Name: "content", Type: "string"},
					},
					ReturnError: true,
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							MessageActions returns a list of possible actions in
							pretty strings that the frontend will use to
							directly display. This method must not do IO.
							
							The string slice returned can be nil or empty.
						`},
						Name: "Actions",
					},
					Returns: []NamedType{{Type: "[]string"}},
				},
				IOMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							DoAction executes a message action on the given
							messageID, which would be taken from
							MessageHeader.ID(). This method is allowed to do
							IO; the frontend should take care of running it
							asynchronously.
						`},
						Name: "DoAction",
					},
					Parameters: []NamedType{
						{Name: "action", Type: "string"},
						{Name: "id", Type: "ID"},
					},
					ReturnError: true,
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
					RegularMethod: RegularMethod{Name: "Nickname"},
					HasContext:    true,
					ContainerType: "LabelContainer",
					HasStopFn:     true,
				},
			},
		}, {
			Comment: Comment{`
				Backlogger adds message history capabilities into a message
				container. The frontend should typically call this method when
				the user scrolls to the top.
				
				As there is no stop callback, if the backend needs to fetch
				messages asynchronously, it is expected to use the context to
				know when to cancel.
				
				The frontend should usually call this method when the user
				scrolls to the top. It is expected to guarantee not to call
				MessagesBefore more than once on the same ID. This can usually
				be done by deactivating the UI.
				
				Note: Although backends might rely on this context, the frontend
				is still expected to invalidate the given container when the
				channel is changed.
			`},
			Name: "Backlogger",
			Methods: []Method{
				IOMethod{
					RegularMethod: RegularMethod{
						Comment: Comment{`
							MessagesBefore fetches messages before the given
							message ID into the MessagesContainer.
						`},
						Name: "MessagesBefore",
					},
					Parameters: []NamedType{
						{Name: "ctx", Type: "context.Context"},
						{Name: "before", Type: "ID"},
						{Name: "c", Type: "MessagePrepender"},
					},
					ReturnError: true,
				},
			},
		}, {
			Comment: Comment{`
				MemberLister adds a member list into a message server.
			`},
			Name: "MemberLister",
			Methods: []Method{
				ContainerMethod{
					RegularMethod: RegularMethod{
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
				UnreadIndicator adds an unread state API for frontends to use.
			`},
			Name: "UnreadIndicator",
			Methods: []Method{
				ContainerMethod{
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							Typing is called by the client to indicate that the
							user is typing. This function can do IO calls, and
							the client must take care of calling it in a
							goroutine (or an asynchronous queue) as well as
							throttling it to TypingTimeout.
						`},
						Name: "Typing",
					},
					ReturnError: true,
				},
				GetterMethod{
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{
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
					RegularMethod: RegularMethod{Name: "UpdateServer"},
					Parameters:    []NamedType{{Type: "ServerUpdate"}},
				},
			},
		}, {
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
					RegularMethod: RegularMethod{
						Comment: Comment{`
							PreviousID returns the ID of the item before this
							server.
						`},
						Name: "PreviousID",
					},
					Returns: []NamedType{{Type: "ID"}},
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
				SetterMethod{},
			},
		}},
	},
}
