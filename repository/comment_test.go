package repository

import (
	"testing"

	"github.com/go-test/deep"
)

const _comment = `
The authenticator interface allows for a multistage initial authentication API
that the backend could use. Multistage is done by calling AuthenticateForm then
Authenticate again forever until no errors are returned.

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
	}`

// Trim away the prefix new line.
var comment = _comment[1:]

func TestComment(t *testing.T) {
	var authenticator = Main["cchat"].Interface("Authenticator")
	var authDoc = authenticator.Comment.GoString()

	if eq := deep.Equal(comment, authDoc); eq != nil {
		t.Fatal("Comment inequality:", eq)
	}
}
