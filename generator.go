package cchat

import "context"

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

// CtxCallbacks binds a set of given callbacks to the given context. This is
// useful for disconnecting handlers when the context expires.
func CtxCallbacks(ctx context.Context, fns ...func()) {
	go func() {
		<-ctx.Done()
		for _, fn := range fns {
			fn()
		}
	}()
}
