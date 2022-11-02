package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// frontend serves the "/" route on port 80
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081 (80 before)")
	// err := http.ListenAndServe(":8081", nil)
	// k8s
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

// render renders the frontend templates
func render(w http.ResponseWriter, t string) {

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}
	// partials := []string{
	// 	"./cmd/web/templates/base.layout.gohtml",
	// 	"./cmd/web/templates/header.partial.gohtml",
	// 	"./cmd/web/templates/footer.partial.gohtml",
	// }
	// add t to templateSlice
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	// next add partials to templateSlice too
	// should replace loop with
	// templateSlice = append(templateSlice, partials...) (S1011)go-staticcheck

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	// swarm
	// tmpl, err := template.ParseFiles(templateSlice...)
	tmpl, err := template.ParseFS(templateFS, templateSlice...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}

	data.BrokerURL = os.Getenv("BROKER_URL")
	// k8s outside try
	// data.BrokerURL = "http://localhost:8080"
	// swarm
	// if err := tmpl.Execute(w, nil); err != nil {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
