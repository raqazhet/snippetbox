package main

import (
	"alex/pkg/models/mock"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies
func newTestApplication(t *testing.T) *application {
	//Create a session manager instance, with the same settings as production
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	return &application{
		errprlog: log.New(io.Discard, "", 0),
		infolog:  log.New(io.Discard, "", 0),
		session:  session,
		snippets: &mock.SnippetModel{},
		users:    &mock.UserModel{},
	}
}

// Define a custom testSErver type which anonymously embeds a httptest.Server
// instance.
type testSEerver struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of or custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testSEerver {
	ts := httptest.NewTLSServer(h)
	//Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	//And the cookie jar to the client, so that response cookies are stored
	//and then sent with subsqent requests.
	ts.Client().Jar = jar
	//Disable redirect-following for the client. Essentially this functon
	//is called after a 3xx response is recived by the client, and returning
	//the http.ErrUseLastResponse error forces it to immediately return th
	//received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testSEerver{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body
func (ts *testSEerver) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
