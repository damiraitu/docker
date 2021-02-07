package main

import (
	"errors"
	"fmt"
	_ "html/template"
	"net/http"
	"se03.com/pkg/forms"
	"se03.com/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

//data := &templateData{Snippets: s}
//files := []string{
//	"./ui/html/home.page.tmpl",
//	"./ui/html/base.layout.tmpl",
//	"./ui/html/footer.partial.tmpl",
//}

//ts, err := template.ParseFiles(files...)
//if err != nil {
//	app.errorLog.Println(err.Error())
//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	return
//}
//err = ts.Execute(w, data)
//if err != nil {
//		app.serverError(w, err)
//	}

//	for _, snippet := range s {
//		fmt.Fprintf(w, "%v\n", snippet)
//	}
//}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippet.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})

}
