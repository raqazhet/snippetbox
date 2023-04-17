package main

import (
	"alex/pkg"
	"alex/pkg/models"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		pkg.Errorhandler(w, 500)
		return
	}
	// inintialize a new buffer
	// buf := new(bytes.Buffer)
	//Write the template to the buffer, instead of straight to the
	//http.REsponse.Writer If there's an error, call our serverError helper and
	//return
	// Execute the template set passing in any dynamic data.
	err := ts.Execute(w, app.addDefaultData(data, r))
	if err != nil {
		pkg.Errorhandler(w, 500)
		return
	}

}

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the ponter. Again, we're not using the *http.Request parametr at the
// moment, but we will do later in the book

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	//Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	return td
}

// The authenticateUser method returns the Id of the current user from the
// session, or zero if the request is from an unauthenticated user
func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user

}
