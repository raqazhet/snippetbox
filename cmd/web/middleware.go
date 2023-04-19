package main

import (
	"alex/pkg"
	"alex/pkg/models"
	"context"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Execute our middleware logic here..
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if auserId value exists in the session. If this *isnt
		//present* then call the next handler in the chain as normal
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		//Fetch the details of the current user from the database. If
		//no matching record is found,remove the (invalid) userId from
		//their session and call the next handler in the chain as normal
		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			log.Fatal("middleware", err)
			pkg.Errorhandler(w, http.StatusInternalServerError)

			return
		}
		//Otherwise, we know that the request is coming from a valid,
		//authenticated (logged in)user. we create a new copy of the
		//request with the user information added to the request context, and
		//call the next handler in the chain *using this new copy of the
		//request*.
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Create a deffered function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				//Set a connection: close header on the response
				w.Header().Set("Connection", "close")
				//Call the pkg.Errorhandler helper method to return a 500
				//Internal server response
				pkg.Errorhandler(w, http.StatusInternalServerError)

			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check that the authenticatedUser helper doesn`t return nil
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound) //перенаправления 302
			return
		}
		// otherwise call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRf cookie with
// the secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
