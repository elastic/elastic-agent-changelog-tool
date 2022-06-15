package githubtest

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// getHttpClient instantiate a http.Client backed by a recorder.Recorder to be used in testing
// scenarios.
// As GitHub may require authentication, this function leverages AuthToken functionality to load
// GitHub token from env or file; when missing instantiate an unauthenticated client.
// NOTE: always remember to call (recorder.Recorder).Stop() in your test case.
func GetHttpClient(t *testing.T) (*recorder.Recorder, *http.Client) {
	t.Helper()
	var err error

	testFs := afero.NewOsFs()
	tkloc, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")

	tk := github.NewAuthToken(testFs, tkloc)
	token, err := tk.AuthToken()

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("unexpected failure getting GitHub token: %v", err)
	}

	var rec *recorder.Recorder
	var hc *http.Client

	if err != nil && errors.Is(err, os.ErrNotExist) {
		// NOTE: the auth token was not present (nor in the env nor in the file) so we
		// initialize an unauthenticated http client and plain recorder.
		rec, err = recorder.New(path.Join("testdata", "fixtures", t.Name()))
		hc = &http.Client{
			Transport: rec,
			Timeout:   2 * time.Second,
		}
		require.NoError(t, err)

		return rec, hc
	}

	if token == "" {
		log.Fatal("GitHub authorization token value was read but it is empty")
	}

	// NOTE: at this point the token is present, so we initialize a oauth2 http client
	// and the corresponding recorder; this is necessary to pass the authentication credentials
	// to the http client transport.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tr := &oauth2.Transport{
		Base:   http.DefaultTransport,
		Source: oauth2.ReuseTokenSource(nil, ts),
	}

	rec, err = recorder.NewAsMode(path.Join("testdata", "fixtures", t.Name()), recorder.ModeReplaying, tr)
	require.NoError(t, err)

	// filter out dynamic & sensitive data/headers.
	// NOTE: your test code will continue to see (and use) the real access token and
	// it is redacted before the recorded interactions are saved.
	rec.AddSaveFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		// i.Request.Headers["Authorization"] = []string{"Basic REDACTED"}

		return nil
	})

	hc = &http.Client{
		Transport: rec,
		Timeout:   2 * time.Second,
	}

	return rec, hc
}
