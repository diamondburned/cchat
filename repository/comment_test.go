package repository

import (
	"testing"

	"github.com/go-test/deep"
)

const _goComment = `
// The authenticator interface allows for a multistage initial authentication
// API that the backend could use. Multistage is done by calling Authenticate
// and check for AuthenticateError's NextStage method.`

// Trim away the prefix new line.
var goComment = _goComment[1:]

func TestComment(t *testing.T) {
	var authenticator = Main[RootPath].Interface("Authenticator")

	t.Run("godoc", func(t *testing.T) {
		godoc := authenticator.Comment.GoString(0)

		if eq := deep.Equal(goComment, godoc); eq != nil {
			t.Fatal("go comment inequality:", eq)
		}
	})
}
