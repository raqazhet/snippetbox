package main

import (
	"alex/pkg/forms"
	"alex/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	CSRFToken         string
	AuthenticatedUser *models.User
	Flash             string
	Form              *forms.Form
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	CurrentYear       int
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// first we need initialize a new mapto act as the cashe
	cashe := map[string]*template.Template{}
	//Use the filefath.Glob function to get a slice of all filefaths with
	// the extension '.html'. This essentially gives us a slice of all the
	//'page' templates for the application
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	//Loop through the pages one-by-one
	for _, page := range pages {
		//Extract the filename (like 'home.page.tmpl') from the full file path
		//and assign it to the name variable.
		name := filepath.Base(page)
		//Parse the page template file in to a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		//use the Parseglob method to add any 'partial' templates to the
		//template set(in our case, it's just the 'footer' partial at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		cashe[name] = ts
	}
	return cashe, nil

}
