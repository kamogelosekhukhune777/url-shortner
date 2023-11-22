package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/kamogelosekhukhune777/url-shortner/pkg/models"
)

const (
	shortURLLength = 8
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLHandler struct {
	DB models.URLRepository
}

func NewURLHandler(db models.URLRepository) *URLHandler {
	return &URLHandler{
		DB: db,
	}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LongURL string `json:"long_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	mapping := &models.URLMapping{
		ShortURL: generateShortURL(),
		LongURL:  requestBody.LongURL,
	}

	if err := h.DB.Save(mapping); err != nil {
		http.Error(w, "failed to save mapping", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(mapping)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *URLHandler) Redirect(w http.Response, r *http.Request) {}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {}

func (h *URLHandler) Update(w http.ResponseWriter, r *http.Request) {}

func generateShortURL() string {
	shortURL := make([]byte, shortURLLength)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortURL)
}
