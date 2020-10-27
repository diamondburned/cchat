package cchat

//go:generate go run ./cmd/internal/cchat-generator

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
