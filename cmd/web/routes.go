package main

import (
	"net/http"

	"github.com/bmizerany/pat" //new import
	"github.com/justinas/alice"
)

func (app *application) Routes() http.Handler {
	// Create a middleware chain containning our standard middleware
	//Which will be used for every request our application receives
	standardMiddlewre := alice.New(app.recoverPanic, app.logRequest)
	// Create a new middleware chain containning the middleware specific to
	//our dynamic application routes. For now, this chain will only contain
	//the session middleware but we'll add more to it later
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	//use the noSurf middleware on all our 'dynamic' routes.
	//Add the authenticate() middleware to the chain
	mux := pat.New()

	//Update these routes to use the new dynamic middleware chain followed
	//by the appropriate handler function
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet)) //moved down
	//Add the five new routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))
	//Register the ping handler function as the handler for the GET /ping
	//route.
	mux.Get("/ping", http.HandlerFunc(ping))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddlewre.Then(mux)
}
