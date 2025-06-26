package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joaofreitas21/waggis/views" // Replace 'yourmodule' with your actual module name
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join(".", "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve landing page at "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Landing().Render(r.Context(), w)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}