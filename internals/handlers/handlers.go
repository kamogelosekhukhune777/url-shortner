package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/kamogelosekhukhune777/url-shortner/internals/storage"
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

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	//short url from the request path
	shortURL := r.URL.Path[1:]

	mapping, err := h.DB.Get(shortURL)
	if err != nil {
		if err == storage.ErrMappingNotFound {
			http.Error(w, "short url not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch mapping", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, mapping.LongURL, http.StatusFound)
}

func (h *URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ShortURL string `json:"shortURL"`
		LongURL  string `json:"longURL"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "inalid request body", http.StatusBadRequest)
		return
	}

	//check if url checks
	_, err := h.DB.Get(requestBody.ShortURL)
	if err != nil {
		if err == storage.ErrMappingNotFound {
			http.Error(w, "short URL not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch mapping", http.StatusNotFound)
		}
		return
	}

	//update
	mapping := &models.URLMapping{
		ShortURL: requestBody.ShortURL,
		LongURL:  requestBody.LongURL,
	}

	if err := h.DB.Update(mapping); err != nil {
		http.Error(w, "failed to update mapping", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(mapping)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	shortenURL := r.URL.Query().Get("short_url")

	//check if the specified short url exists
	_, err := h.DB.Get(shortenURL)
	if err != nil {
		if err == storage.ErrMappingNotFound {
			http.Error(w, " short URl not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch mapping", http.StatusInternalServerError)
		}
		return
	}

	//delete the mapping
	if err := h.DB.Delete(shortenURL); err != nil {
		http.Error(w, "failed to delete mapping", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func generateShortURL() string {
	shortURL := make([]byte, shortURLLength)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortURL)
}
