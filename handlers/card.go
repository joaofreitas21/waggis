package handlers

import (
	"encoding/json"
	"net/http"

	"path/filepath"
	"os"

	"github.com/joaofreitas21/waggis/views"
)

type CardItem = views.CardItem

//loads all cards, returns the whole map
func loadCards() (map[string]struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Items       map[string]CardItem `json:"items"`
}, error) {
	file, err := os.Open(filepath.Join("data", "cards.json"))
	if err != nil {
		return nil,err
	}
	defer file.Close()

	var all map[string]struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Items       map[string]CardItem `json:"items"`
	}
	if err := json.NewDecoder(file).Decode(&all); err != nil {
	
		return nil,err
	}
	return all,nil
}

// Loads a single card from the JSON file by key
func ServeCard(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}
	
	all, err := loadCards()

	if err != nil {
		http.Error(w, "malformed cards.json", 500)
		return
	}
	card, ok := all[key]
	if !ok {
		http.Error(w, "card not found", 404)
		return
	}

	// Render using Templ!
	views.Card(card.Title, card.Description, card.Items).Render(r.Context(), w)
	
}

func ServeCardItem(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	id := r.URL.Query().Get("id")

	if key == "" || id == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	all, err := loadCards()

	if err != nil {
		http.Error(w, "malformed cards.json", 500)
		return
	}
	card, ok := all[key]
	if !ok {
		http.Error(w, "card not found", 404)
		return
	}

	item, ok := card.Items[id]
	if !ok {
		http.Error(w, "card item not found", 404)
		return
	}

	views.ItemCard(item.Title, item.Description, item.Github, item.Image).Render(r.Context(), w)

}