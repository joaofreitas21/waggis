package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/joaofreitas21/waggis/views"
)

type ctfSearchResponse struct {
	Count int    `json:"count"`
	HTML  string `json:"html"`
}


var (
	ctfOnce     sync.Once
	ctfMetadata []views.CtfMetadata
	ctfLoadErr  error
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,        
		extension.Footnote,  
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), 
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(), 
		html.WithXHTML(),
		html.WithUnsafe(), 
	),
)

func getCTFMetadata() ([]views.CtfMetadata, error) {
	ctfOnce.Do(func() {
		f, err := os.Open(filepath.Join("data", "metadata.json"))
		if err != nil {
			ctfLoadErr = err
			return
		}
		defer f.Close()
		ctfLoadErr = json.NewDecoder(f).Decode(&ctfMetadata)
	})
	return ctfMetadata, ctfLoadErr
}


func LoadCtfMetadata() {
	if _, err := getCTFMetadata(); err != nil {
		log.Printf("warning: could not preload CTF metadata: %v", err)
	}
}

func formatDateLabel(dateStr string) string {
	s := strings.TrimSpace(dateStr)
	if s == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return dateStr
	}
	return t.Format("02 Jan 2006")
}

func ServeCTFPage(w http.ResponseWriter, r *http.Request) {
	metadata, err := getCTFMetadata()
	if err != nil {
		http.Error(w, "failed to load CTF metadata", http.StatusInternalServerError)
		return
	}
	views.CTFPage(len(metadata)).Render(r.Context(), w)
}

func SearchCTFWriteups(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	tagLower := strings.ToLower(tag)

	if q == "" && tagLower == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(ctfSearchResponse{Count: 0, HTML: ""})
		return
	}

	metadata, err := getCTFMetadata()
	if err != nil {
		http.Error(w, "failed to load CTF metadata", http.StatusInternalServerError)
		return
	}

	cards := make([]views.CtfCard, 0)

	for _, m := range metadata {
		if !matchesSearch(m, q) || !matchesTag(m, tagLower) {
			continue
		}

		renderTags := make([]views.CtfRender, len(m.Tags))
		for i, t := range m.Tags {
			renderTags[i] = views.CtfRender{
				Name:   t,
				Active: tagLower != "" && strings.ToLower(t) == tagLower,
			}
		}

		cards = append(cards, views.CtfCard{
			Slug:          m.Slug,
			Title:         m.Title,
			Platform:      m.Platform,
			Difficulty:    strings.ToLower(m.Difficulty),
			DateFormatted: formatDateLabel(m.Date),
			Tags:          renderTags,
			Description:   m.Description,
		})
	}

	var buf bytes.Buffer
	if err := views.CTFWriteupCards(cards, tag).Render(r.Context(), &buf); err != nil {
		http.Error(w, "failed to render cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(ctfSearchResponse{
		Count: len(cards),
		HTML:  buf.String(),
	})
}

func matchesSearch(m views.CtfMetadata, q string) bool {
	if q == "" {
		return true
	}
	if strings.Contains(strings.ToLower(m.Title), q) ||
		strings.Contains(strings.ToLower(m.Platform), q) {
		return true
	}
	for _, t := range m.Tags {
		if strings.Contains(strings.ToLower(t), q) {
			return true
		}
	}
	return false
}

func matchesTag(m views.CtfMetadata, tagLower string) bool {
	if tagLower == "" {
		return true
	}
	for _, t := range m.Tags {
		if strings.ToLower(t) == tagLower {
			return true
		}
	}
	return false
}

func ServeCTFReport(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(r.URL.Query().Get("slug"))
	if slug == "" {
		http.Error(w, "missing slug", http.StatusBadRequest)
		return
	}
 
	if strings.ContainsAny(slug, "/\\.") {
		http.Error(w, "invalid slug", http.StatusBadRequest)
		return
	}
 
	src, err := os.ReadFile(filepath.Join("data", slug, "report.md"))
	if err != nil {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}
 
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		http.Error(w, "failed to render report", http.StatusInternalServerError)
		return
	}
 
	views.CTFReportPage(slug, buf.String()).Render(r.Context(), w)
}
