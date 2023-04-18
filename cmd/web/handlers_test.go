package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	//Create a new instance of our application struct. For now, this just
	//contains a couple of mock loggers
	app := &application{
		errprlog: *log.New(io.Discard, "", 0),
		infolog:  *log.New(io.Discard, "", 0),
	}
	//When theb use the httpset.NewTlsServer()function to create a new test
	//server,passing in the value returned by our app.routes()method as the
	//handler for the server. This starts up a HTTPS server which listens on a
	//randomly-chosen port of your local machine for the duration of the test
	//Notice that we defer a call to ts.Close() to shutdown the server when
	//the test finishes.
	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()
	//The network addres that the test server is listening on is contained
	//in the ts.Url field. We can use this along with the ts.Client().Get()
	//method to make a GET /ping request against the rest server. This
	//returns a http.Response struct containing the response
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	//We can then check value of the response status code and body using
	//the same code as before
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
