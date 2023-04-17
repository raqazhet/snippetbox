package main

import (
	"alex/pkg"
	"alex/pkg/forms"
	"alex/pkg/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		fmt.Println("func home", err)
		pkg.Errorhandler(w, http.StatusInternalServerError)
		return
	}
	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.Paresform() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT for PATCH
	// requests. If there are any errors, we use our pkg.ErrorHandler() helper to
	// a 400 Bad Request response to the user
	err := r.ParseForm()
	if err != nil {
		pkg.Errorhandler(w, http.StatusBadRequest)
		return
	}
	// use the r.PostForm.Get() method to retrive the relevant data fields
	// from the r.PostForm map.
	form := forms.New(r.PostForm)
	form.Required([]string{"title", "content", "expires"})
	form.MaxLength("title", 10)
	form.PermittedValues("expires", []string{"10", "7", "2"})
	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{
			Form: form,
		})
	}
	// Pass the data to the SnippetModel.Insert() method, receiving the Id of the new record back
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		fmt.Println("err:", err)
		pkg.Errorhandler(w, http.StatusInternalServerError)
		return
	}
	// Use the Put()method to add a string value ("You snippet was saved
	// succesfuly!") and the corresponding key ("flash") to the session data
	// Note that if theres no exsiting session for the current user
	// or their session has expired then a new, empty, session for them
	// will automatically be created by the session middleware
	app.session.Put(r, "flash", "Snippet succesfully created!!!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ':id' from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		pkg.Errorhandler(w, http.StatusNotFound)
		return
	}
	// Use the SnippetModel object`s  Get method to retrieve the data for a
	// Specific record based on its ID. If no matching record is found,
	// reurn a 404 Not Found response
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		pkg.Errorhandler(w, http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
	}
	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})
}

// Add the new five new functions
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{Form: forms.New(nil)})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		pkg.Errorhandler(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required([]string{"name", "email", "password"})
	form.MatchesPattern("email", forms.EmailRx)
	form.MinLenght("password", 8)
	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	}
	// Try to create a new user record in the database. If the email already exexted
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Addres is already in use")
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	} else if err != nil {
		fmt.Println("func: signUpuser")
		pkg.Errorhandler(w, http.StatusInternalServerError)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming
	// their signup worked and asking them to login
	app.session.Put(r, "flash", "Your signup was succesful. Please log in")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		pkg.Errorhandler(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentails {
		form.Errors.Add("generic", "Email or password is incorrect")
		app.render(w, r, "login.page.html", &templateData{Form: form})
		return
	} else if err != nil {
		fmt.Println("err", err)
		pkg.Errorhandler(w, http.StatusInternalServerError)
		return
	}
	// Add the id of the current user to the session, so that they are now 'login'
	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the "userId" from the session data so that the user is 'logged out'
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've put
	app.session.Put(r, "flash", "You've been logged out succesfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
