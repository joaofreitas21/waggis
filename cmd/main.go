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

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env")
	}

	handlers.LoadCtfMetadata()
	
	fs := http.FileServer(http.Dir(filepath.Join(".", "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	http.HandleFunc("/ctf", handlers.ServeCTFPage)
	http.HandleFunc("/api/ctf/search", handlers.SearchCTFWriteups)
	http.HandleFunc("/ctf/report", handlers.ServeCTFReport)

	http.HandleFunc("/api/ip", handlers.GetIP)

	http.HandleFunc("/card", handlers.ServeCard)

	http.HandleFunc("/api/email", handlers.SendEmail)

	http.HandleFunc("/card/item", handlers.ServeCardItem);

	http.HandleFunc("/email-form/success", handlers.ServeEmailSuccessCard)

	http.HandleFunc("/email-form/error", handlers.ServeEmailErrorCard)

	http.HandleFunc("/email-form", func(w http.ResponseWriter, r *http.Request) {
		views.EmailForm().Render(r.Context(), w)
	})

	// Serve landing page at "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			views.Landing().Render(r.Context(), w)
			return
		} 

		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(r.Context(), w)
		
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Fatal(http.ListenAndServe(":"+port, securityMiddleware(http.DefaultServeMux)))
}