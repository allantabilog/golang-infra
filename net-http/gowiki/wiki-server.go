package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// data structures for a wiki
type Page struct {
	Title string 
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/view/"):]
	p, err := loadPage(fileName)
	if err != nil {
		http.Redirect(w, r, "/edit/"+fileName, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/edit/"):]
	p, err := loadPage(fileName)
	if err != nil {
		p = &Page{Title: fileName}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: fileName, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+fileName, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}