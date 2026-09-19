package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// create an newTemplateData() helper, which returns a pointer to a templateData struct initialized with the current year.
// note that we're not using the  *http.Request parameter here at the moment, but we will do later in the book
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// the serverError writes  an error r message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad Request" when there's a problem with the request that the user sent
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
    // Retrieve the appropriate template set from the cache based on the page
    // name (like 'home.tmpl'). If no entry exists in the cache with the
    // provided name, then create a new error and call the serverError() helper
    // method that we made earlier and return.
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.serverError(w, err)
        return
    }

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write the template to the buffer, instead of straight to http.ResponseWriter.
	// if there's an error, call our serverError() helper and then return
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

    // Write out the provided HTTP status code ('200 OK', '400 Bad Request'
    // etc).
    w.WriteHeader(status)
    
	buf.WriteTo(w)
}
