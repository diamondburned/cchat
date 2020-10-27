package repository

import (
	"testing"

	"github.com/go-test/deep"
)

const _goComment = `
// The authenticator interface allows for a multistage initial authentication
// API that the backend could use. Multistage is done by calling
// AuthenticateForm then Authenticate again forever until no errors are
// returned.
// 
//    var s *cchat.Session
//    var err error
// 
//    for {
//        // Pseudo-function to render the form and return the results of those
//        // forms when the user confirms it.
//        outputs := renderAuthForm(svc.AuthenticateForm())
// 
//        s, err = svc.Authenticate(outputs)
//        if err != nil {
//            renderError(errors.Wrap(err, "Error while authenticating"))
//            continue // retry
//        }
// 
//        break // success
//    }`

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
