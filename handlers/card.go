package handlers

import (
	"encoding/json"
	"net/http"

	"path/filepath"
	"os"

	"github.com/joaofreitas21/waggis/views"
)

// Loads a single card from the JSON file by key
func ServeCard(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}
	// Open cards.json and decode
	file, err := os.Open(filepath.Join("data", "cards.json"))
	if err != nil {
		http.Error(w, "missing cards.json", 500)
		return
	}
	defer file.Close()

	var all map[string]struct {
		Title       string   `json:"title"`
		Github      string   `json:"github"`
		Description string   `json:"description"`
		Icons       []string `json:"icons"`
	}
	if err := json.NewDecoder(file).Decode(&all); err != nil {
		http.Error(w, "malformed cards.json", 500)
		return
	}
	card, ok := all[key]
	if !ok {
		http.Error(w, "card not found", 404)
		return
	}
	// Render using Templ!
	views.Card(card.Title, card.Description, card.Icons, card.Github).Render(r.Context(), w)
}