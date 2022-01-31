package main

import (
	"log"
	"net/http"
	"text/template"
)

var templates = make(map[string]*template.Template)

const ColorGrading = "ColorGrading"

func main() {
	templates[ColorGrading] = loadTemplate(ColorGrading)
	http.HandleFunc("/", handleIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := templates[ColorGrading].Execute(w, struct {
		Title string
	}{
		Title: ColorGrading,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"templates/"+name+".template",
		"templates/_header.template",
		"templates/_footer.template",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}
