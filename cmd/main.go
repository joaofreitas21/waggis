package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/joaofreitas21/waggis/handlers"
	"github.com/joaofreitas21/waggis/views" 
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env")
	}

	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join(".", "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api/ip", handlers.GetIP)

	// Serve landing page at "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Landing().Render(r.Context(), w)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}