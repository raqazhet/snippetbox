package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid Id", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative Id", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal Id", "/snippet/1.34", http.StatusNotFound, nil},
		{"Empty Id", "/snippet/", http.StatusNotFound, nil},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			code, _, body := ts.get(t, v.urlPath)
			if code != v.wantCode {
				t.Errorf("want %d; got %d", v.wantCode, code)
			}
			if !bytes.Contains(body, v.wantBody) {
				t.Errorf("want body to contain %q", v.wantBody)
			}
		})
	}
}
func TestSignupUser(t *testing.T) {
	//Create the application struct containing our mocked dependencies and set
	//up the test server for running and end-to-end test
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	//Make a Get /user/signup request and then extract the CSRF token from the
	//response body
	_, _, body := ts.get(t, "/user/signup")
	csrftoken := extractCSRFToken(t, body)
	//Log the CSRF token value in our test output
	t.Log(csrftoken)
}
