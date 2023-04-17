package pkg

import (
	"fmt"
	"net/http"
	"text/template"
)

type Errors struct {
	Message string
	Status  int
}

func Errorhandler(w http.ResponseWriter, status int) {
	html, err := template.ParseFiles("ui/html/error.html")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Internal Server Errror 500")
		return
	}
	var Eror Errors
	Eror.Message = http.StatusText(status) // func StatusText()
	Eror.Status = status
	err = html.Execute(w, Eror)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Internal Server Errror 500")
		return
	}
}
