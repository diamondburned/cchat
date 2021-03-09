package cchat

//go:generate go run ./cmd/internal/cchat-generator ./
//go:generate go run ./cmd/internal/cchat-empty-gen ./utils/empty/

type authenticateError struct{ error }

func (authenticateError) NextStage() []Authenticator { return nil }

// WrapAuthenticateError wraps the given error to become an AuthenticateError.
// Its NextStage method returns nil. If the given err is nil, then nil is
// returned.
func WrapAuthenticateError(err error) AuthenticateError {
	if err == nil {
		return nil
	}
	return authenticateError{err}
}

// ServerColumn is a convenient function to get the column of a server. It
// returns 0 if server does not implement Columnator.
func ServerColumn(server Server) int {
	var column int
	if columnator := server.AsColumnator(); columnator != nil {
		column = columnator.Column()
	}

	if column < 1 {
		return 0
	}

	return column
}
