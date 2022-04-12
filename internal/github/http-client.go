package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/afero"
	"golang.org/x/oauth2"
)

var ErrEmptyGitHubAuthToken = errors.New("GitHub authorization token value was read but it is empty")

// GetHTTPClient collect GitHub authentication token information and initialize an appropriate
// client (authenticated or unauthenticated) depending on auth information availability.
func GetHTTPClient(fs afero.Fs) (*http.Client, error) {
	tkloc, err := TokenLocation()
	if err != nil {
		return nil, fmt.Errorf("cannot get GitHub token location: %w", err)
	}

	tk := NewAuthToken(fs, tkloc)

	token, err := tk.AuthToken()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("unexpected failure getting GitHub token: %v", err)
	}

	var hc *http.Client

	if err != nil && errors.Is(err, os.ErrNotExist) {
		// NOTE: the auth token was not present (nor in the env nor in the file) so we
		// initialize an unauthenticated http client and plain recorder.
		hc = &http.Client{}

		return hc, nil
	}

	if token == "" {
		return nil, ErrEmptyGitHubAuthToken
	}

	// NOTE: at this point the token is present, so we initialize a oauth2 http client
	// and the corresponding recorder; this is necessary to pass the authentication credentials
	// to the http client transport.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return tc, nil
}
